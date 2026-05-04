package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID      uuid.UUID         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID       uuid.UUID         `json:"order_id" gorm:"type:uuid;not null;index"`
	CustomerID    uuid.UUID         `json:"customer_id" gorm:"type:uuid;not null;index"`
	Amount        float64           `json:"amount" gorm:"not null"`
	Status        TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentMethod string            `json:"payment_method"` // promptpay, credit_card
	Reference     string            `json:"reference"`      // external payment ref
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`

	Order    Order  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}
