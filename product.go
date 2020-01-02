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

// ProductService is an interface for interfacing with the product endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/product
type ProductService interface {
	List(interface{}) ([]Product, error)
	ListWithPagination(interface{}) ([]Product, Pagination, error)
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
	NextPageOptions		*ListOptions
	PreviousPageOptions *ListOptions
}

// List products
func (s *ProductServiceOp) List(options interface{}) ([]Product, error) {
	products, _, err := s.ListWithPagination(options)
	return products, err
}

// List products and return pagination to retrieve next/previous results
func (s *ProductServiceOp) ListWithPagination(options interface{}) ([]Product, Pagination, error) {
	path := fmt.Sprintf("%s/%s.json", globalApiPathPrefix, productsBasePath)
	resource := new(ProductsResource)
	headers := http.Header{}

	err := s.client.CreateAndDo("GET", path, nil, options, resource, &headers)

	// Extract pagination from header
	linkHeader := ""
	_, ok := headers["Link"]
	if ok {
		linkHeader = headers["Link"][0]
	}
	pagination, err := extractPagination(linkHeader)

	return resource.Products, pagination, err
}

// Extract pagination from the Link header 
// Details on the format are here:
// https://help.shopify.com/en/api/guides/paginated-rest-results
func extractPagination(linkHeader string) (Pagination, error) {
	pagination := Pagination{}
	
	if linkHeader == "" {
		return pagination, nil
	}

	linkRegex := regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

	for _, link := range strings.Split(linkHeader, ",") {
		match := linkRegex.FindStringSubmatch(link); 
		if len(match) == 0 {
			err := ResponseDecodingError{
				Message: "could not extract pagination link header",
			}
			return pagination, err
		}

		rel, err := url.Parse(match[1])
		if err != nil {
			err = ResponseDecodingError{
				Message: "pagination does not contain a valid URL",
			}
			return pagination, err 
		}

		params, _ := url.ParseQuery(rel.RawQuery)

		paginationListOptions := ListOptions{}

		paginationListOptions.PageInfo = getSingleValueParam(params, "page_info")
		if paginationListOptions.PageInfo == "" {
			err = ResponseDecodingError{
				Message: "page_info is missing",
			}
			return pagination, err
		}

		paginationListOptions.Limit, _ = strconv.Atoi(getSingleValueParam(params, "limit"))

		// 'rel' is either next or previous
		if match[2] == "next" {
			pagination.NextPageOptions = &paginationListOptions
		} else {
			pagination.PreviousPageOptions = &paginationListOptions
		}
	}

	return pagination, nil
}

// Gets the value from a param that has only 1 value, defaults to empty string
func getSingleValueParam(params url.Values, name string) string {
	values, ok := params[name]
	if ok {
		return values[0]
	}

	return ""
}

// Count products
func (s *ProductServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%s/count.json", globalApiPathPrefix, productsBasePath)
	return s.client.Count(path, options)
}

// Get individual product
func (s *ProductServiceOp) Get(productID int64, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%s/%d.json", globalApiPathPrefix, productsBasePath, productID)
	resource := new(ProductResource)
	err := s.client.Get(path, resource, options)
	return resource.Product, err
}

// Create a new product
func (s *ProductServiceOp) Create(product Product) (*Product, error) {
	path := fmt.Sprintf("%s/%s.json", globalApiPathPrefix, productsBasePath)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Product, err
}

// Update an existing product
func (s *ProductServiceOp) Update(product Product) (*Product, error) {
	path := fmt.Sprintf("%s/%s/%d.json", globalApiPathPrefix, productsBasePath, product.ID)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Product, err
}

// Delete an existing product
func (s *ProductServiceOp) Delete(productID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%s/%d.json", globalApiPathPrefix, productsBasePath, productID))
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
