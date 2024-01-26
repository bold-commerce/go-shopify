package goshopify

import (
	"context"

	"github.com/shopspring/decimal"
)

// ShippingZoneService is an interface for interfacing with the shipping zones endpoint
// of the Shopify API.
// See: https://help.shopify.com/api/reference/store-properties/shippingzone
type ShippingZoneService interface {
	List(context.Context) ([]ShippingZone, error)
}

// ShippingZoneServiceOp handles communication with the shipping zone related methods
// of the Shopify API.
type ShippingZoneServiceOp struct {
	client *Client
}

// ShippingZone represents a Shopify shipping zone
type ShippingZone struct {
	Id                           int64                         `json:"id,omitempty"`
	Name                         string                        `json:"name,omitempty"`
	ProfileId                    string                        `json:"profile_id,omitempty"`
	LocationGroupId              string                        `json:"location_group_id,omitempty"`
	AdminGraphqlApiId            string                        `json:"admin_graphql_api_id,omitempty"`
	Countries                    []ShippingCountry             `json:"countries,omitempty"`
	WeightBasedShippingRates     []WeightBasedShippingRate     `json:"weight_based_shipping_rates,omitempty"`
	PriceBasedShippingRates      []PriceBasedShippingRate      `json:"price_based_shipping_rates,omitempty"`
	CarrierShippingRateProviders []CarrierShippingRateProvider `json:"carrier_shipping_rate_providers,omitempty"`
}

// ShippingCountry represents a Shopify shipping country
type ShippingCountry struct {
	Id             int64              `json:"id,omitempty"`
	ShippingZoneId int64              `json:"shipping_zone_id,omitempty"`
	Name           string             `json:"name,omitempty"`
	Tax            *decimal.Decimal   `json:"tax,omitempty"`
	Code           string             `json:"code,omitempty"`
	TaxName        string             `json:"tax_name,omitempty"`
	Provinces      []ShippingProvince `json:"provinces,omitempty"`
}

// ShippingProvince represents a Shopify shipping province
type ShippingProvince struct {
	Id             int64            `json:"id,omitempty"`
	CountryId      int64            `json:"country_id,omitempty"`
	ShippingZoneId int64            `json:"shipping_zone_id,omitempty"`
	Name           string           `json:"name,omitempty"`
	Code           string           `json:"code,omitempty"`
	Tax            *decimal.Decimal `json:"tax,omitempty"`
	TaxName        string           `json:"tax_name,omitempty"`
	TaxType        string           `json:"tax_type,omitempty"`
	TaxPercentage  *decimal.Decimal `json:"tax_percentage,omitempty"`
}

// WeightBasedShippingRate represents a Shopify weight-constrained shipping rate
type WeightBasedShippingRate struct {
	Id             int64            `json:"id,omitempty"`
	ShippingZoneId int64            `json:"shipping_zone_id,omitempty"`
	Name           string           `json:"name,omitempty"`
	Price          *decimal.Decimal `json:"price,omitempty"`
	WeightLow      *decimal.Decimal `json:"weight_low,omitempty"`
	WeightHigh     *decimal.Decimal `json:"weight_high,omitempty"`
}

// PriceBasedShippingRate represents a Shopify subtotal-constrained shipping rate
type PriceBasedShippingRate struct {
	Id               int64            `json:"id,omitempty"`
	ShippingZoneId   int64            `json:"shipping_zone_id,omitempty"`
	Name             string           `json:"name,omitempty"`
	Price            *decimal.Decimal `json:"price,omitempty"`
	MinOrderSubtotal *decimal.Decimal `json:"min_order_subtotal,omitempty"`
	MaxOrderSubtotal *decimal.Decimal `json:"max_order_subtotal,omitempty"`
}

// CarrierShippingRateProvider represents a Shopify carrier-constrained shipping rate
type CarrierShippingRateProvider struct {
	Id               int64             `json:"id,omitempty"`
	CarrierServiceId int64             `json:"carrier_service_id,omitempty"`
	ShippingZoneId   int64             `json:"shipping_zone_id,omitempty"`
	FlatModifier     *decimal.Decimal  `json:"flat_modifier,omitempty"`
	PercentModifier  *decimal.Decimal  `json:"percent_modifier,omitempty"`
	ServiceFilter    map[string]string `json:"service_filter,omitempty"`
}

// Represents the result from the shipping_zones.json endpoint
type ShippingZonesResource struct {
	ShippingZones []ShippingZone `json:"shipping_zones"`
}

// List shipping zones
func (s *ShippingZoneServiceOp) List(ctx context.Context) ([]ShippingZone, error) {
	resource := new(ShippingZonesResource)
	err := s.client.Get(ctx, "shipping_zones.json", resource, nil)
	return resource.ShippingZones, err
}
