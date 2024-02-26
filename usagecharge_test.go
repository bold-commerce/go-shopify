package goshopify

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func usageChargeTests(t *testing.T, usageCharge UsageCharge) {
	price := decimal.NewFromFloat(1.0)
	createdAt, _ := time.Parse(time.RFC3339, "2018-07-05T13:05:43-04:00")
	billingOn, _ := time.Parse("2006-01-02", "2018-08-04")
	balanceUsed := decimal.NewFromFloat(11.0)
	balancedRemaining := decimal.NewFromFloat(89.0)
	riskLevel := decimal.NewFromFloat(0.08)

	expected := UsageCharge{
		Id:               1034618208,
		Description:      "Super Mega Plan 1000 emails",
		Price:            &price,
		CreatedAt:        &createdAt,
		BillingOn:        &billingOn,
		BalanceRemaining: &balancedRemaining,
		BalanceUsed:      &balanceUsed,
		RiskLevel:        &riskLevel,
	}

	if usageCharge.Id != expected.Id {
		t.Errorf("UsageCharge.Id returned %v, expected %v", usageCharge.Id, expected.Id)
	}
	if usageCharge.Description != expected.Description {
		t.Errorf("UsageCharge.Description returned %v, expected %v", usageCharge.Description, expected.Description)
	}
	if !usageCharge.Price.Equal(*expected.Price) {
		t.Errorf("UsageCharge.Price returned %v, expected %v", usageCharge.Price, expected.Price)
	}
	if !usageCharge.CreatedAt.Equal(*expected.CreatedAt) {
		t.Errorf("UsageCharge.CreatedAt returned %v, expected %v", usageCharge.CreatedAt, expected.CreatedAt)
	}
	if !usageCharge.BillingOn.Equal(*expected.BillingOn) {
		t.Errorf("UsageCharge.BillingOn returned %v, expected %v", usageCharge.BillingOn, expected.BillingOn)
	}
	if !usageCharge.BalanceRemaining.Equal(*expected.BalanceRemaining) {
		t.Errorf("UsageCharge.BalanceRemaining returned %v, expected %v", usageCharge.BalanceRemaining, expected.BalanceRemaining)
	}
	if !usageCharge.BalanceUsed.Equal(*expected.BalanceUsed) {
		t.Errorf("UsageCharge.BalanceUsed returned %v, expected %v", usageCharge.BalanceUsed, expected.BalanceUsed)
	}
	if !usageCharge.RiskLevel.Equal(*expected.RiskLevel) {
		t.Errorf("UsageCharge.RiskLevel returned %v, expected %v", usageCharge.RiskLevel, expected.RiskLevel)
	}
}

func TestUsageChargeServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharge.json"),
		),
	)

	p := decimal.NewFromFloat(1.0)
	charge := UsageCharge{
		Description: "Super Mega Plan 1000 emails",
		Price:       &p,
	}

	returnedCharge, err := client.UsageCharge.Create(context.Background(), 455696195, charge)
	if err != nil {
		t.Errorf("UsageCharge.Create returned an error: %v", err)
	}
	usageChargeTests(t, *returnedCharge)
}

func TestUsageChargeServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharge.json"),
		),
	)

	charge, err := client.UsageCharge.Get(context.Background(), 455696195, 1034618210, nil)
	if err != nil {
		t.Errorf("UsageCharge.Get returned an error: %v", err)
	}

	usageChargeTests(t, *charge)
}

func TestUsageChargeServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharges.json"),
		),
	)

	charges, err := client.UsageCharge.List(context.Background(), 455696195, nil)
	if err != nil {
		t.Errorf("UsageCharge.List returned an error: %v", err)
	}

	// Check that usage charges were parsed
	if len(charges) != 1 {
		t.Errorf("UsageCharage.List got %v usage charges, expected: 1", len(charges))
	}

	usageChargeTests(t, charges[0])
}

func TestUsageChargeServiceOp_GetBadFields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200, `{"usage_charge":{"id":"wrong_id_type"}}`,
		),
	)

	if _, err := client.UsageCharge.Get(context.Background(), 455696195, 1034618210, nil); err == nil {
		t.Errorf("UsageCharge.Get should have returned an error")
	}

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200, `{"usage_charge":{"billing_on":"2018-14-01"}}`,
		),
	)
	if _, err := client.UsageCharge.Get(context.Background(), 455696195, 1034618210, nil); err == nil {
		t.Errorf("UsageCharge.Get should have returned an error")
	}
}
