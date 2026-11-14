package goshopify

import (
	"context"
	"fmt"
	"time"
)

const articlesBasePath = "articles"

// The ArticlesService allows you to create, publish, and edit articles on a shop's blog
// See: https://shopify.dev/docs/api/admin-rest/stable/resources/article
type ArticlesService interface {
	List(context.Context, uint64, interface{}) ([]Article, error)
	ListWithPagination(ctx context.Context, blogId uint64, options interface{}) ([]Article, *Pagination, error)
	ListAll(ctx context.Context, blogId uint64, options interface{}) ([]Article, error)
	Create(context.Context, uint64, Article) (*Article, error)
	Get(context.Context, uint64, uint64) (*Article, error)
	Update(context.Context, uint64, uint64, Article) (*Article, error)
	Delete(context.Context, uint64, uint64) error
	Count(context.Context, uint64, interface{}) (int, error)
	ListTags(context.Context, interface{}) ([]string, error)
	ListTagsWithPagination(ctx context.Context, options interface{}) ([]string, *Pagination, error)
	ListTagsAll(ctx context.Context, options interface{}) ([]string, error)
	ListBlogTags(context.Context, uint64, interface{}) ([]string, error)
	ListBlogTagsWithPagination(ctx context.Context, blogId uint64, options interface{}) ([]string, *Pagination, error)
	ListBlogTagsAll(ctx context.Context, blogId uint64, options interface{}) ([]string, error)
}

type ArticleResource struct {
	Article *Article `json:"article"`
}

type ArticlesResource struct {
	Articles []Article `json:"articles"`
}

// ArticlesServiceOp handles communication with the articles related methods of
// the Shopify API.
type ArticlesServiceOp struct {
	client *Client
}

type ArticleTagsResource struct {
	Tags []string `json:"tags,omitempty"`
}

type ArticleImage struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Alt       string     `json:"alt,omitempty"`
	Width     int        `json:"width,omitempty"`
	Height    int        `json:"height,omitempty"`
	Src       string     `json:"src,omitempty"`
}

