package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const usageChargesPath = "usage_charges"

// UsageChargeService is an interface for interacting with the
// UsageCharge endpoints of the Shopify API.
// See https://help.shopify.com/en/api/reference/billing/usagecharge#endpoints
type UsageChargeService interface {
	Create(context.Context, uint64, UsageCharge) (*UsageCharge, error)
	Get(context.Context, uint64, uint64, interface{}) (*UsageCharge, error)
	List(context.Context, uint64, interface{}) ([]UsageCharge, error)
}

// UsageChargeServiceOp handles communication with the
// UsageCharge related methods of the Shopify API.
type UsageChargeServiceOp struct {
	client *Client
}

// UsageCharge represents a Shopify UsageCharge.
type UsageCharge struct {
	BalanceRemaining *decimal.Decimal `json:"balance_remaining,omitempty"`
	BalanceUsed      *decimal.Decimal `json:"balance_used,omitempty"`
	CreatedAt        *time.Time       `json:"created_at,omitempty"`
	Description      string           `json:"description,omitempty"`
	Id               uint64           `json:"id,omitempty"`
	Price            *decimal.Decimal `json:"price,omitempty"`
	RiskLevel        *decimal.Decimal `json:"risk_level,omitempty"`
}

// UsageChargeResource represents the result from the
// /admin/recurring_application_charges/X/usage_charges/X.json endpoints
type UsageChargeResource struct {
	Charge *UsageCharge `json:"usage_charge"`
}

// UsageChargesResource represents the result from the
// admin/recurring_application_charges/X/usage_charges.json endpoint.
type UsageChargesResource struct {
	Charges []UsageCharge `json:"usage_charges"`
}

// Create creates new usage charge given a recurring charge. *required fields: price and description
func (r *UsageChargeServiceOp) Create(ctx context.Context, chargeId uint64, usageCharge UsageCharge) (
	*UsageCharge, error,
) {
	path := fmt.Sprintf("%s/%d/%s.json", recurringApplicationChargesBasePath, chargeId, usageChargesPath)
	wrappedData := UsageChargeResource{Charge: &usageCharge}
	resource := &UsageChargeResource{}
	err := r.client.Post(ctx, path, wrappedData, resource)
	return resource.Charge, err
}

// Get gets individual usage charge.
func (r *UsageChargeServiceOp) Get(ctx context.Context, chargeId uint64, usageChargeId uint64, options interface{}) (
	*UsageCharge, error,
) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", recurringApplicationChargesBasePath, chargeId, usageChargesPath, usageChargeId)
	resource := &UsageChargeResource{}
	err := r.client.Get(ctx, path, resource, options)
	return resource.Charge, err
}

// List gets all usage charges associated with the recurring charge.
func (r *UsageChargeServiceOp) List(ctx context.Context, chargeId uint64, options interface{}) (
	[]UsageCharge, error,
) {
	path := fmt.Sprintf("%s/%d/%s.json", recurringApplicationChargesBasePath, chargeId, usageChargesPath)
	resource := &UsageChargesResource{}
	err := r.client.Get(ctx, path, resource, options)
	return resource.Charges, err
}
