package goshopify

import (
	"context"
	"testing"
)

func TestLocationsForMoveAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "fulfillment_orders/11/locations_for_move.json",
			response: `{"locations_for_move":[{"location":{"id":7,"name":"Warehouse"},"movable":false,"message":"Not stocked here."},{"location":{"id":8},"movable":true,"message":null}]}`,
			want:     []LocationForMove{{Location: &Location{Id: 7, Name: "Warehouse"}, Movable: false, Message: "Not stocked here."}, {Location: &Location{Id: 8}, Movable: true}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.LocationsForMove.List(ctx, 11, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/locationsformove
func TestLocationsForMoveDocumentedResponse(t *testing.T) {
	var value LocationsForMoveResource
	testFeatureResponse(t, "fixtures/locations_for_move_documented.json", &value, true)
	if !(len(value.LocationsForMove) > 0 && value.LocationsForMove[0].Location.Id != 0) {
		t.Fatalf("incomplete response: %#v", value)
	}
}
