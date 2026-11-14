package goshopify

import (
	"context"
	"testing"
)

func TestPriceRuleAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Count",
			method:   "GET",
			path:     "price_rules/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.PriceRule.Count(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
