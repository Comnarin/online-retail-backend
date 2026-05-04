package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
)

type customerPaymentMethodService struct {
	repo repository.ICustomerPaymentMethodRepository
}

func NewCustomerPaymentMethodService(repo repository.ICustomerPaymentMethodRepository) ICustomerPaymentMethodService {
	return &customerPaymentMethodService{repo: repo}
}

func (s *customerPaymentMethodService) List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerPaymentMethod, error) {
	return s.repo.List(ctx, tenantID, customerID)
}

func (s *customerPaymentMethodService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerPaymentMethod, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *customerPaymentMethodService) Create(ctx context.Context, cpm *domain.CustomerPaymentMethod, asDefault bool) error {
	if err := s.repo.Create(ctx, cpm); err != nil {
		return err
	}
	if asDefault {
		return s.repo.SetDefault(ctx, cpm.TenantID, cpm.CustomerID, cpm.ID)
	}
	return nil
}

func (s *customerPaymentMethodService) Update(ctx context.Context, cpm *domain.CustomerPaymentMethod, asDefault bool) error {
	if err := s.repo.Update(ctx, cpm); err != nil {
		return err
	}
	if asDefault {
		return s.repo.SetDefault(ctx, cpm.TenantID, cpm.CustomerID, cpm.ID)
	}
	return nil
}

func (s *customerPaymentMethodService) Delete(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	// Security check
	method, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if method.CustomerID != customerID {
		return repository.ErrForbidden
	}
	return s.repo.Delete(ctx, tenantID, id)
}
