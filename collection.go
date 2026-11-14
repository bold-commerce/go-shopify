package goshopify

import (
	"context"
	"fmt"
	"time"
)

const collectionsBasePath = "collections"

// CollectionService is an interface for interfacing with the collection endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/collection
type CollectionService interface {
	MetafieldsService
	Get(ctx context.Context, collectionId uint64, options interface{}) (*Collection, error)
	ListProducts(ctx context.Context, collectionId uint64, options interface{}) ([]Product, error)
	ListProductsWithPagination(ctx context.Context, collectionId uint64, options interface{}) ([]Product, *Pagination, error)
}

// CollectionServiceOp handles communication with the collection related methods of
// the Shopify API.
type CollectionServiceOp struct {
	client *Client
}

// Collection represents a Shopify collection
type Collection struct {
	Id             uint64     `json:"id"`
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
func (s *CollectionServiceOp) Get(ctx context.Context, collectionId uint64, options interface{}) (*Collection, error) {
	path := fmt.Sprintf("%s/%d.json", collectionsBasePath, collectionId)
	resource := new(CollectionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collection, err
}

// List products for a collection
func (s *CollectionServiceOp) ListProducts(ctx context.Context, collectionId uint64, options interface{}) ([]Product, error) {
	products, _, err := s.ListProductsWithPagination(ctx, collectionId, options)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// List products for a collection and return pagination to retrieve next/previous results.
func (s *CollectionServiceOp) ListProductsWithPagination(ctx context.Context, collectionId uint64, options interface{}) ([]Product, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/products.json", collectionsBasePath, collectionId)
	resource := new(ProductsResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Products, pagination, nil
}

// ListMetafields for a Collection resource.
func (s *CollectionServiceOp) ListMetafields(ctx context.Context, collectionId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a Collection resource.
func (s *CollectionServiceOp) CountMetafields(ctx context.Context, collectionId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a Collection resource.
func (s *CollectionServiceOp) GetMetafield(ctx context.Context, collectionId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a Collection resource.
func (s *CollectionServiceOp) CreateMetafield(ctx context.Context, collectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a Collection resource.
func (s *CollectionServiceOp) UpdateMetafield(ctx context.Context, collectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a Collection resource.
func (s *CollectionServiceOp) DeleteMetafield(ctx context.Context, collectionId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "collections", resourceId: collectionId}
	return metafieldService.Delete(ctx, metafieldId)
}
