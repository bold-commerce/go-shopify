package goshopify

import (
	"context"
	"fmt"
	"time"
)

const (
	customCollectionsBasePath     = "custom_collections"
	customCollectionsResourceName = "collections"
)

// CustomCollectionService is an interface for interacting with the custom
// collection endpoints of the Shopify API.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/customcollection
type CustomCollectionService interface {
	List(context.Context, interface{}) ([]CustomCollection, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]CustomCollection, *Pagination, error)
	ListAll(ctx context.Context, options interface{}) ([]CustomCollection, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*CustomCollection, error)
	Create(context.Context, CustomCollection) (*CustomCollection, error)
	Update(context.Context, CustomCollection) (*CustomCollection, error)
	Delete(context.Context, uint64) error

	// MetafieldsService used for CustomCollection resource to communicate with Metafields resource
	MetafieldsService
}

// CustomCollectionServiceOp handles communication with the custom collection
// related methods of the Shopify API.
type CustomCollectionServiceOp struct {
	client *Client
}

// CustomCollection represents a Shopify custom collection.
type CustomCollection struct {
	Id             uint64      `json:"id,omitempty"`
	Handle         string      `json:"handle,omitempty"`
	Title          string      `json:"title,omitempty"`
	UpdatedAt      *time.Time  `json:"updated_at,omitempty"`
	BodyHTML       *string     `json:"body_html,omitempty"`
	SortOrder      string      `json:"sort_order,omitempty"`
	TemplateSuffix *string     `json:"template_suffix,omitempty"`
	Image          *Image      `json:"image,omitempty"`
	Published      *bool       `json:"published,omitempty"`
	PublishedAt    *time.Time  `json:"published_at,omitempty"`
	PublishedScope string      `json:"published_scope,omitempty"`
	Metafields     []Metafield `json:"metafields,omitempty"`
}

// CustomCollectionResource represents the result form the custom_collections/X.json endpoint
type CustomCollectionResource struct {
	Collection *CustomCollection `json:"custom_collection"`
}

// CustomCollectionsResource represents the result from the custom_collections.json endpoint
type CustomCollectionsResource struct {
	Collections []CustomCollection `json:"custom_collections"`
}

// List custom collections
func (s *CustomCollectionServiceOp) List(ctx context.Context, options interface{}) ([]CustomCollection, error) {
	path := fmt.Sprintf("%s.json", customCollectionsBasePath)
	resource := new(CustomCollectionsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collections, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *CustomCollectionServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]CustomCollection, *Pagination, error) {
	path := fmt.Sprintf("%s.json", customCollectionsBasePath)
	resource := new(CustomCollectionsResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.Collections, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *CustomCollectionServiceOp) ListAll(ctx context.Context, options interface{}) ([]CustomCollection, error) {
	items := []CustomCollection{}
	for {
		page, pagination, err := s.ListWithPagination(ctx, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
}

// Count custom collections
func (s *CustomCollectionServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", customCollectionsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual custom collection
func (s *CustomCollectionServiceOp) Get(ctx context.Context, collectionId uint64, options interface{}) (*CustomCollection, error) {
	path := fmt.Sprintf("%s/%d.json", customCollectionsBasePath, collectionId)
	resource := new(CustomCollectionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collection, err
}

// Create a new custom collection
// See Image for the details of the Image creation for a collection.
func (s *CustomCollectionServiceOp) Create(ctx context.Context, collection CustomCollection) (*CustomCollection, error) {
	path := fmt.Sprintf("%s.json", customCollectionsBasePath)
	wrappedData := CustomCollectionResource{Collection: &collection}
	resource := new(CustomCollectionResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Collection, err
}

// Update an existing custom collection
func (s *CustomCollectionServiceOp) Update(ctx context.Context, collection CustomCollection) (*CustomCollection, error) {
	path := fmt.Sprintf("%s/%d.json", customCollectionsBasePath, collection.Id)
	wrappedData := CustomCollectionResource{Collection: &collection}
	resource := new(CustomCollectionResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Collection, err
}

// Delete an existing custom collection.
func (s *CustomCollectionServiceOp) Delete(ctx context.Context, collectionId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", customCollectionsBasePath, collectionId))
}

// List metafields for a custom collection
func (s *CustomCollectionServiceOp) ListMetafields(ctx context.Context, customCollectionId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a custom collection
func (s *CustomCollectionServiceOp) CountMetafields(ctx context.Context, customCollectionId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for a custom collection
func (s *CustomCollectionServiceOp) GetMetafield(ctx context.Context, customCollectionId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for a custom collection
func (s *CustomCollectionServiceOp) CreateMetafield(ctx context.Context, customCollectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for a custom collection
func (s *CustomCollectionServiceOp) UpdateMetafield(ctx context.Context, customCollectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Update(ctx, metafield)
}

// // Delete an existing metafield for a custom collection
func (s *CustomCollectionServiceOp) DeleteMetafield(ctx context.Context, customCollectionId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Delete(ctx, metafieldId)
}

// ListMetafieldsWithPagination retrieves one page and its continuation options.
func (s *CustomCollectionServiceOp) ListMetafieldsWithPagination(ctx context.Context, customCollectionId uint64, options interface{}) ([]Metafield, *Pagination, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.ListWithPagination(ctx, options)
}

// ListMetafieldsAll retrieves all pages, returning completed pages on error.
func (s *CustomCollectionServiceOp) ListMetafieldsAll(ctx context.Context, customCollectionId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.ListAll(ctx, options)
}
