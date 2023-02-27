package goshopify

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const productsBasePath = "products"
const productsResourceName = "products"

// linkRegex is used to extract pagination links from product search results.
var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

// ProductService is an interface for interfacing with the product endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/product
type ProductService interface {
	List(interface{}) ([]Product, error)
	ListWithPagination(interface{}) ([]Product, *Pagination, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Product, error)
	Create(Product) (*Product, error)
	Update(Product) (*Product, error)
	Delete(int64) error

	// MetafieldsService used for Product resource to communicate with Metafields resource
	MetafieldsService
}

// ProductServiceOp handles communication with the product related methods of
// the Shopify API.
type ProductServiceOp struct {
	client *Client
}

// Product represents a Shopify product
type Product struct {
	ID                             int64           `json:"id,omitempty"`
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
	Status                         string          `json:"status,omitempty"`
	Options                        []ProductOption `json:"options,omitempty"`
	Variants                       []Variant       `json:"variants,omitempty"`
	Image                          Image           `json:"image,omitempty"`
	Images                         []Image         `json:"images,omitempty"`
	TemplateSuffix                 string          `json:"template_suffix,omitempty"`
	MetafieldsGlobalTitleTag       string          `json:"metafields_global_title_tag,omitempty"`
	MetafieldsGlobalDescriptionTag string          `json:"metafields_global_description_tag,omitempty"`
	Metafields                     []Metafield     `json:"metafields,omitempty"`
	AdminGraphqlAPIID              string          `json:"admin_graphql_api_id,omitempty"`
}

// The options provided by Shopify
type ProductOption struct {
	ID        int64    `json:"id,omitempty"`
	ProductID int64    `json:"product_id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Position  int      `json:"position,omitempty"`
	Values    []string `json:"values,omitempty"`
}

type ProductListOptions struct {
	ListOptions
	CollectionID          int64     `url:"collection_id,omitempty"`
	ProductType           string    `url:"product_type,omitempty"`
	Vendor                string    `url:"vendor,omitempty"`
	Handle                string    `url:"handle,omitempty"`
	PublishedAtMin        time.Time `url:"published_at_min,omitempty"`
	PublishedAtMax        time.Time `url:"published_at_max,omitempty"`
	PublishedStatus       string    `url:"published_status,omitempty"`
	PresentmentCurrencies string    `url:"presentment_currencies,omitempty"`
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
func (s *ProductServiceOp) List(options interface{}) ([]Product, error) {
	products, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (s *ProductServiceOp) ListWithPagination(options interface{}) ([]Product, *Pagination, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	resource := new(ProductsResource)
	headers := http.Header{}

	headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
	if err != nil {
		return nil, nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource.Products, pagination, nil
}

// extractPagination extracts pagination info from linkHeader.
// Details on the format are here:
// https://help.shopify.com/en/api/guides/paginated-rest-results
func extractPagination(linkHeader string) (*Pagination, error) {
	pagination := new(Pagination)

	if linkHeader == "" {
		return pagination, nil
	}

	for _, link := range strings.Split(linkHeader, ",") {
		match := linkRegex.FindStringSubmatch(link)
		// Make sure the link is not empty or invalid
		if len(match) != 3 {
			// We expect 3 values:
			// match[0] = full match
			// match[1] is the URL and match[2] is either 'previous' or 'next'
			err := ResponseDecodingError{
				Message: "could not extract pagination link header",
			}
			return nil, err
		}

		rel, err := url.Parse(match[1])
		if err != nil {
			err = ResponseDecodingError{
				Message: "pagination does not contain a valid URL",
			}
			return nil, err
		}

		params, err := url.ParseQuery(rel.RawQuery)
		if err != nil {
			return nil, err
		}

		paginationListOptions := ListOptions{}

		paginationListOptions.PageInfo = params.Get("page_info")
		if paginationListOptions.PageInfo == "" {
			err = ResponseDecodingError{
				Message: "page_info is missing",
			}
			return nil, err
		}

		limit := params.Get("limit")
		if limit != "" {
			paginationListOptions.Limit, err = strconv.Atoi(params.Get("limit"))
			if err != nil {
				return nil, err
			}
		}

		// 'rel' is either next or previous
		if match[2] == "next" {
			pagination.NextPageOptions = &paginationListOptions
		} else {
			pagination.PreviousPageOptions = &paginationListOptions
		}
	}

	return pagination, nil
}

// Count products
func (s *ProductServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", productsBasePath)
	return s.client.Count(path, options)
}

// Get individual product
func (s *ProductServiceOp) Get(productID int64, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, productID)
	resource := new(ProductResource)
	err := s.client.Get(path, resource, options)
	return resource.Product, err
}

// Create a new product
func (s *ProductServiceOp) Create(product Product) (*Product, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Product, err
}

// Update an existing product
func (s *ProductServiceOp) Update(product Product) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, product.ID)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Product, err
}

// Delete an existing product
func (s *ProductServiceOp) Delete(productID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", productsBasePath, productID))
}

// ListMetafields for a product
func (s *ProductServiceOp) ListMetafields(productID int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.List(options)
}

// Count metafields for a product
func (s *ProductServiceOp) CountMetafields(productID int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Count(options)
}

// GetMetafield for a product
func (s *ProductServiceOp) GetMetafield(productID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Get(metafieldID, options)
}

// CreateMetafield for a product
func (s *ProductServiceOp) CreateMetafield(productID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Create(metafield)
}

// UpdateMetafield for a product
func (s *ProductServiceOp) UpdateMetafield(productID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Update(metafield)
}

// DeleteMetafield for a product
func (s *ProductServiceOp) DeleteMetafield(productID int64, metafieldID int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Delete(metafieldID)
}
