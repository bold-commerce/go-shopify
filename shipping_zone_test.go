package goshopify

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

// shippingZoneTests tests if fields are properly parsed.
func shippingZoneTests(t *testing.T, zone ShippingZone) {
	cases := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", int64(1039932365), zone.ID},
		{"Name", "Some zone", zone.Name},
		{"WeightBasedShippingRates.0.ID", int64(882078075), zone.WeightBasedShippingRates[0].ID},
		{"WeightBasedShippingRates.0.ShippingZoneID", zone.ID, zone.WeightBasedShippingRates[0].ShippingZoneID},
		{"WeightBasedShippingRates.0.Name", "Canada Air Shipping", zone.WeightBasedShippingRates[0].Name},
		{"WeightBasedShippingRates.0.Price", decimal.NewFromFloat(25.00).String(), zone.WeightBasedShippingRates[0].Price.String()},
		{"WeightBasedShippingRates.0.WeightLow", decimal.NewFromFloat(0).String(), zone.WeightBasedShippingRates[0].WeightLow.String()},
		{"WeightBasedShippingRates.0.WeightHigh", decimal.NewFromFloat(11.0231).String(), zone.WeightBasedShippingRates[0].WeightHigh.String()},
		{"PriceBasedShippingRates.0.ID", int64(882078074), zone.PriceBasedShippingRates[0].ID},
		{"PriceBasedShippingRates.0.Name", "$5 Shipping", zone.PriceBasedShippingRates[0].Name},
		{"PriceBasedShippingRates.0.Price", decimal.NewFromFloat(5.05).String(), zone.PriceBasedShippingRates[0].Price.String()},
		{"PriceBasedShippingRates.0.MinOrderSubtotal", decimal.NewFromFloat(40.0).String(), zone.PriceBasedShippingRates[0].MinOrderSubtotal.String()},
		{"PriceBasedShippingRates.0.MaxOrderSubtotal", decimal.NewFromFloat(100.0).String(), zone.PriceBasedShippingRates[0].MaxOrderSubtotal.String()},
		{"CarrierShippingRateProviders.0.ID", int64(882078076), zone.CarrierShippingRateProviders[0].ID},
		{"CarrierShippingRateProviders.0.ShippingZoneID", zone.ID, zone.CarrierShippingRateProviders[0].ShippingZoneID},
		{"CarrierShippingRateProviders.0.CarrierServiceID", int64(770241334), zone.CarrierShippingRateProviders[0].CarrierServiceID},
		{"CarrierShippingRateProviders.0.FlatModifier", decimal.NewFromFloat(0).String(), zone.CarrierShippingRateProviders[0].FlatModifier.String()},
		{"CarrierShippingRateProviders.0.PercentModifier", decimal.NewFromFloat(0).String(), zone.CarrierShippingRateProviders[0].PercentModifier.String()},
	}

	for _, c := range cases {
		if c.expected != c.actual {
			t.Errorf("ShippingZone.%s returned %v, expected %v", c.field, c.actual,
				c.expected)
		}
	}
}

func TestShippingZoneListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shipping_zones.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	shippingZones, err := client.ShippingZone.List(context.Background())
	if shippingZones != nil {
		t.Errorf("ShippingZone.List returned shippingZones, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("ShippingZone.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestShippingZoneList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shipping_zones.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("shipping_zones.json")))

	shippingZones, err := client.ShippingZone.List(context.Background())
	if err != nil {
		t.Errorf("ShippingZone.List returned error: %v", err)
	}

	// Check that shippingZones were parsed
	if len(shippingZones) != 1 {
		t.Errorf("ShippingZone.List got %v shippingZones, expected: 1", len(shippingZones))
	}

	zone := shippingZones[0]
	shippingZoneTests(t, zone)
}
