package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// CountryService provides access to Country endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/country
type CountryService interface {
	List(ctx context.Context, options interface{}) ([]Country, error)
	Get(ctx context.Context, id uint64, options interface{}) (*Country, error)
	Count(ctx context.Context, options interface{}) (int, error)
	Create(ctx context.Context, country Country) (*Country, error)
	Update(ctx context.Context, country Country) (*Country, error)
	Delete(ctx context.Context, id uint64) error
}

// CountryServiceOp implements CountryService.
type CountryServiceOp struct {
	client *Client
}

// Country represents Shopify country.
type Country struct {
	Id        uint64           `json:"id,omitempty"`
	Code      string           `json:"code,omitempty"`
	Name      string           `json:"name,omitempty"`
	Tax       *decimal.Decimal `json:"tax,omitempty"`
	TaxName   string           `json:"tax_name,omitempty"`
	Provinces []Province       `json:"provinces,omitempty"`
}

// CountriesResource represents the countries response.
type CountriesResource struct {
	Countries []Country `json:"countries"`
}

// CountryResource represents the country response.
type CountryResource struct {
	Country *Country `json:"country"`
}

// List retrieves country records.
func (s *CountryServiceOp) List(ctx context.Context, options interface{}) ([]Country, error) {
	resource := new(CountriesResource)
	err := s.client.Get(ctx, "countries.json", resource, options)
	return resource.Countries, err
}

// Get retrieves a country record.
func (s *CountryServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*Country, error) {
	resource := new(CountryResource)
	err := s.client.Get(ctx, fmt.Sprintf("countries/%d.json", id), resource, options)
	return resource.Country, err
}

// Count counts country records.
func (s *CountryServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	return s.client.Count(ctx, "countries/count.json", options)
}

// Create creates a country record.
func (s *CountryServiceOp) Create(ctx context.Context, country Country) (*Country, error) {
	resource := new(CountryResource)
	err := s.client.Post(ctx, "countries.json", CountryResource{Country: &country}, resource)
	return resource.Country, err
}

// Update updates a country record.
func (s *CountryServiceOp) Update(ctx context.Context, country Country) (*Country, error) {
	resource := new(CountryResource)
	err := s.client.Put(ctx, fmt.Sprintf("countries/%d.json", country.Id), CountryResource{Country: &country}, resource)
	return resource.Country, err
}

// Delete deletes a country record.
func (s *CountryServiceOp) Delete(ctx context.Context, id uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("countries/%d.json", id))
}
