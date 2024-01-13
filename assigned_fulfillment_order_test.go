package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func AssignedFulfillmentOrderTests(t *testing.T, assignedFulfillmentOrder AssignedFulfillmentOrder) {
	// Check that ID is assigned to the returned fulfillment
	expectedInt := int64(255858046) // in assigned_fulfillment_orders.json fixture
	if assignedFulfillmentOrder.Id != expectedInt {
		t.Errorf("AssignedFulfillmentOrder.ID returned %+v, expected %+v", assignedFulfillmentOrder.Id, expectedInt)
	}
}

func TestAssignedFulfillmentOrderGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/assigned_fulfillment_orders.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment_orders": [{"id":1},{"id":2}]}`))

	assignedFulfillmentOrderService := &AssignedFulfillmentOrderServiceOp{client: client}

	assignedFulfillmentOrders, err := assignedFulfillmentOrderService.Get(context.Background(), nil)
	if err != nil {
		t.Errorf("AssignedFulfillmentOrder.List returned error: %v", err)
	}

	expected := []AssignedFulfillmentOrder{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(assignedFulfillmentOrders, expected) {
		t.Errorf("AssignedFulfillmentOrder.List returned %+v, expected %+v", assignedFulfillmentOrders, expected)
	}
}

// func TestFulfillmentOrderGet(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	fixture := loadFixture("fulfillment_order.json")
// 	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046.json", client.pathPrefix),
// 		httpmock.NewBytesResponder(200, fixture))

// 	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

// 	fulfillment, err := fulfillmentOrderService.Get(255858046, nil)
// 	if err != nil {
// 		t.Errorf("FulfillmentOrder.Get returned error: %v", err)
// 	}

// 	expected := FulfillmentOrderResource{}
// 	err = json.Unmarshal(fixture, &expected)
// 	if err != nil {
// 		t.Errorf("json.Unmarshall returned error : %v", err)
// 	}

// 	if !reflect.DeepEqual(fulfillment, expected.FulfillmentOrder) {
// 		t.Errorf("FulfillmentOrder.Get returned %+v, expected %+v", fulfillment, expected)
// 	}
// }
