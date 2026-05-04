package repository

import (
	"context"

	"errors"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"gorm.io/gorm"
)

var (
	ErrNotEmpty  = errors.New("resource is not empty")
	ErrForbidden = errors.New("forbidden: you do not have permission to access this resource")
)

// ITenantRepository handles tenant persistence (System-wide)
type ITenantRepository interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	List(ctx context.Context, opts model.ListOptions) ([]domain.Tenant, int64, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	UpdateFeatures(ctx context.Context, id uuid.UUID, features domain.TenantFeature) error
	UpdateAppearance(ctx context.Context, id uuid.UUID, appearance domain.TenantAppearance) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTx(tx *gorm.DB) ITenantRepository
}

// IUserRepository handles administrative user persistence
type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.User, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, user *domain.User) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
}

// ICustomerRepository handles LINE customer persistence
type ICustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Customer, error)
	GetByLineUserID(ctx context.Context, tenantID uuid.UUID, lineUserID string) (*domain.Customer, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Customer, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	WithTx(tx *gorm.DB) ICustomerRepository
}

// IProductRepository handles product persistence with tenant isolation
type IProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Product, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, product *domain.Product) error
	UpdateFields(ctx context.Context, tenantID, id uuid.UUID, fields map[string]interface{}) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	WithTx(tx *gorm.DB) IProductRepository

	CreateCategory(ctx context.Context, cat *domain.ProductCategory) error
	ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.ProductCategory, error)
	UpdateCategory(ctx context.Context, tenantID uuid.UUID, cat *domain.ProductCategory) error
	DeleteCategory(ctx context.Context, tenantID, id uuid.UUID) error
	AddImage(ctx context.Context, image *domain.ProductImage) error
	RemoveImage(ctx context.Context, tenantID, productID, imageID uuid.UUID) error
}

// IOrderRepository handles order persistence with tenant isolation
type IOrderRepository interface {
	Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error)
	UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.OrderStatus) error
	Update(ctx context.Context, order *domain.Order) error
	WithTx(tx *gorm.DB) IOrderRepository
}

// ICouponRepository handles coupon/promotion persistence with tenant isolation
type ICouponRepository interface {
	Create(ctx context.Context, coupon *domain.Coupon) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Coupon, error)
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Coupon, int64, error)
	Update(ctx context.Context, tenantID uuid.UUID, coupon *domain.Coupon) error
	IncrementUsage(ctx context.Context, tenantID, id uuid.UUID) error
	ToggleActive(ctx context.Context, tenantID, id uuid.UUID) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	WithTx(tx *gorm.DB) ICouponRepository
}

// IMembershipRepository handles membership tier persistence with tenant isolation
type IMembershipRepository interface {
	CreateTier(ctx context.Context, tier *domain.MembershipTier) error
	ListTiers(ctx context.Context, tenantID uuid.UUID) ([]domain.MembershipTier, error)
	UpdateTier(ctx context.Context, tenantID uuid.UUID, tier *domain.MembershipTier) error
	DeleteTier(ctx context.Context, tenantID, id uuid.UUID) error

	GetCustomerMembership(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.CustomerMembership, error)
	UpsertCustomerMembership(ctx context.Context, cm *domain.CustomerMembership) error
	WithTx(tx *gorm.DB) IMembershipRepository
}

// IPointRepository handles loyalty point persistence
type IPointRepository interface {
	GetBalance(ctx context.Context, tenantID, customerID uuid.UUID) (int, error)
	AddTransaction(ctx context.Context, tx *domain.PointTransaction) error
	ListTransactions(ctx context.Context, tenantID, customerID uuid.UUID, opts model.ListOptions) ([]domain.PointTransaction, int64, error)
	WithTx(tx *gorm.DB) IPointRepository
}

// IDashboardRepository handles analytical dashboard queries
type IDashboardRepository interface {
	GetStats(ctx context.Context, tenantID uuid.UUID) (*model.DashboardStats, error)
	GetRecentOrders(ctx context.Context, tenantID uuid.UUID, limit int) ([]model.RecentOrderRow, error)
}

// ITransactionRepository handles financial transaction logging
type ITransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Transaction, error)
	Update(ctx context.Context, tx *domain.Transaction) error
	List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Transaction, int64, error)
	WithTx(tx *gorm.DB) ITransactionRepository
}

// ICustomerAddressRepository handles customer addresses
type ICustomerAddressRepository interface {
	List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerAddress, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerAddress, error)
	Create(ctx context.Context, addr *domain.CustomerAddress) error
	Update(ctx context.Context, addr *domain.CustomerAddress) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	SetDefault(ctx context.Context, tenantID, customerID, id uuid.UUID) error
}

type IPaymentMethodRepository interface {
	List(ctx context.Context, tenantID uuid.UUID) ([]domain.PaymentMethod, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.PaymentMethod, error)
	Create(ctx context.Context, pm *domain.PaymentMethod) error
	Update(ctx context.Context, tenantID uuid.UUID, pm *domain.PaymentMethod) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
}

type ICustomerPaymentMethodRepository interface {
	List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerPaymentMethod, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerPaymentMethod, error)
	Create(ctx context.Context, cpm *domain.CustomerPaymentMethod) error
	Update(ctx context.Context, cpm *domain.CustomerPaymentMethod) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	SetDefault(ctx context.Context, tenantID, customerID, id uuid.UUID) error
}

type ICartRepository interface {
	GetByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Cart, error)
	UpsertItem(ctx context.Context, item *domain.CartItem) error
	UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
	ToggleSelectItem(ctx context.Context, itemID uuid.UUID, selected bool) error
	ClearSelectedItems(ctx context.Context, cartID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
	WithTx(tx *gorm.DB) ICartRepository
}
