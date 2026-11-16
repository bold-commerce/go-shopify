package goshopify

import (
	"context"
	"time"
)

// PolicyService provides access to Policy endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/policy
type PolicyService interface {
	List(ctx context.Context, options interface{}) ([]Policy, error)
}

// PolicyServiceOp implements PolicyService.
type PolicyServiceOp struct {
	client *Client
}

// Policy represents Shopify policy.
type Policy struct {
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	URL       string     `json:"url,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Handle    string     `json:"handle,omitempty"`
}

// PoliciesResource represents the policies response.
type PoliciesResource struct {
	Policies []Policy `json:"policies"`
}

// List retrieves policy records.
func (s *PolicyServiceOp) List(ctx context.Context, options interface{}) ([]Policy, error) {
	resource := new(PoliciesResource)
	err := s.client.Get(ctx, "policies.json", resource, options)
	return resource.Policies, err
}
