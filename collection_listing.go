package goshopify

import (
	"context"
	"fmt"
	"time"
)

// CollectionListingService provides access to CollectionListing endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/collectionlisting
type CollectionListingService interface {
	GetProductIdsWithPagination(ctx context.Context, collectionId uint64, options interface{}) ([]uint64, *Pagination, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]CollectionListing, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]CollectionListing, error)
	Get(ctx context.Context, id uint64, options interface{}) (*CollectionListing, error)
	Delete(ctx context.Context, id uint64) error
	Publish(ctx context.Context, collectionId uint64) (*CollectionListing, error)
	GetProductIds(ctx context.Context, collectionId uint64, options interface{}) ([]uint64, error)
}

// CollectionListingServiceOp implements CollectionListingService.
type CollectionListingServiceOp struct {
	client *Client
}

// CollectionListing represents Shopify collection listing.
type CollectionListing struct {
	CollectionId        uint64     `json:"collection_id,omitempty"`
	BodyHTML            string     `json:"body_html,omitempty"`
	DefaultProductImage *Image     `json:"default_product_image,omitempty"`
	Handle              string     `json:"handle,omitempty"`
	Image               *Image     `json:"image,omitempty"`
	Title               string     `json:"title,omitempty"`
	SortOrder           string     `json:"sort_order,omitempty"`
	PublishedAt         *time.Time `json:"published_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

// CollectionListingsResource represents the collection_listings response.
type CollectionListingsResource struct {
	CollectionListings []CollectionListing `json:"collection_listings"`
}

// CollectionListingResource represents the collection_listing response.
type CollectionListingResource struct {
	CollectionListing *CollectionListing `json:"collection_listing"`
}

// CollectionListingProductIdsResource represents the product_ids response.
type CollectionListingProductIdsResource struct {
	ProductIds []uint64 `json:"product_ids"`
}

// List retrieves collection listing records.
func (s *CollectionListingServiceOp) List(ctx context.Context, options interface{}) ([]CollectionListing, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a collection listing record.
func (s *CollectionListingServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*CollectionListing, error) {
	resource := new(CollectionListingResource)
	err := s.client.Get(ctx, fmt.Sprintf("collection_listings/%d.json", id), resource, options)
	return resource.CollectionListing, err
}

// Delete deletes a collection listing record.
func (s *CollectionListingServiceOp) Delete(ctx context.Context, id uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("collection_listings/%d.json", id))
}

// Publish publishes a collection to the app's sales channel.
func (s *CollectionListingServiceOp) Publish(ctx context.Context, collectionId uint64) (*CollectionListing, error) {
	resource := new(CollectionListingResource)
	err := s.client.Put(ctx, fmt.Sprintf("collection_listings/%d.json", collectionId), CollectionListingResource{CollectionListing: &CollectionListing{CollectionId: collectionId}}, resource)
	return resource.CollectionListing, err
}

// GetProductIds retrieves the published product IDs in a collection.
func (s *CollectionListingServiceOp) GetProductIds(ctx context.Context, collectionId uint64, options interface{}) ([]uint64, error) {
	items, _, err := s.GetProductIdsWithPagination(ctx, collectionId, options)
	return items, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *CollectionListingServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]CollectionListing, *Pagination, error) {
	resource := new(CollectionListingsResource)
	pagination, err := s.client.ListWithPagination(ctx, "collection_listings.json", resource, options)
	return resource.CollectionListings, pagination, err
}

// GetProductIdsWithPagination returns a page and its next/previous page options.
func (s *CollectionListingServiceOp) GetProductIdsWithPagination(ctx context.Context, collectionId uint64, options interface{}) ([]uint64, *Pagination, error) {
	resource := new(CollectionListingProductIdsResource)
	pagination, err := s.client.ListWithPagination(ctx, fmt.Sprintf("collection_listings/%d/product_ids.json", collectionId), resource, options)
	return resource.ProductIds, pagination, err
}
