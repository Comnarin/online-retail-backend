package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/platform/cache"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo         repository.IOrderRepository
	productRepo       repository.IProductRepository
	pointRepo         repository.IPointRepository
	couponRepo        repository.ICouponRepository
	membershipRepo    repository.IMembershipRepository
	transactionRepo   repository.ITransactionRepository
	addressRepo       repository.ICustomerAddressRepository
	paymentMethodRepo repository.IPaymentMethodRepository
	tenantRepo        repository.ITenantRepository
	cartRepo          repository.ICartRepository
	redis             *cache.RedisClient
	db                *gorm.DB
}

func NewOrderService(
	orderRepo repository.IOrderRepository,
	productRepo repository.IProductRepository,
	pointRepo repository.IPointRepository,
	couponRepo repository.ICouponRepository,
	membershipRepo repository.IMembershipRepository,
	transactionRepo repository.ITransactionRepository,
	addressRepo repository.ICustomerAddressRepository,
	paymentMethodRepo repository.IPaymentMethodRepository,
	tenantRepo repository.ITenantRepository,
	cartRepo repository.ICartRepository,
	redis *cache.RedisClient,
	db *gorm.DB,
) IOrderService {
	return &OrderService{orderRepo, productRepo, pointRepo, couponRepo, membershipRepo, transactionRepo, addressRepo, paymentMethodRepo, tenantRepo, cartRepo, redis, db}
}

