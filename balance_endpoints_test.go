package goshopify

import (
	"context"
	"testing"
)

func TestBalanceAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "shopify_payments/balance.json",
			response: `{"balance":[{"currency":"USD"}]}`,
			want:     []Balance{{Currency: "USD"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Balance.List(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/balance
func TestBalanceDocumentedResponse(t *testing.T) {
	var value BalancesResource
	testFeatureResponse(t, "fixtures/balance_documented.json", &value, true)
	if !(len(value.Balances) == 1 && value.Balances[0].Amount.String() == "53.99") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
