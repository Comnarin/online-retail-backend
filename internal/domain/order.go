package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending             OrderStatus = "pending"
	OrderStatusPendingVerification OrderStatus = "pending_verification"
	OrderStatusConfirmed           OrderStatus = "confirmed"
	OrderStatusProcessing          OrderStatus = "processing"
	OrderStatusShipped             OrderStatus = "shipped"
	OrderStatusCompleted           OrderStatus = "completed"
	OrderStatusCancelled           OrderStatus = "cancelled"
	OrderStatusExpired             OrderStatus = "expired"
)


type Order struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID       uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderNumber    string         `json:"order_number" gorm:"uniqueIndex;not null"`
	CustomerID     uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null;index" validate:"required"`
	Subtotal       float64        `json:"subtotal" gorm:"not null"`
	DiscountAmount float64        `json:"discount_amount" gorm:"default:0"`
	Total          float64        `json:"total" gorm:"not null"`
	Status         OrderStatus    `json:"status" gorm:"default:'pending'"`
	CouponID        *uuid.UUID     `json:"coupon_id" gorm:"type:uuid"`
	PointsUsed      int            `json:"points_used" gorm:"default:0"`
	PaymentMethodID *uuid.UUID     `json:"payment_method_id" gorm:"type:uuid"`
	PaymentDeadline *time.Time     `json:"payment_deadline"`
	SlipImageURL    string         `json:"slip_image_url" gorm:"type:text;default:''"`
	Note            string         `json:"note"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	Items           []OrderItem  `json:"items" gorm:"foreignKey:OrderID"`
	Customer        Customer     `json:"customer" gorm:"foreignKey:CustomerID"`
	ShippingAddress OrderAddress `json:"shipping_address" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index" validate:"required"`
	NameTh    string    `json:"name_th" gorm:"not null"`
	NameEn    string    `json:"name_en"`
	Price     float64   `json:"price" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Subtotal  float64   `json:"subtotal" gorm:"not null"`

	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

type OrderAddress struct {
	OrderID  uuid.UUID `json:"order_id" gorm:"type:uuid;primaryKey"`
	Name     string    `json:"name" validate:"required"`
	Phone    string    `json:"phone" validate:"required"`
	Address  string    `json:"address" validate:"required"`
	District string    `json:"district" validate:"required"`
	Province string    `json:"province" validate:"required"`
	ZipCode  string    `json:"zip_code" validate:"required"`
}
