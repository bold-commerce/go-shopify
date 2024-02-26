package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func FulfillmentTests(t *testing.T, fulfillment Fulfillment) {
	// Check that Id is assigned to the returned fulfillment
	expectedInt := uint64(1022782888)
	if fulfillment.Id != expectedInt {
		t.Errorf("Fulfillment.Id returned %+v, expected %+v", fulfillment.Id, expectedInt)
	}
}

func TestFulfillmentList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillments": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillments, err := fulfillmentService.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Fulfillment.List returned error: %v", err)
	}

	expected := []Fulfillment{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(fulfillments, expected) {
		t.Errorf("Fulfillment.List returned %+v, expected %+v", fulfillments, expected)
	}
}

func TestFulfillmentCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	cnt, err := fulfillmentService.Count(context.Background(), nil)
	if err != nil {
		t.Errorf("Fulfillment.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Fulfillment.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = fulfillmentService.Count(context.Background(), CountOptions{CreatedAtMin: date})
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment": {"id":1}}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillment, err := fulfillmentService.Get(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Fulfillment.Get returned error: %v", err)
	}

	expected := &Fulfillment{Id: 1}
	if !reflect.DeepEqual(fulfillment, expected) {
		t.Errorf("Fulfillment.Get returned %+v, expected %+v", fulfillment, expected)
	}
}

func TestFulfillmentCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillment := Fulfillment{
		LocationId:     905684977,
		TrackingNumber: "123456789",
		TrackingUrls: []string{
			"https://shipping.xyz/track.php?num=123456789",
			"https://anothershipper.corp/track.php?code=abc",
		},
		NotifyCustomer: true,
	}

	returnedFulfillment, err := fulfillmentService.Create(context.Background(), fulfillment)
	if err != nil {
		t.Errorf("Fulfillment.Create returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1022782888.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillment := Fulfillment{
		Id:             1022782888,
		TrackingNumber: "987654321",
	}

	returnedFulfillment, err := fulfillmentService.Update(context.Background(), fulfillment)
	if err != nil {
		t.Errorf("Fulfillment.Update returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentComplete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/complete.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	returnedFulfillment, err := fulfillmentService.Complete(context.Background(), 1)
	if err != nil {
		t.Errorf("Fulfillment.Complete returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentTransition(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/open.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	returnedFulfillment, err := fulfillmentService.Transition(context.Background(), 1)
	if err != nil {
		t.Errorf("Fulfillment.Transition returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestFulfillmentCancel(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1/cancel.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	returnedFulfillment, err := fulfillmentService.Cancel(context.Background(), 1)
	if err != nil {
		t.Errorf("Fulfillment.Cancel returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}
