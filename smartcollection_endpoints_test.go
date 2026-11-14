package goshopify

import (
	"context"
	"testing"
)

func TestSmartCollectionAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:   "Order",
			method: "PUT",
			path:   "smart_collections/11/order.json",
			query:  "products[]=7&products[]=8",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.SmartCollection.Order(ctx, 11, SmartCollectionOrderOptions{Products: []uint64{7, 8}})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestSmartCollectionSortOrder(t *testing.T) {
	testFeatureEndpoint(t, featureEndpoint{
		name:   "Order",
		method: "PUT",
		path:   "smart_collections/11/order.json",
		query:  "sort_order=alpha-desc",
		want:   nil,
		call: func(ctx context.Context) (interface{}, error) {
			return nil, client.SmartCollection.Order(ctx, 11, SmartCollectionOrderOptions{SortOrder: "alpha-desc"})
		},
	})
}
