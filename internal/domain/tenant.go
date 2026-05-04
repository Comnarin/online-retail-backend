package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null" validate:"required,lowercase"`
	OrderCode string         `json:"order_code" gorm:"size:10;default:'ORD'"`
	Status    string         `json:"status" gorm:"default:'active'"` // active, inactive, suspended
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Features   TenantFeature    `json:"features" gorm:"foreignKey:TenantID"`
	Appearance TenantAppearance `json:"appearance" gorm:"foreignKey:TenantID"`
}

type TenantFeature struct {
	TenantID           uuid.UUID `json:"tenant_id" gorm:"type:uuid;primaryKey"`
	MaxProducts        int       `json:"max_products"`
	EnableLiff         bool      `json:"enable_liff"`
	EnableCoupons      bool      `json:"enable_coupons"`
	EnableMembership   bool      `json:"enable_membership"`
	EnablePoints       bool      `json:"enable_points"`
	EnableInventory    bool      `json:"enable_inventory"`
	EnableReviews      bool      `json:"enable_reviews"`
	EnableDelivery     bool      `json:"enable_delivery"`
	MultiBranchSupport bool      `json:"multi_branch_support"`
	PointExchangeRate  int       `json:"point_exchange_rate" gorm:"default:100"`
}

type TenantAppearance struct {
	TenantID        uuid.UUID `json:"tenant_id" gorm:"type:uuid;primaryKey"`
	LogoURL         string    `json:"logo_url"`
	PrimaryColor    string    `json:"primary_color" validate:"omitempty,min=4,max=10"`
	BannerURL       string    `json:"banner_url"`
	Theme           string    `json:"theme"` // default, dark, curator
	ShopTitle       string    `json:"shop_title"`
	ShopDescription string    `json:"shop_description"`
	ShopHeroURL     string    `json:"shop_hero_url"`
}
