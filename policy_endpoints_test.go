package goshopify

import (
	"context"
	"testing"
)

func TestPolicyAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "policies.json",
			response: `{"policies":[{"title":"Refund Policy","body":"30 days"}]}`,
			want:     []Policy{{Title: "Refund Policy", Body: "30 days"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Policy.List(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/policy
func TestPolicyDocumentedResponse(t *testing.T) {
	var value PoliciesResource
	testFeatureResponse(t, "fixtures/policy_documented.json", &value, true)
	if !(len(value.Policies) > 0 && value.Policies[0].Title != "" && value.Policies[0].URL != "") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
