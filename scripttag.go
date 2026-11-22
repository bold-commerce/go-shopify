package goshopify

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const scriptTagsBasePath = "script_tags"

// ErrScriptTagWriteUnsupported is returned when creating or updating a script tag.
var ErrScriptTagWriteUnsupported = errors.New("Shopify no longer accepts script tag writes. Use a theme app embed block")

// ScriptTagService is an interface for interfacing with the ScriptTag endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/scripttag
type ScriptTagService interface {
	List(context.Context, interface{}) ([]ScriptTag, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]ScriptTag, *Pagination, error)
	ListAll(ctx context.Context, options interface{}) ([]ScriptTag, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*ScriptTag, error)
	// Deprecated: Shopify disabled script tag writes on October 1, 2026. Use an app embed block.
	Create(context.Context, ScriptTag) (*ScriptTag, error)
	// Deprecated: Shopify disabled script tag writes on October 1, 2026. Use an app embed block.
	Update(context.Context, ScriptTag) (*ScriptTag, error)
	Delete(context.Context, uint64) error
}

// ScriptTagServiceOp handles communication with the shop related methods of the
// Shopify API.
type ScriptTagServiceOp struct {
	client *Client
}

// ScriptTag represents a Shopify ScriptTag.
type ScriptTag struct {
	CreatedAt    *time.Time `json:"created_at"`
	Event        string     `json:"event"`
	Id           uint64     `json:"id"`
	Src          string     `json:"src"`
	DisplayScope string     `json:"display_scope"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Cache        *bool      `json:"cache,omitempty"`
}

// The options provided by Shopify.
type ScriptTagOption struct {
	Limit        int       `url:"limit,omitempty"`
	Page         int       `url:"page,omitempty"`
	SinceId      uint64    `url:"since_id,omitempty"`
	CreatedAtMin time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax time.Time `url:"updated_at_max,omitempty"`
	Src          string    `url:"src,omitempty"`
	Fields       string    `url:"fields,omitempty"`
}

// ScriptTagsResource represents the result from the admin/script_tags.json
// endpoint.
type ScriptTagsResource struct {
	ScriptTags []ScriptTag `json:"script_tags"`
}

// ScriptTagResource represents the result from the
// admin/script_tags/{#script_tag_id}.json endpoint.
type ScriptTagResource struct {
	ScriptTag *ScriptTag `json:"script_tag"`
}

// List script tags
func (s *ScriptTagServiceOp) List(ctx context.Context, options interface{}) ([]ScriptTag, error) {
	path := fmt.Sprintf("%s.json", scriptTagsBasePath)
	resource := &ScriptTagsResource{}
	err := s.client.Get(ctx, path, resource, options)
	return resource.ScriptTags, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *ScriptTagServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]ScriptTag, *Pagination, error) {
	path := fmt.Sprintf("%s.json", scriptTagsBasePath)
	resource := &ScriptTagsResource{}
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.ScriptTags, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *ScriptTagServiceOp) ListAll(ctx context.Context, options interface{}) ([]ScriptTag, error) {
	items := []ScriptTag{}
	err := forEachPage(options, func(options interface{}) (*Pagination, error) {
		page, pagination, err := s.ListWithPagination(ctx, options)
		if err == nil {
			items = append(items, page...)
		}
		return pagination, err
	})
	return items, err
}

// Count script tags
func (s *ScriptTagServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", scriptTagsBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual script tag
func (s *ScriptTagServiceOp) Get(ctx context.Context, tagId uint64, options interface{}) (*ScriptTag, error) {
	path := fmt.Sprintf("%s/%d.json", scriptTagsBasePath, tagId)
	resource := &ScriptTagResource{}
	err := s.client.Get(ctx, path, resource, options)
	return resource.ScriptTag, err
}

// Create returns ErrScriptTagWriteUnsupported.
//
// Deprecated: Shopify disabled script tag writes on October 1, 2026. Use an app embed block.
func (s *ScriptTagServiceOp) Create(ctx context.Context, tag ScriptTag) (*ScriptTag, error) {
	return nil, ErrScriptTagWriteUnsupported
}

// Update returns ErrScriptTagWriteUnsupported.
//
// Deprecated: Shopify disabled script tag writes on October 1, 2026. Use an app embed block.
func (s *ScriptTagServiceOp) Update(ctx context.Context, tag ScriptTag) (*ScriptTag, error) {
	return nil, ErrScriptTagWriteUnsupported
}

// Delete an existing script tag
func (s *ScriptTagServiceOp) Delete(ctx context.Context, tagId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", scriptTagsBasePath, tagId))
}
