package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipTier struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID     uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name         string         `json:"name" gorm:"not null" validate:"required"`
	Level        int            `json:"level" gorm:"not null" validate:"required,gte=0"`
	MinPoints    int            `json:"min_points" gorm:"not null" validate:"required,gte=0"`
	DiscountRate float64        `json:"discount_rate" gorm:"type:decimal(5,2);default:0"`
	Color        string         `json:"color" gorm:"type:varchar(10);default:'#E8572A'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	Benefits []TierBenefit `json:"benefits" gorm:"foreignKey:TierID"`
}

type TierBenefit struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TierID    uuid.UUID `json:"tier_id" gorm:"type:uuid;not null;index"`
	Benefit   string    `json:"benefit" gorm:"not null"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
}

type CustomerMembership struct {
	TenantID     uuid.UUID      `json:"tenant_id" gorm:"type:uuid;primaryKey"`
	CustomerID   uuid.UUID      `json:"customer_id" gorm:"type:uuid;primaryKey"`
	TierID       uuid.UUID      `json:"tier_id" gorm:"type:uuid;not null" validate:"required"`
	Points       int            `json:"points" gorm:"default:0" validate:"gte=0"`
	TotalSpent   float64        `json:"total_spent" gorm:"default:0"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	Tier MembershipTier `json:"tier" gorm:"foreignKey:TierID"`
	Customer Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

type PointTransaction struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CustomerID  uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null;index" validate:"required"`
	OrderID     *uuid.UUID     `json:"order_id" gorm:"type:uuid;index"`
	Points      int            `json:"points" gorm:"not null" validate:"required,ne=0"`
	Type        string         `json:"type" gorm:"not null"` // earn, redeem
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
}
