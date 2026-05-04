package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID      uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name          string         `json:"name" gorm:"not null"`
	Code          string         `json:"code" gorm:"not null;index"` // promptpay, credit_card, cod
	Description   string         `json:"description"`
	Icon          string         `json:"icon" gorm:"default:'CreditCard'"` // Lucide icon name
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	ExpiryMinutes int            `json:"expiry_minutes" gorm:"default:30"`
	QRCodeURL     string         `json:"qr_code_url" gorm:"type:text;default:''"` // MinIO URL for QR code image
	Config        string         `json:"config" gorm:"type:text"` // JSON config for keys, etc.
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
