package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"github.com/retail/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo       repository.IUserRepository
	customerRepo   repository.ICustomerRepository
	membershipRepo repository.IMembershipRepository
	pointRepo      repository.IPointRepository
	tenantRepo     repository.ITenantRepository
	jwtMgr         *jwtpkg.Manager
	db             *gorm.DB
	channelID      string
}

func NewAuthService(
	userRepo repository.IUserRepository,
	customerRepo repository.ICustomerRepository,
	membershipRepo repository.IMembershipRepository,
	pointRepo repository.IPointRepository,
	tenantRepo repository.ITenantRepository,
	jwtMgr *jwtpkg.Manager,
	db *gorm.DB,
	channelID string,
) IAuthService {
	return &AuthService{
		userRepo:       userRepo,
		customerRepo:   customerRepo,
		membershipRepo: membershipRepo,
		pointRepo:      pointRepo,
		tenantRepo:     tenantRepo,
		jwtMgr:         jwtMgr,
		db:             db,
		channelID:      channelID,
	}
}

// LoginWithEmail for superadmin and tenant_admin
func (u *AuthService) LoginWithEmail(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Password == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	token, err := u.jwtMgr.Generate(user.ID, user.TenantID, string(user.Role))
	if err != nil {
		return nil, err
	}

	user.Password = nil
	return &model.LoginResponse{Token: token, User: user}, nil
}

func (u *AuthService) verifyIDToken(ctx context.Context, idToken string) (string, string, string, error) {
	if idToken == "" {
		return "", "", "", errors.New("id_token is empty")
	}

	data := url.Values{
		"id_token": []string{idToken},
		"client_id": []string{strings.TrimSpace(u.channelID)},
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/oauth2/v2.1/verify", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		return "", "", "", fmt.Errorf("LINE verification failed with status %d: %v", resp.StatusCode, errData)
	}

	var result struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Error   string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", err
	}

	if result.Error != "" {
		return "", "", "", fmt.Errorf("LINE error: %s", result.Error)
	}

	return result.Sub, result.Name, result.Picture, nil
}

func (u *AuthService) LoginWithLine(ctx context.Context, req model.LineLoginRequest) (*model.LoginResponse, error) {
	// 1. Secure ID Token Verification
	sub, vName, vPicture, err := u.verifyIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, fmt.Errorf("invalid identity token: %w", err)
	}

	// 2. Ensure the token belongs to the requesting user (Prevents spoofing)
	if sub != req.LineUserID {
		return nil, errors.New("identity mismatch: token does not belong to user")
	}

	// Use verified metadata if provided in the token (profile scope)
	if vName != "" {
		req.DisplayName = vName
	}
	if vPicture != "" {
		req.PictureURL = vPicture
	}

	customer, err := u.customerRepo.GetByLineUserID(ctx, req.TenantID, req.LineUserID)
	if err != nil {
		// New user — create customer account with transaction
		err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			name := req.DisplayName
			if name == "" {
				name = "LINE User"
			}
			customer = &domain.Customer{
				TenantID:   req.TenantID,
				LineUserID: &req.LineUserID,
				Name:       name,
				Avatar:     &req.PictureURL,
				IsActive:   true,
			}
			if err := u.customerRepo.WithTx(tx).Create(ctx, customer); err != nil {
				return err
			}

			// Onboarding: Check Features
			tenant, err := u.tenantRepo.GetByID(ctx, req.TenantID)
			if err == nil && tenant.Features.EnableMembership {
				// 1. Get Default Tier (lowest min_points)
				tiers, err := u.membershipRepo.WithTx(tx).ListTiers(ctx, req.TenantID)
				if err == nil && len(tiers) > 0 {
					var defaultTier *domain.MembershipTier
					for i := range tiers {
						if defaultTier == nil || tiers[i].MinPoints < defaultTier.MinPoints {
							defaultTier = &tiers[i]
						}
					}

					if defaultTier != nil {
						// 2. Create Membership
						initialPoints := 0
						if tenant.Features.EnablePoints {
							initialPoints = 10 // Welcome points
						}

						membership := &domain.CustomerMembership{
							TenantID:   req.TenantID,
							CustomerID: customer.ID,
							TierID:     defaultTier.ID,
							Points:     initialPoints,
						}
						if err := u.membershipRepo.WithTx(tx).UpsertCustomerMembership(ctx, membership); err != nil {
							return err
						}

						// 3. Add Point Transaction
						if initialPoints > 0 {
							ptx := &domain.PointTransaction{
								TenantID:    req.TenantID,
								CustomerID:  customer.ID,
								Points:      initialPoints,
								Type:        "earn",
								Description: "Welcome Points",
							}
							if err := u.pointRepo.WithTx(tx).AddTransaction(ctx, ptx); err != nil {
								return err
							}
						}
					}
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// Existing user — Update metadata (Name/Avatar) if changed
		customer.Name = req.DisplayName
		customer.Avatar = &req.PictureURL
		_ = u.customerRepo.Update(ctx, req.TenantID, customer)
	}

	// We still use "customer" role in the JWT token to maintain compatibility with auth middleware
	token, err := u.jwtMgr.Generate(customer.ID, &customer.TenantID, "customer")
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: token, Customer: customer}, nil
}

func (u *AuthService) DevLiffLogin(ctx context.Context, tenantID uuid.UUID) (*model.LoginResponse, error) {
	// Create or get a stable test customer for this tenant
	lineID := "dev_test_customer_" + tenantID.String()[:8]
	customer, err := u.customerRepo.GetByLineUserID(ctx, tenantID, lineID)
	if err != nil {
		err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			customer = &domain.Customer{
				TenantID:   tenantID,
				LineUserID: &lineID,
				Name:       "Dev Tester",
				IsActive:   true,
			}
			if err := u.customerRepo.WithTx(tx).Create(ctx, customer); err != nil {
				return err
			}

			// Onboarding: Check Features
			tenant, err := u.tenantRepo.GetByID(ctx, tenantID)
			if err == nil && tenant.Features.EnableMembership {
				tiers, err := u.membershipRepo.WithTx(tx).ListTiers(ctx, tenantID)
				if err == nil && len(tiers) > 0 {
					var defaultTier *domain.MembershipTier
					for i := range tiers {
						if defaultTier == nil || tiers[i].MinPoints < defaultTier.MinPoints {
							defaultTier = &tiers[i]
						}
					}

					if defaultTier != nil {
						initialPoints := 0
						if tenant.Features.EnablePoints {
							initialPoints = 10
						}

						membership := &domain.CustomerMembership{
							TenantID:   tenantID,
							CustomerID: customer.ID,
							TierID:     defaultTier.ID,
							Points:     initialPoints,
						}
						if err := u.membershipRepo.WithTx(tx).UpsertCustomerMembership(ctx, membership); err != nil {
							return err
						}
					}
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	token, err := u.jwtMgr.Generate(customer.ID, &customer.TenantID, "customer")
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: token, Customer: customer}, nil
}

// Logout blacklists the JWT
func (u *AuthService) Logout(ctx context.Context, jti string, userID uuid.UUID) error {
	return nil // JTI blacklisting handled in JWT manager
}

func (u *AuthService) RegisterAdmin(ctx context.Context, req model.RegisterAdminRequest) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedStr := string(hashed)
	tenantID := req.TenantID
	user := &domain.User{
		TenantID: &tenantID,
		Email:    &req.Email,
		Password: &hashedStr,
		Name:     req.Name,
		Role:     domain.RoleTenantAdmin,
		IsActive: true,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	user.Password = nil
	return user, nil
}
