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

func TestDeprecatedAPICallSupplementalResponse(t *testing.T) {
	at := featureTime(t, "2026-11-01T00:00:00Z")
	deadline := featureTime(t, "2027-01-01T00:00:00Z")
	testFeatureEndpoint(t, featureEndpoint{
		name:     "PopulatedOptionalFields",
		method:   "GET",
		path:     "deprecated_api_calls.json",
		response: `{"data_updated_at":"2026-11-01T00:00:00Z","deprecated_api_calls":[{"api_type":"GraphQL","graphql_schema_name":"QueryRoot.deprecatedField","description":"A deprecated field was used.","documentation_url":"https://shopify.dev/docs/api/usage/versioning","endpoint":"graphql","last_call_at":"2026-11-01T00:00:00Z","migration_deadline":"2027-01-01T00:00:00Z","version":"2026-10"}]}`,
		want: &DeprecatedAPICallsResource{
			DataUpdatedAt: at,
			DeprecatedAPICalls: []DeprecatedAPICall{
				{
					APIType:           "GraphQL",
					GraphQLSchemaName: "QueryRoot.deprecatedField",
					Description:       "A deprecated field was used.",
					DocumentationURL:  "https://shopify.dev/docs/api/usage/versioning",
					Endpoint:          "graphql",
					LastCallAt:        at,
					MigrationDeadline: deadline,
					Version:           "2026-10",
				},
			},
		},
		call: func(ctx context.Context) (interface{}, error) { return client.DeprecatedAPICall.List(ctx, nil) },
	})
}
