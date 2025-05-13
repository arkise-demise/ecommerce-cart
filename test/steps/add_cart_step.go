package steps

import (
	"bytes"
	"ecommerce-cart/data"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iHaveTheFollowingProductsInStock(products *godog.Table) error {
	for _, row := range products.Rows[1:] { 
		productID, err := strconv.Atoi(row.Cells[0].Value)
		if err != nil {
			return fmt.Errorf("invalid ProductID: %v", err)
		}
		quantity, err := strconv.Atoi(row.Cells[1].Value)
		if err != nil {
			return fmt.Errorf("invalid Quantity: %v", err)
		}
		err = c.iAddItemToCart(int32(productID), int32(quantity))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FeatureContext) iAddTheseItemsToMyCart() error {
	if c.response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(c.response.Body)
		return fmt.Errorf("failed to add items to cart, status: %v, response: %s", c.response.StatusCode, string(body))
	}
	return nil
}

func (c *FeatureContext) iShouldReceive(successMessage string) error {

	body, err := io.ReadAll(c.response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if !bytes.Contains(body, []byte(successMessage)) {
		return fmt.Errorf("expected success message %q not found in response: %s", successMessage, string(body))
	}

	return nil
}

func (c *FeatureContext) iAlreadyHaveTheFollowingUniqueItemsInMyCart(products *godog.Table) error {
	if c.cartProductIDs == nil {
		c.cartProductIDs = make(map[int32]bool)
	}

	for _, row := range products.Rows[1:] { 
		productID, err := strconv.Atoi(row.Cells[0].Value)
		if err != nil {
			return fmt.Errorf("invalid ProductID: %v", err)
		}
		quantity, err := strconv.Atoi(row.Cells[1].Value)
		if err != nil {
			return fmt.Errorf("invalid Quantity: %v", err)
		}
		err = c.iAddItemToCart(int32(productID), int32(quantity))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *FeatureContext) iTryToAddOneMoreUniqueItemToTheCart(table *godog.Table) error {
    for _, row := range table.Rows[1:] { 
        productID, err := strconv.Atoi(row.Cells[0].Value)
        if err != nil {
            return fmt.Errorf("invalid ProductID: %v", err)
        }
        quantity, err := strconv.Atoi(row.Cells[1].Value)
        if err != nil {
            return fmt.Errorf("invalid Quantity: %v", err)
        }

        err = c.iAddItemToCart(int32(productID), int32(quantity))
		if err.Error() == "can't add more than 10 unique items to the cart" {
            return nil
        }
    }
    return nil
}

func (c *FeatureContext) iShouldSeeAnError(expectedError string) error {

    body, err := io.ReadAll(c.response.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %v", err)
    }

    if !bytes.Contains(body, []byte(expectedError)) {
        return fmt.Errorf("expected error %q not found in response: %s", expectedError, string(body))
    }

    return nil
}

func (c *FeatureContext) iAddItemToCart(productID int32, quantity int32) error {
	if c.cartProductIDs[productID] {
		return fmt.Errorf("product %d is already in the cart", productID)
	}

	payload := data.RequestForAddItem{
		ProductID: productID,
		Quantity:  quantity,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %v", err)
	}

	resp, err := http.Post(c.server.URL+"/cart/add_items", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to add item to cart: %v", err)
	}

	c.response = resp

	if len(c.cartProductIDs) >= 10 {
		return fmt.Errorf("can't add more than 10 unique items to the cart")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add item to cart, status: %v, response: %s", resp.StatusCode, string(body))
	}

	if c.cartProductIDs == nil {
		c.cartProductIDs = make(map[int32]bool)
	}

	if !c.cartProductIDs[productID] {
		c.cartItems += 1
		c.cartProductIDs[productID] = true
	}

	return nil
}

func CartInitializeScenario(ctx *godog.ScenarioContext) {
    feature := &FeatureContext{}

    ctx.BeforeScenario(func(*godog.Scenario) {
        err := feature.initialize()
        if err != nil {
            log.Fatalf("Failed to initialize scenario: %v", err)
        }
    })

    ctx.AfterScenario(func(*godog.Scenario, error) {
        feature.tearDown()
    })

    ctx.Step(`^I have the following products in stock:$`, feature.iHaveTheFollowingProductsInStock)
    ctx.Step(`^I add these items to my cart$`, feature.iAddTheseItemsToMyCart)
    ctx.Step(`^I should receive "([^"]*)"$`, feature.iShouldReceive)
    ctx.Step(`^I already have the following unique items in my cart:$`, feature.iAlreadyHaveTheFollowingUniqueItemsInMyCart)
    ctx.Step(`^I try to add one more unique item to the cart:$`, feature.iTryToAddOneMoreUniqueItemToTheCart)
    ctx.Step(`^I should see an error "([^"]*)"$`, feature.iShouldSeeAnError)
}
