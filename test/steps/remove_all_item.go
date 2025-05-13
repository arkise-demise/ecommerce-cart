package steps

import (
	"bytes"
	"ecommerce-cart/data"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iHaveCartState(cartState string) error {
    switch cartState {
    case "multiple":
        return c.iHaveMultipleItemsInMyCart()
    case "empty":
        return c.iHaveAnEmptyCart()
    default:
        return fmt.Errorf("unknown cart state: %s", cartState)
    }
}

func (c *FeatureContext) iHaveMultipleItemsInMyCart() error {
	items := []data.RequestForAddItem{
		{ProductID: 1001, Quantity: 5},
		{ProductID: 1002, Quantity: 3},
		{ProductID: 1003, Quantity: 2},
		{ProductID: 1004, Quantity: 4},
	}

	for _, item := range items {
		err := c.iAddItemToCart(item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %v", err)
		}
	}

	return nil
}

func (c *FeatureContext) iHaveAnEmptyCart() error {
	req, err := http.NewRequest(http.MethodDelete, c.server.URL+"/cart/remove", nil)
	if err != nil {
		return fmt.Errorf("failed to create request to clear the cart: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to clear the cart: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to clear the cart, received status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *FeatureContext) iRemoveAllItemsFromTheCart() error {
    req, err := http.NewRequest(http.MethodDelete, c.server.URL+"/cart/remove", nil)
    if err != nil {
        return fmt.Errorf("failed to create request to remove all items: %v", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to remove all items from the cart: %v", err)
    }
    c.response = resp

    return nil
}

func (c *FeatureContext) iShouldReceiveConfirmationMessages(expectedMessage string) error {
    if c.response == nil {
        return fmt.Errorf("no response received")
    }

    body, err := io.ReadAll(c.response.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %v", err)
    }

    if !bytes.Contains(body, []byte(expectedMessage)) {
        return fmt.Errorf("expected message %q not found in response: %s", expectedMessage, string(body))
    }

    return nil
}

func (c *FeatureContext) theCartShouldBeEmpty() error {
    req, err := http.NewRequest(http.MethodGet, c.server.URL+"/cart/view", nil)
    if err != nil {
        return fmt.Errorf("failed to create request to view the cart: %v", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to view cart: %v", err)
    }
    defer resp.Body.Close()

    var cartResponse struct {
        CartItems []data.Items `json:"cart_items"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&cartResponse); err != nil {
        return fmt.Errorf("failed to decode cart items: %v", err)
    }

    if len(cartResponse.CartItems) != 0 {
        return fmt.Errorf("cart is not empty, contains %d items", len(cartResponse.CartItems))
    }

    return nil
}

func RemoveAllItemFromCartInitializeScenario(ctx *godog.ScenarioContext) {
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

	ctx.Step(`^I have a (.*) cart$`, feature.iHaveCartState)
	ctx.Step(`^I remove all items from the cart$`, feature.iRemoveAllItemsFromTheCart)
	ctx.Step(`^I should receive a confirmation message "([^"]*)"$`, feature.iShouldReceiveConfirmationMessages)
	ctx.Step(`^The cart should be empty$`, feature.theCartShouldBeEmpty)
}
