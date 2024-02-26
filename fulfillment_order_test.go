package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func FulfillmentOrderTests(t *testing.T, fulfillmentOrder FulfillmentOrder) {
	// Check that Id is assigned to the returned fulfillment
	expectedInt := uint64(255858046) // in fulfillment_order.json fixture
	if fulfillmentOrder.Id != expectedInt {
		t.Errorf("FulfillmentOrder.Id returned %+v, expected %+v", fulfillmentOrder.Id, expectedInt)
	}
}

func TestFulfillmentOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillment_orders.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment_orders": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentOrderServiceOp{client: client}

	fulfillmentOrders, err := fulfillmentService.List(context.Background(), 123, nil)
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

	fixture := loadFixture("fulfillment_order.json")
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, fixture))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	fulfillment, err := fulfillmentOrderService.Get(context.Background(), 255858046, nil)
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
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Cancel(context.Background(), 1)
	if err != nil {
		t.Errorf("FulfillmentOrder.Cancel returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderClose(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/close.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Close(context.Background(), 1, "test")
	if err != nil {
		t.Errorf("FulfillmentOrder.Close returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderHold(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/hold.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Hold(context.Background(), 1, false, HoldReasonOutOfStock, "test")
	if err != nil {
		t.Errorf("FulfillmentOrder.Hold returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderMove(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1046000818/move.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order_move.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	req := FulfillmentOrderMoveRequest{
		NewLocationId: 655441491,
		LineItems: []FulfillmentOrderLineItemQuantity{
			{Id: 1058737594, Quantity: 1},
		},
	}

	result, err := fulfillmentOrderService.Move(context.Background(), 1046000818, req)
	if err != nil {
		t.Errorf("FulfillmentOrder.Move returned error: %v", err)
	}
	if result.MovedFulfillmentOrder.AssignedLocationId != 655441491 {
		t.Errorf("FulfillmentOrder.Move result AssignedLocation is is %d, expected %d",
			result.MovedFulfillmentOrder.AssignedLocationId, 655441491)
	}
}

func TestFulfillmentOrderOpen(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/open.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.Open(context.Background(), fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentOrder.Open returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderReleaseHold(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/release_hold.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.ReleaseHold(context.Background(), fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentOrder.ReleaseHold returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderReschedule(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/reschedule.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.Reschedule(context.Background(), fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentOrder.Reschedule returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderSetDeadline(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/set_fulfillment_orders_deadline.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)
	newDeadline := time.Now().Add(time.Hour * 24 * 7)
	err := fulfillmentOrderService.SetDeadline(context.Background(), []uint64{fulfillmentId}, newDeadline)
	if err != nil {
		t.Errorf("FulfillmentOrder.SetDeadline returned error: %v", err)
	}
}
