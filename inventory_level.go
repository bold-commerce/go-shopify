package goshopify

import (
	"fmt"
	"time"
)

const inventoryLevelBasePath = "inventory_levels"
const inventoryLevelAdjustmentBasePath = "adjust"
const inventoryLevelConnectBasePath = "connect"
const inventoryLevelSetBasePath = "set"

// InventoryLevelService is an interface for interacting with the
// inventory level endpoints of the Shopify API
// See https://help.shopify.com/en/api/reference/inventory/inventorylevel
type InventoryLevelService interface {
	Get(interface{}) ([]InventoryLevel, error)
	Adjust(OptionsInventoryLevel) (InventoryLevel, error)
	Delete(interface{}) error
	Connect(OptionsInventoryLevel) (InventoryLevel, error)
	Set(OptionsInventoryLevel) (InventoryLevel, error)
}

// InventoryLevelServiceOp is the default implementation of the InventoryLevelService interface
type InventoryLevelServiceOp struct {
	client *Client
}

// InventoryLevel represents a Shopify inventory Level
type InventoryLevel struct {
	Available         int64      `json:"available,omitempty"`
	InventoryItemID   int64      `json:"inventory_item_id,omitempty"`
	LocationID        int64      `json:"location_id,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	AdminGraphqlAPIID string     `json:"admin_graphql_api_id,omitempty"`
}

// OptionsInventoryLevel represents a adjusts requirement
type OptionsInventoryLevel struct {
	InventoryItemID       int64 `json:"inventory_item_id,omitempty"`
	LocationID            int64 `json:"location_id,omitempty"`
	AvailableAdjustment   int64 `json:"available_adjustment,omitempty"`
	Available             int64 `json:"available,omitempty"`
	DisconnectIfNecessary bool  `json:"disconnect_if_necessary,omitempty"`
}

// InventoryLevelResource is used for handling single item requests and responses
type InventoryLevelResource struct {
	InventoryLevel InventoryLevel `json:"inventory_level"`
}

// InventoryLevelsResource is used for handling multiple level responsees
type InventoryLevelsResource struct {
	InventoryLevels []InventoryLevel `json:"inventory_levels"`
}

// Get inventory leves by inventoryItemId and / or location_id
func (s *InventoryLevelServiceOp) Get(options interface{}) ([]InventoryLevel, error) {
	path := fmt.Sprintf("%s/%s.json", globalApiPathPrefix, inventoryLevelBasePath)
	resource := new(InventoryLevelsResource)
	err := s.client.Get(path, resource, options)
	return resource.InventoryLevels, err
}

// Adjust the inventory level of an inventory item at a single location
// Parameters required from OptionsInventoryLevel {InventoryItemID, LocationID, AvailableAdjustment}
func (s *InventoryLevelServiceOp) Adjust(o OptionsInventoryLevel) (InventoryLevel, error) {
	path := fmt.Sprintf("%s/%s/%s.json", globalApiPathPrefix, inventoryLevelBasePath, inventoryLevelAdjustmentBasePath)
	resource := new(InventoryLevelResource)
	fmt.Println(path)
	fmt.Println(o)
	err := s.client.Post(path, o, resource)
	return resource.InventoryLevel, err
}

// Delete an inventory level of an inventory item at a location
// options interface LIKE
// optionsDelete := struct {
//	 InventoryItemID int64 `url:"inventory_item_id,omitempty"`
//	 LocationID      int64 `url:"location_id,omitempty"`
// }{InventoryItemID: xxxxxxxxxxxxxx, LocationID: xxxxxxxxxxx}
func (s *InventoryLevelServiceOp) Delete(options interface{}) error {
	return s.client.DeleteWithOptions(fmt.Sprintf("%s/%s.json", globalApiPathPrefix, inventoryLevelBasePath), options)
}

// Connect an inventory item to a location by creating an inventory level at that location.
func (s *InventoryLevelServiceOp) Connect(o OptionsInventoryLevel) (InventoryLevel, error) {
	path := fmt.Sprintf("%s/%s/%s.json", globalApiPathPrefix, inventoryLevelBasePath, inventoryLevelConnectBasePath)
	resource := new(InventoryLevelResource)
	err := s.client.Post(path, o, resource)
	fmt.Println(path)
	return resource.InventoryLevel, err
}

// Set the inventory level for an inventory item at a location
func (s *InventoryLevelServiceOp) Set(o OptionsInventoryLevel) (InventoryLevel, error) {
	path := fmt.Sprintf("%s/%s/%s.json", globalApiPathPrefix, inventoryLevelBasePath, inventoryLevelSetBasePath)
	resource := new(InventoryLevelResource)
	err := s.client.Post(path, o, resource)
	fmt.Println(path)
	return resource.InventoryLevel, err
}
