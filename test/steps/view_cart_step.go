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

func (c *FeatureContext) iAddItemsToCartAndViewThem(products *godog.Table) error {
	if c.cartProductIDs == nil {
		c.cartProductIDs = make(map[int32]bool)
	}

	for _, row := range products.Rows[1:] {
		itemID, err := strconv.Atoi(row.Cells[0].Value)
		if err != nil {
			return fmt.Errorf("invalid ItemID %q: %v", row.Cells[0].Value, err)
		}

		quantity, err := strconv.Atoi(row.Cells[1].Value)
		if err != nil {
			return fmt.Errorf("invalid Quantity %q: %v", row.Cells[1].Value, err)
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
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status code %d, response: %s", resp.StatusCode, string(body))
		}

		c.cartProductIDs[int32(itemID)] = true
	}

	resp, err := http.Get(c.server.URL + "/cart/view")
	if err != nil {
		return fmt.Errorf("failed to view cart: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	c.response = resp
	c.responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	return nil
}


func (c *FeatureContext) iShouldSeeTheFollowingItemsInTheCart(items *godog.Table) error {
    var response map[string]interface{}
    
    err := json.Unmarshal(c.responseBody, &response)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    cartItems, ok := response["items"].([]interface{})
    if !ok {
        return fmt.Errorf("items field not found in response or invalid type")
    }

    for i, row := range items.Rows[1:] { 
        expectedItemID, err := strconv.Atoi(row.Cells[0].Value)
        if err != nil {
            return fmt.Errorf("invalid ItemID: %v", err)
        }
        expectedItemName := row.Cells[1].Value
        expectedQuantity, err := strconv.Atoi(row.Cells[2].Value)
        if err != nil {
            return fmt.Errorf("invalid Quantity: %v", err)
        }
        expectedPricePerItem, err := strconv.ParseFloat(row.Cells[3].Value, 64)
        if err != nil {
            return fmt.Errorf("invalid Price: %v", err)
        }

        if i >= len(cartItems) {
            return fmt.Errorf("cart item index %d out of range", i)
        }

        item, ok := cartItems[i].(map[string]interface{})
        if !ok {
            return fmt.Errorf("cart item is not in the expected format")
        }

        actualItemID := int(item["item_id"].(float64))
        actualItemName := item["item_name"].(string)
        actualQuantity := int(item["quantity"].(float64))
        actualPricePerItem := item["price"].(float64)

        if actualItemID != expectedItemID || actualItemName != expectedItemName || actualQuantity != expectedQuantity || actualPricePerItem != expectedPricePerItem {
            return fmt.Errorf("expected itemID %d with name %q, quantity %d, and price %.2f, but got itemID %d with name %q, quantity %d, and price %.2f",
                expectedItemID, expectedItemName, expectedQuantity, expectedPricePerItem,
                actualItemID, actualItemName, actualQuantity, actualPricePerItem)
        }
    }

    return nil
}


func (c *FeatureContext) theTotalPriceOfTheCartShouldBe(expectedTotalPriceStr string) error {
	expectedTotalPrice, err := strconv.ParseFloat(expectedTotalPriceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid total price: %v", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(c.responseBody, &responseBody)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	totalPrice, ok := responseBody["totalPrice"].(float64)
	if !ok {
		return fmt.Errorf("total Price field not found in response or invalid type")
	}

	if totalPrice != expectedTotalPrice {
		return fmt.Errorf("expected total price %.2f, but got %.2f", expectedTotalPrice, totalPrice)
	}

	return nil
}

func CartViewScenario(ctx *godog.ScenarioContext) {
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

	ctx.Step(`^I add the following items to my cart and view them:$`, feature.iAddItemsToCartAndViewThem)
	ctx.Step(`^I should see the following items:$`, feature.iShouldSeeTheFollowingItemsInTheCart)
	ctx.Step(`^the total price of the cart should be "([^"]*)"$`, feature.theTotalPriceOfTheCartShouldBe)
}

