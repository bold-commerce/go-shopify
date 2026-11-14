package goshopify

import (
	"context"
	"testing"
)

func TestCurrencyAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "currencies.json",
			response: `{"currencies":[{"currency":"USD"}]}`,
			want:     []Currency{{Currency: "USD"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Currency.List(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/currency
func TestCurrencyDocumentedResponse(t *testing.T) {
	var value CurrenciesResource
	testFeatureResponse(t, "fixtures/currency_documented.json", &value, true)
	if !(len(value.Currencies) > 0 && value.Currencies[0].RateUpdatedAt != nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}
