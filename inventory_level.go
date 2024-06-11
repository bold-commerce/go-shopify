package goshopify

import (
	"context"
	"fmt"
	"time"
)

const inventoryLevelsBasePath = "inventory_levels"

// InventoryLevelService is an interface for interacting with the
// inventory items endpoints of the Shopify API
// See https://help.shopify.com/en/api/reference/inventory/inventorylevel
type InventoryLevelService interface {
	List(context.Context, interface{}) ([]InventoryLevel, error)
	Adjust(context.Context, interface{}) (*InventoryLevel, error)
	Delete(context.Context, uint64, uint64) error
	Connect(context.Context, InventoryLevel) (*InventoryLevel, error)
	Set(context.Context, InventoryLevel) (*InventoryLevel, error)
}

// InventoryLevelServiceOp is the default implementation of the InventoryLevelService interface
type InventoryLevelServiceOp struct {
	client *Client
}

// InventoryLevel represents a Shopify inventory level
type InventoryLevel struct {
	InventoryItemId   uint64     `json:"inventory_item_id,omitempty"`
	LocationId        uint64     `json:"location_id,omitempty"`
	Available         int        `json:"available"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	AdminGraphqlApiId string     `json:"admin_graphql_api_id,omitempty"`
}

// InventoryLevelResource is used for handling single level requests and responses
type InventoryLevelResource struct {
	InventoryLevel *InventoryLevel `json:"inventory_level"`
}

// InventoryLevelsResource is used for handling multiple item responsees
type InventoryLevelsResource struct {
	InventoryLevels []InventoryLevel `json:"inventory_levels"`
}

// InventoryLevelListOptions is used for get list
type InventoryLevelListOptions struct {
	InventoryItemIds []uint64  `url:"inventory_item_ids,omitempty,comma"`
	LocationIds      []uint64  `url:"location_ids,omitempty,comma"`
	Limit            int       `url:"limit,omitempty"`
	UpdatedAtMin     time.Time `url:"updated_at_min,omitempty"`
}

// InventoryLevelAdjustOptions is used for Adjust inventory levels
type InventoryLevelAdjustOptions struct {
	InventoryItemId uint64 `json:"inventory_item_id"`
	LocationId      uint64 `json:"location_id"`
	Adjust          int    `json:"available_adjustment"`
}

// List inventory levels
func (s *InventoryLevelServiceOp) List(ctx context.Context, options interface{}) ([]InventoryLevel, error) {
	path := fmt.Sprintf("%s.json", inventoryLevelsBasePath)
	resource := new(InventoryLevelsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.InventoryLevels, err
}

// Delete an inventory level
func (s *InventoryLevelServiceOp) Delete(ctx context.Context, itemId, locationId uint64) error {
	path := fmt.Sprintf("%s.json?inventory_item_id=%v&location_id=%v",
		inventoryLevelsBasePath, itemId, locationId)
	return s.client.Delete(ctx, path)
}

// Connect an inventory level
func (s *InventoryLevelServiceOp) Connect(ctx context.Context, level InventoryLevel) (*InventoryLevel, error) {
	return s.post(ctx, fmt.Sprintf("%s/connect.json", inventoryLevelsBasePath), level)
}

// Set an inventory level
func (s *InventoryLevelServiceOp) Set(ctx context.Context, level InventoryLevel) (*InventoryLevel, error) {
	return s.post(ctx, fmt.Sprintf("%s/set.json", inventoryLevelsBasePath), level)
}

// Adjust the inventory level of an inventory item at a single location
func (s *InventoryLevelServiceOp) Adjust(ctx context.Context, options interface{}) (*InventoryLevel, error) {
	return s.post(ctx, fmt.Sprintf("%s/adjust.json", inventoryLevelsBasePath), options)
}

func (s *InventoryLevelServiceOp) post(ctx context.Context, path string, options interface{}) (*InventoryLevel, error) {
	resource := new(InventoryLevelResource)
	err := s.client.Post(ctx, path, options, resource)
	return resource.InventoryLevel, err
}
