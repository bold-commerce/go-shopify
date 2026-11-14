package goshopify

import (
	"context"
	"testing"
)

func TestAddedMetafieldOwners(t *testing.T) {
	for _, owner := range []struct{ name, resource string }{
		{"Article", "articles"},
		{"Blog", "blogs"},
		{"Image", "product_images"},
		{"Collection", "collections"},
	} {
		t.Run(owner.name, func(t *testing.T) {
			// setup installs a fresh client, so select services after setup in each call.
			service := func() MetafieldsService {
				switch owner.name {
				case "Article":
					return client.Article
				case "Blog":
					return client.Blog
				case "Image":
					return client.Image
				default:
					return client.Collection
				}
			}
			for _, tc := range []featureEndpoint{
				{name: "List", method: "GET", path: owner.resource + "/11/metafields.json", response: `{"metafields":[{"id":7,"value":"hello"}]}`, want: []Metafield{{Id: 7, Value: "hello"}}, call: func(ctx context.Context) (interface{}, error) { return service().ListMetafields(ctx, 11, nil) }},
				{name: "Count", method: "GET", path: owner.resource + "/11/metafields/count.json", response: `{"count":3}`, want: 3, call: func(ctx context.Context) (interface{}, error) { return service().CountMetafields(ctx, 11, nil) }},
				{name: "Get", method: "GET", path: owner.resource + "/11/metafields/7.json", response: `{"metafield":{"id":7,"value":"hello"}}`, want: &Metafield{Id: 7, Value: "hello"}, call: func(ctx context.Context) (interface{}, error) { return service().GetMetafield(ctx, 11, 7, nil) }},
				{name: "Create", method: "POST", path: owner.resource + "/11/metafields.json", request: `{"metafield":{"key":"summary","namespace":"custom","type":"single_line_text_field","value":"hello"}}`, response: `{"metafield":{"id":7,"value":"hello"}}`, want: &Metafield{Id: 7, Value: "hello"}, call: func(ctx context.Context) (interface{}, error) {
					return service().CreateMetafield(ctx, 11, Metafield{Key: "summary", Namespace: "custom", Type: MetafieldTypeSingleLineTextField, Value: "hello"})
				}},
				{name: "Update", method: "PUT", path: owner.resource + "/11/metafields/7.json", request: `{"metafield":{"id":7,"value":"hello"}}`, response: `{"metafield":{"id":7,"value":"hello"}}`, want: &Metafield{Id: 7, Value: "hello"}, call: func(ctx context.Context) (interface{}, error) {
					return service().UpdateMetafield(ctx, 11, Metafield{Id: 7, Value: "hello"})
				}},
				{name: "Delete", method: "DELETE", path: owner.resource + "/11/metafields/7.json", want: nil, call: func(ctx context.Context) (interface{}, error) { return nil, service().DeleteMetafield(ctx, 11, 7) }},
			} {
				t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
			}
		})
	}
}
