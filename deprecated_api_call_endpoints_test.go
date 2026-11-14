package goshopify

import (
	"context"
	"testing"
)

func TestDeprecatedAPICallAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "deprecated_api_calls.json",
			response: `{"deprecated_api_calls":[{"endpoint":"Product","version":"2026-10"}]}`,
			want:     &DeprecatedAPICallsResource{DeprecatedAPICalls: []DeprecatedAPICall{{Endpoint: "Product", Version: "2026-10"}}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DeprecatedAPICall.List(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/deprecated-api-calls
func TestDeprecatedApiCallDocumentedResponse(t *testing.T) {
	var value DeprecatedAPICallsResource
	testFeatureResponse(t, "fixtures/deprecated_api_call_documented.json", &value, true)
	if !(value.DataUpdatedAt != nil && len(value.DeprecatedAPICalls) > 0 && value.DeprecatedAPICalls[0].MigrationDeadline != nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}
