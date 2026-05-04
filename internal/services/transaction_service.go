package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type TransactionService struct {
	repo repository.ITransactionRepository
}

func NewTransactionService(repo repository.ITransactionRepository) ITransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Transaction, int64, error) {
	return s.repo.List(ctx, tenantID, opts)
}
