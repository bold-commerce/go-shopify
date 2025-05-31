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
	List(context.Context, interface{}) ([]Blog, error)
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
	Id                 uint64     `json:"id"`
	Title              string     `json:"title"`
	Commentable        string     `json:"commentable"`
	Feedburner         string     `json:"feedburner"`
	FeedburnerLocation string     `json:"feedburner_location"`
	Handle             string     `json:"handle"`
	Metafield          Metafield  `json:"metafield"`
	Tags               string     `json:"tags"`
	TemplateSuffix     string     `json:"template_suffix"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	AdminGraphqlApiId  string     `json:"admin_graphql_api_id,omitempty"`
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
	wrappedData := BlogResource{Blog: &blog}
	resource := new(BlogResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Blog, err
}

// Delete an blog
func (s *BlogServiceOp) Delete(ctx context.Context, blogId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", blogsBasePath, blogId))
}
