package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type customerRepo struct{ db *gorm.DB }

func NewCustomerRepository(db *gorm.DB) repository.ICustomerRepository {
	return &customerRepo{db: db}
}

func (r *customerRepo) WithTx(tx *gorm.DB) repository.ICustomerRepository {
	return &customerRepo{db: tx}
}

func (r *customerRepo) Create(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepo) GetByLineUserID(ctx context.Context, tenantID uuid.UUID, lineUserID string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND line_user_id = ?", tenantID, lineUserID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var total int64

	// Base query with subqueries for aggregates
	subQueryOrders := r.db.Model(&domain.Order{}).
		Select("customer_id, COUNT(id) as order_count, SUM(total) as total_spent").
		Where("tenant_id = ? AND status != 'cancelled'", tenantID).
		Group("customer_id")

	q := r.db.WithContext(ctx).Model(&domain.Customer{}).
		Select("customers.*, COALESCE(ord.order_count, 0) as order_count, COALESCE(ord.total_spent, 0) as total_spent").
		Joins("LEFT JOIN (?) ord ON ord.customer_id = customers.id", subQueryOrders).
		Where("customers.tenant_id = ?", tenantID)

	if opts.Search != "" {
		q = q.Where("customers.name ILIKE ? OR customers.phone ILIKE ? OR customers.email ILIKE ?", 
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}
	
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (opts.Page - 1) * opts.Limit
	err := q.Offset(offset).Limit(opts.Limit).
		Order("customers.created_at DESC").
		Preload("Membership.Tier").
		Find(&customers).Error

	// Map TotalPoints from Membership after loading
	for i := range customers {
		if customers[i].Membership != nil {
			customers[i].TotalPoints = int64(customers[i].Membership.Points)
		}
	}

	return customers, total, err
}

func (r *customerRepo) Update(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Model(&domain.Customer{}).
		Where("id = ? AND tenant_id = ?", customer.ID, tenantID).
		Updates(customer).Error
}

func (r *customerRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.Customer{}).Error
}
