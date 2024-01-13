package goshopify

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestPriceRuleGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("PriceRule.Get returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Get returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/list.json"),
		),
	)

	rules, err := client.PriceRule.List(context.Background())
	if err != nil {
		t.Errorf("PriceRule.List returned error: %v", err)
	}

	expected := []PriceRule{{ID: 1}}
	if expected[0].ID != rules[0].ID {
		t.Errorf("PriceRule.List returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Create(context.Background(), PriceRule{})
	if err != nil {
		t.Errorf("PriceRule.Create returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Create returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Update(context.Background(), PriceRule{ID: 1})
	if err != nil {
		t.Errorf("PriceRule.Update returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Update returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	err := client.PriceRule.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("PriceRule.Delete returned error: %v", err)
	}
}

func TestPriceRuleSetters(t *testing.T) {
	pr := PriceRule{}
	prereqSubtotalRange := "1.5"
	prereqQuantityRange := 2
	prereqShippingPrice := "5.5"
	prereqRatioQuantity := 1
	prereqRatioEntitledQuantity := 1
	badMoneyString := "dog"

	// Test bad money strings
	err := pr.SetPrerequisiteSubtotalRange(&badMoneyString)
	if err == nil {
		t.Errorf("Expected error from setting bad string as prerequisite subtotal range: %s", badMoneyString)
	}

	err = pr.SetPrerequisiteShippingPriceRange(&badMoneyString)
	if err == nil {
		t.Errorf("Expected error from setting bad string as prerequisite shipping price: %s", badMoneyString)
	}

	// Test populating values
	err = pr.SetPrerequisiteSubtotalRange(&prereqSubtotalRange)
	if err != nil {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteQuantityRange(&prereqQuantityRange)
	err = pr.SetPrerequisiteShippingPriceRange(&prereqShippingPrice)
	if err != nil {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(&prereqRatioQuantity, &prereqRatioEntitledQuantity)

	if pr.PrerequisiteSubtotalRange.GreaterThanOrEqualTo != prereqSubtotalRange {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	if pr.PrerequisiteQuantityRange.GreaterThanOrEqualTo != prereqQuantityRange {
		t.Errorf("Failed to set prerequisite quantity range: %d", prereqQuantityRange)
	}

	if pr.PrerequisiteShippingPriceRange.LessThanOrEqualTo != prereqShippingPrice {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqShippingPrice)
	}

	if pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != prereqRatioQuantity {
		t.Errorf("Failed to set prerequisite ratio quantity: %d", prereqRatioQuantity)
	}

	if pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != prereqRatioEntitledQuantity {
		t.Errorf("Failed to set prerequisite ratio entitled quantity: %d", prereqRatioEntitledQuantity)
	}

	// Test clearing values by setting nil
	err = pr.SetPrerequisiteSubtotalRange(nil)
	if err != nil {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteQuantityRange(nil)
	err = pr.SetPrerequisiteShippingPriceRange(nil)
	if err != nil {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqSubtotalRange)
	}

	if pr.PrerequisiteSubtotalRange != nil {
		t.Errorf("Failed to clear prerequisite subtotal range")
	}

	if pr.PrerequisiteQuantityRange != nil {
		t.Errorf("Failed to clear prerequisite quantity range")
	}

	if pr.PrerequisiteShippingPriceRange != nil {
		t.Errorf("Failed to clear prerequisite shipping price")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(nil, &prereqRatioEntitledQuantity)
	if pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != 0 || pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != prereqRatioEntitledQuantity {
		t.Errorf("Failed to clear prerequisite-to-entitlement-quantity-ratio prerequisite quantity")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(&prereqRatioQuantity, nil)
	if pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != 0 || pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != prereqRatioQuantity {
		t.Errorf("Failed to clear prerequisite-to-entitlement-quantity-ratio entitled quantity")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(nil, nil)
	if pr.PrerequisiteToEntitlementQuantityRatio != nil {
		t.Errorf("Failed to clear wholly prerequisite to entitlement quantity ratio")
	}
}
