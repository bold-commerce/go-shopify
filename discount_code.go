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
	Count(ctx context.Context, options interface{}) (int, error)
	Lookup(ctx context.Context, code string) (*PriceRuleDiscountCode, error)
	CreateBatch(ctx context.Context, priceRuleId uint64, codes []PriceRuleDiscountCode) (*DiscountCodeCreation, error)
	GetBatch(ctx context.Context, priceRuleId, batchId uint64) (*DiscountCodeCreation, error)
	ListBatchCodes(ctx context.Context, priceRuleId, batchId uint64, options interface{}) ([]BatchDiscountCode, error)
	Create(context.Context, uint64, PriceRuleDiscountCode) (*PriceRuleDiscountCode, error)
	Update(context.Context, uint64, PriceRuleDiscountCode) (*PriceRuleDiscountCode, error)
	List(context.Context, uint64) ([]PriceRuleDiscountCode, error)
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

// DiscountCodeCreation represents Shopify discount code creation.
type DiscountCodeCreation struct {
	Id            uint64     `json:"id,omitempty"`
	PriceRuleId   uint64     `json:"price_rule_id,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Status        string     `json:"status,omitempty"`
	CodesCount    int        `json:"codes_count,omitempty"`
	ImportedCount int        `json:"imported_count,omitempty"`
	FailedCount   int        `json:"failed_count,omitempty"`
	Logs          []string   `json:"logs,omitempty"`
}

// BatchDiscountCode represents Shopify batch discount code.
type BatchDiscountCode struct {
	Id     *uint64             `json:"id,omitempty"`
	Code   string              `json:"code,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

// Count counts discount code records.
func (s *DiscountCodeServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	return s.client.Count(ctx, "discount_codes/count.json", options)
}

// Lookup finds a discount code, following Shopify's redirect to the code.
func (s *DiscountCodeServiceOp) Lookup(ctx context.Context, code string) (*PriceRuleDiscountCode, error) {
	resource := new(DiscountCodeResource)
	err := s.client.Get(ctx, "discount_codes/lookup.json", resource, struct {
		Code string `url:"code"`
	}{code})
	return resource.PriceRuleDiscountCode, err
}

// CreateBatch starts an asynchronous job to create discount codes.
func (s *DiscountCodeServiceOp) CreateBatch(ctx context.Context, priceRuleId uint64, codes []PriceRuleDiscountCode) (*DiscountCodeCreation, error) {
	var resource struct {
		Creation *DiscountCodeCreation `json:"discount_code_creation"`
	}
	err := s.client.Post(ctx, fmt.Sprintf("price_rules/%d/batch.json", priceRuleId), DiscountCodesResource{DiscountCodes: codes}, &resource)
	return resource.Creation, err
}

// GetBatch retrieves the status of a discount code creation job.
func (s *DiscountCodeServiceOp) GetBatch(ctx context.Context, priceRuleId, batchId uint64) (*DiscountCodeCreation, error) {
	var resource struct {
		Creation *DiscountCodeCreation `json:"discount_code_creation"`
	}
	err := s.client.Get(ctx, fmt.Sprintf("price_rules/%d/batch/%d.json", priceRuleId, batchId), &resource, nil)
	return resource.Creation, err
}

// ListBatchCodes retrieves the codes and errors from a discount code creation job.
func (s *DiscountCodeServiceOp) ListBatchCodes(ctx context.Context, priceRuleId, batchId uint64, options interface{}) ([]BatchDiscountCode, error) {
	var resource struct {
		Codes []BatchDiscountCode `json:"discount_codes"`
	}
	err := s.client.Get(ctx, fmt.Sprintf("price_rules/%d/batch/%d/discount_codes.json", priceRuleId, batchId), &resource, options)
	return resource.Codes, err
}
