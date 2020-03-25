package goshopify

import (
	"fmt"
	"time"
)

const pagesBasePath = "pages"
const pagesResourceName = "pages"

// PagesPageService is an interface for interacting with the pages
// endpoints of the Shopify API.
// See https://help.shopify.com/api/reference/online_store/page
type PageService interface {
	List(interface{}) ([]Page, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Page, error)
	Create(Page) (*Page, error)
	Update(Page) (*Page, error)
	Delete(int64) error

	// MetafieldsService used for Pages resource to communicate with Metafields
	// resource
	MetafieldsService
}

// PageServiceOp handles communication with the page related methods of the
// Shopify API.
type PageServiceOp struct {
	client *Client
}

// Page represents a Shopify page.
type Page struct {
	ID             int64       `json:"id,omitempty"`
	Author         string      `json:"author,omitempty"`
	Handle         string      `json:"handle,omitempty"`
	Title          string      `json:"title,omitempty"`
	CreatedAt      *time.Time  `json:"created_at,omitempty"`
	UpdatedAt      *time.Time  `json:"updated_at,omitempty"`
	BodyHTML       string      `json:"body_html,omitempty"`
	TemplateSuffix string      `json:"template_suffix,omitempty"`
	PublishedAt    *time.Time  `json:"published_at,omitempty"`
	ShopID         int64       `json:"shop_id,omitempty"`
	Metafields     []Metafield `json:"metafields,omitempty"`
}

// PageResource represents the result from the pages/X.json endpoint
type PageResource struct {
	Page *Page `json:"page"`
}

// PagesResource represents the result from the pages.json endpoint
type PagesResource struct {
	Pages []Page `json:"pages"`
}

// List pages
func (s *PageServiceOp) List(options interface{}) ([]Page, error) {
	path := fmt.Sprintf("%s.json", pagesBasePath)
	resource := new(PagesResource)
	err := s.client.Get(path, resource, options)
	return resource.Pages, err
}

// Count pages
func (s *PageServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", pagesBasePath)
	return s.client.Count(path, options)
}

// Get individual page
func (s *PageServiceOp) Get(pageID int64, options interface{}) (*Page, error) {
	path := fmt.Sprintf("%s/%d.json", pagesBasePath, pageID)
	resource := new(PageResource)
	err := s.client.Get(path, resource, options)
	return resource.Page, err
}

// Create a new page
func (s *PageServiceOp) Create(page Page) (*Page, error) {
	path := fmt.Sprintf("%s.json", pagesBasePath)
	wrappedData := PageResource{Page: &page}
	resource := new(PageResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Page, err
}

// Update an existing page
func (s *PageServiceOp) Update(page Page) (*Page, error) {
	path := fmt.Sprintf("%s/%d.json", pagesBasePath, page.ID)
	wrappedData := PageResource{Page: &page}
	resource := new(PageResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Page, err
}

// Delete an existing page.
func (s *PageServiceOp) Delete(pageID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", pagesBasePath, pageID))
}

// List metafields for a page
func (s *PageServiceOp) ListMetafields(pageID int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.List(options)
}

// Count metafields for a page
func (s *PageServiceOp) CountMetafields(pageID int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.Count(options)
}

// Get individual metafield for a page
func (s *PageServiceOp) GetMetafield(pageID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.Get(metafieldID, options)
}

// Create a new metafield for a page
func (s *PageServiceOp) CreateMetafield(pageID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.Create(metafield)
}

// Update an existing metafield for a page
func (s *PageServiceOp) UpdateMetafield(pageID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.Update(metafield)
}

// Delete an existing metafield for a page
func (s *PageServiceOp) DeleteMetafield(pageID int64, metafieldID int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: pagesResourceName, resourceID: pageID}
	return metafieldService.Delete(metafieldID)
}
