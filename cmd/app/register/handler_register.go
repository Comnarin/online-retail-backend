package registry

import (
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/handler/admin"
	"github.com/retail/backend/internal/handler/auth"
	"github.com/retail/backend/internal/handler/liff"
	"github.com/retail/backend/internal/handler/superadmin"
	"github.com/retail/backend/internal/platform/storage"
)

type Handlers struct {
	Auth                handler.IAuthHandler
	Tenant              handler.ITenantHandler
	Dashboard           handler.IDashboardHandler
	Appearance          handler.IAppearanceHandler
	Product             handler.IProductHandler
	Order               handler.IOrderHandler
	Customer            handler.ICustomerHandler
	Coupon              handler.ICouponHandler
	Membership          handler.IMembershipHandler
	LiffProduct         handler.ILiffProductHandler
	LiffOrder           handler.ILiffOrderHandler
	LiffCustomer        handler.ILiffCustomerHandler
	LiffAddress         handler.ILiffAddressHandler
	LiffPayment         handler.ILiffPaymentHandler
	LiffCart            handler.ILiffCartHandler
	AdminPaymentMethod  handler.IAdminPaymentMethodHandler
	Transaction         handler.ITransactionHandler
	Storage             *storage.MinIOClient
}

func NewHandlers(services *Services, minioClient *storage.MinIOClient) *Handlers {
	return &Handlers{
		Auth:               auth.NewAuthHandler(services.Auth),
		Tenant:             superadmin.NewTenantHandler(services.Tenant, services.Auth),
		Dashboard:          admin.NewDashboardHandler(services.Dashboard),
		Appearance:         admin.NewAppearanceHandler(services.Dashboard),
		Product:            admin.NewProductHandler(services.Product, minioClient),
		Order:              admin.NewOrderHandler(services.Order),
		Customer:           admin.NewCustomerHandler(services.Customer),
		Coupon:             admin.NewCouponHandler(services.Coupon),
		Membership:         admin.NewMembershipHandler(services.Membership),
		LiffProduct:        liff.NewLiffProductHandler(services.Product),
		LiffOrder:          liff.NewLiffOrderHandler(services.Order, minioClient),
		LiffCustomer:       liff.NewLiffCustomerHandler(services.Customer, services.Membership),
		LiffAddress:        liff.NewLiffAddressHandler(services.Address),
		LiffPayment:        liff.NewLiffPaymentHandler(services.PaymentMethod, services.CustomerPayment),
		LiffCart:           liff.NewLiffCartHandler(services.Cart),
		Transaction:        admin.NewTransactionHandler(services.Transaction),
		AdminPaymentMethod: admin.NewPaymentMethodHandler(services.PaymentMethod, minioClient),
		Storage:            minioClient,
	}
}

