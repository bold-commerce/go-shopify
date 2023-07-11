package goshopify

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

func FulfillmentOrderTests(t *testing.T, fulfillmentOrder FulfillmentOrder) {
	// Check that ID is assigned to the returned fulfillment
	expectedInt := int64(255858046) // in fulfillmentorder.json fixture
	if fulfillmentOrder.Id != expectedInt {
		t.Errorf("FulfillmentOrder.ID returned %+v, expected %+v", fulfillmentOrder.Id, expectedInt)
	}
}

func TestFulfillmentOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillment_orders.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment_orders": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentOrderServiceOp{client: client}

	fulfillmentOrders, err := fulfillmentService.List(123, nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.List returned error: %v", err)
	}

	expected := []FulfillmentOrder{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(fulfillmentOrders, expected) {
		t.Errorf("FulfillmentOrder.List returned %+v, expected %+v", fulfillmentOrders, expected)
	}
}

func TestFulfillmentOrderGet(t *testing.T) {
	setup()
	defer teardown()

	fixture := loadFixture("fulfillmentorder.json")
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, fixture))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	fulfillment, err := fulfillmentOrderService.Get(255858046, nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.Get returned error: %v", err)
	}

	expected := FulfillmentOrderResource{}
	err = json.Unmarshal(fixture, &expected)
	if err != nil {
		t.Errorf("json.Unmarshall returned error : %v", err)
	}

	if !reflect.DeepEqual(fulfillment, expected.FulfillmentOrder) {
		t.Errorf("FulfillmentOrder.Get returned %+v, expected %+v", fulfillment, expected)
	}
}

func TestFulfillmentOrderCancel(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/cancel.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillmentorder.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Cancel(1)
	if err != nil {
		t.Errorf("FulfillmentOrder.Cancel returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderClose(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/close.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillmentorder.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Close(1, "test")
	if err != nil {
		t.Errorf("FulfillmentOrder.Close returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderHold(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/hold.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillmentorder.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Hold(1, false, HoldReasonOutOfStock, "test")
	if err != nil {
		t.Errorf("FulfillmentOrder.Hold returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

// TODO
func TestFulfillmentOrderMove(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1046000818/move.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillmentorder_move.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	req := FulfillmentOrderMoveRequest{
		NewLocationId: 655441491,
		LineItems: []FulfillmentOrderLineItemQuantity{
			{Id: 1058737594, Quantity: 1},
		},
	}

	result, err := fulfillmentOrderService.Move(1046000818, req, nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.Move returned error: %v", err)
	}
	if result.MovedFulfillmentOrder.AssignedLocationId != 655441491 {
		t.Errorf("FulfillmentOrder.Move result AssignedLocation is is %d, expected %d",
			result.MovedFulfillmentOrder.AssignedLocationId, 655441491)
	}
}

// TODO
func TestFulfillmentOrderOpen(t *testing.T) {
	t.Skip("TestFulfillmentOrderOpen is not implemented")
}

// TODO
func TestFulfillmentOrderReleaseHold(t *testing.T) {
	t.Skip("TestFulfillmentOrderReleaseHold is not implemented")
}

// TODO
func TestFulfillmentOrderReschedule(t *testing.T) {
	t.Skip("TestFulfillmentOrderReschedule is not implemented")
}

// TODO
func TestFulfillmentOrderSetDeadline(t *testing.T) {
	t.Skip("TestFulfillmentOrderSetDeadline is not implemented")
}
