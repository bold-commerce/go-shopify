package goshopify

import (
	"context"
	"testing"
)

func TestPartialWrites(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{
			name: "set customer consent", method: "PUT", path: "customers/7.json", request: `{"customer":{"id":7,"email_marketing_consent":{"state":"subscribed"},"sms_marketing_consent":{"state":"not_subscribed"}}}`, response: `{"customer":{"id":7}}`, want: &Customer{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Customer.Update(ctx, Customer{Id: 7, EmailMarketingConsent: &EmailMarketingConsent{State: "subscribed"}, SMSMarketingConsent: &SMSMarketingConsent{State: "not_subscribed"}})
			},
		},
		{
			name: "update product", method: "PUT", path: "products/7.json", request: `{"product":{"id":7,"title":"Mug"}}`, response: `{"product":{"id":7}}`, want: &Product{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Product.Update(ctx, Product{Id: 7, Title: "Mug", Image: Image{VariantIds: []uint64{}}})
			},
		},
		{
			name: "update product image variants", method: "PUT", path: "products/7.json", request: `{"product":{"id":7,"image":{"variant_ids":[9]}}}`, response: `{"product":{"id":7}}`, want: &Product{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Product.Update(ctx, Product{Id: 7, Image: Image{VariantIds: []uint64{9}}})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
