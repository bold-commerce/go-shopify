package goshopify

import (
	"context"
	"testing"
)

func TestArticlesAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "ListAuthors",
			method:   "GET",
			path:     "articles/authors.json",
			response: `{"authors":["Alice","Bob"]}`,
			want:     []string{"Alice", "Bob"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Article.ListAuthors(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
