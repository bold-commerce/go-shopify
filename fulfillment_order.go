package goshopify

import (
	"fmt"
)

// FulfillmentOrderService is an interface for interfacing with the fulfillment endpoints
// of the Shopify API.
// https://shopify.dev/api/admin-rest/2022-10/resources/fulfillmentorder
type FulfillmentOrderService interface {
	List(interface{}) ([]FulfillmentOrder, error)
}

// FulfillmentOrdersService is an interface for other Shopify resources
// to interface with the fulfillment endpoints of the Shopify API.
// https://help.shopify.com/api/reference/fulfillment
type FulfillmentOrdersService interface {
	ListFulfillmentOrders(int64, interface{}) ([]FulfillmentOrder, error)
}

// FulfillmentOrderServiceOp handles communication with the fulfillment order
// related methods of the Shopify API.
type FulfillmentOrderServiceOp struct {
	client     *Client
	resource   string
	resourceID int64
}

// FulfillmentOrder represents a Shopify fulfillment order.
type FulfillmentOrder struct {
	ID      int64  `json:"id,omitempty"`
	ShopID  int64  `json:"shop_id,omitempty"`
	OrderID int64  `json:"order_id,omitempty"`
	Status  string `json:"status,omitempty"`
	// Note: There are a lot of other possible fields however current we only care about ID
	// https://shopify.dev/api/admin-rest/2022-10/resources/fulfillmentorder#resource-object
}

// FulfillmentOrdersResource represents the result from the orders/{order_id}/fulfillment_orders.json endpoint
type FulfillmentOrdersResource struct {
	FulfillmentOrders []FulfillmentOrder `json:"fulfillment_orders"`
}

// List fulfillment orders
func (s *FulfillmentOrderServiceOp) List(options interface{}) ([]FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s.json", prefix)
	resource := new(FulfillmentOrdersResource)
	err := s.client.Get(path, resource, options)
	return resource.FulfillmentOrders, err
}
