package goshopify

import (
	"context"
	"fmt"
	"time"
)

const paymentsTransactionsBasePath = "shopify_payments/balance/transactions"

// PaymentsTransactionsService is an interface for interfacing with the Transactions endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-01/resources/transactions
type PaymentsTransactionsService interface {
	List(context.Context, interface{}) ([]PaymentsTransactions, error)
	ListAll(context.Context, interface{}) ([]PaymentsTransactions, error)
	ListWithPagination(context.Context, interface{}) ([]PaymentsTransactions, *Pagination, error)
	// Deprecated: Shopify does not document an individual balance transaction endpoint. Use List.
	Get(context.Context, uint64, interface{}) (*PaymentsTransactions, error)
}

// PaymentsTransactionsServiceOp handles communication with the transactions related methods of
// the Payment methods of Shopify API.
type PaymentsTransactionsServiceOp struct {
	client *Client
}

// PaymentsTransactionsListOptions filters balance transactions.
// https://shopify.dev/docs/api/admin-rest/latest/resources/transactions#get-shopify-payments-balance-transactions
type PaymentsTransactionsListOptions struct {
	PageInfo     string       `url:"page_info,omitempty"`
	Limit        int          `url:"limit,omitempty"`
	Fields       string       `url:"fields,omitempty"`
	LastId       uint64       `url:"last_id,omitempty"`
	SinceId      uint64       `url:"since_id,omitempty"`
	PayoutId     uint64       `url:"payout_id,omitempty"`
	PayoutStatus PayoutStatus `url:"payout_status,omitempty"`
	Test         *bool        `url:"test,omitempty"`
	// Deprecated: Shopify does not document date_min for balance transactions.
	DateMin *OnlyDate `url:"date_min,omitempty"`
	// Deprecated: Shopify does not document date_max for balance transactions.
	DateMax *OnlyDate `url:"date_max,omitempty"`
	// Deprecated: Shopify does not document processed_at as a balance transaction filter.
	ProcessedAt *OnlyDate `url:"processed_at,omitempty"`
}

// PaymentsTransactions represents a Shopify Transactions
type PaymentsTransactions struct {
	Id                       uint64                    `json:"id,omitempty"`
	Type                     PaymentsTransactionsTypes `json:"type,omitempty"`
	Test                     bool                      `json:"test,omitempty"`
	PayoutId                 int                       `json:"payout_id,omitempty"`
	PayoutStatus             PayoutStatus              `json:"payout_status,omitempty"`
	Currency                 string                    `json:"currency,omitempty"`
	Amount                   string                    `json:"amount,omitempty"`
	Fee                      string                    `json:"fee,omitempty"`
	Net                      string                    `json:"net,omitempty"`
	SourceId                 int                       `json:"source_id,omitempty"`
	SourceType               string                    `json:"source_type,omitempty"`
	SourceOrderTransactionId int                       `json:"source_order_transaction_id,omitempty"`
	SourceOrderId            int                       `json:"source_order_id,omitempty"`
	ProcessedAt              time.Time                 `json:"processed_at,omitempty"`

	AdjustmentReason            string                               `json:"adjustment_reason,omitempty"`
	AdjustmentOrderTransactions []PaymentsAdjustmentOrderTransaction `json:"adjustment_order_transactions,omitempty"`
}

// PaymentsAdjustmentOrderTransaction links a balance adjustment to an order transaction.
type PaymentsAdjustmentOrderTransaction struct {
	Id     uint64                   `json:"id,omitempty"`
	Amount string                   `json:"amount,omitempty"`
	Fee    string                   `json:"fee,omitempty"`
	Net    string                   `json:"net,omitempty"`
	Order  *PaymentsAdjustmentOrder `json:"order,omitempty"`
}

type PaymentsAdjustmentOrder struct {
	Id   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PaymentsTransactionsTypes string

const (
	PaymentsTransactionsCharge             PaymentsTransactionsTypes = "charge"
	PaymentsTransactionsRefund             PaymentsTransactionsTypes = "refund"
	PaymentsTransactionsDispute            PaymentsTransactionsTypes = "dispute"
	PaymentsTransactionsReserve            PaymentsTransactionsTypes = "reserve"
	PaymentsTransactionsAdjustment         PaymentsTransactionsTypes = "adjustment"
	PaymentsTransactionsCredit             PaymentsTransactionsTypes = "credit"
	PaymentsTransactionsDebit              PaymentsTransactionsTypes = "debit"
	PaymentsTransactionsPayout             PaymentsTransactionsTypes = "payout"
	PaymentsTransactionsPayoutFailure      PaymentsTransactionsTypes = "payout_failure"
	PaymentsTransactionsPayoutCancellation PaymentsTransactionsTypes = "payout_cancellation"
)

// Represents the result from the PaymentsTransactions/X.json endpoint
type PaymentsTransactionResource struct {
	PaymentsTransaction *PaymentsTransactions `json:"transaction"`
}

// Represents the result from the PaymentsTransactions.json endpoint
type PaymentsTransactionsResource struct {
	PaymentsTransactions []PaymentsTransactions `json:"transactions"`
}

// List PaymentsTransactions
func (s *PaymentsTransactionsServiceOp) List(ctx context.Context, options interface{}) ([]PaymentsTransactions, error) {
	PaymentsTransactions, _, err := s.ListWithPagination(ctx, options)
	if err != nil {
		return nil, err
	}
	return PaymentsTransactions, nil
}

// ListAll Lists all PaymentsTransactions, iterating over pages
func (s *PaymentsTransactionsServiceOp) ListAll(ctx context.Context, options interface{}) ([]PaymentsTransactions, error) {
	collector := []PaymentsTransactions{}
	err := forEachPage(options, func(options interface{}) (*Pagination, error) {
		entities, pagination, err := s.ListWithPagination(ctx, options)
		if err == nil {
			collector = append(collector, entities...)
		}
		return pagination, err
	})
	return collector, err
}

func (s *PaymentsTransactionsServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]PaymentsTransactions, *Pagination, error) {
	path := fmt.Sprintf("%s.json", paymentsTransactionsBasePath)
	resource := new(PaymentsTransactionsResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.PaymentsTransactions, pagination, nil
}

// Get retrieves an individual balance transaction through an undocumented endpoint.
//
// Deprecated: Shopify documents only the list endpoint. Use List instead.
func (s *PaymentsTransactionsServiceOp) Get(ctx context.Context, payoutId uint64, options interface{}) (*PaymentsTransactions, error) {
	path := fmt.Sprintf("%s/%d.json", paymentsTransactionsBasePath, payoutId)
	resource := new(PaymentsTransactionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.PaymentsTransaction, err
}
