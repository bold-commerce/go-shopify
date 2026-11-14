package goshopify

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

func TestProvinceAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "countries/11/provinces.json",
			response: `{"provinces":[{"id":7}]}`,
			want:     []Province{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Province.List(ctx, 11, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "countries/11/provinces/7.json",
			response: `{"province":{"id":7}}`,
			want:     &Province{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Province.Get(ctx, 11, 7, nil)
			},
		},
		{
			name:     "Count",
			method:   "GET",
			path:     "countries/11/provinces/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.Province.Count(ctx, 11, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/province
func TestProvinceDocumentedResponse(t *testing.T) {
	var value ProvinceResource
	testFeatureResponse(t, "fixtures/province_documented.json", &value, true)
	if !(value.Province.Id == 224293623 && value.Province.Tax.String() == "0.09") {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestProvinceUpdateZeroTax(t *testing.T) {
	zero := decimal.Zero
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Update",
		method:   "PUT",
		path:     "countries/11/provinces/7.json",
		request:  `{"province":{"id":7,"tax":"0","tax_name":"Tax"}}`,
		response: `{"province":{"id":7}}`,
		want:     &Province{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Province.Update(ctx, 11, Province{Id: 7, Tax: &zero, TaxName: "Tax"})
		},
	})
}
