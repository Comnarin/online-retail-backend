package model

// ListOptions defines common pagination and filtering parameters for list endpoints
type ListOptions struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	Category string `json:"category"`
	Status   string `json:"status"`
}
