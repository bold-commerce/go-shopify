package goshopify

import (
	"context"
	"fmt"
	"time"
)

const (
	locationsBasePath     = "locations"
	locationsResourceName = "locations"
)

// LocationService is an interface for interfacing with the location endpoints
// of the Shopify API.
// See: https://help.shopify.com/en/api/reference/inventory/location
type LocationService interface {
	// Retrieves a list of locations
	List(ctx context.Context, options interface{}) ([]Location, error)
	// Retrieves a single location by its Id
	Get(ctx context.Context, id uint64, options interface{}) (*Location, error)
	// Retrieves a count of locations
	Count(ctx context.Context, options interface{}) (int, error)

	// MetafieldsService used for Location resource to communicate with Metafields resource
	MetafieldsService
}

type Location struct {
	// Whether the location is active.If true, then the location can be used to sell products,
	// stock inventory, and fulfill orders.Merchants can deactivate locations from the Shopify admin.
	// Deactivated locations don't contribute to the shop's location limit.
	Active bool `json:"active"`

	// The first line of the address.
	Address1 string `json:"address1"`

	// The second line of the address.
	Address2 string `json:"address2"`

	// The city the location is in.
	City string `json:"city"`

	// The country the location is in.
	Country string `json:"country"`

	// The two-letter code (ISO 3166-1 alpha-2 format) corresponding to country the location is in.
	CountryCode string `json:"country_code"`

	CountryName string `json:"country_name"`

	// The date and time (ISO 8601 format) when the location was created.
	CreatedAt time.Time `json:"created_at"`

	// The Id for the location.
	Id uint64 `json:"id"`

	// Whether this is a fulfillment service location.
	// If true, then the location is a fulfillment service location.
	// If false, then the location was created by the merchant and isn't tied to a fulfillment service.
	Legacy bool `json:"legacy"`

	// The name of the location.
	Name string `json:"name"`

	// The phone number of the location.This value can contain special characters like - and +.
	Phone string `json:"phone"`

	// The province the location is in.
	Province string `json:"province"`

	// The two-letter code corresponding to province or state the location is in.
	ProvinceCode string `json:"province_code"`

	// The date and time (ISO 8601 format) when the location was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// The zip or postal code.
	Zip string `json:"zip"`

	AdminGraphqlApiId string `json:"admin_graphql_api_id"`
}

// LocationServiceOp handles communication with the location related methods of
// the Shopify API.
type LocationServiceOp struct {
	client *Client
}

func (s *LocationServiceOp) List(ctx context.Context, options interface{}) ([]Location, error) {
	path := fmt.Sprintf("%s.json", locationsBasePath)
	resource := new(LocationsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Locations, err
}

func (s *LocationServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*Location, error) {
	path := fmt.Sprintf("%s/%d.json", locationsBasePath, id)
	resource := new(LocationResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Location, err
}

func (s *LocationServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", locationsBasePath)
	return s.client.Count(ctx, path, options)
}

// ListMetafields for a Location resource.
func (s *LocationServiceOp) ListMetafields(ctx context.Context, locationId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a Location resource.
func (s *LocationServiceOp) CountMetafields(ctx context.Context, locationId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a Location resource.
func (s *LocationServiceOp) GetMetafield(ctx context.Context, locationId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a Location resource.
func (s *LocationServiceOp) CreateMetafield(ctx context.Context, locationId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a Location resource.
func (s *LocationServiceOp) UpdateMetafield(ctx context.Context, locationId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a Location resource.
func (s *LocationServiceOp) DeleteMetafield(ctx context.Context, locationId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: locationsResourceName, resourceId: locationId}
	return metafieldService.Delete(ctx, metafieldId)
}

// Represents the result from the locations/X.json endpoint
type LocationResource struct {
	Location *Location `json:"location"`
}

// Represents the result from the locations.json endpoint
type LocationsResource struct {
	Locations []Location `json:"locations"`
}
