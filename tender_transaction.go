package goshopify

import (
	"context"

	"github.com/shopspring/decimal"
	"time"
)

// TenderTransactionService provides access to TenderTransaction endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/tendertransaction
type TenderTransactionService interface {
	ListWithPagination(ctx context.Context, options interface{}) ([]TenderTransaction, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]TenderTransaction, error)
}

// TenderTransactionServiceOp implements TenderTransactionService.
type TenderTransactionServiceOp struct {
	client *Client
}

// TenderTransaction represents Shopify tender transaction.
type TenderTransaction struct {
	Id              uint64           `json:"id,omitempty"`
	OrderId         uint64           `json:"order_id,omitempty"`
	Amount          *decimal.Decimal `json:"amount,omitempty"`
	Currency        string           `json:"currency,omitempty"`
	UserId          *uint64          `json:"user_id,omitempty"`
	Test            bool             `json:"test,omitempty"`
	ProcessedAt     *time.Time       `json:"processed_at,omitempty"`
	RemoteReference string           `json:"remote_reference,omitempty"`
	PaymentDetails  *PaymentDetails  `json:"payment_details,omitempty"`
	PaymentMethod   string           `json:"payment_method,omitempty"`
}

// TenderTransactionsResource represents the tender_transactions response.
type TenderTransactionsResource struct {
	TenderTransactions []TenderTransaction `json:"tender_transactions"`
}

// TenderTransactionListOptions controls filtering and pagination.
type TenderTransactionListOptions struct {
	ListOptions
	ProcessedAt    time.Time `url:"processed_at,omitempty"`
	ProcessedAtMin time.Time `url:"processed_at_min,omitempty"`
	ProcessedAtMax time.Time `url:"processed_at_max,omitempty"`
}

// List retrieves tender transaction records.
func (s *TenderTransactionServiceOp) List(ctx context.Context, options interface{}) ([]TenderTransaction, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *TenderTransactionServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]TenderTransaction, *Pagination, error) {
	resource := new(TenderTransactionsResource)
	pagination, err := s.client.ListWithPagination(ctx, "tender_transactions.json", resource, options)
	return resource.TenderTransactions, pagination, err
}
