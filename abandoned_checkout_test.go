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
