package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
)

type IAuthService interface {
	LoginWithEmail(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	LoginWithLine(ctx context.Context, req model.LineLoginRequest) (*model.LoginResponse, error)
	DevLiffLogin(ctx context.Context, tenantID uuid.UUID) (*model.LoginResponse, error)
	Logout(ctx context.Context, jti string, userID uuid.UUID) error
	RegisterAdmin(ctx context.Context, req model.RegisterAdminRequest) (*domain.User, error)
}

type IDashboardService interface {
	GetStats(ctx context.Context, tenantID uuid.UUID) (*model.DashboardStats, error)
	GetRecentOrders(ctx context.Context, tenantID uuid.UUID, limit int) ([]model.RecentOrderRow, error)
	GetAppearance(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error)
	UpdateAppearance(ctx context.Context, tenantID uuid.UUID, name, orderCode string, appearance domain.TenantAppearance, features domain.TenantFeature) error
}

type ICustomerService interface {
	ListCustomers(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Customer, int64, error)
	GetCustomer(ctx context.Context, tenantID, id uuid.UUID) (*domain.Customer, error)
}

type ITenantService interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	List(ctx context.Context, opts model.ListOptions) ([]domain.Tenant, int64, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	UpdateFeatures(ctx context.Context, id uuid.UUID, features domain.TenantFeature) error
	UpdateAppearance(ctx context.Context, id uuid.UUID, appearance domain.TenantAppearance) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAdmins(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.User, int64, error)
	DeleteAdmin(ctx context.Context, tenantID, adminID uuid.UUID) error
}

type ICouponService interface {
	Create(ctx context.Context, coupon *domain.Coupon) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error)
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Coupon, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, coupon *domain.Coupon) error
	Toggle(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error)
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
}

type IMembershipService interface {
	CreateTier(ctx context.Context, tier *domain.MembershipTier) error
	ListTiers(ctx context.Context, tenantID uuid.UUID) ([]domain.MembershipTier, error)
	UpdateTier(ctx context.Context, tenantID uuid.UUID, tier *domain.MembershipTier) error
	DeleteTier(ctx context.Context, tenantID, id uuid.UUID) error
	GetCustomerMembership(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.CustomerMembership, error)
	ListTransactions(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, opts model.ListOptions) ([]domain.PointTransaction, int64, error)
}

type IProductService interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *model.CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Product, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req *model.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	CreateCategory(ctx context.Context, cat *domain.ProductCategory) error
	ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.ProductCategory, error)
	UpdateCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, cat *domain.ProductCategory) error
	DeleteCategory(ctx context.Context, tenantID, id uuid.UUID) error
	AddImage(ctx context.Context, image *domain.ProductImage) error
	RemoveImage(ctx context.Context, tenantID, productID, imageID uuid.UUID) error
}

type CreateOrderRequest struct {
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Items      []struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	} `json:"items"`
	CouponCode      string              `json:"coupon_code"`
	PointsToRedeem  int                 `json:"points_to_redeem"`
	ShippingAddress domain.OrderAddress `json:"shipping_address"`
	SaveAddress     bool                `json:"save_address"`
	AddressLabel    string              `json:"address_label"`
	PaymentMethodID uuid.UUID           `json:"payment_method_id"`
	Note            string              `json:"note"`
}

type IOrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, tenantID, id uuid.UUID) (*domain.Order, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error)
	UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.OrderStatus) error
	UploadSlip(ctx context.Context, tenantID, orderID, customerID uuid.UUID, slipURL string) error
	ConfirmPayment(ctx context.Context, tenantID, orderID uuid.UUID) error
	RejectPayment(ctx context.Context, tenantID, orderID uuid.UUID) error
	CheckExpiration(ctx context.Context, orderID uuid.UUID) error
	GetActivePendingOrder(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Order, error)
	CancelOrder(ctx context.Context, tenantID, orderID, customerID uuid.UUID) error
	UpdatePaymentMethod(ctx context.Context, tenantID, orderID, customerID, paymentMethodID uuid.UUID) (*domain.Order, error)
}

type ITransactionService interface {
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Transaction, int64, error)
}

type ICustomerAddressService interface {
	List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerAddress, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerAddress, error)
	Create(ctx context.Context, addr *domain.CustomerAddress, asDefault bool) error
	Update(ctx context.Context, addr *domain.CustomerAddress, asDefault bool) error
	Delete(ctx context.Context, tenantID, customerID, id uuid.UUID) error
}

type IPaymentMethodService interface {
	List(ctx context.Context, tenantID uuid.UUID) ([]domain.PaymentMethod, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.PaymentMethod, error)
	Update(ctx context.Context, tenantID uuid.UUID, pm *domain.PaymentMethod) error
}

type ICustomerPaymentMethodService interface {
	List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerPaymentMethod, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerPaymentMethod, error)
	Create(ctx context.Context, cpm *domain.CustomerPaymentMethod, asDefault bool) error
	Update(ctx context.Context, cpm *domain.CustomerPaymentMethod, asDefault bool) error
	Delete(ctx context.Context, tenantID, customerID, id uuid.UUID) error
}

type ICartService interface {
	GetCart(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Cart, error)
	AddItem(ctx context.Context, tenantID, customerID uuid.UUID, productID uuid.UUID, quantity int) error
	UpdateItemQty(ctx context.Context, tenantID, customerID, itemID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, tenantID, customerID, itemID uuid.UUID) error
	ToggleSelection(ctx context.Context, tenantID, customerID, itemID uuid.UUID, selected bool) error
	ClearCart(ctx context.Context, tenantID, customerID uuid.UUID) error
}
