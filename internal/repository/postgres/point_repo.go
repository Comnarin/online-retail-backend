package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type pointRepo struct{ db *gorm.DB }

func NewPointRepository(db *gorm.DB) repository.IPointRepository {
	return &pointRepo{db: db}
}

func (r *pointRepo) WithTx(tx *gorm.DB) repository.IPointRepository {
	return &pointRepo{db: tx}
}

func (r *pointRepo) GetBalance(ctx context.Context, tenantID, customerID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.PointTransaction{}).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Select("COALESCE(SUM(points), 0)").Scan(&total).Error
	return int(total), err
}

func (r *pointRepo) AddTransaction(ctx context.Context, tx *domain.PointTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *pointRepo) ListTransactions(ctx context.Context, tenantID, customerID uuid.UUID, opts model.ListOptions) ([]domain.PointTransaction, int64, error) {
	var txs []domain.PointTransaction
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.PointTransaction{}).Where("tenant_id = ?", tenantID)
	if customerID != uuid.Nil {
		q = q.Where("customer_id = ?", customerID)
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&txs).Error
	return txs, total, err
}
