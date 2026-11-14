package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestCommentAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "comments.json",
			response: `{"comments":[{"id":7,"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}]}`,
			want:     []Comment{{Id: 7, Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "comments/7.json",
			response: `{"comment":{"id":7,"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}}`,
			want:     &Comment{Id: 7, Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Get(ctx, 7, nil)
			},
		},
		{
			name:     "Count",
			method:   "GET",
			path:     "comments/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Count(ctx, nil)
			},
		},
		{
			name:     "Create",
			method:   "POST",
			path:     "comments.json",
			request:  `{"comment":{"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}}`,
			response: `{"comment":{"id":7,"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}}`,
			want:     &Comment{Id: 7, Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Create(ctx, Comment{Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11})
			},
		},
		{
			name:     "Update",
			method:   "PUT",
			path:     "comments/7.json",
			request:  `{"comment":{"id":7,"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}}`,
			response: `{"comment":{"id":7,"body":"Great post!","author":"Reader","email":"reader@example.com","article_id":12,"blog_id":11}}`,
			want:     &Comment{Id: 7, Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Update(ctx, Comment{Id: 7, Body: "Great post!", Author: "Reader", Email: "reader@example.com", ArticleId: 12, BlogId: 11})
			},
		},
		{
			name:     "Approve",
			method:   "POST",
			path:     "comments/7/approve.json",
			request:  `{}`,
			response: `{"id":7,"status":"published"}`,
			want:     &Comment{Id: 7, Status: "published"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Approve(ctx, 7)
			},
		},
		{
			name:     "NotSpam",
			method:   "POST",
			path:     "comments/7/not_spam.json",
			request:  `{}`,
			response: `{"id":7,"status":"published"}`,
			want:     &Comment{Id: 7, Status: "published"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.NotSpam(ctx, 7)
			},
		},
		{
			name:     "Remove",
			method:   "POST",
			path:     "comments/7/remove.json",
			request:  `{}`,
			response: `{"id":7,"status":"removed"}`,
			want:     &Comment{Id: 7, Status: "removed"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Remove(ctx, 7)
			},
		},
		{
			name:     "Restore",
			method:   "POST",
			path:     "comments/7/restore.json",
			request:  `{}`,
			response: `{"id":7,"status":"published"}`,
			want:     &Comment{Id: 7, Status: "published"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Restore(ctx, 7)
			},
		},
		{
			name:     "Spam",
			method:   "POST",
			path:     "comments/7/spam.json",
			request:  `{}`,
			response: `{"id":7,"status":"spam"}`,
			want:     &Comment{Id: 7, Status: "spam"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Comment.Spam(ctx, 7)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestCommentListWithPagination(t *testing.T) {
	testFeaturePagination(t, "comments.json", `{"comments":[{"id":7}]}`, []Comment{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Comment.ListWithPagination(ctx, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/comment
func TestCommentDocumentedResponse(t *testing.T) {
	var value CommentResource
	testFeatureResponse(t, "fixtures/comment_documented.json", &value, true)
	if !(value.Comment.Id == 757536352 && value.Comment.BodyHtml != "" && value.Comment.PublishedAt == nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/comment
func TestCommentActionDocumentedResponse(t *testing.T) {
	var value Comment
	testFeatureResponse(t, "fixtures/comment_action_documented.json", &value, true)
	if !(value.Id == 653537639 && value.Status == "removed" && value.PublishedAt == nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestCommentQueryFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "comments.json",
		query:    "article_id=7&blog_id=11&status=published&published_at_min=2026-01-15T00%3A00%3A00Z",
		response: `{"comments":[]}`,
		want:     []Comment{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Comment.List(ctx, CommentListOptions{ArticleId: 7, BlogId: 11, Status: "published", PublishedAtMin: at})
		},
	})
}