func (u *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	// 0. Check Redis for existing pending order
	activeKey := fmt.Sprintf("retail:order:active:%s:%s", req.TenantID, req.CustomerID)
	existingID, _ := u.redis.Get(ctx, activeKey)
	if existingID != "" {
		return nil, fmt.Errorf("you have a pending order (#%s). please complete or cancel it first", existingID[:8])
	}

	pm, err := u.paymentMethodRepo.GetByID(ctx, req.TenantID, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment method")
	}
	deadline := time.Now().Add(time.Duration(pm.ExpiryMinutes) * time.Minute)

	var orderItems []domain.OrderItem
	var subtotal float64
	var order *domain.Order

	if err := u.db.Transaction(func(tx *gorm.DB) error {
		// New transaction-aware repos
		pRepo := u.productRepo.WithTx(tx)
		cRepo := u.couponRepo.WithTx(tx)
		mRepo := u.membershipRepo.WithTx(tx)
		oRepo := u.orderRepo.WithTx(tx)
		cartRepo := u.cartRepo.WithTx(tx)
		
		for _, item := range req.Items {
			// Fetch product in transaction
			product, err := pRepo.GetByID(ctx, req.TenantID, item.ProductID)
			if err != nil {
				return err
			}

			// 1. Check & Deduct Stock
			if product.Inventory < item.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", product.NameTh)
			}

			product.Inventory -= item.Quantity
			if err := pRepo.Update(ctx, req.TenantID, product); err != nil {
				return err
			}

			lineTotal := product.Price * float64(item.Quantity)
			subtotal += lineTotal
			orderItems = append(orderItems, domain.OrderItem{
				ProductID: item.ProductID,
				NameTh:    product.NameTh,
				NameEn:    product.NameEn,
				Price:     product.Price,
				Quantity:  item.Quantity,
				Subtotal:  lineTotal,
			})
		}

		discount := 0.0
		var couponID *uuid.UUID

		if req.CouponCode != "" {
			coupon, err := cRepo.GetByCode(ctx, req.TenantID, req.CouponCode)
			if err == nil && coupon.IsActive {
				if coupon.TierID != nil {
					membership, err := mRepo.GetCustomerMembership(ctx, req.TenantID, req.CustomerID)
					if err != nil || membership == nil || membership.TierID != *coupon.TierID {
						return fmt.Errorf("coupon %s is restricted to a different membership tier", req.CouponCode)
					}
				}

				if coupon.DiscountType == domain.DiscountTypePercent {
					discount = subtotal * coupon.DiscountValue / 100
					if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
						discount = coupon.MaxDiscount
					}
				} else {
					discount = coupon.DiscountValue
				}

				if subtotal < coupon.MinOrderValue {
					return fmt.Errorf("order subtotal must be at least ฿%.2f to use this coupon", coupon.MinOrderValue)
				}

				_ = cRepo.IncrementUsage(ctx, req.TenantID, coupon.ID)
				couponID = &coupon.ID
			}
		}

		order = &domain.Order{
			TenantID:        req.TenantID,
			CustomerID:      req.CustomerID,
			Subtotal:        subtotal,
			DiscountAmount:  discount,
			Total:           subtotal - discount,
			CouponID:        couponID,
			PointsUsed:      req.PointsToRedeem,
			ShippingAddress: req.ShippingAddress,
			Note:            req.Note,
			Status:          domain.OrderStatusPending,
			PaymentMethodID: &req.PaymentMethodID,
			PaymentDeadline: &deadline,
		}

		if req.PointsToRedeem > 0 {
			membership, err := mRepo.GetCustomerMembership(ctx, req.TenantID, req.CustomerID)
			if err != nil || membership == nil || membership.Points < req.PointsToRedeem {
				return fmt.Errorf("insufficient points to redeem")
			}
			order.Total -= float64(req.PointsToRedeem)
			if order.Total < 0 {
				order.Total = 0
			}

			// Record Point Transaction
			pointTx := &domain.PointTransaction{
				TenantID:    order.TenantID,
				CustomerID:  order.CustomerID,
				OrderID:     &order.ID,
				Points:      -req.PointsToRedeem,
				Type:        "redeem",
				Description: fmt.Sprintf("Points redeemed for order #%s", order.OrderNumber),
			}
			if err := u.pointRepo.WithTx(tx).AddTransaction(ctx, pointTx); err != nil {
				return err
			}

			// Update Balance
			membership.Points -= req.PointsToRedeem
			if err := mRepo.UpsertCustomerMembership(ctx, membership); err != nil {
				return err
			}
		}

		// Create Order
		if err := oRepo.Create(ctx, order, orderItems); err != nil {
			return err
		}

		// Clear selected items from persistent cart
		cart, err := cartRepo.GetByCustomerID(ctx, order.TenantID, order.CustomerID)
		if err == nil && cart != nil {
			_ = cartRepo.ClearSelectedItems(ctx, cart.ID)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	_ = u.transactionRepo.Create(ctx, &domain.Transaction{
		TenantID:      order.TenantID,
		OrderID:       order.ID,
		CustomerID:    order.CustomerID,
		Amount:        order.Total,
		Status:        domain.TransactionStatusPending,
		PaymentMethod: "promptpay",
	})

	// Optional: Save address for future use
	if req.SaveAddress {
		label := req.AddressLabel
		if label == "" {
			label = "Saved Address"
		}
		_ = u.addressRepo.Create(ctx, &domain.CustomerAddress{
			TenantID:   order.TenantID,
			CustomerID: order.CustomerID,
			Label:      label,
			Name:       req.ShippingAddress.Name,
			Phone:      req.ShippingAddress.Phone,
			Address:    req.ShippingAddress.Address,
			District:   req.ShippingAddress.District,
			Province:   req.ShippingAddress.Province,
			ZipCode:    req.ShippingAddress.ZipCode,
		})
	}

	// 5. Track in Redis
	_ = u.redis.Set(ctx, activeKey, order.ID.String(), time.Duration(pm.ExpiryMinutes)*time.Minute)

	return order, nil
}

func (u *OrderService) GetOrder(ctx context.Context, tenantID, id uuid.UUID) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Auto-repair disabled as per user request (no timing/expiring)
	// Auto-check expiration disabled
	
	// Re-fetch in case we need it, though no mutation happened
	order, err = u.orderRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderService) CheckExpiration(ctx context.Context, orderID uuid.UUID) error {
	// Expiration logic disabled as per user request
	return nil
}

func (u *OrderService) ListByTenant(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error) {
	return u.orderRepo.ListByTenant(ctx, tenantID, opts)
}

func (u *OrderService) ListByCustomer(ctx context.Context, customerID uuid.UUID, opts model.ListOptions) ([]domain.Order, int64, error) {
	return u.orderRepo.ListByCustomer(ctx, customerID, opts)
}

func (u *OrderService) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.OrderStatus) error {
	err := u.orderRepo.UpdateStatus(ctx, tenantID, id, status)
	if err == nil && (status == domain.OrderStatusCompleted || status == domain.OrderStatusCancelled || status == domain.OrderStatusExpired) {
		order, _ := u.orderRepo.GetByID(ctx, tenantID, id)
		if order != nil {
			activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, order.CustomerID)
			_ = u.redis.Del(ctx, activeKey)
		}
	}
	return err
}

