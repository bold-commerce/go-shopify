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
// See https://help.shopify.com/api/reference/customcollection
type CustomCollectionService interface {
	List(context.Context, interface{}) ([]CustomCollection, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, int64, interface{}) (*CustomCollection, error)
	Create(context.Context, CustomCollection) (*CustomCollection, error)
	Update(context.Context, CustomCollection) (*CustomCollection, error)
	Delete(context.Context, int64) error

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
	Id             int64       `json:"id"`
	Handle         string      `json:"handle"`
	Title          string      `json:"title"`
	UpdatedAt      *time.Time  `json:"updated_at"`
	BodyHTML       string      `json:"body_html"`
	SortOrder      string      `json:"sort_order"`
	TemplateSuffix string      `json:"template_suffix"`
	Image          Image       `json:"image"`
	Published      bool        `json:"published"`
	PublishedAt    *time.Time  `json:"published_at"`
	PublishedScope string      `json:"published_scope"`
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

// Count custom collections
func (s *CustomCollectionServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", customCollectionsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual custom collection
func (s *CustomCollectionServiceOp) Get(ctx context.Context, collectionId int64, options interface{}) (*CustomCollection, error) {
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
func (s *CustomCollectionServiceOp) Delete(ctx context.Context, collectionId int64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", customCollectionsBasePath, collectionId))
}

// List metafields for a custom collection
func (s *CustomCollectionServiceOp) ListMetafields(ctx context.Context, customCollectionId int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a custom collection
func (s *CustomCollectionServiceOp) CountMetafields(ctx context.Context, customCollectionId int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for a custom collection
func (s *CustomCollectionServiceOp) GetMetafield(ctx context.Context, customCollectionId int64, metafieldId int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for a custom collection
func (s *CustomCollectionServiceOp) CreateMetafield(ctx context.Context, customCollectionId int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for a custom collection
func (s *CustomCollectionServiceOp) UpdateMetafield(ctx context.Context, customCollectionId int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Update(ctx, metafield)
}

// // Delete an existing metafield for a custom collection
func (s *CustomCollectionServiceOp) DeleteMetafield(ctx context.Context, customCollectionId int64, metafieldId int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customCollectionsResourceName, resourceId: customCollectionId}
	return metafieldService.Delete(ctx, metafieldId)
}
