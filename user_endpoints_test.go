package goshopify

import (
	"context"
	"testing"
)

func TestUserAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "users.json",
			response: `{"users":[{"id":7}]}`,
			want:     []User{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.User.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "users/7.json",
			response: `{"user":{"id":7}}`,
			want:     &User{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.User.Get(ctx, 7, nil)
			},
		},
		{
			name:     "Current",
			method:   "GET",
			path:     "users/current.json",
			response: `{"user":{"id":7,"account_owner":true}}`,
			want:     &User{Id: 7, AccountOwner: true},
			call: func(ctx context.Context) (interface{}, error) {
				return client.User.Current(ctx, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestUserListWithPagination(t *testing.T) {
	testFeaturePagination(t, "users.json", `{"users":[{"id":7}]}`, []User{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.User.ListWithPagination(ctx, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/user
func TestUserDocumentedResponse(t *testing.T) {
	var value UsersResource
	testFeatureResponse(t, "fixtures/user_documented.json", &value, true)
	if !(len(value.Users) == 2 && value.Users[0].AccountOwner && value.Users[0].AdminGraphqlApiId != "") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
