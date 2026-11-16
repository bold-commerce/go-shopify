package goshopify

import (
	"context"
	"fmt"
)

// UserService provides access to User endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/user
type UserService interface {
	ListWithPagination(ctx context.Context, options interface{}) ([]User, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]User, error)
	Get(ctx context.Context, id uint64, options interface{}) (*User, error)
	Current(ctx context.Context, options interface{}) (*User, error)
}

// UserServiceOp implements UserService.
type UserServiceOp struct {
	client *Client
}

// User represents Shopify user.
type User struct {
	AdminGraphqlApiId    string   `json:"admin_graphql_api_id,omitempty"`
	TFAEnabled           bool     `json:"tfa_enabled?,omitempty"`
	Id                   uint64   `json:"id,omitempty"`
	AccountOwner         bool     `json:"account_owner,omitempty"`
	Bio                  string   `json:"bio,omitempty"`
	Email                string   `json:"email,omitempty"`
	FirstName            string   `json:"first_name,omitempty"`
	LastName             string   `json:"last_name,omitempty"`
	IM                   string   `json:"im,omitempty"`
	Permissions          []string `json:"permissions,omitempty"`
	Phone                string   `json:"phone,omitempty"`
	ReceiveAnnouncements int      `json:"receive_announcements,omitempty"`
	ScreenName           string   `json:"screen_name,omitempty"`
	URL                  string   `json:"url,omitempty"`
	Locale               string   `json:"locale,omitempty"`
	UserType             string   `json:"user_type,omitempty"`
}

// UsersResource represents the users response.
type UsersResource struct {
	Users []User `json:"users"`
}

// UserResource represents the user response.
type UserResource struct {
	User *User `json:"user"`
}

// List retrieves user records.
func (s *UserServiceOp) List(ctx context.Context, options interface{}) ([]User, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a user record.
func (s *UserServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*User, error) {
	resource := new(UserResource)
	err := s.client.Get(ctx, fmt.Sprintf("users/%d.json", id), resource, options)
	return resource.User, err
}

// Current retrieves the currently authenticated user.
func (s *UserServiceOp) Current(ctx context.Context, options interface{}) (*User, error) {
	resource := new(UserResource)
	err := s.client.Get(ctx, "users/current.json", resource, options)
	return resource.User, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *UserServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]User, *Pagination, error) {
	resource := new(UsersResource)
	pagination, err := s.client.ListWithPagination(ctx, "users.json", resource, options)
	return resource.Users, pagination, err
}
