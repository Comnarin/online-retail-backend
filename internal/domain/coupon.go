package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypeFixed   DiscountType = "fixed"
	DiscountTypePercent DiscountType = "percent"
)

type Coupon struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID      uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Code          string         `json:"code" gorm:"not null" validate:"required,uppercase"`
	Description   string         `json:"description"`
	DiscountType  DiscountType   `json:"discount_type" gorm:"not null" validate:"required,oneof=fixed percent"`
	DiscountValue float64        `json:"discount_value" gorm:"not null" validate:"required,gt=0"`
	MinOrderValue float64        `json:"min_order_value" gorm:"default:0"`
	MaxDiscount   float64        `json:"max_discount" gorm:"default:0"`
	StartDate     time.Time      `json:"start_date" validate:"required"`
	EndDate       time.Time      `json:"end_date" validate:"required,gtfield=StartDate"`
	UsageLimit    int            `json:"usage_limit" gorm:"default:0"`
	UsedCount     int            `json:"used_count" gorm:"default:0"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	TierID        *uuid.UUID     `json:"tier_id" gorm:"type:uuid"`
	Tier          *MembershipTier `json:"tier" gorm:"foreignKey:TierID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
