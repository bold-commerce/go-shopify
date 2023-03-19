package goshopify

import (
	"fmt"
	"time"
)

const blogsBasePath = "blogs"

// BlogService is an interface for interfacing with the blogs endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/online_store/blog
type BlogService interface {
	List(interface{}) ([]Blog, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Blog, error)
	Create(Blog) (*Blog, error)
	Update(Blog) (*Blog, error)
	Delete(int64) error
}

// BlogServiceOp handles communication with the blog related methods of
// the Shopify API.
type BlogServiceOp struct {
	client *Client
}

// Blog represents a Shopify blog
type Blog struct {
	ID                 int64      `json:"id"`
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
	AdminGraphqlAPIID  string     `json:"admin_graphql_api_id,omitempty"`
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
func (s *BlogServiceOp) List(options interface{}) ([]Blog, error) {
	path := fmt.Sprintf("%s.json", blogsBasePath)
	resource := new(BlogsResource)
	err := s.client.Get(path, resource, options)
	return resource.Blogs, err
}

// Count blogs
func (s *BlogServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", blogsBasePath)
	return s.client.Count(path, options)
}

// Get single blog
func (s *BlogServiceOp) Get(blogId int64, options interface{}) (*Blog, error) {
	path := fmt.Sprintf("%s/%d.json", blogsBasePath, blogId)
	resource := new(BlogResource)
	err := s.client.Get(path, resource, options)
	return resource.Blog, err
}

// Create a new blog
func (s *BlogServiceOp) Create(blog Blog) (*Blog, error) {
	path := fmt.Sprintf("%s.json", blogsBasePath)
	wrappedData := BlogResource{Blog: &blog}
	resource := new(BlogResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Blog, err
}

// Update an existing blog
func (s *BlogServiceOp) Update(blog Blog) (*Blog, error) {
	path := fmt.Sprintf("%s/%d.json", blogsBasePath, blog.ID)
	wrappedData := BlogResource{Blog: &blog}
	resource := new(BlogResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Blog, err
}

// Delete an blog
func (s *BlogServiceOp) Delete(blogId int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", blogsBasePath, blogId))
}
