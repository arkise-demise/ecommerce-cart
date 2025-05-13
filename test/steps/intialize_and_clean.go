package steps

import (
	"context"
	"database/sql"
	"ecommerce-cart/config"
	"ecommerce-cart/data"
	"ecommerce-cart/handler"
	"ecommerce-cart/repository"
	"ecommerce-cart/routes"
	"ecommerce-cart/service"

	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

type FeatureContext struct {
	cartHandler    *handler.Handler
	ctx            context.Context
	server         *httptest.Server
	response       *http.Response
	responseBody   []byte 
	cartError      error
	tx             *sql.Tx
	products       []data.Products 
	db             *sql.DB
	cartProductIDs map[int32]bool
	productID      int32
	quantity       int32
	cartItems      int32
	cartRepository *repository.Repository
	errMessage     string
	viewCartResponseBody []byte
}

func setupTestRouter(tx *sql.Tx) (*gin.Engine, error) {
	_, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	queries := data.New(tx)
	cartRepo := repository.NewRepository(queries)
	cartService := service.NewService(cartRepo)
	cartHandler := handler.NewHandler(cartService)

	r := gin.Default()
	routes.SetupRouter(cartHandler,r)

	return r, nil
}

func (c *FeatureContext) initialize() error {
	c.ctx = context.Background()
	db, err := config.Connection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	c.tx, err = db.BeginTx(c.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	queries := data.New(c.tx) 
    c.cartRepository = repository.NewRepository(queries) 

	router, err := setupTestRouter(c.tx)
	if err != nil {
		return err
	}
	c.cartProductIDs = make(map[int32]bool)

	c.server = httptest.NewServer(router)

	if c.server == nil {
		return fmt.Errorf("server is not initialized")
	}

	fmt.Println("Server and transaction successfully initialized.")
	return nil
}

func (c *FeatureContext) tearDown() {
	if c.server != nil {
		c.server.Close()
	}

	if c.tx != nil {
		if err := c.tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("Failed to rollback transaction: %v\n", err)
		} else {
			fmt.Println("Transaction rolled back successfully.")
		}
	}
}

func (c *FeatureContext) BeforeScenario(sc *godog.Scenario) {
	if err := c.initialize(); err != nil {
		fmt.Printf("Failed to initialize before scenario: %v\n", err)
	}
}

func (c *FeatureContext) AfterScenario(sc *godog.Scenario) {
	c.tearDown()
}