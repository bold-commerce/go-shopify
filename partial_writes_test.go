package goshopify

import (
	"context"
	"testing"
)

func TestPartialWrites(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{
			name: "update blog", method: "PUT", path: "blogs/7.json", request: `{"blog": {"id": 7, "title": "News"}}`, response: `{"blog": {"id": 7}}`, want: &Blog{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Blog.Update(ctx, Blog{Id: 7, Title: "News"})
			},
		},
		{
			name: "reset blog template", method: "PUT", path: "blogs/7.json",
			request: `{"blog":{"id":7,"title":"News","template_suffix":""}}`, response: `{"blog":{"id":7,"template_suffix":null}}`, want: &Blog{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				blog := Blog{Id: 7, Title: "News", TemplateSuffix: "custom"}
				blog.SetTemplateSuffix("")
				return client.Blog.Update(ctx, blog)
			},
		},
		{
			name: "set blog template", method: "PUT", path: "blogs/7.json",
			request: `{"blog":{"id":7,"template_suffix":"custom"}}`, response: `{"blog":{"id":7,"template_suffix":"custom"}}`, want: &Blog{Id: 7, TemplateSuffix: "custom"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Blog.Update(ctx, Blog{Id: 7, TemplateSuffix: "custom"})
			},
		},
		{
			name: "update webhook", method: "PUT", path: "webhooks/7.json", request: `{"webhook": {"id": 7, "address": "https://example.com/events"}}`, response: `{"webhook": {"id": 7}}`, want: &Webhook{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Webhook.Update(ctx, Webhook{Id: 7, Address: "https://example.com/events"})
			},
		},
		{
			name: "clear webhook filters", method: "PUT", path: "webhooks/7.json", request: `{"webhook": {"id": 7, "fields": [], "metafield_namespaces": [], "private_metafield_namespaces": []}}`, response: `{"webhook": {"id": 7}}`, want: &Webhook{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Webhook.Update(ctx, Webhook{Id: 7, Fields: []string{}, MetafieldNamespaces: []string{}, PrivateMetafieldNamespaces: []string{}})
			},
		},
		{
			name: "update redirect", method: "PUT", path: "redirects/7.json", request: `{"redirect": {"id": 7, "target": "/new"}}`, response: `{"redirect": {"id": 7}}`, want: &Redirect{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Redirect.Update(ctx, Redirect{Id: 7, Target: "/new"})
			},
		},
		{
			name: "update theme", method: "PUT", path: "themes/7.json", request: `{"theme": {"id": 7, "name": "New"}}`, response: `{"theme": {"id": 7}}`, want: &Theme{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Theme.Update(ctx, Theme{Id: 7, Name: "New", Previewable: true})
			},
		},
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
