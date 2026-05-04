package model

import (
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LineLoginRequest struct {
	TenantID    uuid.UUID `json:"tenant_id" validate:"required"`
	LineUserID  string    `json:"line_user_id" validate:"required"`
	IDToken     string    `json:"id_token" validate:"required"`
	DisplayName string    `json:"display_name"`
	PictureURL  string    `json:"picture_url"`
}

type RegisterAdminRequest struct {
	TenantID uuid.UUID `json:"tenant_id"` // Set by URL params
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token    string           `json:"token"`
	User     *domain.User     `json:"user,omitempty"`
	Customer *domain.Customer `json:"customer,omitempty"`
}
