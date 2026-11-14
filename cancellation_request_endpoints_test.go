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

func TestCancellationRequestWithoutMessage(t *testing.T) {
	for _, tc := range []struct {
		name, path, status, requestStatus string
		call                              func(context.Context) (interface{}, error)
	}{
		{"Send", "cancellation_request", "in_progress", "cancellation_requested", func(ctx context.Context) (interface{}, error) {
			return client.CancellationRequest.Send(ctx, 11, "")
		}},
		{"Accept", "cancellation_request/accept", "cancelled", "cancellation_accepted", func(ctx context.Context) (interface{}, error) {
			return client.CancellationRequest.Accept(ctx, 11, "")
		}},
		{"Reject", "cancellation_request/reject", "in_progress", "cancellation_rejected", func(ctx context.Context) (interface{}, error) {
			return client.CancellationRequest.Reject(ctx, 11, "")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testFeatureEndpoint(t, featureEndpoint{
				method:   "POST",
				path:     "fulfillment_orders/11/" + tc.path + ".json",
				request:  `{"cancellation_request":{}}`,
				response: `{"fulfillment_order":{"id":11,"status":"` + tc.status + `","request_status":"` + tc.requestStatus + `"}}`,
				want:     &FulfillmentOrder{Id: 11, Status: tc.status, RequestStatus: tc.requestStatus},
				call:     tc.call,
			})
		})
	}
}
