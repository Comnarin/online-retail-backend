package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
)

type paymentMethodService struct {
	repo repository.IPaymentMethodRepository
}

func NewPaymentMethodService(repo repository.IPaymentMethodRepository) IPaymentMethodService {
	return &paymentMethodService{repo: repo}
}

func (s *paymentMethodService) List(ctx context.Context, tenantID uuid.UUID) ([]domain.PaymentMethod, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *paymentMethodService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.PaymentMethod, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *paymentMethodService) Update(ctx context.Context, tenantID uuid.UUID, pm *domain.PaymentMethod) error {
	return s.repo.Update(ctx, tenantID, pm)
}
