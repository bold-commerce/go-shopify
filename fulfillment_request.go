package goshopify

import (
	"context"
	"fmt"
)

const (
	fulfillmentRequestBasePath = "fulfillment_orders"
)

// FulfillmentRequestService is an interface for interfacing with the fulfillment request endpoints of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/fulfillmentrequest
type FulfillmentRequestService interface {
	Send(context.Context, uint64, FulfillmentRequest) (*FulfillmentOrder, error)
	SendWithResponse(context.Context, uint64, FulfillmentRequest) (*FulfillmentRequestResource, error)
	Accept(context.Context, uint64, FulfillmentRequest) (*FulfillmentOrder, error)
	Reject(context.Context, uint64, FulfillmentRequest) (*FulfillmentOrder, error)
}

type FulfillmentRequest struct {
	NotifyCustomer            *bool                        `json:"notify_customer,omitempty"`
	Message                   string                       `json:"message,omitempty"`
	FulfillmentOrderLineItems []FulfillmentOrderLineItem   `json:"fulfillment_order_line_items,omitempty"`
	Reason                    string                       `json:"reason,omitempty"`
	LineItems                 []FulfillmentRequestLineItem `json:"line_items,omitempty"`
}

type FulfillmentRequestOrderLineItem struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
}

type FulfillmentRequestLineItem struct {
	FulfillmentOrderLineItemId uint64 `json:"fulfillment_order_line_item_id,omitempty"`
	Message                    string `json:"message,omitempty"`
}

type FulfillmentRequestResource struct {
	FulfillmentOrder            *FulfillmentOrder  `json:"fulfillment_order,omitempty"`
	FulfillmentRequest          FulfillmentRequest `json:"fulfillment_request,omitempty"`
	OriginalFulfillmentOrder    *FulfillmentOrder  `json:"original_fulfillment_order,omitempty"`
	SubmittedFulfillmentOrder   *FulfillmentOrder  `json:"submitted_fulfillment_order,omitempty"`
	UnsubmittedFulfillmentOrder *FulfillmentOrder  `json:"unsubmitted_fulfillment_order,omitempty"`
}

// FulfillmentRequestServiceOp handles communication with the fulfillment request related methods of the Shopify API.
type FulfillmentRequestServiceOp struct {
	client *Client
}

// Send sends a fulfillment request and returns the original fulfillment order.
// Use SendWithResponse to also retrieve the submitted and unsubmitted orders.
func (s *FulfillmentRequestServiceOp) Send(ctx context.Context, fulfillmentOrderId uint64, request FulfillmentRequest) (*FulfillmentOrder, error) {
	resource, err := s.SendWithResponse(ctx, fulfillmentOrderId, request)
	return resource.OriginalFulfillmentOrder, err
}

// SendWithResponse sends a fulfillment request and returns all resulting fulfillment orders.
func (s *FulfillmentRequestServiceOp) SendWithResponse(ctx context.Context, fulfillmentOrderId uint64, request FulfillmentRequest) (*FulfillmentRequestResource, error) {
	path := fmt.Sprintf("%s/%d/fulfillment_request.json", fulfillmentRequestBasePath, fulfillmentOrderId)
	wrappedData := FulfillmentRequestResource{FulfillmentRequest: request}
	resource := new(FulfillmentRequestResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource, err
}

// Accept accepts a fulfillment request sent to a fulfillment service for a fulfillment order.
func (s *FulfillmentRequestServiceOp) Accept(ctx context.Context, fulfillmentOrderId uint64, request FulfillmentRequest) (*FulfillmentOrder, error) {
	path := fmt.Sprintf("%s/%d/fulfillment_request/accept.json", fulfillmentRequestBasePath, fulfillmentOrderId)
	wrappedData := map[string]interface{}{"fulfillment_request": request}
	resource := new(FulfillmentRequestResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.FulfillmentOrder, err
}

// Reject rejects a fulfillment request sent to a fulfillment service for a fulfillment order.
func (s *FulfillmentRequestServiceOp) Reject(ctx context.Context, fulfillmentOrderId uint64, request FulfillmentRequest) (*FulfillmentOrder, error) {
	path := fmt.Sprintf("%s/%d/fulfillment_request/reject.json", fulfillmentRequestBasePath, fulfillmentOrderId)
	wrappedData := map[string]interface{}{"fulfillment_request": request}
	resource := new(FulfillmentRequestResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.FulfillmentOrder, err
}
