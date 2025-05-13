package repository

import (
	"context"
	"ecommerce-cart/data"
	"errors"
	"fmt"
)



type Repository struct {
	dbQueries *data.Queries
}

func NewRepository(dbQueries *data.Queries) *Repository {
	return &Repository{
		dbQueries: dbQueries,
	}
}

func (r *Repository) AddProductRepo(ctx context.Context, productID int32, productName string, price float64, stock int32) error {
	params := data.InsertProductParams{
		ProductID:   productID,
		ProductName: productName,
		Price:       price,
		Stock:       stock,
	}

	err := r.dbQueries.InsertProduct(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetProduct(ctx context.Context, productID int32) (*data.Products, error) {
	product, err := r.dbQueries.FetchProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *Repository) AddCartItemRepo(ctx context.Context, productID int32, quantity int32) (int, error) {
    product, err := r.dbQueries.FetchProductByID(ctx, productID)
    if err != nil {
        return 0, errors.New("product does not exist")
    }

    count, err := r.dbQueries.GetUniqueItemCountInCart(ctx)
    if err != nil {
        return int(count), err
    }

 
    if count >= 10 {
        return int(count), errors.New("can't add more than 10 unique items to the cart")
    }

    if quantity > product.Stock {
        return int(count), errors.New("requested quantity exceeds available stock")
    }

    err = r.dbQueries.InsertCartItem(ctx, data.InsertCartItemParams{
        ItemID:   product.ProductID,
        ItemName: product.ProductName,
        Price:    product.Price,
        Quantity: quantity,
    })
    if err != nil {
        return int(count), err
    }

    err = r.dbQueries.UpdateProductStockLevel(ctx, data.UpdateProductStockLevelParams{
        ProductID: productID,
        Stock:     product.Stock - quantity,
    })
    if err != nil {
        return int(count + 1), err
    }

    return int(count + 1), nil
}




func (r *Repository) RemoveItemRepo(ctx context.Context, itemID int) error {
	err := r.dbQueries.DeleteCartItem(ctx, int32(itemID))
	if err != nil {
		return errors.New("failed to remove item from the cart")
	}
	return nil
}

func (r *Repository) RemoveAllItemRepo(ctx context.Context) error {
	err := r.dbQueries.ClearCartItems(ctx)
	if err != nil {
		return errors.New("failed to remove all items from the cart")
	}
	return nil
}

func (r *Repository) FindItemRepo(ctx context.Context, itemID int) (bool, error) {
	exists, err := r.dbQueries.CheckItemInCart(ctx, int32(itemID))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdateItemQuantityRepo(ctx context.Context, itemID int32, quantity int32) error {
	existingItem, err := r.dbQueries.GetItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("error checking if item exists: %w", err)
	}

	if existingItem.ItemID == 0 {
		return fmt.Errorf("no items in the cart to update")
	}

	params := data.ModifyItemQuantityParams{
		ItemID:   itemID,
		Quantity: quantity,
	}

	err = r.dbQueries.ModifyItemQuantity(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update the item quantity for itemID %d: %w", itemID, err)
	}

	return nil
}

func (r *Repository) ApplyPercentDiscountRepo(ctx context.Context, discount float64, itemID int32) (float64, error) {
	params := data.ApplyPercentageDiscountParams{
		Price: discount,
		ItemID:   itemID,
	}

	err := r.dbQueries.ApplyPercentageDiscount(ctx, params)
	if err != nil {
		return 0, errors.New("failed to apply percentage discount to the cart")
	}

	updatedItem, err := r.dbQueries.GetItemByID(ctx, itemID)
	if err != nil {
		return 0, errors.New("failed to retrieve updated item")
	}

	return updatedItem.Price, nil
}


func (r *Repository) ApplyFlatDiscountRepo(ctx context.Context, discount float64, itemID int32) (float64, error) {
	params := data.ApplyFlatDiscountParams{
		Price: discount,
		ItemID: itemID,
	}

	err := r.dbQueries.ApplyFlatDiscount(ctx, params)
	if err != nil {
		return 0, errors.New("failed to apply flat discount to the item")
	}

	updatedItem, err := r.dbQueries.GetItemByID(ctx, itemID)
	if err != nil {
		return 0, errors.New("failed to retrieve updated item")
	}

	return updatedItem.Price, nil
}


func (r *Repository) ViewCartRepo(ctx context.Context) ([]data.Items, error) {
    items, err := r.dbQueries.GetCartItems(ctx)
    if err != nil {
        return nil, errors.New("failed to retrieve cart items")
    }

    return items, nil
}

func (r *Repository) Checkout(ctx context.Context) error {
    count, err := r.dbQueries.CheckCartIsEmpty(ctx)
    if err != nil {
        return errors.New("failed to check cart status")
    }
    if count == 0 {
        return fmt.Errorf("cart doesn't contain any items")
    }

    items, err := r.dbQueries.GetCartItems(ctx)
    if err != nil {
        return errors.New("failed to retrieve cart items")
    }

    for _, item := range items {
        product, err := r.dbQueries.FetchProductByID(ctx, item.ItemID)
        if err != nil {
            return fmt.Errorf("product with ID %d does not exist", item.ItemID)
        }

        if product.Stock < item.Quantity {
            return fmt.Errorf("insufficient stock for product with ID %d", item.ItemID)
        }
    }

    if err := r.dbQueries.CompleteCheckout(ctx); err != nil {
        return errors.New("failed to complete checkout")
    }

    return nil
}

func (r *Repository) FetchProductByID(ctx context.Context, productID int32) (data.Products, error) {
    return r.dbQueries.FetchProductByID(ctx, productID)
}