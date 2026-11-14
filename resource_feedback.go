package goshopify

import (
	"context"
	"time"
)

// ResourceFeedbackService provides access to ResourceFeedback endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/resourcefeedback
type ResourceFeedbackService interface {
	List(ctx context.Context, options interface{}) ([]ResourceFeedback, error)
	Create(ctx context.Context, feedback ResourceFeedback) (*ResourceFeedback, error)
}

// ResourceFeedbackServiceOp implements ResourceFeedbackService.
type ResourceFeedbackServiceOp struct {
	client *Client
}

// ResourceFeedback represents Shopify resource feedback.
type ResourceFeedback struct {
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	ResourceId          uint64     `json:"resource_id,omitempty"`
	ResourceType        string     `json:"resource_type,omitempty"`
	ResourceUpdatedAt   *time.Time `json:"resource_updated_at,omitempty"`
	Messages            []string   `json:"messages,omitempty"`
	FeedbackGeneratedAt *time.Time `json:"feedback_generated_at,omitempty"`
	State               string     `json:"state,omitempty"`
}

// ResourceFeedbacksResource represents the resource_feedback response.
type ResourceFeedbacksResource struct {
	ResourceFeedbacks []ResourceFeedback `json:"resource_feedback"`
}

// ResourceFeedbackResource represents the resource_feedback response.
type ResourceFeedbackResource struct {
	ResourceFeedback *ResourceFeedback `json:"resource_feedback"`
}

// List retrieves resource feedback records.
func (s *ResourceFeedbackServiceOp) List(ctx context.Context, options interface{}) ([]ResourceFeedback, error) {
	resource := new(ResourceFeedbacksResource)
	err := s.client.Get(ctx, "resource_feedback.json", resource, options)
	return resource.ResourceFeedbacks, err
}

// Create creates a resource feedback record.
func (s *ResourceFeedbackServiceOp) Create(ctx context.Context, feedback ResourceFeedback) (*ResourceFeedback, error) {
	resource := new(ResourceFeedbackResource)
	err := s.client.Post(ctx, "resource_feedback.json", ResourceFeedbackResource{ResourceFeedback: &feedback}, resource)
	return resource.ResourceFeedback, err
}
