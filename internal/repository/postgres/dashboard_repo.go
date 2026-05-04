package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) repository.IDashboardRepository {
	return &dashboardRepo{db: db}
}

func (r *dashboardRepo) GetStats(ctx context.Context, tenantID uuid.UUID) (*model.DashboardStats, error) {
	var stats model.DashboardStats

	// We use the safe WithContext wrapper for executing analytical queries without breaking architectural boundaries
	r.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(total), 0) as total_revenue FROM orders WHERE tenant_id = ? AND status != 'cancelled'`, tenantID).Scan(&stats.TotalRevenue)
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders WHERE tenant_id = ?`, tenantID).Scan(&stats.TotalOrders)
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM users WHERE tenant_id = ? AND role = 'customer'`, tenantID).Scan(&stats.TotalCustomers)
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM products WHERE tenant_id = ? AND status = 'active'`, tenantID).Scan(&stats.TotalProducts)
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM coupons WHERE tenant_id = ? AND is_active = true`, tenantID).Scan(&stats.ActiveCoupons)
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders WHERE tenant_id = ? AND status = 'pending'`, tenantID).Scan(&stats.PendingOrders)

	return &stats, nil
}

func (r *dashboardRepo) GetRecentOrders(ctx context.Context, tenantID uuid.UUID, limit int) ([]model.RecentOrderRow, error) {
	var orders []model.RecentOrderRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT o.id, o.order_number, u.name as customer_name, o.status, o.total,
		       to_char(o.created_at, 'DD Mon YYYY') as created_at
		FROM orders o
		JOIN users u ON u.id = o.customer_id
		WHERE o.tenant_id = ?
		ORDER BY o.created_at DESC
		LIMIT ?`, tenantID, limit).Scan(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}
