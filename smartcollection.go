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
	Id             uint64      `json:"id"`
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
	Rules          []Rule      `json:"rules"`
	Disjunctive    bool        `json:"disjunctive"`
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
