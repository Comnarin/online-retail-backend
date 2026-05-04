package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_cart_tenant_customer" json:"tenant_id"`
	CustomerID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_cart_tenant_customer" json:"customer_id"`
	Items      []CartItem     `gorm:"foreignKey:CartID" json:"items"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index;uniqueIndex:idx_cart_tenant_customer"`
}

type CartItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CartID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Selected  bool           `gorm:"default:true" json:"selected"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
