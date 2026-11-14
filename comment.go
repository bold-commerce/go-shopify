package goshopify

import (
	"context"
	"fmt"
	"time"
)

// CommentService provides access to Comment endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/comment
type CommentService interface {
	ListWithPagination(ctx context.Context, options interface{}) ([]Comment, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]Comment, error)
	Get(ctx context.Context, id uint64, options interface{}) (*Comment, error)
	Count(ctx context.Context, options interface{}) (int, error)
	Create(ctx context.Context, comment Comment) (*Comment, error)
	Update(ctx context.Context, comment Comment) (*Comment, error)
	Approve(ctx context.Context, id uint64) (*Comment, error)
	NotSpam(ctx context.Context, id uint64) (*Comment, error)
	Remove(ctx context.Context, id uint64) (*Comment, error)
	Restore(ctx context.Context, id uint64) (*Comment, error)
	Spam(ctx context.Context, id uint64) (*Comment, error)
}

// CommentServiceOp implements CommentService.
type CommentServiceOp struct {
	client *Client
}

// Comment represents Shopify comment.
type Comment struct {
	Id          uint64     `json:"id,omitempty"`
	ArticleId   uint64     `json:"article_id,omitempty"`
	BlogId      uint64     `json:"blog_id,omitempty"`
	Author      string     `json:"author,omitempty"`
	Body        string     `json:"body,omitempty"`
	BodyHtml    string     `json:"body_html,omitempty"`
	Email       string     `json:"email,omitempty"`
	IP          string     `json:"ip,omitempty"`
	Status      string     `json:"status,omitempty"`
	UserAgent   string     `json:"user_agent,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// CommentsResource represents the comments response.
type CommentsResource struct {
	Comments []Comment `json:"comments"`
}

// CommentResource represents the comment response.
type CommentResource struct {
	Comment *Comment `json:"comment"`
}

// CommentListOptions controls filtering and pagination.
type CommentListOptions struct {
	ListOptions
	ArticleId      uint64    `url:"article_id,omitempty"`
	BlogId         uint64    `url:"blog_id,omitempty"`
	Status         string    `url:"status,omitempty"`
	PublishedAtMin time.Time `url:"published_at_min,omitempty"`
	PublishedAtMax time.Time `url:"published_at_max,omitempty"`
}

// List retrieves comment records.
func (s *CommentServiceOp) List(ctx context.Context, options interface{}) ([]Comment, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a comment record.
func (s *CommentServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*Comment, error) {
	resource := new(CommentResource)
	err := s.client.Get(ctx, fmt.Sprintf("comments/%d.json", id), resource, options)
	return resource.Comment, err
}

// Count counts comment records.
func (s *CommentServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	return s.client.Count(ctx, "comments/count.json", options)
}

// Create creates a comment record.
func (s *CommentServiceOp) Create(ctx context.Context, comment Comment) (*Comment, error) {
	resource := new(CommentResource)
	err := s.client.Post(ctx, "comments.json", CommentResource{Comment: &comment}, resource)
	return resource.Comment, err
}

// Update updates a comment record.
func (s *CommentServiceOp) Update(ctx context.Context, comment Comment) (*Comment, error) {
	resource := new(CommentResource)
	err := s.client.Put(ctx, fmt.Sprintf("comments/%d.json", comment.Id), CommentResource{Comment: &comment}, resource)
	return resource.Comment, err
}

// Approve approves and publishes a comment.
func (s *CommentServiceOp) Approve(ctx context.Context, id uint64) (*Comment, error) {
	resource := new(Comment)
	err := s.client.Post(ctx, fmt.Sprintf("comments/%d/approve.json", id), struct{}{}, resource)
	return resource, err
}

// NotSpam marks a comment as not spam.
func (s *CommentServiceOp) NotSpam(ctx context.Context, id uint64) (*Comment, error) {
	resource := new(Comment)
	err := s.client.Post(ctx, fmt.Sprintf("comments/%d/not_spam.json", id), struct{}{}, resource)
	return resource, err
}

// Remove removes a comment.
func (s *CommentServiceOp) Remove(ctx context.Context, id uint64) (*Comment, error) {
	resource := new(Comment)
	err := s.client.Post(ctx, fmt.Sprintf("comments/%d/remove.json", id), struct{}{}, resource)
	return resource, err
}

// Restore restores a removed comment.
func (s *CommentServiceOp) Restore(ctx context.Context, id uint64) (*Comment, error) {
	resource := new(Comment)
	err := s.client.Post(ctx, fmt.Sprintf("comments/%d/restore.json", id), struct{}{}, resource)
	return resource, err
}

// Spam marks a comment as spam.
func (s *CommentServiceOp) Spam(ctx context.Context, id uint64) (*Comment, error) {
	resource := new(Comment)
	err := s.client.Post(ctx, fmt.Sprintf("comments/%d/spam.json", id), struct{}{}, resource)
	return resource, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *CommentServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Comment, *Pagination, error) {
	resource := new(CommentsResource)
	pagination, err := s.client.ListWithPagination(ctx, "comments.json", resource, options)
	return resource.Comments, pagination, err
}