func (u *OrderService) GetActivePendingOrder(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Order, error) {
	activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, customerID)
	orderIDStr, _ := u.redis.Get(ctx, activeKey)
	
	var activeOrder *domain.Order
	if orderIDStr != "" {
		if orderID, err := uuid.Parse(orderIDStr); err == nil {
			// GetOrder will automatically run CheckExpiration
			activeOrder, _ = u.GetOrder(ctx, tenantID, orderID)
		}
	}

	// If we have an order but it's no longer pending, it's not our "active" session anymore
	// Return the order if it's pending OR if it's recent enough to show the status screen
	if activeOrder != nil && (activeOrder.Status == domain.OrderStatusPending || activeOrder.Status == domain.OrderStatusExpired || activeOrder.Status == domain.OrderStatusCancelled) {
		return activeOrder, nil
	}

	// Session expired or no session in Redis: Fallback to searching DB for any pending order
	if activeOrder != nil {
		_ = u.redis.Del(ctx, activeKey)
	}

	// Fetch up to 10 pending orders to clean up any old ones
	opts := model.ListOptions{Limit: 10, Status: string(domain.OrderStatusPending)}
	orders, _, err := u.orderRepo.ListByCustomer(ctx, customerID, opts)
	if err != nil || len(orders) == 0 {
		return nil, fmt.Errorf("no active pending order")
	}

	var validOrder *domain.Order
	for _, o := range orders {
		// GetOrder will automatically run CheckExpiration
		order, _ := u.GetOrder(ctx, tenantID, o.ID)
		if order != nil && order.Status == domain.OrderStatusPending {
			validOrder = order
			break
		}
	}

	if validOrder == nil {
		return nil, fmt.Errorf("no active pending order")
	}

	return validOrder, nil
}

func (u *OrderService) CancelOrder(ctx context.Context, tenantID, orderID, customerID uuid.UUID) error {
	order, err := u.orderRepo.GetByID(ctx, tenantID, orderID)
	if err != nil {
		return err
	}

	if order.CustomerID != customerID {
		return fmt.Errorf("unauthorized")
	}

	if order.Status != domain.OrderStatusPending {
		return fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update Status
		if err := u.orderRepo.WithTx(tx).UpdateStatus(ctx, tenantID, orderID, domain.OrderStatusCancelled); err != nil {
			return err
		}

		// 2. Release Stock
		var o domain.Order
		if err := tx.Preload("Items").First(&o, "id = ?", orderID).Error; err == nil {
			for _, item := range o.Items {
				_ = tx.Model(&domain.Product{}).Where("id = ?", item.ProductID).
					Update("inventory", gorm.Expr("inventory + ?", item.Quantity))
			}
		}

		// 3. Clear Redis
		activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, customerID)
		_ = u.redis.Del(ctx, activeKey)

		return nil
	})
}

func (u *OrderService) UploadSlip(ctx context.Context, tenantID, orderID, customerID uuid.UUID, slipURL string) error {
	order, err := u.orderRepo.GetByID(ctx, tenantID, orderID)
	if err != nil {
		return err
	}

	if order.CustomerID != customerID {
		return fmt.Errorf("unauthorized")
	}

	if order.Status != domain.OrderStatusPending {
		return fmt.Errorf("order is not in pending status")
	}

	order.SlipImageURL = slipURL
	order.Status = domain.OrderStatusPendingVerification

	if err := u.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Clear the active pending key — customer has submitted proof
	activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, customerID)
	_ = u.redis.Del(ctx, activeKey)

	return nil
}

