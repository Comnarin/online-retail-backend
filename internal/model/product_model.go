package model

import "github.com/google/uuid"

type CreateProductRequest struct {
	CategoryID    *uuid.UUID `json:"category_id" validate:"omitempty"`
	NameTh        string     `json:"name_th" validate:"required,min=2"`
	NameEn        string     `json:"name_en" validate:"omitempty"`
	SKU           *string    `json:"sku" validate:"omitempty"`
	DescriptionTh string     `json:"description_th" validate:"omitempty"`
	DescriptionEn string     `json:"description_en" validate:"omitempty"`
	Price         float64    `json:"price" validate:"required,gte=0"`
	Inventory     int        `json:"inventory" validate:"gte=0"`
	Status        string     `json:"status" validate:"omitempty"`
}

type UpdateProductRequest struct {
	CategoryID    *uuid.UUID `json:"category_id" validate:"omitempty"`
	NameTh        string     `json:"name_th" validate:"omitempty,min=2"`
	NameEn        string     `json:"name_en" validate:"omitempty"`
	SKU           *string    `json:"sku" validate:"omitempty"`
	DescriptionTh string     `json:"description_th" validate:"omitempty"`
	DescriptionEn string     `json:"description_en" validate:"omitempty"`
	Price         *float64   `json:"price" validate:"omitempty,gte=0"`
	Inventory     *int       `json:"inventory" validate:"omitempty,gte=0"`
	Status        string     `json:"status" validate:"omitempty"`
}
