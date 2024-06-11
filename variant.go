package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	variantsBasePath     = "variants"
	variantsResourceName = "variants"
)

// VariantService is an interface for interacting with the variant endpoints
// of the Shopify API.
// See https://help.shopify.com/api/reference/product_variant
type VariantService interface {
	List(context.Context, uint64, interface{}) ([]Variant, error)
	Count(context.Context, uint64, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Variant, error)
	Create(context.Context, uint64, Variant) (*Variant, error)
	Update(context.Context, Variant) (*Variant, error)
	Delete(context.Context, uint64, uint64) error

	// MetafieldsService used for Variant resource to communicate with Metafields resource
	MetafieldsService
}

// VariantServiceOp handles communication with the variant related methods of
// the Shopify API.
type VariantServiceOp struct {
	client *Client
}

type variantInventoryPolicy string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/product-variant#resource-object
const (
	// Customers are not allowed to place orders for the product variant if it's out
	// of stock. This is the default value.
	VariantInventoryPolicyDeny variantInventoryPolicy = "deny"

	// Customers are allowed to place orders for the product variant if it's out of
	// stock.
	VariantInventoryPolicyContinue variantInventoryPolicy = "continue"
)

// Variant represents a Shopify variant
type Variant struct {
	Id                   uint64                 `json:"id,omitempty"`
	ProductId            uint64                 `json:"product_id,omitempty"`
	Title                string                 `json:"title,omitempty"`
	Sku                  string                 `json:"sku,omitempty"`
	Position             int                    `json:"position,omitempty"`
	Grams                int                    `json:"grams,omitempty"`
	InventoryPolicy      variantInventoryPolicy `json:"inventory_policy,omitempty"`
	Price                *decimal.Decimal       `json:"price,omitempty"`
	CompareAtPrice       *decimal.Decimal       `json:"compare_at_price,omitempty"`
	FulfillmentService   string                 `json:"fulfillment_service,omitempty"`
	InventoryManagement  string                 `json:"inventory_management,omitempty"`
	InventoryItemId      uint64                 `json:"inventory_item_id,omitempty"`
	Option1              string                 `json:"option1,omitempty"`
	Option2              string                 `json:"option2,omitempty"`
	Option3              string                 `json:"option3,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	Taxable              bool                   `json:"taxable,omitempty"`
	TaxCode              string                 `json:"tax_code,omitempty"`
	Barcode              string                 `json:"barcode,omitempty"`
	ImageId              uint64                 `json:"image_id,omitempty"`
	InventoryQuantity    int                    `json:"inventory_quantity,omitempty"`
	Weight               *decimal.Decimal       `json:"weight,omitempty"`
	WeightUnit           string                 `json:"weight_unit,omitempty"`
	OldInventoryQuantity int                    `json:"old_inventory_quantity,omitempty"`
	RequireShipping      bool                   `json:"requires_shipping"`
	AdminGraphqlApiId    string                 `json:"admin_graphql_api_id,omitempty"`
	Metafields           []Metafield            `json:"metafields,omitempty"`
	PresentmentPrices    []presentmentPrices    `json:"presentment_prices,omitempty"`
}

type presentmentPrices struct {
	Price          *AmountSetEntry `json:"price,omitempty"`
	CompareAtPrice *AmountSetEntry `json:"compare_at_price,omitempty"`
}

// VariantResource represents the result from the variants/X.json endpoint
type VariantResource struct {
	Variant *Variant `json:"variant"`
}

// VariantsResource represents the result from the products/X/variants.json endpoint
type VariantsResource struct {
	Variants []Variant `json:"variants"`
}

// List variants
func (s *VariantServiceOp) List(ctx context.Context, productId uint64, options interface{}) ([]Variant, error) {
	path := fmt.Sprintf("%s/%d/variants.json", productsBasePath, productId)
	resource := new(VariantsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Variants, err
}

// Count variants
func (s *VariantServiceOp) Count(ctx context.Context, productId uint64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/variants/count.json", productsBasePath, productId)
	return s.client.Count(ctx, path, options)
}

// Get individual variant
func (s *VariantServiceOp) Get(ctx context.Context, variantId uint64, options interface{}) (*Variant, error) {
	path := fmt.Sprintf("%s/%d.json", variantsBasePath, variantId)
	resource := new(VariantResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Variant, err
}

// Create a new variant
func (s *VariantServiceOp) Create(ctx context.Context, productId uint64, variant Variant) (*Variant, error) {
	path := fmt.Sprintf("%s/%d/variants.json", productsBasePath, productId)
	wrappedData := VariantResource{Variant: &variant}
	resource := new(VariantResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Variant, err
}

// Update existing variant
func (s *VariantServiceOp) Update(ctx context.Context, variant Variant) (*Variant, error) {
	path := fmt.Sprintf("%s/%d.json", variantsBasePath, variant.Id)
	wrappedData := VariantResource{Variant: &variant}
	resource := new(VariantResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Variant, err
}

// Delete an existing variant
func (s *VariantServiceOp) Delete(ctx context.Context, productId uint64, variantId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d/variants/%d.json", productsBasePath, productId, variantId))
}

// ListMetafields for a variant
func (s *VariantServiceOp) ListMetafields(ctx context.Context, variantId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.List(ctx, options)
}

// CountMetafields for a variant
func (s *VariantServiceOp) CountMetafields(ctx context.Context, variantId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.Count(ctx, options)
}

// GetMetafield for a variant
func (s *VariantServiceOp) GetMetafield(ctx context.Context, variantId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// CreateMetafield for a variant
func (s *VariantServiceOp) CreateMetafield(ctx context.Context, variantId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.Create(ctx, metafield)
}

// UpdateMetafield for a variant
func (s *VariantServiceOp) UpdateMetafield(ctx context.Context, variantId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.Update(ctx, metafield)
}

// DeleteMetafield for a variant
func (s *VariantServiceOp) DeleteMetafield(ctx context.Context, variantId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantsResourceName, resourceId: variantId}
	return metafieldService.Delete(ctx, metafieldId)
}
