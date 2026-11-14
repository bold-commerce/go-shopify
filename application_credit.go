package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// ApplicationCreditService provides access to ApplicationCredit endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/applicationcredit
type ApplicationCreditService interface {
	List(ctx context.Context, options interface{}) ([]ApplicationCredit, error)
	Get(ctx context.Context, id uint64, options interface{}) (*ApplicationCredit, error)
}

// ApplicationCreditServiceOp implements ApplicationCreditService.
type ApplicationCreditServiceOp struct {
	client *Client
}

// ApplicationCredit represents Shopify application credit.
type ApplicationCredit struct {
	Id          uint64           `json:"id,omitempty"`
	Description string           `json:"description,omitempty"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Currency    string           `json:"currency,omitempty"`
	Test        *bool            `json:"test,omitempty"`
}

// ApplicationCreditsResource represents the application_credits response.
type ApplicationCreditsResource struct {
	ApplicationCredits []ApplicationCredit `json:"application_credits"`
}

// ApplicationCreditResource represents the application_credit response.
type ApplicationCreditResource struct {
	ApplicationCredit *ApplicationCredit `json:"application_credit"`
}

// List retrieves application credit records.
func (s *ApplicationCreditServiceOp) List(ctx context.Context, options interface{}) ([]ApplicationCredit, error) {
	resource := new(ApplicationCreditsResource)
	err := s.client.Get(ctx, "application_credits.json", resource, options)
	return resource.ApplicationCredits, err
}

// Get retrieves a application credit record.
func (s *ApplicationCreditServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*ApplicationCredit, error) {
	resource := new(ApplicationCreditResource)
	err := s.client.Get(ctx, fmt.Sprintf("application_credits/%d.json", id), resource, options)
	return resource.ApplicationCredit, err
}
