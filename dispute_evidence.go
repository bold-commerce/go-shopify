package goshopify

import (
	"context"
	"fmt"
	"time"
)

// DisputeEvidenceService provides access to DisputeEvidence endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/dispute-evidence
type DisputeEvidenceService interface {
	Get(ctx context.Context, disputeId uint64) (*DisputeEvidence, error)
	Update(ctx context.Context, disputeId uint64, evidence DisputeEvidenceUpdate) (*DisputeEvidence, error)
}

// DisputeEvidenceServiceOp implements DisputeEvidenceService.
type DisputeEvidenceServiceOp struct {
	client *Client
}

// DisputeEvidence represents Shopify dispute evidence.
type DisputeEvidence struct {
	Id                           uint64                       `json:"id,omitempty"`
	PaymentsDisputeId            uint64                       `json:"payments_dispute_id,omitempty"`
	ProductDescription           string                       `json:"product_description,omitempty"`
	BillingAddress               *Address                     `json:"billing_address,omitempty"`
	ShippingAddress              *Address                     `json:"shipping_address,omitempty"`
	Fulfillments                 []DisputeEvidenceFulfillment `json:"fulfillments,omitempty"`
	DisputeEvidenceFiles         *DisputeEvidenceFiles        `json:"dispute_evidence_files,omitempty"`
	CreatedAt                    *time.Time                   `json:"created_at,omitempty"`
	UpdatedAt                    *time.Time                   `json:"updated_at,omitempty"`
	SubmittedByMerchantOn        *string                      `json:"submitted_by_merchant_on,omitempty"`
	AccessActivityLog            string                       `json:"access_activity_log,omitempty"`
	CancellationPolicyDisclosure string                       `json:"cancellation_policy_disclosure,omitempty"`
	CancellationRebuttal         string                       `json:"cancellation_rebuttal,omitempty"`
	CustomerEmailAddress         string                       `json:"customer_email_address,omitempty"`
	CustomerFirstName            string                       `json:"customer_first_name,omitempty"`
	CustomerLastName             string                       `json:"customer_last_name,omitempty"`
	RefundPolicyDisclosure       string                       `json:"refund_policy_disclosure,omitempty"`
	RefundRefusalExplanation     string                       `json:"refund_refusal_explanation,omitempty"`
	UncategorizedText            string                       `json:"uncategorized_text,omitempty"`
}

// DisputeEvidenceUpdate represents Shopify dispute evidence update.
type DisputeEvidenceUpdate struct {
	AccessActivityLog            *string                           `json:"access_activity_log,omitempty"`
	CancellationPolicyDisclosure *string                           `json:"cancellation_policy_disclosure,omitempty"`
	CancellationRebuttal         *string                           `json:"cancellation_rebuttal,omitempty"`
	CustomerEmailAddress         *string                           `json:"customer_email_address,omitempty"`
	CustomerFirstName            *string                           `json:"customer_first_name,omitempty"`
	CustomerLastName             *string                           `json:"customer_last_name,omitempty"`
	RefundPolicyDisclosure       *string                           `json:"refund_policy_disclosure,omitempty"`
	RefundRefusalExplanation     *string                           `json:"refund_refusal_explanation,omitempty"`
	UncategorizedText            *string                           `json:"uncategorized_text,omitempty"`
	ShippingAddressAttributes    *Address                          `json:"shipping_address_attributes,omitempty"`
	FulfillmentsAttributes       []DisputeEvidenceFulfillmentInput `json:"fulfillments_attributes,omitempty"`
}

// DisputeEvidenceFulfillment represents Shopify dispute evidence fulfillment.
type DisputeEvidenceFulfillment struct {
	ShippingCarrier        string    `json:"shipping_carrier,omitempty"`
	ShippingTrackingNumber string    `json:"shipping_tracking_number,omitempty"`
	ShippingDate           *OnlyDate `json:"shipping_date,omitempty"`
}

// DisputeEvidenceFulfillmentInput represents Shopify dispute evidence fulfillment input.
type DisputeEvidenceFulfillmentInput struct {
	ShippingCarrier        string     `json:"shipping_carrier,omitempty"`
	ShippingTrackingNumber string     `json:"shipping_tracking_number,omitempty"`
	ShippingDate           *time.Time `json:"shipping_date,omitempty"`
}

// DisputeEvidenceFiles represents Shopify dispute evidence files.
type DisputeEvidenceFiles struct {
	CancellationPolicyFileId    *uint64 `json:"cancellation_policy_file_id,omitempty"`
	CustomerCommunicationFileId *uint64 `json:"customer_communication_file_id,omitempty"`
	RefundPolicyFileId          *uint64 `json:"refund_policy_file_id,omitempty"`
	ServiceDocumentationFileId  *uint64 `json:"service_documentation_file_id,omitempty"`
	ShippingDocumentationFileId *uint64 `json:"shipping_documentation_file_id,omitempty"`
	UncategorizedFileId         *uint64 `json:"uncategorized_file_id,omitempty"`
}

// DisputeEvidenceResource represents the dispute_evidence response.
type DisputeEvidenceResource struct {
	DisputeEvidence *DisputeEvidence `json:"dispute_evidence"`
}

// Get retrieves a dispute evidence record.
func (s *DisputeEvidenceServiceOp) Get(ctx context.Context, disputeId uint64) (*DisputeEvidence, error) {
	resource := new(DisputeEvidenceResource)
	err := s.client.Get(ctx, fmt.Sprintf("shopify_payments/disputes/%d/dispute_evidences.json", disputeId), resource, nil)
	return resource.DisputeEvidence, err
}

// Update saves evidence. The merchant must submit it in the Shopify admin.
func (s *DisputeEvidenceServiceOp) Update(ctx context.Context, disputeId uint64, evidence DisputeEvidenceUpdate) (*DisputeEvidence, error) {
	resource := new(DisputeEvidenceResource)
	err := s.client.Put(ctx, fmt.Sprintf("shopify_payments/disputes/%d/dispute_evidences.json", disputeId), struct {
		Evidence DisputeEvidenceUpdate `json:"dispute_evidence"`
	}{evidence}, resource)
	return resource.DisputeEvidence, err
}
