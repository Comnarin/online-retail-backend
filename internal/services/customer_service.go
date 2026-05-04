package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type CustomerService struct {
	customerRepo repository.ICustomerRepository
}

func NewCustomerService(customerRepo repository.ICustomerRepository) ICustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

func (s *CustomerService) ListCustomers(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Customer, int64, error) {
	return s.customerRepo.ListByTenant(ctx, tenantID, opts)
}

func (s *CustomerService) GetCustomer(ctx context.Context, tenantID, id uuid.UUID) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, tenantID, id)
}
