package goshopify

import (
	"context"
	"fmt"
	"time"
)

const productListingBasePath = "product_listings"

// ProductListingService is an interface for interfacing with the product listing endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/admin-api/rest/reference/sales-channels/productlisting
type ProductListingService interface {
	List(context.Context, interface{}) ([]ProductListing, error)
	ListWithPagination(context.Context, interface{}) ([]ProductListing, *Pagination, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*ProductListing, error)
	GetProductIds(context.Context, interface{}) ([]uint64, error)
	Publish(context.Context, uint64) (*ProductListing, error)
	Delete(context.Context, uint64) error
}

// ProductListingServiceOp handles communication with the product related methods of
// the Shopify API.
type ProductListingServiceOp struct {
	client *Client
}

// ProductListing represents a Shopify product published to your sales channel app
type ProductListing struct {
	Id          uint64          `json:"product_id,omitempty"`
	Title       string          `json:"title,omitempty"`
	BodyHTML    string          `json:"body_html,omitempty"`
	Vendor      string          `json:"vendor,omitempty"`
	ProductType string          `json:"product_type,omitempty"`
	Handle      string          `json:"handle,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
	PublishedAt *time.Time      `json:"published_at,omitempty"`
	Tags        string          `json:"tags,omitempty"`
	Options     []ProductOption `json:"options,omitempty"`
	Variants    []Variant       `json:"variants,omitempty"`
	Images      []Image         `json:"images,omitempty"`
}

// Represents the result from the product_listings/X.json endpoint
type ProductListingResource struct {
	ProductListing *ProductListing `json:"product_listing"`
}

// Represents the result from the product_listings.json endpoint
type ProductsListingsResource struct {
	ProductListings []ProductListing `json:"product_listings"`
}

// Represents the result from the product_listings/product_ids.json endpoint
type ProductListingIdsResource struct {
	ProductIds []uint64 `json:"product_ids"`
}

// Resource which create product_listing endpoint expects in request body
// e.g.
// PUT /admin/api/2020-07/product_listings/921728736.json
//
//	{
//	  "product_listing": {
//	    "product_id": 921728736
//	  }
//	}
type ProductListingPublishResource struct {
	ProductListing struct {
		ProductId uint64 `json:"product_id"`
	} `json:"product_listing"`
}

// List products
func (s *ProductListingServiceOp) List(ctx context.Context, options interface{}) ([]ProductListing, error) {
	products, _, err := s.ListWithPagination(ctx, options)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (s *ProductListingServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]ProductListing, *Pagination, error) {
	path := fmt.Sprintf("%s.json", productListingBasePath)
	resource := new(ProductsListingsResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.ProductListings, pagination, nil
}

// Count products listings published to your sales channel app
func (s *ProductListingServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", productListingBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual product_listing by product Id
func (s *ProductListingServiceOp) Get(ctx context.Context, productId uint64, options interface{}) (*ProductListing, error) {
	path := fmt.Sprintf("%s/%d.json", productListingBasePath, productId)
	resource := new(ProductListingResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.ProductListing, err
}

// GetProductIds lists all product Ids that are published to your sales channel
func (s *ProductListingServiceOp) GetProductIds(ctx context.Context, options interface{}) ([]uint64, error) {
	path := fmt.Sprintf("%s/product_ids.json", productListingBasePath)
	resource := new(ProductListingIdsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.ProductIds, err
}

// Publish an existing product listing to your sales channel app
func (s *ProductListingServiceOp) Publish(ctx context.Context, productId uint64) (*ProductListing, error) {
	path := fmt.Sprintf("%s/%v.json", productListingBasePath, productId)
	wrappedData := new(ProductListingPublishResource)
	wrappedData.ProductListing.ProductId = productId
	resource := new(ProductListingResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.ProductListing, err
}

// Delete unpublishes an existing product from your sales channel app.
func (s *ProductListingServiceOp) Delete(ctx context.Context, productId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", productListingBasePath, productId))
}
