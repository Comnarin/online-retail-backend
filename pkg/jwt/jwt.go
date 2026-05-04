package jwtpkg

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/platform/cache"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID *uuid.UUID `json:"tenant_id"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret      []byte
	expiry      time.Duration
	refreshDays int
	redis       *cache.RedisClient
}

func NewManager(secret string, expiryHours, refreshDays int, redis *cache.RedisClient) *Manager {
	return &Manager{
		secret:      []byte(secret),
		expiry:      time.Duration(expiryHours) * time.Hour,
		refreshDays: refreshDays,
		redis:       redis,
	}
}

func (m *Manager) Generate(userID uuid.UUID, tenantID *uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// Blacklist adds the JWT ID to Redis (logout / revocation)
func (m *Manager) Blacklist(ctx context.Context, jti string, expiry time.Time) error {
	ttl := time.Until(expiry)
	return m.redis.Set(ctx, "blacklist:"+jti, "1", ttl)
}

func (m *Manager) IsBlacklisted(ctx context.Context, jti string) bool {
	exists, _ := m.redis.Exists(ctx, "blacklist:"+jti)
	return exists
}
