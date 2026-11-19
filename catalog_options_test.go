package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestCatalogRequestOptions(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{
			name: "create unpublished product", method: "POST", path: "products.json",
			request:  `{"product":{"title":"Mug","published":false,"image":{"src":"https://example.com/mug.jpg"}}}`,
			response: `{"product":{"id":1}}`, want: &Product{Id: 1},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Product.Create(ctx, Product{Title: "Mug", Published: BoolPtr(false), Image: Image{Src: "https://example.com/mug.jpg"}})
			},
		},
		{
			name: "create collection with product", method: "POST", path: "custom_collections.json",
			request:  `{"custom_collection":{"title":"Sale","collects":[{"product_id":1}]}}`,
			response: `{"custom_collection":{"id":2}}`, want: &CustomCollection{Id: 2},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CustomCollection.Create(ctx, CustomCollection{Title: "Sale", Collects: []Collect{{ProductId: 1}}})
			},
		},
		{
			name: "create image with metafield", method: "POST", path: "products/1/images.json",
			request:  `{"image":{"src":"https://example.com/mug.jpg","metafields":[{"namespace":"custom","key":"caption","type":"single_line_text_field","value":"Mug"}]}}`,
			response: `{"image":{"id":3}}`, want: &Image{Id: 3},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Image.Create(ctx, 1, Image{Src: "https://example.com/mug.jpg", Metafields: []Metafield{{Namespace: "custom", Key: "caption", Type: "single_line_text_field", Value: "Mug"}}})
			},
		},
		{
			name: "article with image attachment", method: "POST", path: "blogs/4/articles.json",
			request:  `{"article":{"title":"News","updated_at":"2026-10-01T00:00:00Z","image":{"attachment":"YWJj"}}}`,
			response: `{"article":{"id":5}}`, want: &Article{Id: 5},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Article.Create(ctx, 4, Article{Title: "News", UpdatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Image: &ArticleImage{Attachment: "YWJj"}})
			},
		},
		{
			name: "select article fields", method: "GET", path: "blogs/4/articles/5.json", query: "fields=id,title",
			response: `{"article":{"id":5,"title":"News"}}`, want: &Article{Id: 5, Title: "News"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Article.GetWithOptions(ctx, 4, 5, ListOptions{Fields: "id,title"})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
