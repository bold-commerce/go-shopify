package goshopify

import (
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

	rules, err := client.PriceRule.Get(1)
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

	rules, err := client.PriceRule.List()
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

	rules, err := client.PriceRule.Create(PriceRule{})
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

	rules, err := client.PriceRule.Update(PriceRule{ID: 1})
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

	err := client.PriceRule.Delete(1)
	if err != nil {
		t.Errorf("PriceRule.Delete returned error: %v", err)
	}
}
