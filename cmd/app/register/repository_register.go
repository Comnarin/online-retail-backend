package registry

import (
	"github.com/retail/backend/internal/repository"
	pgRepo "github.com/retail/backend/internal/repository/postgres"
	"gorm.io/gorm"
)

// Repositories holds all repository instances.
type Repositories struct {
	Tenant     repository.ITenantRepository
	User       repository.IUserRepository
	Product    repository.IProductRepository
	Order      repository.IOrderRepository
	Coupon     repository.ICouponRepository
	Point      repository.IPointRepository
	Membership      repository.IMembershipRepository
	Dashboard       repository.IDashboardRepository
	Transaction     repository.ITransactionRepository
	Customer        repository.ICustomerRepository
	Address         repository.ICustomerAddressRepository
	PaymentMethod   repository.IPaymentMethodRepository
	CustomerPayment repository.ICustomerPaymentMethodRepository
	Cart            repository.ICartRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Tenant:      pgRepo.NewTenantRepository(db),
		User:        pgRepo.NewUserRepository(db),
		Product:     pgRepo.NewProductRepository(db),
		Order:       pgRepo.NewOrderRepository(db),
		Coupon:      pgRepo.NewCouponRepository(db),
		Point:       pgRepo.NewPointRepository(db),
		Membership:      pgRepo.NewMembershipRepository(db),
		Dashboard:       pgRepo.NewDashboardRepository(db),
		Transaction:     pgRepo.NewTransactionRepository(db),
		Customer:        pgRepo.NewCustomerRepository(db),
		Address:         pgRepo.NewCustomerAddressRepository(db),
		PaymentMethod:   pgRepo.NewPaymentMethodRepository(db),
		CustomerPayment: pgRepo.NewCustomerPaymentMethodRepository(db),
		Cart:            pgRepo.NewCartRepository(db),
	}
}
