package goshopify

import (
	"fmt"
	"time"
)

const collectionsBasePath = "collections"

// CollectionService is an interface for interfacing with the collection endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/products/collection
type CollectionService interface {
	Get(collectionID int64, options interface{}) (*Collection, error)
	ListProducts(collectionID int64, options interface{}) ([]Product, error)
	ListProductsWithPagination(collectionID int64, options interface{}) ([]Product, *Pagination, error)
}

// CollectionServiceOp handles communication with the collection related methods of
// the Shopify API.
type CollectionServiceOp struct {
	client *Client
}

// Collection represents a Shopify collection
type Collection struct {
	ID             int64      `json:"id"`
	Handle         string     `json:"handle"`
	Title          string     `json:"title"`
	UpdatedAt      *time.Time `json:"updated_at"`
	BodyHTML       string     `json:"body_html"`
	SortOrder      string     `json:"sort_order"`
	TemplateSuffix string     `json:"template_suffix"`
	Image          Image      `json:"image"`
	PublishedAt    *time.Time `json:"published_at"`
	PublishedScope string     `json:"published_scope"`
}

// Represents the result from the collections/X.json endpoint
type CollectionResource struct {
	Collection *Collection `json:"collection"`
}

// Get individual collection
func (s *CollectionServiceOp) Get(collectionID int64, options interface{}) (*Collection, error) {
	path := fmt.Sprintf("%s/%d.json", collectionsBasePath, collectionID)
	resource := new(CollectionResource)
	err := s.client.Get(path, resource, options)
	return resource.Collection, err
}

// List products for a collection
func (s *CollectionServiceOp) ListProducts(collectionID int64, options interface{}) ([]Product, error) {
	products, _, err := s.ListProductsWithPagination(collectionID, options)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// List products for a collection and return pagination to retrieve next/previous results.
func (s *CollectionServiceOp) ListProductsWithPagination(collectionID int64, options interface{}) ([]Product, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/products.json", collectionsBasePath, collectionID)
	resource := new(ProductsResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Products, pagination, nil
}
