package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
)

type cartService struct {
	cartRepo    repository.ICartRepository
	productRepo repository.IProductRepository
}

func NewCartService(cartRepo repository.ICartRepository, productRepo repository.IProductRepository) ICartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *cartService) GetCart(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Cart, error) {
	return s.cartRepo.GetByCustomerID(ctx, tenantID, customerID)
}

func (s *cartService) AddItem(ctx context.Context, tenantID, customerID uuid.UUID, productID uuid.UUID, quantity int) error {
	// 1. Check Inventory
	product, err := s.productRepo.GetByID(ctx, tenantID, productID)
	if err != nil {
		return err
	}
	if product.Inventory < quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.GetByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		return err
	}

	item := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		Selected:  true, // Default selected by user request
	}

	return s.cartRepo.UpsertItem(ctx, item)
}

func (s *cartService) UpdateItemQty(ctx context.Context, tenantID, customerID, itemID uuid.UUID, quantity int) error {
	// 1. Fetch Cart to find the product ID for this item
	cart, err := s.cartRepo.GetByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		return err
	}

	// Find the item to get ProductID
	var productID uuid.UUID
	found := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			productID = item.ProductID
			found = true
			break
		}
	}

	if !found {
		return errors.New("item not found")	
	}

	if quantity > 0 {
		// 2. Check Inventory
		product, err := s.productRepo.GetByID(ctx, tenantID, productID)
		if err != nil {
			return err
		}
		if product.Inventory < quantity {
			return errors.New("insufficient stock")
		}
	}
	
	if quantity <= 0 {
		return s.cartRepo.DeleteItem(ctx, itemID)
	}
	return s.cartRepo.UpdateItemQuantity(ctx, itemID, quantity)
}

func (s *cartService) RemoveItem(ctx context.Context, tenantID, customerID, itemID uuid.UUID) error {
	return s.cartRepo.DeleteItem(ctx, itemID)
}

func (s *cartService) ToggleSelection(ctx context.Context, tenantID, customerID, itemID uuid.UUID, selected bool) error {
	return s.cartRepo.ToggleSelectItem(ctx, itemID, selected)
}

func (s *cartService) ClearCart(ctx context.Context, tenantID, customerID uuid.UUID) error {
	cart, err := s.cartRepo.GetByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		return err
	}
	return s.cartRepo.ClearCart(ctx, cart.ID)
}
