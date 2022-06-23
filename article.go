package goshopify

import (
	"fmt"
)

const articlesBasePath = "articles"
const articleCountBasePath = "articles/count"

// ArticleService is an interface for interfacing with the article endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/article
type ArticleService interface {
	ListBlog(string, interface{}) ([]Article, error)
	Count(string) (int, error)
}

// ArticleServiceOp handles communication with the Article related methods of
// the Shopify API.
type ArticleServiceOp struct {
	client *Client
}

// Article represents a Shopify article
type Article struct {
	ID                int64  `json:"id,omitempty"`
	Title             string `json:"title,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	BodyHtml          string `json:"body_html,omitempty"`
	ArticleId         int64  `json:"blog_id,omitempty"`
	Author            string `json:"author,omitempty"`
	UserId            int    `json:"user_id,omitempty"`
	PublishedAt       string `json:"published_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
	SummaryHtml       string `json:"summary_html,omitempty"`
	TemplateSuffix    string `json:"template_suffix,omitempty"`
	Handle            string `json:"handle,omitempty"`
	Tags              string `json:"tags,omitempty"`
	AdminGraphqlApiId string `json:"admin_graphql_api_id,omitempty"`
}

type ArticlesResource struct {
	Articles []Article `json:"articles"`
}

type ArticleCountResource struct {
	Count int `json:"count"`
}

type ArticleListOptions struct {
	ListOptions
}

// Retrieves a list of all articles from a blog
func (s *ArticleServiceOp) ListBlog(blogID string, options interface{}) ([]Article, error) {
	path := fmt.Sprintf("blogs/%s/%s.json", blogID, articlesBasePath)
	resource := new(ArticlesResource)
	err := s.client.Get(path, resource, options)
	return resource.Articles, err
}

// Retrieves a count of all articles from a blog
func (s *ArticleServiceOp) Count(blogID string) (int, error) {
	path := fmt.Sprintf("blogs/%s/%s.json", blogID, articleCountBasePath)
	resource := new(ArticleCountResource)
	err := s.client.Get(path, resource, nil)
	return resource.Count, err
}
