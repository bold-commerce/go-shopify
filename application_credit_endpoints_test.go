package goshopify

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplicationCreditAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "application_credits.json",
			response: `{"application_credits":[{"id":7}]}`,
			want:     []ApplicationCredit{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ApplicationCredit.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "application_credits/7.json",
			response: `{"application_credit":{"id":7}}`,
			want:     &ApplicationCredit{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ApplicationCredit.Get(ctx, 7, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestApplicationCreditNullableFields(t *testing.T) {
	amount := decimal.RequireFromString("5.25")
	enabled, disabled := true, false
	testFeatureEndpoint(t, featureEndpoint{
		method:   "GET",
		path:     "application_credits.json",
		response: `{"application_credits":[{"id":7,"amount":"5.25","test":true},{"id":8,"amount":5.25,"test":false},{"id":9,"amount":null,"test":null}]}`,
		want: []ApplicationCredit{
			{Id: 7, Amount: &amount, Test: &enabled},
			{Id: 8, Amount: &amount, Test: &disabled},
			{Id: 9},
		},
		invalidResponses: []string{`{"application_credits":[{"test":"true"}]}`},
		call: func(ctx context.Context) (interface{}, error) {
			return client.ApplicationCredit.List(ctx, nil)
		},
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/applicationcredit
func TestApplicationCreditDocumentedResponse(t *testing.T) {
	var value ApplicationCreditResource
	testFeatureResponse(t, "fixtures/application_credit_documented.json", &value, true)
	if !(value.ApplicationCredit.Id == 140583599 && value.ApplicationCredit.Amount.String() == "5") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
