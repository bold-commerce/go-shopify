package goshopify

import (
	"fmt"

	"github.com/shopspring/decimal"
)

const payoutsBasePath = "shopify_payments/payouts"

// PayoutsService is an interface for interfacing with the payouts endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-01/resources/payouts
type PayoutsService interface {
	List(interface{}) ([]Payout, error)
	ListWithPagination(interface{}) ([]Payout, *Pagination, error)
	Get(int64, interface{}) (*Payout, error)
}

// PayoutsServiceOp handles communication with the payout related methods of the
// Shopify API.
type PayoutsServiceOp struct {
	client *Client
}

// A struct for all available payout list options
type PayoutsListOptions struct {
	PageInfo string       `url:"page_info,omitempty"`
	Limit    int          `url:"limit,omitempty"`
	Fields   string       `url:"fields,omitempty"`
	LastId   int64        `url:"last_id,omitempty"`
	SinceId  int64        `url:"since_id,omitempty"`
	Status   PayoutStatus `url:"status,omitempty"`
	DateMin  *OnlyDate    `url:"date_min,omitempty"`
	DateMax  *OnlyDate    `url:"date_max,omitempty"`
	Date     *OnlyDate    `url:"date,omitempty"`
}

// Payout represents a Shopify payout
type Payout struct {
	Id       int64           `json:"id,omitempty"`
	Date     OnlyDate        `json:"date,omitempty"`
	Currency string          `json:"currency,omitempty"`
	Amount   decimal.Decimal `json:"amount,omitempty"`
	Status   PayoutStatus    `json:"status,omitempty"`
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
func (s *PayoutsServiceOp) List(options interface{}) ([]Payout, error) {
	payouts, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return payouts, nil
}

func (s *PayoutsServiceOp) ListWithPagination(options interface{}) ([]Payout, *Pagination, error) {
	path := fmt.Sprintf("%s.json", payoutsBasePath)
	resource := new(PayoutsResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Payouts, pagination, nil
}

// Get individual payout
func (s *PayoutsServiceOp) Get(id int64, options interface{}) (*Payout, error) {
	path := fmt.Sprintf("%s/%d.json", payoutsBasePath, id)
	resource := new(PayoutResource)
	err := s.client.Get(path, resource, options)
	return resource.Payout, err
}
