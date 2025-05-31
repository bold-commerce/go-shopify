package goshopify

import (
	"context"
	"fmt"
	"time"
)

const themesBasePath = "themes"

// Options for theme list
type ThemeListOptions struct {
	Role   string `url:"role,omitempty"`
	Fields string `url:"fields,omitempty"`
}

// ThemeService is an interface for interfacing with the themes endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/theme
type ThemeService interface {
	List(context.Context, interface{}) ([]Theme, error)
	Create(context.Context, Theme) (*Theme, error)
	Get(context.Context, uint64, interface{}) (*Theme, error)
	Update(context.Context, Theme) (*Theme, error)
	Delete(context.Context, uint64) error
}

// ThemeServiceOp handles communication with the theme related methods of
// the Shopify API.
type ThemeServiceOp struct {
	client *Client
}

// Theme represents a Shopify theme
type Theme struct {
	Id                uint64     `json:"id"`
	Name              string     `json:"name"`
	Previewable       bool       `json:"previewable"`
	Processing        bool       `json:"processing"`
	Role              string     `json:"role"`
	ThemeStoreId      uint64     `json:"theme_store_id"`
	AdminGraphqlApiId string     `json:"admin_graphql_api_id"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// ThemesResource is the result from the themes/X.json endpoint
type ThemeResource struct {
	Theme *Theme `json:"theme"`
}

// ThemesResource is the result from the themes.json endpoint
type ThemesResource struct {
	Themes []Theme `json:"themes"`
}

// List all themes
func (s *ThemeServiceOp) List(ctx context.Context, options interface{}) ([]Theme, error) {
	path := fmt.Sprintf("%s.json", themesBasePath)
	resource := new(ThemesResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Themes, err
}

// Update a theme
func (s *ThemeServiceOp) Create(ctx context.Context, theme Theme) (*Theme, error) {
	path := fmt.Sprintf("%s.json", themesBasePath)
	wrappedData := ThemeResource{Theme: &theme}
	resource := new(ThemeResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Theme, err
}

// Get a theme
func (s *ThemeServiceOp) Get(ctx context.Context, themeId uint64, options interface{}) (*Theme, error) {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, themeId)
	resource := new(ThemeResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Theme, err
}

// Update a theme
func (s *ThemeServiceOp) Update(ctx context.Context, theme Theme) (*Theme, error) {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, theme.Id)
	wrappedData := ThemeResource{Theme: &theme}
	resource := new(ThemeResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Theme, err
}

// Delete a theme
func (s *ThemeServiceOp) Delete(ctx context.Context, themeId uint64) error {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, themeId)
	return s.client.Delete(ctx, path)
}
