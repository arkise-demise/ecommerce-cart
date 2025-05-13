package steps

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (c *FeatureContext) iHaveTheFollowingItemsInMyCart(items *godog.Table) error {
	for _, row := range items.Rows[1:] {
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

func (c *FeatureContext) iPerformActionOnItem(action string, productIDStr string) error {
	_, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid ProductID: %v", err)
	}

	switch action {
	case "remove":
		return c.iRemoveTheItemWithProductID(productIDStr)
	case "try to remove":
		return c.iTryToRemoveTheItemWithProductID(productIDStr)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func (c *FeatureContext) iRemoveTheItemWithProductID(productIDStr string) error {
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid ProductID: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/cart/remove/%d", c.server.URL, productID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request to remove item: %v", err)
	}

	c.response, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to remove item: %v", err)
	}
	return nil
}

func (c *FeatureContext) iShouldSeeMessage(message string) error {
	if c.response == nil {
		return fmt.Errorf("no response received")
	}

	body, err := io.ReadAll(c.response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if !bytes.Contains(body, []byte(message)) {
		return fmt.Errorf("expected message %q not found in response: %s", message, string(body))
	}

	return nil
}

func (c *FeatureContext) iTryToRemoveTheItemWithProductID(productIDStr string) error {
	return c.iRemoveTheItemWithProductID(productIDStr)
}

func (c *FeatureContext) iShouldSee(resultMessage string) error {
	if c.response == nil {
		return fmt.Errorf("no response received")
	}

	body, err := io.ReadAll(c.response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	defer c.response.Body.Close()

	if !bytes.Contains(body, []byte(resultMessage)) {
		return fmt.Errorf("expected message %q not found in response: %s", resultMessage, string(body))
	}

	return nil
}

func RemoveItemInitializeScenario(ctx *godog.ScenarioContext) {
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

	ctx.Step(`^I have the following items in my cart:$`, feature.iHaveTheFollowingItemsInMyCart)
	ctx.Step(`^I (remove|try to remove) the item with ProductID (\d+)$`, feature.iPerformActionOnItem)
	ctx.Step(`^I should see the message "([^"]*)"$`, feature.iShouldSeeMessage)
	ctx.Step(`^I should see "([^"]*)"$`, feature.iShouldSee)
}
