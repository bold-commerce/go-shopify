package goshopify

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestFulfillmentOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillment_orders.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment_orders": [{"id":1},{"id":2}]}`))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client, resource: ordersResourceName, resourceID: 123}

	fulfillmentOrders, err := fulfillmentOrderService.List(nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.List returned error: %v", err)
	}

	expected := []FulfillmentOrder{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(fulfillmentOrders, expected) {
		t.Errorf("FulfillmentOrder.List returned %+v, expected %+v", fulfillmentOrders, expected)
	}
}
