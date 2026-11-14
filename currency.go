package goshopify

import (
	"context"
	"time"
)

// CurrencyService provides access to Currency endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/currency
type CurrencyService interface {
	List(ctx context.Context, options interface{}) ([]Currency, error)
}

// CurrencyServiceOp implements CurrencyService.
type CurrencyServiceOp struct {
	client *Client
}

// Currency represents Shopify currency.
type Currency struct {
	Enabled       bool       `json:"enabled,omitempty"`
	Currency      string     `json:"currency,omitempty"`
	RateUpdatedAt *time.Time `json:"rate_updated_at,omitempty"`
}

// CurrenciesResource represents the currencies response.
type CurrenciesResource struct {
	Currencies []Currency `json:"currencies"`
}

// List retrieves currency records.
func (s *CurrencyServiceOp) List(ctx context.Context, options interface{}) ([]Currency, error) {
	resource := new(CurrenciesResource)
	err := s.client.Get(ctx, "currencies.json", resource, options)
	return resource.Currencies, err
}
