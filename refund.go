package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// RefundService provides access to Refund endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/refund
type RefundService interface {
	ListWithPagination(ctx context.Context, orderId uint64, options interface{}) ([]Refund, *Pagination, error)
	List(ctx context.Context, orderId uint64, options interface{}) ([]Refund, error)
	Get(ctx context.Context, orderId uint64, id uint64, options interface{}) (*Refund, error)
	Create(ctx context.Context, orderId uint64, refund Refund) (*Refund, error)
	Calculate(ctx context.Context, orderId uint64, refund Refund) (*Refund, error)
}

// RefundServiceOp implements RefundService.
type RefundServiceOp struct {
	client *Client
}

// RefundShipping represents Shopify refund shipping.
type RefundShipping struct {
	FullRefund        *bool            `json:"full_refund,omitempty"`
	Amount            *decimal.Decimal `json:"amount,omitempty"`
	Tax               *decimal.Decimal `json:"tax,omitempty"`
	MaximumRefundable *decimal.Decimal `json:"maximum_refundable,omitempty"`
}

// RefundDuty represents Shopify refund duty.
type RefundDuty struct {
	DutyId     uint64 `json:"duty_id"`
	RefundType string `json:"refund_type"`
}

// RefundsResource represents the refunds response.
type RefundsResource struct {
	Refunds []Refund `json:"refunds"`
}

// RefundResource represents the refund response.
type RefundResource struct {
	Refund *Refund `json:"refund"`
}

// List retrieves refund records.
func (s *RefundServiceOp) List(ctx context.Context, orderId uint64, options interface{}) ([]Refund, error) {
	items, _, err := s.ListWithPagination(ctx, orderId, options)
	return items, err
}

// Get retrieves a refund record.
func (s *RefundServiceOp) Get(ctx context.Context, orderId uint64, id uint64, options interface{}) (*Refund, error) {
	resource := new(RefundResource)
	err := s.client.Get(ctx, fmt.Sprintf("orders/%d/refunds/%d.json", orderId, id), resource, options)
	return resource.Refund, err
}

// Create creates a refund record.
func (s *RefundServiceOp) Create(ctx context.Context, orderId uint64, refund Refund) (*Refund, error) {
	resource := new(RefundResource)
	err := s.client.Post(ctx, fmt.Sprintf("orders/%d/refunds.json", orderId), RefundResource{Refund: &refund}, resource)
	return resource.Refund, err
}

// Calculate calculates a refund without issuing it.
// Set each returned transaction's Kind to "refund" before passing it to Create.
func (s *RefundServiceOp) Calculate(ctx context.Context, orderId uint64, refund Refund) (*Refund, error) {
	resource := new(RefundResource)
	err := s.client.Post(ctx, fmt.Sprintf("orders/%d/refunds/calculate.json", orderId), RefundResource{Refund: &refund}, resource)
	return resource.Refund, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *RefundServiceOp) ListWithPagination(ctx context.Context, orderId uint64, options interface{}) ([]Refund, *Pagination, error) {
	resource := new(RefundsResource)
	pagination, err := s.client.ListWithPagination(ctx, fmt.Sprintf("orders/%d/refunds.json", orderId), resource, options)
	return resource.Refunds, pagination, err
}
