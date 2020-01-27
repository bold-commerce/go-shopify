package goshopify

import (
	"fmt"
	"reflect"
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
		ID:               1034618208,
		Description:      "Super Mega Plan 1000 emails",
		Price:            &price,
		CreatedAt:        &createdAt,
		BillingOn:        &billingOn,
		BalanceRemaining: &balancedRemaining,
		BalanceUsed:      &balanceUsed,
		RiskLevel:        &riskLevel,
	}

	if !reflect.DeepEqual(usageCharge, expected) {
		t.Errorf("UsageCharge.Create returned %+v, expected %+v", usageCharge, expected)
	}
}
func TestUsageChargeServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharge.json"),
		),
	)

	p := decimal.NewFromFloat(1.0)
	charge := UsageCharge{
		Description: "Super Mega Plan 1000 emails",
		Price:       &p,
	}

	returnedCharge, err := client.UsageCharge.Create(455696195, charge)
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
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharge.json"),
		),
	)

	charge, err := client.UsageCharge.Get(455696195, 1034618210, nil)
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
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(
			200, loadFixture("usagecharges.json"),
		),
	)

	charges, err := client.UsageCharge.List(455696195, nil)
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
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200, `{"usage_charge":{"id":"wrong_id_type"}}`,
		),
	)

	if _, err := client.UsageCharge.Get(455696195, 1034618210, nil); err == nil {
		t.Errorf("UsageCharge.Get should have returned an error")
	}

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/recurring_application_charges/455696195/usage_charges/1034618210.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200, `{"usage_charge":{"billing_on":"2018-14-01"}}`,
		),
	)
	if _, err := client.UsageCharge.Get(455696195, 1034618210, nil); err == nil {
		t.Errorf("UsageCharge.Get should have returned an error")
	}

}
