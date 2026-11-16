package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"time"
)

// DisputeService provides access to Dispute endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/dispute
type DisputeService interface {
	ListWithPagination(ctx context.Context, options interface{}) ([]Dispute, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]Dispute, error)
	Get(ctx context.Context, id uint64, options interface{}) (*Dispute, error)
}

// DisputeServiceOp implements DisputeService.
type DisputeServiceOp struct {
	client *Client
}

// Dispute represents Shopify dispute.
type Dispute struct {
	Id                uint64           `json:"id,omitempty"`
	OrderId           *uint64          `json:"order_id,omitempty"`
	Type              string           `json:"type,omitempty"`
	Amount            *decimal.Decimal `json:"amount,omitempty"`
	Currency          string           `json:"currency,omitempty"`
	Reason            string           `json:"reason,omitempty"`
	NetworkReasonCode string           `json:"network_reason_code,omitempty"`
	Status            string           `json:"status,omitempty"`
	EvidenceDueBy     *time.Time       `json:"evidence_due_by,omitempty"`
	EvidenceSentOn    *time.Time       `json:"evidence_sent_on,omitempty"`
	FinalizedOn       *time.Time       `json:"finalized_on,omitempty"`
	InitiatedAt       *time.Time       `json:"initiated_at,omitempty"`
}

// DisputesResource represents the disputes response.
type DisputesResource struct {
	Disputes []Dispute `json:"disputes"`
}

// DisputeResource represents the dispute response.
type DisputeResource struct {
	Dispute *Dispute `json:"dispute"`
}

// DisputeListOptions controls filtering and pagination.
type DisputeListOptions struct {
	ListOptions
	InitiatedAt *OnlyDate `url:"initiated_at,omitempty"`
	LastId      uint64    `url:"last_id,omitempty"`
	Status      string    `url:"status,omitempty"`
}

// List retrieves dispute records.
func (s *DisputeServiceOp) List(ctx context.Context, options interface{}) ([]Dispute, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a dispute record.
func (s *DisputeServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*Dispute, error) {
	resource := new(DisputeResource)
	err := s.client.Get(ctx, fmt.Sprintf("shopify_payments/disputes/%d.json", id), resource, options)
	return resource.Dispute, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *DisputeServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Dispute, *Pagination, error) {
	resource := new(DisputesResource)
	pagination, err := s.client.ListWithPagination(ctx, "shopify_payments/disputes.json", resource, options)
	return resource.Disputes, pagination, err
}
