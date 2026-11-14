package goshopify

import (
	"context"
	"time"
)

// DeprecatedAPICallService provides access to DeprecatedAPICall endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/deprecated-api-calls
type DeprecatedAPICallService interface {
	List(ctx context.Context, options interface{}) (*DeprecatedAPICallsResource, error)
}

// DeprecatedAPICallServiceOp implements DeprecatedAPICallService.
type DeprecatedAPICallServiceOp struct {
	client *Client
}

// DeprecatedAPICall represents Shopify deprecated apicall.
type DeprecatedAPICall struct {
	APIType           string     `json:"api_type,omitempty"`
	Description       string     `json:"description,omitempty"`
	DocumentationURL  string     `json:"documentation_url,omitempty"`
	Endpoint          string     `json:"endpoint,omitempty"`
	LastCallAt        *time.Time `json:"last_call_at,omitempty"`
	MigrationDeadline *time.Time `json:"migration_deadline,omitempty"`
	GraphQLSchemaName string     `json:"graphql_schema_name,omitempty"`
	Version           string     `json:"version,omitempty"`
}

// DeprecatedAPICallsResource represents Shopify deprecated apicalls resource.
type DeprecatedAPICallsResource struct {
	DataUpdatedAt      *time.Time          `json:"data_updated_at,omitempty"`
	DeprecatedAPICalls []DeprecatedAPICall `json:"deprecated_api_calls,omitempty"`
}

// List retrieves deprecated apicall records.
func (s *DeprecatedAPICallServiceOp) List(ctx context.Context, options interface{}) (*DeprecatedAPICallsResource, error) {
	resource := new(DeprecatedAPICallsResource)
	err := s.client.Get(ctx, "deprecated_api_calls.json", resource, options)
	return resource, err
}
