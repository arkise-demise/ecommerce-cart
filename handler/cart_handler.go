package handler

import (
	"context"
	"ecommerce-cart/data"
	"ecommerce-cart/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{service: service}
}

// AddProduct godoc
// @Summary Add a new product
// @Description Adds a new product to the store
// @Tags Products
// @Accept json
// @Produce json
// @Param product body data.Products true "Product details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /products/add_product [post]
func (h *Handler) AddProduct(c *gin.Context) {
	var product data.Products

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err := h.service.AddProduct(context.Background(), product.ProductID, product.ProductName, product.Price, product.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "products": product})
}

// AddItem godoc
// @Summary Add items to the cart
// @Description Adds a specified quantity of a product to the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param item body data.RequestForAddItem true "Item to add"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /cart/add_items [post]
func (h *Handler) AddItem(c *gin.Context) {
	var itemAdded data.RequestForAddItem

	if err := c.ShouldBindJSON(&itemAdded); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	cartItems, err := h.service.AddItem(context.Background(), itemAdded.ProductID, itemAdded.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "items successfully added to the cart",
		"total_items_added": cartItems,
		"added_product":     itemAdded,
	})
}

// RemoveItem godoc
// @Summary Remove item from the cart
// @Description Removes a specific item from the cart by item ID
// @Tags Cart
// @Param itemID path int true "Item ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /cart/remove/{itemID} [delete]
func (h *Handler) RemoveItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid itemID"})
		return
	}

	err = h.service.RemoveItem(context.Background(), int32(itemID))
	if err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// RemoveAllItemsHandler godoc
// @Summary Remove all items from the cart
// @Description Clears all items from the cart
// @Tags Cart
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/remove [delete]
func (h *Handler) RemoveAllItemsHandler(c *gin.Context) {
	err := h.service.RemoveAllItems(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove all items from the cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All items removed from the cart"})
}

// UpdateItemQuantity godoc
// @Summary Update item quantity in the cart
// @Description Updates the quantity of a specific item in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param updateRequest body data.RequestForUpdate true "Update item quantity"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /cart/update [put]
func (h *Handler) UpdateItemQuantity(c *gin.Context) {
	var updateRequest data.RequestForUpdate

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, please check the item ID and quantity format"})
		return
	}

	if updateRequest.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than zero"})
		return
	}

	err := h.service.UpdateItemQuantity(c.Request.Context(), updateRequest.ItemID, updateRequest.Quantity)
	if err != nil {
		if err.Error() == "no items in the cart to update" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in the cart"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item's quantity is updated!"})
}

// ApplyDiscount godoc
// @Summary Apply a discount to an item
// @Description Applies a percentage or fixed discount to a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param discountRequest body data.RequestForDiscount true "Discount details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /cart/discount [post]
func (h *Handler) ApplyDiscount(c *gin.Context) {
	var discountRequest data.RequestForDiscount

	if err := c.ShouldBindJSON(&discountRequest); err != nil {
		fmt.Printf("Invalid discount request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	updatedPrice, err := h.service.ApplyDiscount(c.Request.Context(), discountRequest.DiscountType, discountRequest.Discount, discountRequest.ItemID)
	if err != nil {
		fmt.Printf("Error applying discount: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("%.2f %s discount applied to item %d", discountRequest.Discount, discountRequest.DiscountType, discountRequest.ItemID),
		"newPrice": updatedPrice,
	})
}

// ViewCart godoc
// @Summary View all items in the cart
// @Description Retrieves all items currently in the cart and calculates the total price
// @Tags Cart
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/view [get]
func (h *Handler) ViewCart(c *gin.Context) {
	items, totalPrice, err := h.service.ViewCart(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalPrice": totalPrice,
		"items":      items,
	})
}

// Checkout godoc
// @Summary Checkout the cart
// @Description Finalizes the cart and proceeds with the checkout process
// @Tags Cart
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/checkout [post]
func (h *Handler) Checkout(c *gin.Context) {
	items, totalPrice, err := h.service.ViewCart(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Checkout(context.Background())
	if err != nil {
		if err.Error() == "cart doesn't contain any items" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Checkout is successful!",
		"items":      items,
		"totalPrice": totalPrice,
	})
}
