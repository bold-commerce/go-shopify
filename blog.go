package goshopify

import (
	"context"
	"fmt"
	"time"
)

const blogsBasePath = "blogs"

// BlogService is an interface for interfacing with the blogs endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/blog
type BlogService interface {
	MetafieldsService
	List(context.Context, interface{}) ([]Blog, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]Blog, *Pagination, error)
	ListAll(ctx context.Context, options interface{}) ([]Blog, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Blog, error)
	Create(context.Context, Blog) (*Blog, error)
	Update(context.Context, Blog) (*Blog, error)
	Delete(context.Context, uint64) error
}

// BlogServiceOp handles communication with the blog related methods of
// the Shopify API.
type BlogServiceOp struct {
	client *Client
}

// Blog represents a Shopify blog
type Blog struct {
	Id                 uint64      `json:"id,omitempty"`
	Title              string      `json:"title,omitempty"`
	Commentable        string      `json:"commentable,omitempty"`
	Feedburner         string      `json:"feedburner,omitempty"`
	FeedburnerLocation string      `json:"feedburner_location,omitempty"`
	Handle             string      `json:"handle,omitempty"`
	Metafields         []Metafield `json:"metafields,omitempty"`
	Tags               string      `json:"tags,omitempty"`
	TemplateSuffix     string      `json:"template_suffix,omitempty"`
	templateSuffixSet  bool
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	AdminGraphqlApiId  string     `json:"admin_graphql_api_id,omitempty"`
}

// SetTemplateSuffix selects the template for Update. Pass "" to restore the default.
func (b *Blog) SetTemplateSuffix(suffix string) {
	b.TemplateSuffix = suffix
	b.templateSuffixSet = true
}

// BlogsResource is the result from the blogs.json endpoint
type BlogsResource struct {
	Blogs []Blog `json:"blogs"`
}

// Represents the result from the blogs/X.json endpoint
type BlogResource struct {
	Blog *Blog `json:"blog"`
}

// List all blogs
func (s *BlogServiceOp) List(ctx context.Context, options interface{}) ([]Blog, error) {
	path := fmt.Sprintf("%s.json", blogsBasePath)
	resource := new(BlogsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Blogs, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *BlogServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Blog, *Pagination, error) {
	path := fmt.Sprintf("%s.json", blogsBasePath)
	resource := new(BlogsResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.Blogs, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *BlogServiceOp) ListAll(ctx context.Context, options interface{}) ([]Blog, error) {
	items := []Blog{}
	err := forEachPage(options, func(options interface{}) (*Pagination, error) {
		page, pagination, err := s.ListWithPagination(ctx, options)
		if err == nil {
			items = append(items, page...)
		}
		return pagination, err
	})
	return items, err
}

// Count blogs
func (s *BlogServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", blogsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get single blog
func (s *BlogServiceOp) Get(ctx context.Context, blogId uint64, options interface{}) (*Blog, error) {
	path := fmt.Sprintf("%s/%d.json", blogsBasePath, blogId)
	resource := new(BlogResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Blog, err
}

// Create a new blog
func (s *BlogServiceOp) Create(ctx context.Context, blog Blog) (*Blog, error) {
	path := fmt.Sprintf("%s.json", blogsBasePath)
	wrappedData := BlogResource{Blog: &blog}
	resource := new(BlogResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Blog, err
}

// Update an existing blog
func (s *BlogServiceOp) Update(ctx context.Context, blog Blog) (*Blog, error) {
	path := fmt.Sprintf("%s/%d.json", blogsBasePath, blog.Id)
	fields := struct {
		*Blog
		TemplateSuffix *string `json:"template_suffix,omitempty"`
	}{Blog: &blog}
	if blog.TemplateSuffix != "" || blog.templateSuffixSet {
		fields.TemplateSuffix = &blog.TemplateSuffix
	}
	wrappedData := map[string]interface{}{"blog": fields}
	resource := new(BlogResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Blog, err
}

// Delete an blog
func (s *BlogServiceOp) Delete(ctx context.Context, blogId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", blogsBasePath, blogId))
}

// ListMetafields for a Blog resource.
func (s *BlogServiceOp) ListMetafields(ctx context.Context, blogId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a Blog resource.
func (s *BlogServiceOp) CountMetafields(ctx context.Context, blogId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a Blog resource.
func (s *BlogServiceOp) GetMetafield(ctx context.Context, blogId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a Blog resource.
func (s *BlogServiceOp) CreateMetafield(ctx context.Context, blogId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a Blog resource.
func (s *BlogServiceOp) UpdateMetafield(ctx context.Context, blogId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a Blog resource.
func (s *BlogServiceOp) DeleteMetafield(ctx context.Context, blogId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.Delete(ctx, metafieldId)
}

// ListMetafieldsWithPagination retrieves one page and its continuation options.
func (s *BlogServiceOp) ListMetafieldsWithPagination(ctx context.Context, blogId uint64, options interface{}) ([]Metafield, *Pagination, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.ListWithPagination(ctx, options)
}

// ListMetafieldsAll retrieves all pages, returning completed pages on error.
func (s *BlogServiceOp) ListMetafieldsAll(ctx context.Context, blogId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: "blogs", resourceId: blogId}
	return metafieldService.ListAll(ctx, options)
}
