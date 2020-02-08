package goshopify

import (
	"fmt"
	"time"
)

const scriptTagsBasePath = "script_tags"

// ScriptTagService is an interface for interfacing with the ScriptTag endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/scripttag
type ScriptTagService interface {
	List(interface{}) ([]ScriptTag, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*ScriptTag, error)
	Create(ScriptTag) (*ScriptTag, error)
	Update(ScriptTag) (*ScriptTag, error)
	Delete(int64) error
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
	ID           int64      `json:"id"`
	Src          string     `json:"src"`
	DisplayScope string     `json:"display_scope"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

// The options provided by Shopify.
type ScriptTagOption struct {
	Limit        int       `url:"limit,omitempty"`
	Page         int       `url:"page,omitempty"`
	SinceID      int64     `url:"since_id,omitempty"`
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
func (s *ScriptTagServiceOp) List(options interface{}) ([]ScriptTag, error) {
	path := fmt.Sprintf("%s.json", scriptTagsBasePath)
	resource := &ScriptTagsResource{}
	err := s.client.Get(path, resource, options)
	return resource.ScriptTags, err
}

// Count script tags
func (s *ScriptTagServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", scriptTagsBasePath)
	return s.client.Count(path, options)
}

// Get individual script tag
func (s *ScriptTagServiceOp) Get(tagID int64, options interface{}) (*ScriptTag, error) {
	path := fmt.Sprintf("%s/%d.json", scriptTagsBasePath, tagID)
	resource := &ScriptTagResource{}
	err := s.client.Get(path, resource, options)
	return resource.ScriptTag, err
}

// Create a new script tag
func (s *ScriptTagServiceOp) Create(tag ScriptTag) (*ScriptTag, error) {
	path := fmt.Sprintf("%s.json", scriptTagsBasePath)
	wrappedData := ScriptTagResource{ScriptTag: &tag}
	resource := &ScriptTagResource{}
	err := s.client.Post(path, wrappedData, resource)
	return resource.ScriptTag, err
}

// Update an existing script tag
func (s *ScriptTagServiceOp) Update(tag ScriptTag) (*ScriptTag, error) {
	path := fmt.Sprintf("%s/%d.json", scriptTagsBasePath, tag.ID)
	wrappedData := ScriptTagResource{ScriptTag: &tag}
	resource := &ScriptTagResource{}
	err := s.client.Put(path, wrappedData, resource)
	return resource.ScriptTag, err
}

// Delete an existing script tag
func (s *ScriptTagServiceOp) Delete(tagID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", scriptTagsBasePath, tagID))
}
