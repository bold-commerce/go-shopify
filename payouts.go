package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const payoutsBasePath = "shopify_payments/payouts"

// PayoutService is an interface for interfacing with the payout endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-01/resources/payout
type PayoutService interface {
	List(interface{}) ([]Payout, error)
	ListWithPagination(interface{}) ([]Payout, *Pagination, error)
	Get(int64, interface{}) (*Payout, error)
}

// PayoutOp handles communication with the payout related methods of the
// Shopify API.
type PayoutServiceOp struct {
	client *Client
}

// A struct for all available payout list options
type PayoutListOptions struct {
	PageInfo string       `url:"page_info,omitempty"`
	Limit    int          `url:"limit,omitempty"`
	Fields   string       `url:"fields,omitempty"`
	LastID   int64        `url:"last_id,omitempty"`
	SinceID  int64        `url:"since_id,omitempty"`
	Status   PayoutStatus `url:"status,omitempty"`
	DateMin  *time.Time   `url:"date_min,omitempty"`
	DateMax  *time.Time   `url:"date_max,omitempty"`
	Date     *time.Time   `url:"date,omitempty"`
}

// Payout represents a Shopify payout
type Payout struct {
	ID       int64            `json:"id,omitempty"`
	Date     *time.Time       `json:"date,omitempty"`
	Currency string           `json:"currency,omitempty"`
	Amount   *decimal.Decimal `json:"amount,omitempty"`
	Status   PayoutStatus     `json:"status,omitempty"`
}

type PayoutStatus string

const (
	PayoutStatusScheduled PayoutStatus = "scheduled"
	PayoutStatusInTransit PayoutStatus = "in_transit"
	PayoutStatusPaid      PayoutStatus = "paid"
	PayoutStatusFailed    PayoutStatus = "failed"
	PayoutStatusCancelled PayoutStatus = "canceled"
)

// Represents the result from the payouts/X.json endpoint
type PayoutResource struct {
	Payout *Payout `json:"payout"`
}

// Represents the result from the payouts.json endpoint
type PayoutsResource struct {
	Payouts []Payout `json:"payouts"`
}

// List payouts
func (s *PayoutServiceOp) List(options interface{}) ([]Payout, error) {
	payouts, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return payouts, nil
}

func (s *PayoutServiceOp) ListWithPagination(options interface{}) ([]Payout, *Pagination, error) {
	path := fmt.Sprintf("%s.json", payoutsBasePath)
	resource := new(PayoutsResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Payouts, pagination, nil
}

// Get individual payout
func (s *PayoutServiceOp) Get(payoutID int64, options interface{}) (*Payout, error) {
	path := fmt.Sprintf("%s/%d.json", payoutsBasePath, payoutID)
	resource := new(PayoutResource)
	err := s.client.Get(path, resource, options)
	return resource.Payout, err
}
