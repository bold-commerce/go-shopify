package goshopify

import (
	"context"
	"testing"
)

func TestLocationAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "ListInventoryLevels",
			method:   "GET",
			path:     "locations/11/inventory_levels.json",
			response: `{"inventory_levels":[{"inventory_item_id":7,"location_id":11}]}`,
			want:     []InventoryLevel{{InventoryItemId: 7, LocationId: 11}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Location.ListInventoryLevels(ctx, 11, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestLocationListInventoryLevelsWithPagination(t *testing.T) {
	testFeaturePagination(t, "locations/11/inventory_levels.json", `{"inventory_levels":[{"inventory_item_id":7}]}`, []InventoryLevel{{InventoryItemId: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Location.ListInventoryLevelsWithPagination(ctx, 11, options)
	})
}