type Article struct {
	Author      string        `json:"author,omitempty"`
	BlogId      uint64        `json:"blog_id,omitempty"`
	BodyHtml    string        `json:"body_html,omitempty"`
	Id          uint64        `json:"id,omitempty"`
	Handle      string        `json:"handle,omitempty"`
	Image       *ArticleImage `json:"image,omitempty"`
	Metafields  []Metafield   `json:"metafields,omitempty"`
	Published   bool          `json:"published,omitempty"`
	SummaryHtml string        `json:"summary_html,omitempty"`
	Tags        string        `json:"tags,omitempty"`
	Title       string        `json:"title,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	UserId      int           `json:"user_id,omitempty"`
	PublishedAt *time.Time    `json:"published_at,omitempty"`
	CreatedAt   *time.Time    `json:"created_at,omitempty"`
}

// List all the articles in a blog.
func (s *ArticlesServiceOp) List(ctx context.Context, blogId uint64, options interface{}) ([]Article, error) {
	path := fmt.Sprintf("%s/%d/%s.json", blogsBasePath, blogId, articlesBasePath)
	resource := new(ArticlesResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Articles, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *ArticlesServiceOp) ListWithPagination(ctx context.Context, blogId uint64, options interface{}) ([]Article, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/%s.json", blogsBasePath, blogId, articlesBasePath)
	resource := new(ArticlesResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.Articles, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *ArticlesServiceOp) ListAll(ctx context.Context, blogId uint64, options interface{}) ([]Article, error) {
	items := []Article{}
	for {
		page, pagination, err := s.ListWithPagination(ctx, blogId, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
}

// Create a article in a blog.
func (s *ArticlesServiceOp) Create(ctx context.Context, blogId uint64, article Article) (*Article, error) {
	path := fmt.Sprintf("%s/%d/%s.json", blogsBasePath, blogId, articlesBasePath)
	body := ArticleResource{
		Article: &article,
	}
	resource := new(ArticleResource)
	err := s.client.Post(ctx, path, body, resource)
	return resource.Article, err
}

// Get an article by blog id and article id.
func (s *ArticlesServiceOp) Get(ctx context.Context, blogId uint64, articleId uint64) (*Article, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", blogsBasePath, blogId, articlesBasePath, articleId)
	resource := new(ArticleResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.Article, err
}

// Update an article in a blog.
func (s *ArticlesServiceOp) Update(ctx context.Context, blogId uint64, articleId uint64, article Article) (*Article, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", blogsBasePath, blogId, articlesBasePath, articleId)
	wrappedData := ArticleResource{Article: &article}
	resource := new(ArticleResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Article, err
}

// Delete an article in a blog.
func (s *ArticlesServiceOp) Delete(ctx context.Context, blogId uint64, articleId uint64) error {
	path := fmt.Sprintf("%s/%d/%s/%d.json", blogsBasePath, blogId, articlesBasePath, articleId)
	return s.client.Delete(ctx, path)
}

// ListTags Get all tags from all articles.
func (s *ArticlesServiceOp) ListTags(ctx context.Context, options interface{}) ([]string, error) {
	path := fmt.Sprintf("%s/tags.json", articlesBasePath)
	articleTags := new(ArticleTagsResource)
	err := s.client.Get(ctx, path, &articleTags, options)
	return articleTags.Tags, err
}

// ListTagsWithPagination returns one page and its next/previous page options.
func (s *ArticlesServiceOp) ListTagsWithPagination(ctx context.Context, options interface{}) ([]string, *Pagination, error) {
	path := fmt.Sprintf("%s/tags.json", articlesBasePath)
	articleTags := new(ArticleTagsResource)
	pagination, err := s.client.ListWithPagination(ctx, path, articleTags, options)
	if err != nil {
		return nil, nil, err
	}
	return articleTags.Tags, pagination, nil
}

// ListTagsAll retrieves all pages. On error, it returns the items from completed pages.
func (s *ArticlesServiceOp) ListTagsAll(ctx context.Context, options interface{}) ([]string, error) {
	items := []string{}
	for {
		page, pagination, err := s.ListTagsWithPagination(ctx, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
}

// Count Articles from a Blog.
func (s *ArticlesServiceOp) Count(ctx context.Context, blogId uint64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/%s/count.json", blogsBasePath, blogId, articlesBasePath)
	return s.client.Count(ctx, path, options)
}

// ListBlogTags Get all tags from all articles in a blog.
func (s *ArticlesServiceOp) ListBlogTags(ctx context.Context, blogId uint64, options interface{}) ([]string, error) {
	path := fmt.Sprintf("%s/%d/%s/tags.json", blogsBasePath, blogId, articlesBasePath)
	articleTags := new(ArticleTagsResource)
	err := s.client.Get(ctx, path, &articleTags, options)
	return articleTags.Tags, err
}

// ListBlogTagsWithPagination returns one page and its next/previous page options.
func (s *ArticlesServiceOp) ListBlogTagsWithPagination(ctx context.Context, blogId uint64, options interface{}) ([]string, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/%s/tags.json", blogsBasePath, blogId, articlesBasePath)
	articleTags := new(ArticleTagsResource)
	pagination, err := s.client.ListWithPagination(ctx, path, articleTags, options)
	if err != nil {
		return nil, nil, err
	}
	return articleTags.Tags, pagination, nil
}

// ListBlogTagsAll retrieves all pages. On error, it returns the items from completed pages.
func (s *ArticlesServiceOp) ListBlogTagsAll(ctx context.Context, blogId uint64, options interface{}) ([]string, error) {
	items := []string{}
	for {
		page, pagination, err := s.ListBlogTagsWithPagination(ctx, blogId, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
}
