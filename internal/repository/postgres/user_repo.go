package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) repository.IUserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.User{}).Where("tenant_id = ?", tenantID)
	if opts.Search != "" {
		q = q.Where("name ILIKE ? OR email ILIKE ?", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}
	
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (opts.Page - 1) * opts.Limit
	err := q.Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *userRepo) Update(ctx context.Context, tenantID uuid.UUID, user *domain.User) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ? AND tenant_id = ?", user.ID, tenantID).
		Updates(user).Error
}

func (r *userRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.User{}).Error
}
