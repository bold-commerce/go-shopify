package goshopify

import (
	"fmt"
	"time"
)

const locationsBasePath = "locations"

// LocationService is an interface for interfacing with the location endpoints
// of the Shopify API.
// See: https://help.shopify.com/en/api/reference/inventory/location
type LocationService interface {
	// Retrieves a list of locations
	List(options interface{}) ([]Location, error)
	// Retrieves a single location by its ID
	Get(ID int64, options interface{}) (*Location, error)
	// Retrieves a count of locations
	Count(options interface{}) (int, error)
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

	// The ID for the location.
	ID int64 `json:"id"`

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

	AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
}

// LocationServiceOp handles communication with the location related methods of
// the Shopify API.
type LocationServiceOp struct {
	client *Client
}

func (s *LocationServiceOp) List(options interface{}) ([]Location, error) {
	path := fmt.Sprintf("%s.json", locationsBasePath)
	resource := new(LocationsResource)
	err := s.client.Get(path, resource, options)
	return resource.Locations, err
}

func (s *LocationServiceOp) Get(ID int64, options interface{}) (*Location, error) {
	path := fmt.Sprintf("%s/%d.json", locationsBasePath, ID)
	resource := new(LocationResource)
	err := s.client.Get(path, resource, options)
	return resource.Location, err
}

func (s *LocationServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", locationsBasePath)
	return s.client.Count(path, options)
}

// Represents the result from the locations/X.json endpoint
type LocationResource struct {
	Location *Location `json:"location"`
}

// Represents the result from the locations.json endpoint
type LocationsResource struct {
	Locations []Location `json:"locations"`
}
