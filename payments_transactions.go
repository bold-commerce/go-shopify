package goshopify

import (
	"fmt"
)

const paymentsTransactionsBasePath = "shopify_payments/balance/transactions"

// PaymentsTransactionsService is an interface for interfacing with the PaymentsTransactions endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-01/resources/transactions
type PaymentsTransactionsService interface {
	List(interface{}) ([]PaymentsTransactions, error)
	ListWithPagination(interface{}) ([]PaymentsTransactions, *Pagination, error)
	Get(int64, interface{}) (*PaymentsTransactions, error)
}

// PaymentsTransactionsServiceOp handles communication with the payout related methods of the
// Shopify API.
type PaymentsTransactionsServiceOp struct {
	client *Client
}

// A struct for all available payout list options
type PaymentsTransactionsListOptions struct {
	PageInfo     string       `url:"page_info,omitempty"`
	Limit        int          `url:"limit,omitempty"`
	Fields       string       `url:"fields,omitempty"`
	LastId       int64        `url:"last_id,omitempty"`
	SinceId      int64        `url:"since_id,omitempty"`
	PayoutId     int64        `url:"payout_id,omitempty"`
	PayoutStatus PayoutStatus `url:"payout_status,omitempty"`
	DateMin      *OnlyDate    `url:"date_min,omitempty"`
	DateMax      *OnlyDate    `url:"date_max,omitempty"`
	ProcessedAt  *OnlyDate    `json:"processed_at,omitempty"`
}

// PaymentsTransactions represents a Shopify payout
type PaymentsTransactions struct {
	ID                       int64                     `json:"id,omitempty"`
	Type                     PaymentsTransactionsTypes `json:"type,omitempty"`
	Test                     bool                      `json:"test,omitempty"`
	PayoutID                 int                       `json:"payout_id,omitempty"`
	PayoutStatus             PayoutStatus              `json:"payout_status,omitempty"`
	Currency                 string                    `json:"currency,omitempty"`
	Amount                   string                    `json:"amount,omitempty"`
	Fee                      string                    `json:"fee,omitempty"`
	Net                      string                    `json:"net,omitempty"`
	SourceID                 int                       `json:"source_id,omitempty"`
	SourceType               string                    `json:"source_type,omitempty"`
	SourceOrderTransactionID int                       `json:"source_order_transaction_id,omitempty"`
	SourceOrderID            int                       `json:"source_order_id,omitempty"`
	ProcessedAt              OnlyDate                  `json:"processed_at,omitempty"`
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
func (s *PaymentsTransactionsServiceOp) List(options interface{}) ([]PaymentsTransactions, error) {
	PaymentsTransactions, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return PaymentsTransactions, nil
}

func (s *PaymentsTransactionsServiceOp) ListWithPagination(options interface{}) ([]PaymentsTransactions, *Pagination, error) {
	path := fmt.Sprintf("%s.json", paymentsTransactionsBasePath)
	resource := new(PaymentsTransactionsResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.PaymentsTransactions, pagination, nil
}

// Get individual payout
func (s *PaymentsTransactionsServiceOp) Get(payoutID int64, options interface{}) (*PaymentsTransactions, error) {
	path := fmt.Sprintf("%s/%d.json", paymentsTransactionsBasePath, payoutID)
	resource := new(PaymentsTransactionResource)
	err := s.client.Get(path, resource, options)
	return resource.PaymentsTransaction, err
}
