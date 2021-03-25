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
	Get(int64, interface{}) (*Collect, error)
	Create(Collect) (*Collect, error)
	Delete(int64) error
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

// Get individual collect
func (s *CollectServiceOp) Get(collectID int64, options interface{}) (*Collect, error) {
	path := fmt.Sprintf("%s/%d.json", collectsBasePath, collectID)
	resource := new(CollectResource)
	err := s.client.Get(path, resource, options)
	return resource.Collect, err
}

// Create collects
func (s *CollectServiceOp) Create(collect Collect) (*Collect, error) {
	path := fmt.Sprintf("%s.json", collectsBasePath)
	wrappedData := CollectResource{Collect: &collect}
	resource := new(CollectResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Collect, err
}

// Delete an existing collect
func (s *CollectServiceOp) Delete(collectID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", collectsBasePath, collectID))
}
