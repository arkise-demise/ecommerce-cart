package service

import (
	"context"
	"ecommerce-cart/data"
)

type ServiceInterface interface {
	AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error
	AddItem(ctx context.Context, productID int32, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int32) error
	RemoveAllItems(ctx context.Context) error
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discountType string, discount float64, itemID int32) (float64, error)  // Updated signature
	ViewCart(ctx context.Context) ([]data.Items, float64, error)
	Checkout(ctx context.Context) error
}
