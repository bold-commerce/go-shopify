package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestTenderTransactionAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "tender_transactions.json",
			response: `{"tender_transactions":[{"id":7}]}`,
			want:     []TenderTransaction{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.TenderTransaction.List(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestTenderTransactionListWithPagination(t *testing.T) {
	testFeaturePagination(t, "tender_transactions.json", `{"tender_transactions":[{"id":7}]}`, []TenderTransaction{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.TenderTransaction.ListWithPagination(ctx, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/tendertransaction
func TestTenderTransactionDocumentedResponse(t *testing.T) {
	var value TenderTransactionsResource
	testFeatureResponse(t, "fixtures/tender_transaction_documented.json", &value, true)
	if !(len(value.TenderTransactions) == 2 && value.TenderTransactions[0].Amount.String() == "250.94" && value.TenderTransactions[0].UserId == nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestTenderTransactionQueryFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "tender_transactions.json",
		query:    "processed_at=2026-01-15T00%3A00%3A00Z&processed_at_min=2026-01-15T00%3A00%3A00Z&processed_at_max=2026-01-15T00%3A00%3A00Z&order=processed_at+ASC",
		response: `{"tender_transactions":[]}`,
		want:     []TenderTransaction{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.TenderTransaction.List(ctx, TenderTransactionListOptions{ProcessedAt: at, ProcessedAtMin: at, ProcessedAtMax: at, ListOptions: ListOptions{Order: "processed_at ASC"}})
		},
	})
}
