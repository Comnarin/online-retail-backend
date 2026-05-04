package model

type DashboardStats struct {
	TotalRevenue    float64 `json:"total_revenue"`
	TotalOrders      int64   `json:"total_orders"`
	TotalCustomers   int64   `json:"total_customers"`
	TotalProducts    int64   `json:"total_products"`
	ActiveCoupons    int64   `json:"active_coupons"`
	PendingOrders    int64   `json:"pending_orders"`
}

type RecentOrderRow struct {
	ID           string  `json:"id"`
	OrderNumber  string  `json:"order_number"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status"`
	Total        float64 `json:"total"`
	CreatedAt    string  `json:"created_at"`
}
