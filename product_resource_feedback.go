package goshopify

import (
	"context"
	"fmt"
)

// ProductResourceFeedbackService provides access to ProductResourceFeedback endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/product-resourcefeedback
type ProductResourceFeedbackService interface {
	List(ctx context.Context, productId uint64, options interface{}) ([]ResourceFeedback, error)
	Create(ctx context.Context, productId uint64, feedback ResourceFeedback) (*ResourceFeedback, error)
}

// ProductResourceFeedbackServiceOp implements ProductResourceFeedbackService.
type ProductResourceFeedbackServiceOp struct {
	client *Client
}

// List retrieves product resource feedback records.
func (s *ProductResourceFeedbackServiceOp) List(ctx context.Context, productId uint64, options interface{}) ([]ResourceFeedback, error) {
	resource := new(ResourceFeedbacksResource)
	err := s.client.Get(ctx, fmt.Sprintf("products/%d/resource_feedback.json", productId), resource, options)
	return resource.ResourceFeedbacks, err
}

// Create creates a product resource feedback record.
func (s *ProductResourceFeedbackServiceOp) Create(ctx context.Context, productId uint64, feedback ResourceFeedback) (*ResourceFeedback, error) {
	resource := new(ResourceFeedbackResource)
	err := s.client.Post(ctx, fmt.Sprintf("products/%d/resource_feedback.json", productId), ResourceFeedbackResource{ResourceFeedback: &feedback}, resource)
	return resource.ResourceFeedback, err
}
