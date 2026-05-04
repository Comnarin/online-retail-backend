package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerPaymentMethod struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CustomerID      uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null;index"`
	TenantID        uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	PaymentMethodID uuid.UUID      `json:"payment_method_id" gorm:"type:uuid;not null"`
	Provider        string         `json:"provider"` // omise, stripe, mock
	ProviderRef     string         `json:"provider_ref"` // Token or Card ID
	Last4           string         `json:"last4"`
	Brand           string         `json:"brand"` // visa, mastercard
	ExpMonth        int            `json:"exp_month"`
	ExpYear         int            `json:"exp_year"`
	IsDefault       bool           `json:"is_default" gorm:"default:false"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Customer        *Customer       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	PaymentMethod   PaymentMethod  `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
}
