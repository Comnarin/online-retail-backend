package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type transactionRepo struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) repository.ITransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepo) Update(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *transactionRepo) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Transaction, int64, error) {
	var txs []domain.Transaction
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if opts.Search != "" {
		query = query.Joins("LEFT JOIN orders ON orders.id = transactions.order_id").
			Where("orders.order_number ILIKE ?", "%"+opts.Search+"%")
	}

	if err := query.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Order").
		Preload("Customer").
		Offset((opts.Page - 1) * opts.Limit).
		Limit(opts.Limit).
		Order("transactions.created_at DESC").
		Find(&txs).Error

	return txs, total, err
}

func (r *transactionRepo) WithTx(tx *gorm.DB) repository.ITransactionRepository {
	return &transactionRepo{db: tx}
}
