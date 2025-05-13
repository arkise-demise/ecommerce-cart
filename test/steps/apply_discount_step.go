package steps

import (
	"bytes"
	"ecommerce-cart/data"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iHaveAddedTheFollowingItemsToMyCart(products *godog.Table) error {
    if c.cartProductIDs == nil {
        c.cartProductIDs = make(map[int32]bool)
    }

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
            body, _ := io.ReadAll(resp.Body) 
            return fmt.Errorf("unexpected status code: %v, response body: %s", resp.StatusCode, string(body))
        }

        c.response = resp
        c.cartProductIDs[int32(itemID)] = true
    }

    return nil
}


func (c *FeatureContext) iApplyADiscountToTheItemInCart(discountType string, discountValue, itemIDStr string) error {
    itemID, err := strconv.Atoi(itemIDStr)
    if err != nil {
        return fmt.Errorf("invalid ItemID: %v", err)
    }

    discount, err := strconv.ParseFloat(discountValue, 64)
    if err != nil {
        return fmt.Errorf("invalid Discount: %v", err)
    }

    payload := data.RequestForDiscount{
        DiscountType: discountType,
        Discount:     discount,
        ItemID:       int32(itemID),  
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal request data: %v", err)
    }

    resp, err := http.Post(c.server.URL+"/cart/discount", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to send request to apply discount: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body) 
        return fmt.Errorf("unexpected status code: %v, response body: %s", resp.StatusCode, string(body))
    }

    c.response = resp
    c.responseBody, err = io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %v", err)
    }

    return nil
}

func (c *FeatureContext) iShouldReceiveSuccessMessage(successMessage string) error {
    var responseBody map[string]interface{}
    err := json.Unmarshal(c.responseBody, &responseBody)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    actualMessage, ok := responseBody["message"].(string)
    if !ok {
        return fmt.Errorf("expected success message not found in response body")
    }

    if !strings.Contains(actualMessage, successMessage) {
        return fmt.Errorf("expected success message %q not found in response: %s", successMessage, actualMessage)
    }

    return nil
}

func (c *FeatureContext) theNewPriceOfTheItemShouldBe(expectedPriceStr string) error {
    expectedPrice, err := strconv.ParseFloat(expectedPriceStr, 64)
    if err != nil {
        return fmt.Errorf("invalid price: %v", err)
    }

    var responseBody map[string]interface{}
    err = json.Unmarshal(c.responseBody, &responseBody)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    newPrice, ok := responseBody["newPrice"].(float64)
    if !ok {
        return fmt.Errorf("newPrice field not found in response or invalid type")
    }

    if newPrice != expectedPrice {
        return fmt.Errorf("expected new price %v, but got %v", expectedPrice, newPrice)
    }

    return nil
}

func CartDiscountScenario(ctx *godog.ScenarioContext) {
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

    ctx.Step(`^I have added the following items to my cart:$`, feature.iHaveAddedTheFollowingItemsToMyCart)
    ctx.Step(`^I apply a "([^"]*)" discount of "([^"]*)" to the item "([^"]*)" in the cart$`, feature.iApplyADiscountToTheItemInCart)
    ctx.Step(`^I should receive a success message "([^"]*)"$`, feature.iShouldReceiveSuccessMessage)
    ctx.Step(`^the new price of the item should be "([^"]*)"$`, feature.theNewPriceOfTheItemShouldBe)
}
