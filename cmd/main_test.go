package main

import (
	"ecommerce-cart/test/steps"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
		Paths:  []string{"../test/features"},
	}

	status := godog.TestSuite{
		Name:                 "ecommerce-cart",
		ScenarioInitializer:  InitializeScenarios,
		Options:              &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}

	os.Exit(m.Run())
}

func InitializeScenarios(ctx *godog.ScenarioContext) {
	steps.CartInitializeScenario(ctx)
	steps.RemoveItemInitializeScenario(ctx)
	steps.RemoveAllItemFromCartInitializeScenario(ctx)
	steps.CartUpdateQuantityInitializeScenario(ctx)
	steps.CartDiscountScenario(ctx)
	steps.CartViewScenario(ctx)
	steps.CartCheckoutScenario(ctx)

}
