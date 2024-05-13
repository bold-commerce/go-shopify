package goshopify

import (
	"context"
	"fmt"
	"time"
)

const storefrontAccessTokensBasePath = "storefront_access_tokens"

// StorefrontAccessTokenService is an interface for interfacing with the storefront access
// token endpoints of the Shopify API.
// See: https://help.shopify.com/api/reference/access/storefrontaccesstoken
type StorefrontAccessTokenService interface {
	List(context.Context, interface{}) ([]StorefrontAccessToken, error)
	Create(context.Context, StorefrontAccessToken) (*StorefrontAccessToken, error)
	Delete(context.Context, uint64) error
}

// StorefrontAccessTokenServiceOp handles communication with the storefront access token
// related methods of the Shopify API.
type StorefrontAccessTokenServiceOp struct {
	client *Client
}

// StorefrontAccessToken represents a Shopify storefront access token
type StorefrontAccessToken struct {
	Id                uint64     `json:"id,omitempty"`
	Title             string     `json:"title,omitempty"`
	AccessToken       string     `json:"access_token,omitempty"`
	AccessScope       string     `json:"access_scope,omitempty"`
	AdminGraphqlApiId string     `json:"admin_graphql_api_id,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
}

// StorefrontAccessTokenResource represents the result from the admin/storefront_access_tokens.json endpoint
type StorefrontAccessTokenResource struct {
	StorefrontAccessToken *StorefrontAccessToken `json:"storefront_access_token"`
}

// StorefrontAccessTokensResource is the root object for a storefront access tokens get request.
type StorefrontAccessTokensResource struct {
	StorefrontAccessTokens []StorefrontAccessToken `json:"storefront_access_tokens"`
}

// List storefront access tokens
func (s *StorefrontAccessTokenServiceOp) List(ctx context.Context, options interface{}) ([]StorefrontAccessToken, error) {
	path := fmt.Sprintf("%s.json", storefrontAccessTokensBasePath)
	resource := new(StorefrontAccessTokensResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.StorefrontAccessTokens, err
}

// Create a new storefront access token
func (s *StorefrontAccessTokenServiceOp) Create(ctx context.Context, storefrontAccessToken StorefrontAccessToken) (*StorefrontAccessToken, error) {
	path := fmt.Sprintf("%s.json", storefrontAccessTokensBasePath)
	wrappedData := StorefrontAccessTokenResource{StorefrontAccessToken: &storefrontAccessToken}
	resource := new(StorefrontAccessTokenResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.StorefrontAccessToken, err
}

// Delete an existing storefront access token
func (s *StorefrontAccessTokenServiceOp) Delete(ctx context.Context, Id uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", storefrontAccessTokensBasePath, Id))
}
