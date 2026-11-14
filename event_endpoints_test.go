package goshopify

import (
	"context"
	"testing"
)

func TestEventAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "events.json",
			response: `{"events":[{"id":7}]}`,
			want:     []Event{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Event.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "events/7.json",
			response: `{"event":{"id":7}}`,
			want:     &Event{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Event.Get(ctx, 7, nil)
			},
		},
		{
			name:     "Count",
			method:   "GET",
			path:     "events/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.Event.Count(ctx, nil)
			},
		},
		{
			name:     "ListForOrder",
			method:   "GET",
			path:     "orders/11/events.json",
			response: `{"events":[{"id":7}]}`,
			want:     []Event{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Event.ListForOrder(ctx, 11, nil)
			},
		},
		{
			name:     "ListForProduct",
			method:   "GET",
			path:     "products/11/events.json",
			response: `{"events":[{"id":7}]}`,
			want:     []Event{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Event.ListForProduct(ctx, 11, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestEventListWithPagination(t *testing.T) {
	testFeaturePagination(t, "events.json", `{"events":[{"id":7}]}`, []Event{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Event.ListWithPagination(ctx, options)
	})
}

func TestEventListForOrderWithPagination(t *testing.T) {
	testFeaturePagination(t, "orders/11/events.json", `{"events":[{"id":7}]}`, []Event{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Event.ListForOrderWithPagination(ctx, 11, options)
	})
}

func TestEventListForProductWithPagination(t *testing.T) {
	testFeaturePagination(t, "products/11/events.json", `{"events":[{"id":7}]}`, []Event{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Event.ListForProductWithPagination(ctx, 11, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/event
func TestEventDocumentedResponse(t *testing.T) {
	var value EventsResource
	testFeatureResponse(t, "fixtures/event_documented.json", &value, true)
	if !(len(value.Events) > 0 && value.Events[0].SubjectId != 0 && len(value.Events[0].Arguments) > 0) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestEventQueryFields(t *testing.T) {
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "events.json",
		query:    "filter=Product%2COrder&verb=destroy",
		response: `{"events":[]}`,
		want:     []Event{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Event.List(ctx, EventListOptions{Filter: "Product,Order", Verb: "destroy"})
		},
	})
}