func (u *OrderService) ConfirmPayment(ctx context.Context, tenantID, orderID uuid.UUID) error {
	return u.db.Transaction(func(dbTx *gorm.DB) error {
		orderRepo := u.orderRepo.WithTx(dbTx)

		order, err := orderRepo.GetByID(ctx, tenantID, orderID)
		if err != nil {
			return err
		}

		if order.Status != domain.OrderStatusPendingVerification {
			return fmt.Errorf("order is not pending verification, current status: %s", order.Status)
		}

		// 1. Update Order Status to Processing (simplified flow)
		if err := orderRepo.UpdateStatus(ctx, tenantID, orderID, domain.OrderStatusProcessing); err != nil {
			return err
		}

		// 2. Update Transaction
		txRepo := u.transactionRepo.WithTx(dbTx)
		tx, err := txRepo.GetByOrderID(ctx, orderID)
		if err == nil && tx != nil {
			tx.Status = domain.TransactionStatusSuccess
			_ = txRepo.Update(ctx, tx)
		}

		// 3. Award Points
		tenant, err := u.tenantRepo.GetByID(ctx, tenantID)
		if err == nil && tenant.Features.EnablePoints && tenant.Features.PointExchangeRate > 0 {
			pointsEarned := int(order.Total / float64(tenant.Features.PointExchangeRate))
			if pointsEarned > 0 {
				pointTx := &domain.PointTransaction{
					TenantID:    tenantID,
					CustomerID:  order.CustomerID,
					OrderID:     &orderID,
					Points:      pointsEarned,
					Type:        "earn",
					Description: fmt.Sprintf("Points earned from order #%s", order.OrderNumber),
				}
				if err := u.pointRepo.WithTx(dbTx).AddTransaction(ctx, pointTx); err != nil {
					return err
				}

				// Update Membership Points
				membership, err := u.membershipRepo.WithTx(dbTx).GetCustomerMembership(ctx, tenantID, order.CustomerID)
				if err == nil && membership != nil {
					membership.Points += pointsEarned
					_ = u.membershipRepo.WithTx(dbTx).UpsertCustomerMembership(ctx, membership)
				}
			}
		}

		// 4. Clear Redis
		activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, order.CustomerID)
		_ = u.redis.Del(ctx, activeKey)

		return nil
	})
}

func (u *OrderService) RejectPayment(ctx context.Context, tenantID, orderID uuid.UUID) error {
	return u.db.Transaction(func(dbTx *gorm.DB) error {
		orderRepo := u.orderRepo.WithTx(dbTx)

		order, err := orderRepo.GetByID(ctx, tenantID, orderID)
		if err != nil {
			return err
		}

		if order.Status != domain.OrderStatusPendingVerification {
			return fmt.Errorf("order is not pending verification, current status: %s", order.Status)
		}

		// 1. Update Status to Cancelled
		if err := orderRepo.UpdateStatus(ctx, tenantID, orderID, domain.OrderStatusCancelled); err != nil {
			return err
		}

		// 2. Release Stock
		var o domain.Order
		if err := dbTx.Preload("Items").First(&o, "id = ?", orderID).Error; err == nil {
			for _, item := range o.Items {
				_ = dbTx.Model(&domain.Product{}).Where("id = ?", item.ProductID).
					Update("inventory", gorm.Expr("inventory + ?", item.Quantity))
			}
		}

		// 3. Update Transaction
		txRepo := u.transactionRepo.WithTx(dbTx)
		tx, err := txRepo.GetByOrderID(ctx, orderID)
		if err == nil && tx != nil {
			tx.Status = domain.TransactionStatusFailed
			_ = txRepo.Update(ctx, tx)
		}

		// 4. Clear Redis
		activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, order.CustomerID)
		_ = u.redis.Del(ctx, activeKey)

		return nil
	})
}

func (u *OrderService) UpdatePaymentMethod(ctx context.Context, tenantID, orderID, customerID, paymentMethodID uuid.UUID) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, tenantID, orderID)
	if err != nil {
		return nil, err
	}

	if order.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized")
	}

	if order.Status != domain.OrderStatusPending {
		return nil, fmt.Errorf("only pending orders can change payment method")
	}

	// 1. Fetch payment method to get new expiry
	pm, err := u.paymentMethodRepo.GetByID(ctx, tenantID, paymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment method")
	}

	newDeadline := time.Now().Add(time.Duration(pm.ExpiryMinutes) * time.Minute)

	// 2. Update order
	order.PaymentMethodID = &paymentMethodID
	order.PaymentDeadline = &newDeadline

	if err := u.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// 3. Update Redis expiry
	activeKey := fmt.Sprintf("retail:order:active:%s:%s", tenantID, customerID)
	_ = u.redis.Set(ctx, activeKey, order.ID.String(), time.Duration(pm.ExpiryMinutes)*time.Minute)

	return order, nil
}

