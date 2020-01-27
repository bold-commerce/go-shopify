package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func FulfillmentTests(t *testing.T, fulfillment Fulfillment) {
	// Check that ID is assigned to the returned fulfillment
	expectedInt := int64(1022782888)
	if fulfillment.ID != expectedInt {
		t.Errorf("Fulfillment.ID returned %+v, expected %+v", fulfillment.ID, expectedInt)
	}
}

func TestFulfillmentList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillments": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	fulfillments, err := fulfillmentService.List(nil)
	if err != nil {
		t.Errorf("Fulfillment.List returned error: %v", err)
	}

	expected := []Fulfillment{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(fulfillments, expected) {
		t.Errorf("Fulfillment.List returned %+v, expected %+v", fulfillments, expected)
	}
}

func TestFulfillmentCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	cnt, err := fulfillmentService.Count(nil)
	if err != nil {
		t.Errorf("Fulfillment.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Fulfillment.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = fulfillmentService.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Fulfillment.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Fulfillment.Count returned %d, expected %d", cnt, expected)
	}
}

func TestFulfillmentGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment": {"id":1}}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	fulfillment, err := fulfillmentService.Get(1, nil)
	if err != nil {
		t.Errorf("Fulfillment.Get returned error: %v", err)
	}

	expected := &Fulfillment{ID: 1}
	if !reflect.DeepEqual(fulfillment, expected) {
		t.Errorf("Fulfillment.Get returned %+v, expected %+v", fulfillment, expected)
	}
}

func TestFulfillmentCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	fulfillment := Fulfillment{
		LocationID:     905684977,
		TrackingNumber: "123456789",
		TrackingUrls: []string{
			"https://shipping.xyz/track.php?num=123456789",
			"https://anothershipper.corp/track.php?code=abc",
		},
		NotifyCustomer: true,
	}

	returnedFulfillment, err := fulfillmentService.Create(fulfillment)
	if err != nil {
		t.Errorf("Fulfillment.Create returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1022782888.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	fulfillment := Fulfillment{
		ID:             1022782888,
		TrackingNumber: "987654321",
	}

	returnedFulfillment, err := fulfillmentService.Update(fulfillment)
	if err != nil {
		t.Errorf("Fulfillment.Update returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentComplete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/complete.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	returnedFulfillment, err := fulfillmentService.Complete(1)
	if err != nil {
		t.Errorf("Fulfillment.Complete returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentTransition(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/open.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	returnedFulfillment, err := fulfillmentService.Transition(1)
	if err != nil {
		t.Errorf("Fulfillment.Transition returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentCancel(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/cancel.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	returnedFulfillment, err := fulfillmentService.Cancel(1)
	if err != nil {
		t.Errorf("Fulfillment.Cancel returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}
