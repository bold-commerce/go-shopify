package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const abandonedCheckoutsBasePath = "checkouts"

// AbandonedCheckoutService is an interface for interfacing with the abandonedCheckouts endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/abandoned-checkouts
type AbandonedCheckoutService interface {
	List(context.Context, interface{}) ([]AbandonedCheckout, error)
}

// AbandonedCheckoutServiceOp handles communication with the checkout related methods of
// the Shopify API.
type AbandonedCheckoutServiceOp struct {
	client *Client
}

// Represents the result from the checkouts.json endpoint
type AbandonedCheckoutsResource struct {
	AbandonedCheckouts []AbandonedCheckout `json:"checkouts,omitempty"`
}

// AbandonedCheckout represents a Shopify abandoned checkout
type AbandonedCheckout struct {
	Id                       uint64               `json:"id,omitempty"`
	Token                    string               `json:"token,omitempty"`
	CartToken                string               `json:"cart_token,omitempty"`
	Email                    string               `json:"email,omitempty"`
	Gateway                  string               `json:"gateway,omitempty"`
	BuyerAcceptsMarketing    bool                 `json:"buyer_accepts_marketing,omitempty"`
	CreatedAt                *time.Time           `json:"created_at,omitempty"`
	UpdatedAt                *time.Time           `json:"updated_at,omitempty"`
	LandingSite              string               `json:"landing_site,omitempty"`
	Note                     string               `json:"note,omitempty"`
	NoteAttributes           []NoteAttribute      `json:"note_attributes,omitempty"`
	ReferringSite            string               `json:"referring_site,omitempty"`
	ShippingLines            []ShippingLines      `json:"shipping_lines,omitempty"`
	TaxesIncluded            bool                 `json:"taxes_included,omitempty"`
	TotalWeight              int                  `json:"total_weight,omitempty"`
	Currency                 string               `json:"currency,omitempty"`
	CompletedAt              *time.Time           `json:"completed_at,omitempty"`
	ClosedAt                 *time.Time           `json:"closed_at,omitempty"`
	UserId                   uint64               `json:"user_id,omitempty"`
	SourceIdentifier         string               `json:"source_identifier,omitempty"`
	SourceUrl                string               `json:"source_url,omitempty"`
	DeviceId                 uint64               `json:"device_id,omitempty"`
	Phone                    string               `json:"phone,omitempty"`
	CustomerLocale           string               `json:"customer_locale,omitempty"`
	Name                     string               `json:"name,omitempty"`
	Source                   string               `json:"source,omitempty"`
	AbandonedCheckoutUrl     string               `json:"abandoned_checkout_url,omitempty"`
	DiscountCodes            []DiscountCode       `json:"discount_codes,omitempty"`
	TaxLines                 []TaxLine            `json:"tax_lines,omitempty"`
	SourceName               string               `json:"source_name,omitempty"`
	PresentmentCurrency      string               `json:"presentment_currency,omitempty"`
	BuyerAcceptsSmsMarketing bool                 `json:"buyer_accepts_sms_marketing,omitempty"`
	SmsMarketingPhone        string               `json:"sms_marketing_phone,omitempty"`
	TotalDiscounts           *decimal.Decimal     `json:"total_discounts,omitempty"`
	TotalLineItemsPrice      *decimal.Decimal     `json:"total_line_items_price,omitempty"`
	TotalPrice               *decimal.Decimal     `json:"total_price,omitempty"`
	SubtotalPrice            *decimal.Decimal     `json:"subtotal_price,omitempty"`
	TotalDuties              string               `json:"total_duties,omitempty"`
	BillingAddress           *Address             `json:"billing_address,omitempty"`
	ShippingAddress          *Address             `json:"shipping_address,omitempty"`
	Customer                 *Customer            `json:"customer,omitempty"`
	SmsMarketingConsent      *SmsMarketingConsent `json:"sms_marketing_consent,omitempty"`
	AdminGraphqlApiId        string               `json:"admin_graphql_api_id,omitempty"`
	DefaultAddress           *CustomerAddress     `json:"default_address,omitempty"`
}

type SmsMarketingConsent struct {
	State                string     `json:"state,omitempty"`
	OptInLevel           string     `json:"opt_in_level,omitempty"`
	ConsentUpdatedAt     *time.Time `json:"consent_updated_at,omitempty"`
	ConsentCollectedFrom string     `json:"consent_collected_from,omitempty"`
}

// Get abandoned checkout list
func (s *AbandonedCheckoutServiceOp) List(ctx context.Context, options interface{}) ([]AbandonedCheckout, error) {
	path := fmt.Sprintf("/%s.json", abandonedCheckoutsBasePath)
	resource := new(AbandonedCheckoutsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.AbandonedCheckouts, err
}
