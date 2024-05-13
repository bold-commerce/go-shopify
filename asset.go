package goshopify

import (
	"context"
	"fmt"
	"time"
)

const assetsBasePath = "themes"

// AssetService is an interface for interfacing with the asset endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/asset
type AssetService interface {
	List(context.Context, uint64, interface{}) ([]Asset, error)
	Get(context.Context, uint64, string) (*Asset, error)
	Update(context.Context, uint64, Asset) (*Asset, error)
	Delete(context.Context, uint64, string) error
}

// AssetServiceOp handles communication with the asset related methods of
// the Shopify API.
type AssetServiceOp struct {
	client *Client
}

// Asset represents a Shopify asset
type Asset struct {
	Attachment  string     `json:"attachment,omitempty"`
	ContentType string     `json:"content_type,omitempty"`
	Key         string     `json:"key,omitempty"`
	PublicURL   string     `json:"public_url,omitempty"`
	Size        int        `json:"size,omitempty"`
	SourceKey   string     `json:"source_key,omitempty"`
	Src         string     `json:"src,omitempty"`
	ThemeId     uint64     `json:"theme_id,omitempty"`
	Value       string     `json:"value,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// AssetResource is the result from the themes/x/assets.json?asset[key]= endpoint
type AssetResource struct {
	Asset *Asset `json:"asset"`
}

// AssetsResource is the result from the themes/x/assets.json endpoint
type AssetsResource struct {
	Assets []Asset `json:"assets"`
}

type assetGetOptions struct {
	Key     string `url:"asset[key]"`
	ThemeId uint64 `url:"theme_id"`
}

// List the metadata for all assets in the given theme
func (s *AssetServiceOp) List(ctx context.Context, themeId uint64, options interface{}) ([]Asset, error) {
	path := fmt.Sprintf("%s/%d/assets.json", assetsBasePath, themeId)
	resource := new(AssetsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Assets, err
}

// Get an asset by key from the given theme
func (s *AssetServiceOp) Get(ctx context.Context, themeId uint64, key string) (*Asset, error) {
	path := fmt.Sprintf("%s/%d/assets.json", assetsBasePath, themeId)
	options := assetGetOptions{
		Key:     key,
		ThemeId: themeId,
	}
	resource := new(AssetResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Asset, err
}

// Update an asset
func (s *AssetServiceOp) Update(ctx context.Context, themeId uint64, asset Asset) (*Asset, error) {
	path := fmt.Sprintf("%s/%d/assets.json", assetsBasePath, themeId)
	wrappedData := AssetResource{Asset: &asset}
	resource := new(AssetResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Asset, err
}

// Delete an asset
func (s *AssetServiceOp) Delete(ctx context.Context, themeId uint64, key string) error {
	path := fmt.Sprintf("%s/%d/assets.json?asset[key]=%s", assetsBasePath, themeId, key)
	return s.client.Delete(ctx, path)
}
