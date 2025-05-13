package repository

import (
	"context"
	"ecommerce-cart/data"
)

type CartRepoInterface interface {
	AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error
	AddItem(ctx context.Context, productID int32, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int) error
	RemoveAllItem(ctx context.Context) error
	FindItem(ctx context.Context, itemID int) (bool, error)
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discountType string, discount float64, itemID int32) (float64, error)
	ViewCart(ctx context.Context) ([]data.Items, error)
	Checkout(ctx context.Context) error
	FetchProductByID(ctx context.Context, productID int32) (data.Products, error) 

	
}
