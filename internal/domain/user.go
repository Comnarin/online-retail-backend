package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role types
type UserRole string

const (
	RoleSuperAdmin  UserRole = "superadmin"
	RoleTenantAdmin UserRole = "tenant_admin"
)

// User is the core administrative user entity (SuperAdmin, TenantAdmin)
type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID  *uuid.UUID `json:"tenant_id" gorm:"type:uuid;index"`
	Email     *string    `json:"email" gorm:"uniqueIndex"`
	Password  *string    `json:"-"`
	Name      string     `json:"name"`
	Role      UserRole   `json:"role" gorm:"not null"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	Tenant *Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}
