package goshopify

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCountryAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "countries.json",
			response: `{"countries":[{"id":7}]}`,
			want:     []Country{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Country.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "countries/7.json",
			response: `{"country":{"id":7}}`,
			want:     &Country{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Country.Get(ctx, 7, nil)
			},
		},
		{
			name:     "Count",
			method:   "GET",
			path:     "countries/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.Country.Count(ctx, nil)
			},
		},
		{
			name:   "Delete",
			method: "DELETE",
			path:   "countries/7.json",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.Country.Delete(ctx, 7)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/country
func TestCountryDocumentedResponse(t *testing.T) {
	var value CountriesResource
	testFeatureResponse(t, "fixtures/country_documented.json", &value, true)
	if !(len(value.Countries) > 0 && value.Countries[0].Tax.String() == "0.05" && value.Countries[0].Provinces[0].TaxPercentage.String() == "8") {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestCountryWriteFields(t *testing.T) {
	zero := decimal.Zero
	for _, tc := range []featureEndpoint{
		{
			name:     "Create",
			method:   "POST",
			path:     "countries.json",
			request:  `{"country":{"code":"FR"}}`,
			response: `{"country":{"id":7}}`,
			want:     &Country{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Country.Create(ctx, Country{Code: "FR"})
			},
		},
		{
			name:     "UpdateZeroTax",
			method:   "PUT",
			path:     "countries/7.json",
			request:  `{"country":{"id":7,"tax":"0"}}`,
			response: `{"country":{"id":7}}`,
			want:     &Country{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Country.Update(ctx, Country{Id: 7, Tax: &zero})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
