package service

import (
	"context"
	"ecommerce-cart/data"
	"ecommerce-cart/repository"
	"fmt"

	"errors"
)

type Service struct {
	cartRepository *repository.Repository 
}

func NewService(cartRepo *repository.Repository) *Service {
	return &Service{cartRepository: cartRepo}
}

func (s *Service) AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error {
	if price <= 0 {
		return errors.New("price should be greater than zero")
	}
	if productName == "" {
		return errors.New("product name should not be empty")
	}
	if stock <= 0 {
		return errors.New("stock should be greater than zero")
	}

	err := s.cartRepository.AddProductRepo(ctx, productID, productName, price, stock)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddItem(ctx context.Context, productID int32, quantity int32) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be greater than zero")
	}

	count, err := s.cartRepository.AddCartItemRepo(ctx, productID, quantity)
	if err != nil {
		return 0, err
	}

	return count, nil 
}

func (s *Service) RemoveItem(ctx context.Context, itemID int32) error {
	isItemExist, err := s.cartRepository.FindItemRepo(ctx, int(itemID))
	if err != nil {
		return err
	}
	if !isItemExist {
		return errors.New("cart item not found")
	}

	return s.cartRepository.RemoveItemRepo(ctx, int(itemID))
}

func (s *Service) RemoveAllItems(ctx context.Context) error {
	return s.cartRepository.RemoveAllItemRepo(ctx)
}

func (s *Service) UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	itemExists, err := s.cartRepository.FindItemRepo(ctx, int(itemID))
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}
	if !itemExists {
		return fmt.Errorf("no items in the cart to update")
	}

	err = s.cartRepository.UpdateItemQuantityRepo(ctx, itemID, quantity)
	if err != nil {
		return fmt.Errorf("could not update item quantity: %w", err)
	}

	return nil
}


func (s *Service) ApplyDiscount(ctx context.Context, discountType string, discount float64, itemID int32) (float64, error) {
	if discount <= 0 {
		return 0, errors.New("invalid discount value, must be greater than 0")
	}

	var updatedPrice float64
	var err error

	switch discountType {
	case "percentage":
		if discount > 100 {
			return 0, errors.New("invalid discount value, percentage discount cannot exceed 100%")
		}
		updatedPrice, err = s.cartRepository.ApplyPercentDiscountRepo(ctx, discount, itemID)
		if err != nil {
			return 0, err
		}
	case "flat":
		updatedPrice, err = s.cartRepository.ApplyFlatDiscountRepo(ctx, discount, itemID)
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.New("invalid discount type")
	}

	return updatedPrice, nil
}


func (s *Service) ViewCart(ctx context.Context) ([]data.Items, float64, error) {

	items, err := s.cartRepository.ViewCartRepo(ctx)
    if err != nil {
        return nil, 0, errors.New("failed to retrieve cart items")
    }

    totalPrice := 0.0
    for _, item := range items {
        totalPrice += float64(item.Quantity) * item.Price
    }

    return items, totalPrice, nil
}

func (s *Service) Checkout(ctx context.Context) error {
    items, err := s.cartRepository.ViewCartRepo(ctx)
    if err != nil {
        return fmt.Errorf("failed to retrieve cart items: %v", err)
    }

    if len(items) == 0 {
        return fmt.Errorf("cart doesn't contain any items")
    }

    for _, item := range items {
        product, err := s.cartRepository.FetchProductByID(ctx, item.ItemID)
        if err != nil {
            return fmt.Errorf("product with ID %d does not exist", item.ItemID)
        }

        if product.Stock < item.Quantity {
            return fmt.Errorf("insufficient stock for product with ID %d", item.ItemID)
        }
    }

    err = s.cartRepository.Checkout(ctx)
    if err != nil {
        return fmt.Errorf("checkout process failed: %v", err)
    }

    return nil
}