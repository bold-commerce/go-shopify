package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAbandonedCheckoutList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/checkouts.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"checkouts": [{"id":1},{"id":2}]}`,
		),
	)

	abandonedCheckouts, err := client.AbandonedCheckout.List(context.Background(), nil)
	if err != nil {
		t.Errorf("AbandonedCheckout.List returned error: %v", err)
	}

	expected := []AbandonedCheckout{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(abandonedCheckouts, expected) {
		t.Errorf("AbandonedCheckout.List returned %+v, expected %+v", abandonedCheckouts, expected)
	}
}

func TestAbandonedCheckoutDocumentedResponse(t *testing.T) {
	setup()
	defer teardown()
	// Shopify's abandoned checkout response has string shipping IDs and line items.
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/checkouts.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("abandoned_checkout.json")))
	checkouts, err := client.AbandonedCheckout.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(checkouts) != 1 {
		t.Fatalf("got %d checkouts", len(checkouts))
	}
	checkout := checkouts[0]
	if len(checkout.ShippingLines) != 1 || checkout.ShippingLines[0].Id != "5da41c1738454765" {
		t.Fatalf("shipping lines = %#v", checkout.ShippingLines)
	}
	if len(checkout.LineItems) != 2 || checkout.LineItems[0].VariantId != 49148385 || checkout.LineItems[0].Quantity != 1 {
		t.Fatalf("line items = %#v", checkout.LineItems)
	}
}
