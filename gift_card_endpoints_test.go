package goshopify

import (
	"context"
	"testing"
)

func TestGiftCardAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Search",
			method:   "GET",
			path:     "gift_cards/search.json",
			query:    "query=last_characters%3Aabcd",
			response: `{"gift_cards":[{"id":7}]}`,
			want:     []GiftCard{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.GiftCard.Search(ctx, struct {
					Query string `url:"query"`
				}{"last_characters:abcd"})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestGiftCardSearchWithPagination(t *testing.T) {
	testFeaturePagination(t, "gift_cards/search.json", `{"gift_cards":[{"id":7}]}`, []GiftCard{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.GiftCard.SearchWithPagination(ctx, options)
	})
}

func TestGiftCardQueryFields(t *testing.T) {
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "gift_cards/search.json",
		query:    "query=last_characters%3Aabcd&limit=2&order=created_at+ASC",
		response: `{"gift_cards":[]}`,
		want:     []GiftCard{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.GiftCard.Search(ctx, GiftCardSearchOptions{ListOptions: ListOptions{Limit: 2, Order: "created_at ASC"}, Query: "last_characters:abcd"})
		},
	})
}
