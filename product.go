package goshopify

import (
	"context"
	"fmt"
	"regexp"
	"time"
)

const (
	productsBasePath     = "products"
	productsResourceName = "products"
)

// linkRegex is used to extract pagination links from product search results.
var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

// ProductService is an interface for interfacing with the product endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/product
type ProductService interface {
	List(context.Context, interface{}) ([]Product, error)
	ListAll(context.Context, interface{}) ([]Product, error)
	ListWithPagination(context.Context, interface{}) ([]Product, *Pagination, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Product, error)
	Create(context.Context, Product) (*Product, error)
	Update(context.Context, Product) (*Product, error)
	Delete(context.Context, uint64) error

	// MetafieldsService used for Product resource to communicate with Metafields resource
	MetafieldsService
}

// ProductServiceOp handles communication with the product related methods of
// the Shopify API.
type ProductServiceOp struct {
	client *Client
}

// ProductStatus represents a Shopify product status.
type ProductStatus string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/product#resource-object
const (
	// The product is ready to sell and is available to customers on the online store,
	// sales channels, and apps. By default, existing products are set to active.
	ProductStatusActive ProductStatus = "active"

	// The product is no longer being sold and isn't available to customers on sales
	// channels and apps.
	ProductStatusArchived ProductStatus = "archived"

	// The product isn't ready to sell and is unavailable to customers on sales
	// channels and apps. By default, duplicated and unarchived products are set to
	// draft.
	ProductStatusDraft ProductStatus = "draft"
)

// Product represents a Shopify product
type Product struct {
	Id                             uint64          `json:"id,omitempty"`
	Title                          string          `json:"title,omitempty"`
	BodyHTML                       string          `json:"body_html,omitempty"`
	Vendor                         string          `json:"vendor,omitempty"`
	ProductType                    string          `json:"product_type,omitempty"`
	Handle                         string          `json:"handle,omitempty"`
	CreatedAt                      *time.Time      `json:"created_at,omitempty"`
	UpdatedAt                      *time.Time      `json:"updated_at,omitempty"`
	PublishedAt                    *time.Time      `json:"published_at,omitempty"`
	PublishedScope                 string          `json:"published_scope,omitempty"`
	Tags                           string          `json:"tags,omitempty"`
	Status                         ProductStatus   `json:"status,omitempty"`
	Options                        []ProductOption `json:"options,omitempty"`
	Variants                       []Variant       `json:"variants,omitempty"`
	Image                          Image           `json:"image,omitempty"`
	Images                         []Image         `json:"images,omitempty"`
	TemplateSuffix                 string          `json:"template_suffix,omitempty"`
	MetafieldsGlobalTitleTag       string          `json:"metafields_global_title_tag,omitempty"`
	MetafieldsGlobalDescriptionTag string          `json:"metafields_global_description_tag,omitempty"`
	Metafields                     []Metafield     `json:"metafields,omitempty"`
	AdminGraphqlApiId              string          `json:"admin_graphql_api_id,omitempty"`
}

// The options provided by Shopify
type ProductOption struct {
	Id        uint64   `json:"id,omitempty"`
	ProductId uint64   `json:"product_id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Position  int      `json:"position,omitempty"`
	Values    []string `json:"values,omitempty"`
}

type ProductListOptions struct {
	ListOptions
	CollectionId          uint64          `url:"collection_id,omitempty"`
	ProductType           string          `url:"product_type,omitempty"`
	Vendor                string          `url:"vendor,omitempty"`
	Handle                string          `url:"handle,omitempty"`
	PublishedAtMin        time.Time       `url:"published_at_min,omitempty"`
	PublishedAtMax        time.Time       `url:"published_at_max,omitempty"`
	PublishedStatus       string          `url:"published_status,omitempty"`
	PresentmentCurrencies string          `url:"presentment_currencies,omitempty"`
	Status                []ProductStatus `url:"status,omitempty,comma"`
	Title                 string          `url:"title,omitempty"`
}

// Represents the result from the products/X.json endpoint
type ProductResource struct {
	Product *Product `json:"product"`
}

// Represents the result from the products.json endpoint
type ProductsResource struct {
	Products []Product `json:"products"`
}

// Pagination of results
type Pagination struct {
	NextPageOptions     *ListOptions
	PreviousPageOptions *ListOptions
}

// List products
func (s *ProductServiceOp) List(ctx context.Context, options interface{}) ([]Product, error) {
	products, _, err := s.ListWithPagination(ctx, options)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ListAll Lists all products, iterating over pages
func (s *ProductServiceOp) ListAll(ctx context.Context, options interface{}) ([]Product, error) {
	collector := []Product{}

	for {
		entities, pagination, err := s.ListWithPagination(ctx, options)

		if err != nil {
			return collector, err
		}

		collector = append(collector, entities...)

		if pagination.NextPageOptions == nil {
			break
		}

		options = pagination.NextPageOptions
	}

	return collector, nil
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (s *ProductServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Product, *Pagination, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	resource := new(ProductsResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Products, pagination, nil
}

// Count products
func (s *ProductServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", productsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual product
func (s *ProductServiceOp) Get(ctx context.Context, productId uint64, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, productId)
	resource := new(ProductResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Product, err
}

// Create a new product
func (s *ProductServiceOp) Create(ctx context.Context, product Product) (*Product, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Product, err
}

// Update an existing product
func (s *ProductServiceOp) Update(ctx context.Context, product Product) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, product.Id)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Product, err
}

// Delete an existing product
func (s *ProductServiceOp) Delete(ctx context.Context, productId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", productsBasePath, productId))
}

// ListMetafields for a product
func (s *ProductServiceOp) ListMetafields(ctx context.Context, productId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a product
func (s *ProductServiceOp) CountMetafields(ctx context.Context, productId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a product
func (s *ProductServiceOp) GetMetafield(ctx context.Context, productId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a product
func (s *ProductServiceOp) CreateMetafield(ctx context.Context, productId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a product
func (s *ProductServiceOp) UpdateMetafield(ctx context.Context, productId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a product
func (s *ProductServiceOp) DeleteMetafield(ctx context.Context, productId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceId: productId}
	return metafieldService.Delete(ctx, metafieldId)
}
