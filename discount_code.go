package goshopify

import (
	"context"
	"fmt"
	"time"
)

const discountCodeBasePath = "price_rules/%d/discount_codes"

// DiscountCodeService is an interface for interfacing with the discount endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/PriceRuleDiscountCode
type DiscountCodeService interface {
	Create(context.Context, uint64, PriceRuleDiscountCode) (*PriceRuleDiscountCode, error)
	Update(context.Context, uint64, PriceRuleDiscountCode) (*PriceRuleDiscountCode, error)
	List(context.Context, uint64) ([]PriceRuleDiscountCode, error)
	ListWithPagination(ctx context.Context, priceRuleId uint64, options interface{}) ([]PriceRuleDiscountCode, *Pagination, error)
	ListAll(ctx context.Context, priceRuleId uint64, options interface{}) ([]PriceRuleDiscountCode, error)
	Get(context.Context, uint64, uint64) (*PriceRuleDiscountCode, error)
	Delete(context.Context, uint64, uint64) error
}

// DiscountCodeServiceOp handles communication with the discount code
// related methods of the Shopify API.
type DiscountCodeServiceOp struct {
	client *Client
}

// PriceRuleDiscountCode represents a Shopify Discount Code
type PriceRuleDiscountCode struct {
	Id          uint64     `json:"id,omitempty"`
	PriceRuleId uint64     `json:"price_rule_id,omitempty"`
	Code        string     `json:"code,omitempty"`
	UsageCount  int        `json:"usage_count,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// DiscountCodesResource is the result from the discount_codes.json endpoint
type DiscountCodesResource struct {
	DiscountCodes []PriceRuleDiscountCode `json:"discount_codes"`
}

// DiscountCodeResource represents the result from the discount_codes/X.json endpoint
type DiscountCodeResource struct {
	PriceRuleDiscountCode *PriceRuleDiscountCode `json:"discount_code"`
}

// Create a discount code
func (s *DiscountCodeServiceOp) Create(ctx context.Context, priceRuleId uint64, dc PriceRuleDiscountCode) (*PriceRuleDiscountCode, error) {
	path := fmt.Sprintf(discountCodeBasePath+".json", priceRuleId)
	wrappedData := DiscountCodeResource{PriceRuleDiscountCode: &dc}
	resource := new(DiscountCodeResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.PriceRuleDiscountCode, err
}

// Update an existing discount code
func (s *DiscountCodeServiceOp) Update(ctx context.Context, priceRuleId uint64, dc PriceRuleDiscountCode) (*PriceRuleDiscountCode, error) {
	path := fmt.Sprintf(discountCodeBasePath+"/%d.json", priceRuleId, dc.Id)
	wrappedData := DiscountCodeResource{PriceRuleDiscountCode: &dc}
	resource := new(DiscountCodeResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.PriceRuleDiscountCode, err
}

// List of discount codes
func (s *DiscountCodeServiceOp) List(ctx context.Context, priceRuleId uint64) ([]PriceRuleDiscountCode, error) {
	path := fmt.Sprintf(discountCodeBasePath+".json", priceRuleId)
	resource := new(DiscountCodesResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.DiscountCodes, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *DiscountCodeServiceOp) ListWithPagination(ctx context.Context, priceRuleId uint64, options interface{}) ([]PriceRuleDiscountCode, *Pagination, error) {
	path := fmt.Sprintf(discountCodeBasePath+".json", priceRuleId)
	resource := new(DiscountCodesResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.DiscountCodes, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *DiscountCodeServiceOp) ListAll(ctx context.Context, priceRuleId uint64, options interface{}) ([]PriceRuleDiscountCode, error) {
	items := []PriceRuleDiscountCode{}
	for {
		page, pagination, err := s.ListWithPagination(ctx, priceRuleId, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
}

// Get a single discount code
func (s *DiscountCodeServiceOp) Get(ctx context.Context, priceRuleId uint64, discountCodeId uint64) (*PriceRuleDiscountCode, error) {
	path := fmt.Sprintf(discountCodeBasePath+"/%d.json", priceRuleId, discountCodeId)
	resource := new(DiscountCodeResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.PriceRuleDiscountCode, err
}

// Delete a discount code
func (s *DiscountCodeServiceOp) Delete(ctx context.Context, priceRuleId uint64, discountCodeId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf(discountCodeBasePath+"/%d.json", priceRuleId, discountCodeId))
}
