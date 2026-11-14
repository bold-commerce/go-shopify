package goshopify

import (
	"context"
	"testing"
)

func TestCancellationRequestAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Send",
			method:   "POST",
			path:     "fulfillment_orders/11/cancellation_request.json",
			request:  `{"cancellation_request":{"message":"Unable to fulfill."}}`,
			response: `{"fulfillment_order":{"id":11}}`,
			want:     &FulfillmentOrder{Id: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CancellationRequest.Send(ctx, 11, "Unable to fulfill.")
			},
		},
		{
			name:     "Accept",
			method:   "POST",
			path:     "fulfillment_orders/11/cancellation_request/accept.json",
			request:  `{"cancellation_request":{"message":"Unable to fulfill."}}`,
			response: `{"fulfillment_order":{"id":11}}`,
			want:     &FulfillmentOrder{Id: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CancellationRequest.Accept(ctx, 11, "Unable to fulfill.")
			},
		},
		{
			name:     "Reject",
			method:   "POST",
			path:     "fulfillment_orders/11/cancellation_request/reject.json",
			request:  `{"cancellation_request":{"message":"Unable to fulfill."}}`,
			response: `{"fulfillment_order":{"id":11}}`,
			want:     &FulfillmentOrder{Id: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CancellationRequest.Reject(ctx, 11, "Unable to fulfill.")
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/cancellationrequest
func TestCancellationRequestDocumentedResponse(t *testing.T) {
	var value FulfillmentOrderResource
	testFeatureResponse(t, "fixtures/cancellation_request_documented.json", &value, false)
	if !(value.FulfillmentOrder.Id != 0 && value.FulfillmentOrder.RequestStatus != "") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
