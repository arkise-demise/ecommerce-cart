package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iHaveAddedTheFollowingItemsToCart(products *godog.Table) error {
    for _, row := range products.Rows[1:] {
        productID, err := strconv.Atoi(row.Cells[0].Value) 
        if err != nil {
            return fmt.Errorf("invalid itemID format: %v", err)
        }
        quantity, err := strconv.Atoi(row.Cells[1].Value) 
        if err != nil {
            return fmt.Errorf("invalid quantity format: %v", err)
        }

        payload := map[string]interface{}{
            "product_id":  int32(productID),   
            "quantity": int32(quantity), 
        }

        jsonData, err := json.Marshal(payload)
        if err != nil {
            return fmt.Errorf("failed to marshal request data: %v", err)
        }

        resp, err := http.Post(c.server.URL+"/cart/add_items", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            return fmt.Errorf("failed to add item to cart: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            return fmt.Errorf("failed to add item to cart: %s", string(body))
        }
    }

    return nil
}

func (c *FeatureContext) iUpdateTheItemWithIDToQuantity(itemID, quantity int) error {
	payload := map[string]interface{}{
		"item_id":  itemID,
		"quantity": quantity,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.server.URL+"/cart/update", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update item in cart: %v", err)
	}
	defer resp.Body.Close()

	c.response = resp
	c.responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	viewResp, err := http.Get(c.server.URL + "/cart/view")
	if err != nil {
		return fmt.Errorf("failed to view cart after update: %v", err)
	}
	defer viewResp.Body.Close()

	c.viewCartResponseBody, err = io.ReadAll(viewResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read view cart response body: %v", err)
	}

	return nil
}

func (c *FeatureContext) iShouldSeeAConfirmationMessage(expectedMessage string) error {
	if c.response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, but got %d", c.response.StatusCode)
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(c.responseBody, &responseBody)
	if err != nil {
		return err
	}

	actualMessage, ok := responseBody["message"].(string)
	if !ok || expectedMessage != actualMessage {
		return fmt.Errorf("expected message %q, but got %q", expectedMessage, actualMessage)
	}

	return nil
}

func (c *FeatureContext) iShouldSeeAnErrorMessage(expectedMessage string) error {
	if c.response.StatusCode != http.StatusBadRequest && c.response.StatusCode != http.StatusNotFound && c.response.StatusCode != http.StatusConflict {
		return fmt.Errorf("expected status code 400, 404, or 409, but got %d", c.response.StatusCode)
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(c.responseBody, &responseBody)
	if err != nil {
		return err
	}

	actualError, ok := responseBody["error"].(string)
	if !ok || expectedMessage != actualError {
		return fmt.Errorf("expected error message %q, but got %q", expectedMessage, actualError)
	}

	return nil
}

func (c *FeatureContext) theCartShouldContainUnitsOfTheItemWithID(quantity, itemID int) error {
	var viewCartResponse map[string]interface{}
	err := json.Unmarshal(c.viewCartResponseBody, &viewCartResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal view cart response: %v", err)
	}

	items, ok := viewCartResponse["items"].([]interface{})
	if !ok {
		return fmt.Errorf("items not found in the view cart response")
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if int(itemMap["item_id"].(float64)) == itemID {
			actualQuantity := int(itemMap["quantity"].(float64))
			if actualQuantity != quantity {
				return fmt.Errorf("expected %d units of item %d, but got %d", quantity, itemID, actualQuantity)
			}
			return nil
		}
	}

	return fmt.Errorf("item with ID %d not found in the cart", itemID)
}


func CartUpdateQuantityInitializeScenario(ctx *godog.ScenarioContext) {
	feature := &FeatureContext{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		feature.initialize()
	})

	ctx.AfterScenario(func(*godog.Scenario, error) {
		feature.tearDown()
	})

	ctx.Step(`^I have added the following items to my cart$`, feature.iHaveAddedTheFollowingItemsToCart)
	ctx.Step(`^I update the item with ID (\d+) to quantity (\d+)$`, feature.iUpdateTheItemWithIDToQuantity)
	ctx.Step(`^I should see a confirmation message "([^"]*)"$`, feature.iShouldSeeAConfirmationMessage)
	ctx.Step(`^I should see an error message "([^"]*)"$`, feature.iShouldSeeAnErrorMessage)
	ctx.Step(`^the cart should contain (\d+) units of the item with ID (\d+)$`, feature.theCartShouldContainUnitsOfTheItemWithID)
}
