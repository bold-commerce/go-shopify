package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const inventoryItemsBasePath = "inventory_items"

// InventoryItemService is an interface for interacting with the
// inventory items endpoints of the Shopify API
// See https://shopify.dev/docs/api/admin-rest/latest/resources/inventoryitem
type InventoryItemService interface {
	List(context.Context, interface{}) ([]InventoryItem, error)
	Get(context.Context, uint64, interface{}) (*InventoryItem, error)
	Update(context.Context, InventoryItem) (*InventoryItem, error)
}

// InventoryItemServiceOp is the default implementation of the InventoryItemService interface
type InventoryItemServiceOp struct {
	client *Client
}

// CountryHarmonizedSystemCode represents a harmonized system code for a specific country
type CountryHarmonizedSystemCode struct {
	HarmonizedSystemCode string `json:"harmonized_system_code"`
	CountryCode          string `json:"country_code"`
}

// InventoryItem represents a Shopify inventory item
type InventoryItem struct {
	Id                           uint64                        `json:"id,omitempty"`
	SKU                          string                        `json:"sku,omitempty"`
	CreatedAt                    *time.Time                    `json:"created_at,omitempty"`
	UpdatedAt                    *time.Time                    `json:"updated_at,omitempty"`
	Cost                         *decimal.Decimal              `json:"cost,omitempty"`
	Tracked                      *bool                         `json:"tracked,omitempty"`
	AdminGraphqlApiId            string                        `json:"admin_graphql_api_id,omitempty"`
	CountryCodeOfOrigin          *string                       `json:"country_code_of_origin"`
	CountryHarmonizedSystemCodes []CountryHarmonizedSystemCode `json:"country_harmonized_system_codes"`
	HarmonizedSystemCode         *string                       `json:"harmonized_system_code"`
	ProvinceCodeOfOrigin         *string                       `json:"province_code_of_origin"`
}

// InventoryItemResource is used for handling single item requests and responses
type InventoryItemResource struct {
	InventoryItem *InventoryItem `json:"inventory_item"`
}

// InventoryItemsResource is used for handling multiple item responsees
type InventoryItemsResource struct {
	InventoryItems []InventoryItem `json:"inventory_items"`
}

// List inventory items
func (s *InventoryItemServiceOp) List(ctx context.Context, options interface{}) ([]InventoryItem, error) {
	path := fmt.Sprintf("%s.json", inventoryItemsBasePath)
	resource := new(InventoryItemsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.InventoryItems, err
}

// Get a inventory item
func (s *InventoryItemServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*InventoryItem, error) {
	path := fmt.Sprintf("%s/%d.json", inventoryItemsBasePath, id)
	resource := new(InventoryItemResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.InventoryItem, err
}

// Update a inventory item
func (s *InventoryItemServiceOp) Update(ctx context.Context, item InventoryItem) (*InventoryItem, error) {
	path := fmt.Sprintf("%s/%d.json", inventoryItemsBasePath, item.Id)
	wrappedData := InventoryItemResource{InventoryItem: &item}
	resource := new(InventoryItemResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.InventoryItem, err
}
