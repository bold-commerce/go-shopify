package goshopify

import (
	"context"

	"github.com/shopspring/decimal"
)

// BalanceService provides access to Balance endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/balance
type BalanceService interface {
	List(ctx context.Context, options interface{}) ([]Balance, error)
}

// BalanceServiceOp implements BalanceService.
type BalanceServiceOp struct {
	client *Client
}

// Balance represents Shopify balance.
type Balance struct {
	Amount   *decimal.Decimal `json:"amount,omitempty"`
	Currency string           `json:"currency,omitempty"`
}

// BalancesResource represents the balance response.
type BalancesResource struct {
	Balances []Balance `json:"balance"`
}

// List retrieves balance records.
func (s *BalanceServiceOp) List(ctx context.Context, options interface{}) ([]Balance, error) {
	resource := new(BalancesResource)
	err := s.client.Get(ctx, "shopify_payments/balance.json", resource, options)
	return resource.Balances, err
}
