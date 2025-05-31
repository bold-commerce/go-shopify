package goshopify

import (
	"context"
	"fmt"
	"time"
)

const collectsBasePath = "collects"

// CollectService is an interface for interfacing with the collect endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/collect
type CollectService interface {
	List(context.Context, interface{}) ([]Collect, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Collect, error)
	Create(context.Context, Collect) (*Collect, error)
	Delete(context.Context, uint64) error
}

// CollectServiceOp handles communication with the collect related methods of
// the Shopify API.
type CollectServiceOp struct {
	client *Client
}

// Collect represents a Shopify collect
type Collect struct {
	Id           uint64     `json:"id,omitempty"`
	CollectionId uint64     `json:"collection_id,omitempty"`
	ProductId    uint64     `json:"product_id,omitempty"`
	Featured     bool       `json:"featured,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	Position     int        `json:"position,omitempty"`
	SortValue    string     `json:"sort_value,omitempty"`
}

// Represents the result from the collects/X.json endpoint
type CollectResource struct {
	Collect *Collect `json:"collect"`
}

// Represents the result from the collects.json endpoint
type CollectsResource struct {
	Collects []Collect `json:"collects"`
}

// List collects
func (s *CollectServiceOp) List(ctx context.Context, options interface{}) ([]Collect, error) {
	path := fmt.Sprintf("%s.json", collectsBasePath)
	resource := new(CollectsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collects, err
}

// Count collects
func (s *CollectServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", collectsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual collect
func (s *CollectServiceOp) Get(ctx context.Context, collectId uint64, options interface{}) (*Collect, error) {
	path := fmt.Sprintf("%s/%d.json", collectsBasePath, collectId)
	resource := new(CollectResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collect, err
}

// Create collects
func (s *CollectServiceOp) Create(ctx context.Context, collect Collect) (*Collect, error) {
	path := fmt.Sprintf("%s.json", collectsBasePath)
	wrappedData := CollectResource{Collect: &collect}
	resource := new(CollectResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Collect, err
}

// Delete an existing collect
func (s *CollectServiceOp) Delete(ctx context.Context, collectId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", collectsBasePath, collectId))
}
