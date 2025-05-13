package routes

import (
	"ecommerce-cart/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cartHandler *handler.Handler,r *gin.Engine) {

	r.POST("/products/add_product", cartHandler.AddProduct)

	cartRoutes := r.Group("/cart")
	{
		cartRoutes.POST("/add_items", cartHandler.AddItem)
		cartRoutes.DELETE("/remove/:itemID", cartHandler.RemoveItem)
		cartRoutes.DELETE("/remove", cartHandler.RemoveAllItemsHandler)
		cartRoutes.PUT("/update", cartHandler.UpdateItemQuantity)
		cartRoutes.POST("/discount", cartHandler.ApplyDiscount)
		cartRoutes.GET("/view", cartHandler.ViewCart)
		cartRoutes.POST("/checkout", cartHandler.Checkout)
	}

}
