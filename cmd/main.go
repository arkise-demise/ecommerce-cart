package main

import (
	"ecommerce-cart/config"
	"ecommerce-cart/data"
	"ecommerce-cart/handler"
	"ecommerce-cart/repository"
	"ecommerce-cart/routes"
	"ecommerce-cart/service"
	"log"

	_ "ecommerce-cart/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ecommerce Cart API
// @version 1.0
// @description This is an ecommerce cart service API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@ecommerce-cart.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @schemes http https
func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	DB, err := config.Connection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer DB.Close()

	DBQueries := data.New(DB)
	cartRepo := repository.NewRepository(DBQueries)
	Service := service.NewService(cartRepo)
	cartHandler := handler.NewHandler(Service)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRouter(cartHandler, r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
