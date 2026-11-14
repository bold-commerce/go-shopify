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
	testFeaturePagination(t, "users.json", `{"users":[{"id":7}]}`, []User{{Id: 7}}, `{"users":[{"id":8}]}`, []User{{Id: 8}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
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

func TestUserSupplementalResponse(t *testing.T) {

	testFeatureEndpoint(t, featureEndpoint{
		name:     "PopulatedOptionalFields",
		method:   "GET",
		path:     "users/current.json",
		response: `{"user":{"id":9007199254740993,"bio":"Store owner","email":"owner@example.com","first_name":"Zoë","last_name":"Li","im":"owner-im","screen_name":"owner","phone":"+12065550100","url":"https://example.com/about","permissions":["orders","products"],"account_owner":true,"receive_announcements":0,"locale":"fr","user_type":"regular","admin_graphql_api_id":"gid://shopify/User/9007199254740993","tfa_enabled?":true,"new_permission_field":{}}}`,
		want: &User{
			Id:                9007199254740993,
			Bio:               "Store owner",
			Email:             "owner@example.com",
			FirstName:         "Zoë",
			LastName:          "Li",
			IM:                "owner-im",
			ScreenName:        "owner",
			Phone:             "+12065550100",
			URL:               "https://example.com/about",
			Permissions:       []string{"orders", "products"},
			AccountOwner:      true,
			Locale:            "fr",
			UserType:          "regular",
			AdminGraphqlApiId: "gid://shopify/User/9007199254740993",
			TFAEnabled:        true,
		},
		call: func(ctx context.Context) (interface{}, error) { return client.User.Current(ctx, nil) },
	})
}
