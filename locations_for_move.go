package goshopify

import (
	"context"
	"fmt"
)

// LocationsForMoveService provides access to LocationsForMove endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/locationsformove
type LocationsForMoveService interface {
	List(ctx context.Context, fulfillmentOrderId uint64, options interface{}) ([]LocationForMove, error)
}

// LocationsForMoveServiceOp implements LocationsForMoveService.
type LocationsForMoveServiceOp struct {
	client *Client
}

// LocationForMove represents Shopify location for move.
type LocationForMove struct {
	Location *Location `json:"location,omitempty"`
	Movable  bool      `json:"movable"`
	Message  string    `json:"message,omitempty"`
}

// LocationsForMoveResource represents the locations_for_move response.
type LocationsForMoveResource struct {
	LocationsForMove []LocationForMove `json:"locations_for_move"`
}

// List retrieves locations for move records.
func (s *LocationsForMoveServiceOp) List(ctx context.Context, fulfillmentOrderId uint64, options interface{}) ([]LocationForMove, error) {
	resource := new(LocationsForMoveResource)
	err := s.client.Get(ctx, fmt.Sprintf("fulfillment_orders/%d/locations_for_move.json", fulfillmentOrderId), resource, options)
	return resource.LocationsForMove, err
}
