package goshopify

import (
	"context"
	"fmt"
)

const redirectsBasePath = "redirects"

// RedirectService is an interface for interacting with the redirects
// endpoints of the Shopify API.
// See https://help.shopify.com/api/reference/online_store/redirect
type RedirectService interface {
	List(context.Context, interface{}) ([]Redirect, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, int64, interface{}) (*Redirect, error)
	Create(context.Context, Redirect) (*Redirect, error)
	Update(context.Context, Redirect) (*Redirect, error)
	Delete(context.Context, int64) error
}

// RedirectServiceOp handles communication with the redirect related methods of the
// Shopify API.
type RedirectServiceOp struct {
	client *Client
}

// Redirect represents a Shopify redirect.
type Redirect struct {
	ID     int64  `json:"id"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

// RedirectResource represents the result from the redirects/X.json endpoint
type RedirectResource struct {
	Redirect *Redirect `json:"redirect"`
}

// RedirectsResource represents the result from the redirects.json endpoint
type RedirectsResource struct {
	Redirects []Redirect `json:"redirects"`
}

// List redirects
func (s *RedirectServiceOp) List(ctx context.Context, options interface{}) ([]Redirect, error) {
	path := fmt.Sprintf("%s.json", redirectsBasePath)
	resource := new(RedirectsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Redirects, err
}

// Count redirects
func (s *RedirectServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", redirectsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual redirect
func (s *RedirectServiceOp) Get(ctx context.Context, redirectID int64, options interface{}) (*Redirect, error) {
	path := fmt.Sprintf("%s/%d.json", redirectsBasePath, redirectID)
	resource := new(RedirectResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Redirect, err
}

// Create a new redirect
func (s *RedirectServiceOp) Create(ctx context.Context, redirect Redirect) (*Redirect, error) {
	path := fmt.Sprintf("%s.json", redirectsBasePath)
	wrappedData := RedirectResource{Redirect: &redirect}
	resource := new(RedirectResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Redirect, err
}

// Update an existing redirect
func (s *RedirectServiceOp) Update(ctx context.Context, redirect Redirect) (*Redirect, error) {
	path := fmt.Sprintf("%s/%d.json", redirectsBasePath, redirect.ID)
	wrappedData := RedirectResource{Redirect: &redirect}
	resource := new(RedirectResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Redirect, err
}

// Delete an existing redirect.
func (s *RedirectServiceOp) Delete(ctx context.Context, redirectID int64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", redirectsBasePath, redirectID))
}
