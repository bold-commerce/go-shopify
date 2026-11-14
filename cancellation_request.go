package goshopify

import (
	"context"
	"fmt"
)

// CancellationRequestService provides access to CancellationRequest endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/cancellationrequest
type CancellationRequestService interface {
	Send(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error)
	Accept(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error)
	Reject(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error)
}

// CancellationRequestServiceOp implements CancellationRequestService.
type CancellationRequestServiceOp struct {
	client *Client
}

// Send sends a cancellation request for a fulfillment order.
func (s *CancellationRequestServiceOp) Send(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error) {
	var body struct {
		CancellationRequest struct {
			Message string `json:"message,omitempty"`
		} `json:"cancellation_request"`
	}
	body.CancellationRequest.Message = message
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(ctx, fmt.Sprintf("fulfillment_orders/%d/cancellation_request.json", fulfillmentOrderId), body, resource)
	return resource.FulfillmentOrder, err
}

// Accept accepts a cancellation request for a fulfillment order.
func (s *CancellationRequestServiceOp) Accept(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error) {
	var body struct {
		CancellationRequest struct {
			Message string `json:"message,omitempty"`
		} `json:"cancellation_request"`
	}
	body.CancellationRequest.Message = message
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(ctx, fmt.Sprintf("fulfillment_orders/%d/cancellation_request/accept.json", fulfillmentOrderId), body, resource)
	return resource.FulfillmentOrder, err
}

// Reject rejects a cancellation request for a fulfillment order.
func (s *CancellationRequestServiceOp) Reject(ctx context.Context, fulfillmentOrderId uint64, message string) (*FulfillmentOrder, error) {
	var body struct {
		CancellationRequest struct {
			Message string `json:"message,omitempty"`
		} `json:"cancellation_request"`
	}
	body.CancellationRequest.Message = message
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(ctx, fmt.Sprintf("fulfillment_orders/%d/cancellation_request/reject.json", fulfillmentOrderId), body, resource)
	return resource.FulfillmentOrder, err
}
