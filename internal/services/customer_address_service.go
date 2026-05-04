package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
)

type customerAddressService struct {
	repo repository.ICustomerAddressRepository
}

func NewCustomerAddressService(repo repository.ICustomerAddressRepository) ICustomerAddressService {
	return &customerAddressService{repo: repo}
}

func (s *customerAddressService) List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerAddress, error) {
	return s.repo.List(ctx, tenantID, customerID)
}

func (s *customerAddressService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerAddress, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *customerAddressService) Create(ctx context.Context, addr *domain.CustomerAddress, asDefault bool) error {
	if err := s.repo.Create(ctx, addr); err != nil {
		return err
	}
	if asDefault {
		return s.repo.SetDefault(ctx, addr.TenantID, addr.CustomerID, addr.ID)
	}
	return nil
}

func (s *customerAddressService) Update(ctx context.Context, addr *domain.CustomerAddress, asDefault bool) error {
	if err := s.repo.Update(ctx, addr); err != nil {
		return err
	}
	if asDefault {
		return s.repo.SetDefault(ctx, addr.TenantID, addr.CustomerID, addr.ID)
	}
	return nil
}

func (s *customerAddressService) Delete(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	// Security check: ensure address belongs to the customer
	addr, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if addr.CustomerID != customerID {
		return repository.ErrForbidden
	}

	return s.repo.Delete(ctx, tenantID, id)
}
