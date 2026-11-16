package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestResourceFeedbackAddedEndpoints(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "resource_feedback.json",
			response: `{"resource_feedback":[{"state":"success"}]}`,
			want:     []ResourceFeedback{{State: "success"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ResourceFeedback.List(ctx, nil)
			},
		},
		{
			name:     "Create",
			method:   "POST",
			path:     "resource_feedback.json",
			request:  `{"resource_feedback":{"feedback_generated_at":"2026-01-15T00:00:00Z","state":"requires_action","messages":["Needs at least one image."]}}`,
			response: `{"resource_feedback":{"state":"requires_action","messages":["Needs at least one image."]}}`,
			want:     &ResourceFeedback{State: "requires_action", Messages: []string{"Needs at least one image."}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.ResourceFeedback.Create(ctx, ResourceFeedback{FeedbackGeneratedAt: &at, State: "requires_action", Messages: []string{"Needs at least one image."}})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/resourcefeedback
func TestResourceFeedbackDocumentedResponse(t *testing.T) {
	var value ResourceFeedbacksResource
	testFeatureResponse(t, "fixtures/resource_feedback_documented.json", &value, true)
	if !(len(value.ResourceFeedbacks) == 1 && value.ResourceFeedbacks[0].ResourceType == "Shop" && value.ResourceFeedbacks[0].FeedbackGeneratedAt != nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestResourceFeedbackRequestFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Create",
		method:   "POST",
		path:     "resource_feedback.json",
		request:  `{"resource_feedback":{"state":"success","feedback_generated_at":"2026-01-15T00:00:00Z"}}`,
		response: `{"resource_feedback":{"state":"success"}}`,
		want:     &ResourceFeedback{State: "success"},
		call: func(ctx context.Context) (interface{}, error) {
			return client.ResourceFeedback.Create(ctx, ResourceFeedback{State: "success", FeedbackGeneratedAt: &at})
		},
	})
}
