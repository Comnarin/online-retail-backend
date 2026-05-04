package handler

import (
	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	Login(c *fiber.Ctx) error
	LineLogin(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	DevLiffLogin(c *fiber.Ctx) error
}

type ITenantHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UpdateFeatures(c *fiber.Ctx) error
	UpdateAppearance(c *fiber.Ctx) error
	CreateAdmin(c *fiber.Ctx) error
	ListAdmins(c *fiber.Ctx) error
	DeleteAdmin(c *fiber.Ctx) error
}

type IDashboardHandler interface {
	Stats(c *fiber.Ctx) error
	RecentOrders(c *fiber.Ctx) error
}

type IAppearanceHandler interface {
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type IProductHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadImage(c *fiber.Ctx) error
	DeleteImage(c *fiber.Ctx) error
	ListCategories(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type IOrderHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	ConfirmPayment(c *fiber.Ctx) error
	RejectPayment(c *fiber.Ctx) error
}

type ICustomerHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
}

type ICouponHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Toggle(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type IMembershipHandler interface {
	ListTiers(c *fiber.Ctx) error
	CreateTier(c *fiber.Ctx) error
	UpdateTier(c *fiber.Ctx) error
	DeleteTier(c *fiber.Ctx) error
	ListTransactions(c *fiber.Ctx) error
}

type ILiffProductHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	ListCategories(c *fiber.Ctx) error
}

type ILiffOrderHandler interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	UploadSlip(c *fiber.Ctx) error
	GetActivePending(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
	UpdatePaymentMethod(c *fiber.Ctx) error
}

type ILiffCartHandler interface {
	Get(c *fiber.Ctx) error
	AddItem(c *fiber.Ctx) error
	UpdateItem(c *fiber.Ctx) error
	RemoveItem(c *fiber.Ctx) error
	ToggleSelection(c *fiber.Ctx) error
	Clear(c *fiber.Ctx) error
}

type ILiffCustomerHandler interface {
	GetProfile(c *fiber.Ctx) error
}

type ITransactionHandler interface {
	List(c *fiber.Ctx) error
}

type IPaymentWebhookHandler interface {
	Handle(c *fiber.Ctx) error
}

type IAdminPaymentMethodHandler interface {
	List(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	UploadQRCode(c *fiber.Ctx) error
}

type ILiffAddressHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type ILiffPaymentHandler interface {
	ListMethods(c *fiber.Ctx) error
	ListSavedMethods(c *fiber.Ctx) error
	CreateSavedMethod(c *fiber.Ctx) error
	UpdateSavedMethod(c *fiber.Ctx) error
	DeleteSavedMethod(c *fiber.Ctx) error
	SetDefaultMethod(c *fiber.Ctx) error
}