package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestDisputeAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "shopify_payments/disputes.json",
			response: `{"disputes":[{"id":7}]}`,
			want:     []Dispute{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Dispute.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "shopify_payments/disputes/7.json",
			response: `{"dispute":{"id":7}}`,
			want:     &Dispute{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Dispute.Get(ctx, 7, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestDisputeListWithPagination(t *testing.T) {
	testFeaturePagination(t, "shopify_payments/disputes.json", `{"disputes":[{"id":7}]}`, []Dispute{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Dispute.ListWithPagination(ctx, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/dispute
func TestDisputeDocumentedResponse(t *testing.T) {
	var value DisputesResource
	testFeatureResponse(t, "fixtures/dispute_documented.json", &value, true)
	if !(len(value.Disputes) > 0 && value.Disputes[0].InitiatedAt != nil && value.Disputes[0].OrderId == nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestDisputeQueryFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "shopify_payments/disputes.json",
		query:    "initiated_at=2026-01-15&last_id=7&status=won",
		response: `{"disputes":[]}`,
		want:     []Dispute{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Dispute.List(ctx, DisputeListOptions{InitiatedAt: &OnlyDate{Time: at}, LastId: 7, Status: "won"})
		},
	})
}
