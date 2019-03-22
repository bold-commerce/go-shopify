package goshopify

import (
	"fmt"
	"time"
)

const themesBasePath = "admin/themes"

// Options for theme list
type ThemeListOptions struct {
	ListOptions
	Role string `url:"role,omitempty"`
}

// ThemeService is an interface for interfacing with the themes endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/theme
type ThemeService interface {
	List(interface{}) ([]Theme, error)
}

// ThemeServiceOp handles communication with the theme related methods of
// the Shopify API.
type ThemeServiceOp struct {
	client *Client
}

// Theme represents a Shopify theme
type Theme struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Previewable  bool       `json:"previewable"`
	Processing   bool       `json:"processing"`
	Role         string     `json:"role"`
	ThemeStoreID int64      `json:"theme_store_id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
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
