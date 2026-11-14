package goshopify

import (
	"context"
	"fmt"
	"time"
)

const (
	smartCollectionsBasePath     = "smart_collections"
	smartCollectionsResourceName = "collections"
)

// SmartCollectionService is an interface for interacting with the smart
// collection endpoints of the Shopify API.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/smartcollection
type SmartCollectionService interface {
	Order(ctx context.Context, collectionId uint64, options SmartCollectionOrderOptions) error
	List(context.Context, interface{}) ([]SmartCollection, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*SmartCollection, error)
	Create(context.Context, SmartCollection) (*SmartCollection, error)
	Update(context.Context, SmartCollection) (*SmartCollection, error)
	Delete(context.Context, uint64) error

	// MetafieldsService used for SmartCollection resource to communicate with Metafields resource
	MetafieldsService
}

// SmartCollectionServiceOp handles communication with the smart collection
// related methods of the Shopify API.
type SmartCollectionServiceOp struct {
	client *Client
}

type Rule struct {
	Column    string `json:"column"`
	Relation  string `json:"relation"`
	Condition string `json:"condition"`
}

// SmartCollection represents a Shopify smart collection.
type SmartCollection struct {
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
	Rules          []Rule      `json:"rules,omitempty"`
	Disjunctive    *bool       `json:"disjunctive,omitempty"`
	Metafields     []Metafield `json:"metafields,omitempty"`
}

// SmartCollectionResource represents the result from the smart_collections/X.json endpoint
type SmartCollectionResource struct {
	Collection *SmartCollection `json:"smart_collection"`
}

// SmartCollectionsResource represents the result from the smart_collections.json endpoint
type SmartCollectionsResource struct {
	Collections []SmartCollection `json:"smart_collections"`
}

// List smart collections
func (s *SmartCollectionServiceOp) List(ctx context.Context, options interface{}) ([]SmartCollection, error) {
	path := fmt.Sprintf("%s.json", smartCollectionsBasePath)
	resource := new(SmartCollectionsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collections, err
}

// Count smart collections
func (s *SmartCollectionServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", smartCollectionsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual smart collection
func (s *SmartCollectionServiceOp) Get(ctx context.Context, collectionId uint64, options interface{}) (*SmartCollection, error) {
	path := fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collectionId)
	resource := new(SmartCollectionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Collection, err
}

// Create a new smart collection
// See Image for the details of the Image creation for a collection.
func (s *SmartCollectionServiceOp) Create(ctx context.Context, collection SmartCollection) (*SmartCollection, error) {
	path := fmt.Sprintf("%s.json", smartCollectionsBasePath)
	wrappedData := SmartCollectionResource{Collection: &collection}
	resource := new(SmartCollectionResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Collection, err
}

// Update an existing smart collection
func (s *SmartCollectionServiceOp) Update(ctx context.Context, collection SmartCollection) (*SmartCollection, error) {
	path := fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collection.Id)
	wrappedData := SmartCollectionResource{Collection: &collection}
	resource := new(SmartCollectionResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Collection, err
}

// Delete an existing smart collection.
func (s *SmartCollectionServiceOp) Delete(ctx context.Context, collectionId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collectionId))
}

// List metafields for a smart collection
func (s *SmartCollectionServiceOp) ListMetafields(ctx context.Context, smartCollectionId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a smart collection
func (s *SmartCollectionServiceOp) CountMetafields(ctx context.Context, smartCollectionId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for a smart collection
func (s *SmartCollectionServiceOp) GetMetafield(ctx context.Context, smartCollectionId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for a smart collection
func (s *SmartCollectionServiceOp) CreateMetafield(ctx context.Context, smartCollectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for a smart collection
func (s *SmartCollectionServiceOp) UpdateMetafield(ctx context.Context, smartCollectionId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.Update(ctx, metafield)
}

// // Delete an existing metafield for a smart collection
func (s *SmartCollectionServiceOp) DeleteMetafield(ctx context.Context, smartCollectionId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceId: smartCollectionId}
	return metafieldService.Delete(ctx, metafieldId)
}

// SmartCollectionOrderOptions selects a sort order or a manual product order.
type SmartCollectionOrderOptions struct {
	SortOrder string   `url:"sort_order,omitempty"`
	Products  []uint64 `url:"products[],omitempty"`
}

// Order sets a smart collection's sort order or manual product order.
func (s *SmartCollectionServiceOp) Order(ctx context.Context, collectionId uint64, options SmartCollectionOrderOptions) error {
	return s.client.CreateAndDo(ctx, "PUT", fmt.Sprintf("smart_collections/%d/order.json", collectionId), nil, options, nil)
}
