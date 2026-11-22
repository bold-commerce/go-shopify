package goshopify

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
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
	testFeaturePagination(t, "gift_cards/search.json", `{"gift_cards":[{"id":7}]}`, []GiftCard{{Id: 7}}, `{"gift_cards":[{"id":8}]}`, []GiftCard{{Id: 8}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
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

func TestGiftCardCreateForRecipient(t *testing.T) {
	value := decimal.NewFromInt(100)
	sendOn := &OnlyDate{Time: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)}
	testFeatureEndpoint(t, featureEndpoint{
		method: "POST", path: "gift_cards.json",
		request:  `{"gift_card":{"initial_value":"100","recipient_id":42,"message":"Happy birthday!","send_on":"2026-12-31"}}`,
		response: `{"gift_card":{"id":7}}`, want: &GiftCard{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.GiftCard.Create(ctx, GiftCard{InitalValue: &value, RecipientId: 42, Message: "Happy birthday!", SendOn: sendOn})
		},
	})
}
