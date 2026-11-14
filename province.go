package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// ProvinceService provides access to Province endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/province
type ProvinceService interface {
	List(ctx context.Context, countryId uint64, options interface{}) ([]Province, error)
	Get(ctx context.Context, countryId uint64, id uint64, options interface{}) (*Province, error)
	Count(ctx context.Context, countryId uint64, options interface{}) (int, error)
	Update(ctx context.Context, countryId uint64, province Province) (*Province, error)
}

// ProvinceServiceOp implements ProvinceService.
type ProvinceServiceOp struct {
	client *Client
}

// Province represents Shopify province.
type Province struct {
	Id             uint64           `json:"id,omitempty"`
	CountryId      uint64           `json:"country_id,omitempty"`
	ShippingZoneId *uint64          `json:"shipping_zone_id,omitempty"`
	Code           string           `json:"code,omitempty"`
	Name           string           `json:"name,omitempty"`
	Tax            *decimal.Decimal `json:"tax,omitempty"`
	TaxName        string           `json:"tax_name,omitempty"`
	TaxType        string           `json:"tax_type,omitempty"`
	TaxPercentage  *decimal.Decimal `json:"tax_percentage,omitempty"`
}

// ProvincesResource represents the provinces response.
type ProvincesResource struct {
	Provinces []Province `json:"provinces"`
}

// ProvinceResource represents the province response.
type ProvinceResource struct {
	Province *Province `json:"province"`
}

// List retrieves province records.
func (s *ProvinceServiceOp) List(ctx context.Context, countryId uint64, options interface{}) ([]Province, error) {
	resource := new(ProvincesResource)
	err := s.client.Get(ctx, fmt.Sprintf("countries/%d/provinces.json", countryId), resource, options)
	return resource.Provinces, err
}

// Get retrieves a province record.
func (s *ProvinceServiceOp) Get(ctx context.Context, countryId uint64, id uint64, options interface{}) (*Province, error) {
	resource := new(ProvinceResource)
	err := s.client.Get(ctx, fmt.Sprintf("countries/%d/provinces/%d.json", countryId, id), resource, options)
	return resource.Province, err
}

// Count counts province records.
func (s *ProvinceServiceOp) Count(ctx context.Context, countryId uint64, options interface{}) (int, error) {
	return s.client.Count(ctx, fmt.Sprintf("countries/%d/provinces/count.json", countryId), options)
}

// Update updates a province record.
func (s *ProvinceServiceOp) Update(ctx context.Context, countryId uint64, province Province) (*Province, error) {
	resource := new(ProvinceResource)
	err := s.client.Put(ctx, fmt.Sprintf("countries/%d/provinces/%d.json", countryId, province.Id), ProvinceResource{Province: &province}, resource)
	return resource.Province, err
}
