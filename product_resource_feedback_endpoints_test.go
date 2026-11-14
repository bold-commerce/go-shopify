package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestProductResourceFeedbackAddedEndpoints(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "products/11/resource_feedback.json",
			response: `{"resource_feedback":[{"state":"success"}]}`,
			want:     []ResourceFeedback{{State: "success"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ProductResourceFeedback.List(ctx, 11, nil)
			},
		},
		{
			name:     "Create",
			method:   "POST",
			path:     "products/11/resource_feedback.json",
			request:  `{"resource_feedback":{"feedback_generated_at":"2026-01-15T00:00:00Z","state":"requires_action","messages":["Needs at least one image."]}}`,
			response: `{"resource_feedback":{"state":"requires_action","messages":["Needs at least one image."]}}`,
			want:     &ResourceFeedback{State: "requires_action", Messages: []string{"Needs at least one image."}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ProductResourceFeedback.Create(ctx, 11, ResourceFeedback{FeedbackGeneratedAt: &at, State: "requires_action", Messages: []string{"Needs at least one image."}})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/product-resourcefeedback
func TestProductResourceFeedbackDocumentedResponse(t *testing.T) {
	var value ResourceFeedbacksResource
	testFeatureResponse(t, "fixtures/product_resource_feedback_documented.json", &value, true)
	if !(len(value.ResourceFeedbacks) == 1 && value.ResourceFeedbacks[0].ResourceType == "Product" && value.ResourceFeedbacks[0].ResourceUpdatedAt != nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestProductResourceFeedbackRequestFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Create",
		method:   "POST",
		path:     "products/11/resource_feedback.json",
		request:  `{"resource_feedback":{"state":"success","feedback_generated_at":"2026-01-15T00:00:00Z","resource_updated_at":"2026-01-15T00:00:00Z"}}`,
		response: `{"resource_feedback":{"state":"success"}}`,
		want:     &ResourceFeedback{State: "success"},
		call: func(ctx context.Context) (interface{}, error) {
			return client.ProductResourceFeedback.Create(ctx, 11, ResourceFeedback{State: "success", FeedbackGeneratedAt: &at, ResourceUpdatedAt: &at})
		},
	})
}
