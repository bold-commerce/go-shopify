package goshopify

import (
	"fmt"
	"time"
)

const collectsBasePath = "collects"

// CollectService is an interface for interfacing with the collect endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/products/collect
type CollectService interface {
	List(interface{}) ([]Collect, error)
	Count(interface{}) (int, error)
}

// CollectServiceOp handles communication with the collect related methods of
// the Shopify API.
type CollectServiceOp struct {
	client *Client
}

// Collect represents a Shopify collect
type Collect struct {
	ID           int64      `json:"id,omitempty"`
	CollectionID int64      `json:"collection_id,omitempty"`
	ProductID    int64      `json:"product_id,omitempty"`
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
func (s *CollectServiceOp) List(options interface{}) ([]Collect, error) {
	path := fmt.Sprintf("%s.json", collectsBasePath)
	resource := new(CollectsResource)
	err := s.client.Get(path, resource, options)
	return resource.Collects, err
}

// Count collects
func (s *CollectServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", collectsBasePath)
	return s.client.Count(path, options)
}
