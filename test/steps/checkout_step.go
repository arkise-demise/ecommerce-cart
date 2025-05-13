package steps

import (
	"bytes"
	"ecommerce-cart/data"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iAddItemsAndProceedToCheckout(products *godog.Table) error {
    for _, row := range products.Rows[1:] {
        itemID, err := strconv.Atoi(row.Cells[0].Value)
        if err != nil {
            return fmt.Errorf("invalid ItemID: %v", err)
        }

        quantity, err := strconv.Atoi(row.Cells[1].Value)
        if err != nil {
            return fmt.Errorf("invalid Quantity: %v", err)
        }

        payload := data.RequestForAddItem{
            ProductID: int32(itemID),
            Quantity:  int32(quantity),
        }

        jsonData, err := json.Marshal(payload)
        if err != nil {
            return fmt.Errorf("failed to marshal request data: %v", err)
        }

        resp, err := http.Post(c.server.URL+"/cart/add_items", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            return fmt.Errorf("failed to add item to cart: %v", err)
        }

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
        }

        c.cartProductIDs[int32(itemID)] = true
    }

    resp, err := http.Post(c.server.URL+"/cart/checkout", "application/json", nil)
    if err != nil {
        return fmt.Errorf("failed to checkout: %v", err)
    }

    c.response = resp
    c.responseBody, err = io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
    }

    return nil
}

func (c *FeatureContext) theCheckoutShouldBeSuccessful() error {
    var responseBody map[string]interface{}
    err := json.Unmarshal(c.responseBody, &responseBody)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    message, ok := responseBody["message"].(string)
    if !ok || message != "Checkout is successful!" {
        return fmt.Errorf("expected success message not found, got: %v", message)
    }

    return nil
}

func (c *FeatureContext) theTotalPriceShouldBe(expectedTotalStr string) error {
    var responseBody map[string]interface{}
    err := json.Unmarshal(c.responseBody, &responseBody)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    totalPrice, ok := responseBody["totalPrice"].(float64)
    if !ok {
        return fmt.Errorf("total_price field not found in response")
    }

    expectedTotal, err := strconv.Atoi(expectedTotalStr)
    if err != nil {
        return fmt.Errorf("invalid total price: %v", err)
    }

    if int(totalPrice) != expectedTotal {
        return fmt.Errorf("expected total price %d, but got %d", expectedTotal, int(totalPrice))
    }

    return nil
}

func (c *FeatureContext) iHaveNotAddedAnyItemsToMyCart() error {
	c.cartProductIDs = nil
	return nil
}
func (c *FeatureContext) iProceedToCheckout() error {
    resp, err := http.Post(c.server.URL+"/cart/checkout", "application/json", nil)
    if err != nil {
        return fmt.Errorf("failed to checkout: %v", err)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %v", err)
    }

    c.response = resp
    c.responseBody = body

    return nil
}

func (c *FeatureContext) iShouldReceiveAnErrorMessage(expectedErrorMessage string) error {
    var responseBody map[string]interface{}
    err := json.Unmarshal(c.responseBody, &responseBody)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    errorMessage, ok := responseBody["error"].(string)
    if !ok {
        return fmt.Errorf("error message field not found in response")
    }

    if errorMessage != expectedErrorMessage {
        return fmt.Errorf("expected error message %q, but got %q", expectedErrorMessage, errorMessage)
    }

    return nil
}

func CartCheckoutScenario(ctx *godog.ScenarioContext) {
	feature := &FeatureContext{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		err := feature.initialize()
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize scenario: %v", err))
		}
	})

	ctx.AfterScenario(func(*godog.Scenario, error) {
		feature.tearDown()
	})

	ctx.Step(`^I have added the following items to my cart and proceed to checkout:$`, feature.iAddItemsAndProceedToCheckout)
	ctx.Step(`^the checkout should be successful$`, feature.theCheckoutShouldBeSuccessful)
	ctx.Step(`^the total price should be "([^"]*)"$`, feature.theTotalPriceShouldBe)
	ctx.Step(`^I have not added any items to my cart$`, feature.iHaveNotAddedAnyItemsToMyCart)
	ctx.Step(`I proceed to checkout$`,feature.iProceedToCheckout)
	ctx.Step(`^I should receive an error message "([^"]*)"$`, feature.iShouldReceiveAnErrorMessage)
}
