package registry

import (
	"github.com/retail/backend/internal/platform/cache"
	"github.com/retail/backend/internal/platform/config"
	service "github.com/retail/backend/internal/services"
	jwtpkg "github.com/retail/backend/pkg/jwt"
	"gorm.io/gorm"
)

// Services holds all service instances.
type Services struct {
	Auth       service.IAuthService
	Tenant     service.ITenantService
	Product    service.IProductService
	Order      service.IOrderService
	Coupon     service.ICouponService
	Membership service.IMembershipService
	Customer   service.ICustomerService
	Dashboard  service.IDashboardService
	Transaction     service.ITransactionService
	Address         service.ICustomerAddressService
	PaymentMethod   service.IPaymentMethodService
	CustomerPayment service.ICustomerPaymentMethodService
	Cart            service.ICartService
}

func NewServices(db *gorm.DB, repos *Repositories, jwtMgr *jwtpkg.Manager, redis *cache.RedisClient, cfg *config.Config) *Services {
	return &Services{
		Auth:        service.NewAuthService(repos.User, repos.Customer, repos.Membership, repos.Point, repos.Tenant, jwtMgr, db, cfg.LIFF.ChannelID),
		Tenant:      service.NewTenantService(repos.Tenant, repos.User),
		Product:     service.NewProductService(repos.Product),
		Order:       service.NewOrderService(repos.Order, repos.Product, repos.Point, repos.Coupon, repos.Membership, repos.Transaction, repos.Address, repos.PaymentMethod, repos.Tenant, repos.Cart, redis, db),
		Coupon:      service.NewCouponService(repos.Coupon),
		Membership:  service.NewMembershipService(repos.Membership, repos.Point),
		Customer:    service.NewCustomerService(repos.Customer),
		Dashboard:   service.NewDashboardService(repos.Dashboard, repos.Tenant),
		Transaction: service.NewTransactionService(repos.Transaction),
		Address:     service.NewCustomerAddressService(repos.Address),
		PaymentMethod: service.NewPaymentMethodService(repos.PaymentMethod),
		CustomerPayment: service.NewCustomerPaymentMethodService(repos.CustomerPayment),
		Cart:            service.NewCartService(repos.Cart, repos.Product),
	}
}
