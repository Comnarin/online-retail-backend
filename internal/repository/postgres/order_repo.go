package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type orderRepo struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) repository.IOrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Generate order number in backend: {CODE}{YYYYMMDD}{RAND5}
		var tenant domain.Tenant
		if err := tx.Where("id = ?", order.TenantID).First(&tenant).Error; err != nil {
			return err
		}
		
		prefix := tenant.OrderCode
		if prefix == "" {
			prefix = "ORD"
		}
		
		date := time.Now().Format("20060102")
		random := uuid.New().String()[:5] // First 5 chars of a UUID segment
		order.OrderNumber = fmt.Sprintf("%s%s%s", prefix, date, random)

		// Set items and ensure relations are created
		order.Items = items
		return tx.Create(order).Error
	})
}

func (r *orderRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").Preload("Items.Product").Preload("Customer").Preload("ShippingAddress").
		Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error
	return &order, err
}

func (r *orderRepo) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").Preload("Customer").Preload("ShippingAddress").
		Where("order_number = ?", orderNumber).First(&order).Error
	return &order, err
}

func (r *orderRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Order{}).Where("tenant_id = ?", tenantID)
	if opts.Search != "" {
		q = q.Where("order_number ILIKE ?", "%"+opts.Search+"%")
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Preload("Items").Preload("Customer").Preload("ShippingAddress").Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *orderRepo) ListByCustomer(ctx context.Context, customerID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Order{}).Where("customer_id = ?", customerID)
	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Preload("Items").Preload("Items.Product").Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *orderRepo) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("status", status).Error
}

func (r *orderRepo) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepo) WithTx(tx *gorm.DB) repository.IOrderRepository {
	return &orderRepo{db: tx}
}
