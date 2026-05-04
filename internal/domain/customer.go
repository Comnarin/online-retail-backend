package domain

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a LINE LIFF user in the system
type Customer struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	LineUserID *string    `json:"line_user_id" gorm:"uniqueIndex"`
	Name       string     `json:"name" gorm:"not null"`
	Email      *string    `json:"email" gorm:"index"`
	Phone      *string    `json:"phone"`
	Avatar     *string    `json:"avatar"`
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	// Relational / Computed fields
	Tenant      *Tenant             `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Membership  *CustomerMembership `json:"membership,omitempty" gorm:"foreignKey:CustomerID"`
	TotalPoints int64               `json:"total_points" gorm:"->"`
	TotalSpent  float64             `json:"total_spent" gorm:"->"`
	OrderCount  int64               `json:"order_count" gorm:"->"`
}

type CustomerAddress struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CustomerID uuid.UUID  `json:"customer_id" gorm:"type:uuid;not null;index"`
	Label      string     `json:"label"` // Home, Office, etc.
	Name       string     `json:"name" validate:"required"`
	Phone      string     `json:"phone" validate:"required"`
	Address    string     `json:"address" validate:"required"`
	District   string     `json:"district" validate:"required"`
	Province   string     `json:"province" validate:"required"`
	ZipCode    string     `json:"zip_code" validate:"required"`
	IsDefault  bool       `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Customer *Customer `json:"-" gorm:"foreignKey:CustomerID"`
}
