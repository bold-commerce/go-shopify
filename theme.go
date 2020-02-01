package goshopify

import (
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
// See: https://help.shopify.com/api/reference/theme
type ThemeService interface {
	List(interface{}) ([]Theme, error)
	Create(Theme) (*Theme, error)
	Get(int64, interface{}) (*Theme, error)
	Update(Theme) (*Theme, error)
	Delete(int64) error
}

// ThemeServiceOp handles communication with the theme related methods of
// the Shopify API.
type ThemeServiceOp struct {
	client *Client
}

// Theme represents a Shopify theme
type Theme struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Previewable       bool       `json:"previewable"`
	Processing        bool       `json:"processing"`
	Role              string     `json:"role"`
	ThemeStoreID      int64      `json:"theme_store_id"`
	AdminGraphQLApiID string     `json:"admin_graphql_api_id"`
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
func (s *ThemeServiceOp) List(options interface{}) ([]Theme, error) {
	path := fmt.Sprintf("%s.json", themesBasePath)
	resource := new(ThemesResource)
	err := s.client.Get(path, resource, options)
	return resource.Themes, err
}

// Update a theme
func (s *ThemeServiceOp) Create(theme Theme) (*Theme, error) {
	path := fmt.Sprintf("%s.json", themesBasePath)
	wrappedData := ThemeResource{Theme: &theme}
	resource := new(ThemeResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Theme, err
}

// Get a theme
func (s *ThemeServiceOp) Get(themeID int64, options interface{}) (*Theme, error) {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, themeID)
	resource := new(ThemeResource)
	err := s.client.Get(path, resource, options)
	return resource.Theme, err
}

// Update a theme
func (s *ThemeServiceOp) Update(theme Theme) (*Theme, error) {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, theme.ID)
	wrappedData := ThemeResource{Theme: &theme}
	resource := new(ThemeResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Theme, err
}

// Delete a theme
func (s *ThemeServiceOp) Delete(themeID int64) error {
	path := fmt.Sprintf("%s/%d.json", themesBasePath, themeID)
	return s.client.Delete(path)
}
