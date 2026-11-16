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
			response: `{"inventory_levels":[{"inventory_item_id":7,"location_id":11,"available":-2},{"inventory_item_id":8,"location_id":11,"available":0},{"inventory_item_id":9,"location_id":11,"available":25}]}`,
			want:     []InventoryLevel{{InventoryItemId: 7, LocationId: 11, Available: -2}, {InventoryItemId: 8, LocationId: 11}, {InventoryItemId: 9, LocationId: 11, Available: 25}},
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
	testFeaturePagination(t, "locations/11/inventory_levels.json", `{"inventory_levels":[{"inventory_item_id":7}]}`, []InventoryLevel{{InventoryItemId: 7}}, `{"inventory_levels":[{"inventory_item_id":8}]}`, []InventoryLevel{{InventoryItemId: 8}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Location.ListInventoryLevelsWithPagination(ctx, 11, options)
	})
}
