package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID            uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID      uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index" validate:"required"`
	CategoryID    *uuid.UUID       `json:"category_id" gorm:"type:uuid;index" validate:"omitempty"`
	NameTh        string           `json:"name_th" gorm:"not null" validate:"required,min=2"`
	NameEn        string           `json:"name_en" validate:"omitempty"`
	SKU           *string          `json:"sku" gorm:"uniqueIndex:idx_tenant_sku" validate:"omitempty"`
	DescriptionTh string           `json:"description_th"`
	DescriptionEn string           `json:"description_en"`
	Price         float64          `json:"price" gorm:"not null" validate:"required,gte=0"`
	Inventory     int              `json:"inventory" gorm:"default:0" validate:"gte=0"`
	Status        string           `json:"status" gorm:"default:'active'"` // active, inactive
	SortOrder     int              `json:"sort_order" gorm:"default:0"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `json:"-" gorm:"index"`

	Images   []ProductImage   `json:"images" gorm:"foreignKey:ProductID"`
	Category *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	URL       string    `json:"url" gorm:"not null"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
}

type ProductCategory struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID  uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index" validate:"required"`
	NameTh    string         `json:"name_th" gorm:"not null" validate:"required"`
	NameEn    string         `json:"name_en"`
	ParentID  *uuid.UUID     `json:"parent_id" gorm:"type:uuid"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
