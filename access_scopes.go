package goshopify

import (
	"context"
)

type AccessScopesService interface {
	List(context.Context, interface{}) ([]AccessScope, error)
}

type AccessScope struct {
	Handle string `json:"handle,omitempty"`
}

// AccessScopesResource represents the result from the oauth/access_scopes.json endpoint
type AccessScopesResource struct {
	AccessScopes []AccessScope `json:"access_scopes,omitempty"`
}

// AccessScopesServiceOp handles communication with the Access Scopes
// related methods of the Shopify API
type AccessScopesServiceOp struct {
	client *Client
}

// List gets access scopes based on used oauth token
func (s *AccessScopesServiceOp) List(ctx context.Context, options interface{}) ([]AccessScope, error) {
	path := "admin/oauth/access_scopes.json"
	resource := new(AccessScopesResource)
	req, err := s.client.NewRequest(ctx, "GET", path, nil, options)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, resource)
	return resource.AccessScopes, err
}
