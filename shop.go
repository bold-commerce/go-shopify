package goshopify

import (
	"context"
	"time"
)

// The shop resource name is empty because it has no resource id
const shopResourceName = ""

// ShopService is an interface for interfacing with the shop endpoint of the
// Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/shop
type ShopService interface {
	Get(ctx context.Context, options interface{}) (*Shop, error)

	// MetafieldsService used for Shop resource to communicate with Metafields resource
	MetafieldsService
}

// ShopServiceOp handles communication with the shop related methods of the
// Shopify API.
type ShopServiceOp struct {
	client *Client
}

// Shop represents a Shopify shop
type Shop struct {
	Id                              uint64     `json:"id"`
	Name                            string     `json:"name"`
	ShopOwner                       string     `json:"shop_owner"`
	Email                           string     `json:"email"`
	CustomerEmail                   string     `json:"customer_email"`
	CreatedAt                       *time.Time `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
	Address1                        string     `json:"address1"`
	Address2                        string     `json:"address2"`
	City                            string     `json:"city"`
	Country                         string     `json:"country"`
	CountryCode                     string     `json:"country_code"`
	CountryName                     string     `json:"country_name"`
	Currency                        string     `json:"currency"`
	Domain                          string     `json:"domain"`
	Latitude                        float64    `json:"latitude"`
	Longitude                       float64    `json:"longitude"`
	Phone                           string     `json:"phone"`
	Province                        string     `json:"province"`
	ProvinceCode                    string     `json:"province_code"`
	Zip                             string     `json:"zip"`
	MoneyFormat                     string     `json:"money_format"`
	MoneyWithCurrencyFormat         string     `json:"money_with_currency_format"`
	WeightUnit                      string     `json:"weight_unit"`
	MyshopifyDomain                 string     `json:"myshopify_domain"`
	PlanName                        string     `json:"plan_name"`
	PlanDisplayName                 string     `json:"plan_display_name"`
	PasswordEnabled                 bool       `json:"password_enabled"`
	PrimaryLocale                   string     `json:"primary_locale"`
	PrimaryLocationId               uint64     `json:"primary_location_id"`
	Timezone                        string     `json:"timezone"`
	IanaTimezone                    string     `json:"iana_timezone"`
	ForceSSL                        bool       `json:"force_ssl"`
	TaxShipping                     bool       `json:"tax_shipping"`
	TaxesIncluded                   bool       `json:"taxes_included"`
	HasStorefront                   bool       `json:"has_storefront"`
	HasDiscounts                    bool       `json:"has_discounts"`
	HasGiftcards                    bool       `json:"has_gift_cards"`
	SetupRequire                    bool       `json:"setup_required"`
	CountyTaxes                     bool       `json:"county_taxes"`
	CheckoutAPISupported            bool       `json:"checkout_api_supported"`
	Source                          string     `json:"source"`
	GoogleAppsDomain                string     `json:"google_apps_domain"`
	GoogleAppsLoginEnabled          bool       `json:"google_apps_login_enabled"`
	MoneyInEmailsFormat             string     `json:"money_in_emails_format"`
	MoneyWithCurrencyInEmailsFormat string     `json:"money_with_currency_in_emails_format"`
	EligibleForPayments             bool       `json:"eligible_for_payments"`
	RequiresExtraPaymentsAgreement  bool       `json:"requires_extra_payments_agreement"`
	PreLaunchEnabled                bool       `json:"pre_launch_enabled"`
}

// Represents the result from the admin/shop.json endpoint
type ShopResource struct {
	Shop *Shop `json:"shop"`
}

// Get shop
func (s *ShopServiceOp) Get(ctx context.Context, options interface{}) (*Shop, error) {
	resource := new(ShopResource)
	err := s.client.Get(ctx, "shop.json", resource, options)
	return resource.Shop, err
}

// ListMetafields for a shop
func (s *ShopServiceOp) ListMetafields(ctx context.Context, _ uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.List(ctx, options)
}

// CountMetafields for a shop
func (s *ShopServiceOp) CountMetafields(ctx context.Context, _ uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a shop
func (s *ShopServiceOp) GetMetafield(ctx context.Context, _ uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a shop
func (s *ShopServiceOp) CreateMetafield(ctx context.Context, _ uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a shop
func (s *ShopServiceOp) UpdateMetafield(ctx context.Context, _ uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a shop
func (s *ShopServiceOp) DeleteMetafield(ctx context.Context, _ uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: shopResourceName}
	return metafieldService.Delete(ctx, metafieldId)
}
