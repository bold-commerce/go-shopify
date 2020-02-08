package goshopify

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const usageChargesPath = "usage_charges"

// UsageChargeService is an interface for interacting with the
// UsageCharge endpoints of the Shopify API.
// See https://help.shopify.com/en/api/reference/billing/usagecharge#endpoints
type UsageChargeService interface {
	Create(int64, UsageCharge) (*UsageCharge, error)
	Get(int64, int64, interface{}) (*UsageCharge, error)
	List(int64, interface{}) ([]UsageCharge, error)
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
	BillingOn        *time.Time       `json:"billing_on,omitempty"`
	CreatedAt        *time.Time       `json:"created_at,omitempty"`
	Description      string           `json:"description,omitempty"`
	ID               int64            `json:"id,omitempty"`
	Price            *decimal.Decimal `json:"price,omitempty"`
	RiskLevel        *decimal.Decimal `json:"risk_level,omitempty"`
}

func (r *UsageCharge) UnmarshalJSON(data []byte) error {
	// This is a workaround for the API returning BillingOn date in the format of "YYYY-MM-DD"
	// https://help.shopify.com/en/api/reference/billing/usagecharge#endpoints
	// For a longer explanation of the hack check:
	// http://choly.ca/post/go-json-marshalling/
	type alias UsageCharge
	aux := &struct {
		BillingOn *string `json:"billing_on"`
		*alias
	}{alias: (*alias)(r)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if err := parse(&r.BillingOn, aux.BillingOn); err != nil {
		return err
	}
	return nil
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
func (r *UsageChargeServiceOp) Create(chargeID int64, usageCharge UsageCharge) (
	*UsageCharge, error) {

	path := fmt.Sprintf("%s/%d/%s.json", recurringApplicationChargesBasePath, chargeID, usageChargesPath)
	wrappedData := UsageChargeResource{Charge: &usageCharge}
	resource := &UsageChargeResource{}
	err := r.client.Post(path, wrappedData, resource)
	return resource.Charge, err
}

// Get gets individual usage charge.
func (r *UsageChargeServiceOp) Get(chargeID int64, usageChargeID int64, options interface{}) (
	*UsageCharge, error) {

	path := fmt.Sprintf("%s/%d/%s/%d.json", recurringApplicationChargesBasePath, chargeID, usageChargesPath, usageChargeID)
	resource := &UsageChargeResource{}
	err := r.client.Get(path, resource, options)
	return resource.Charge, err
}

// List gets all usage charges associated with the recurring charge.
func (r *UsageChargeServiceOp) List(chargeID int64, options interface{}) (
	[]UsageCharge, error) {

	path := fmt.Sprintf("%s/%d/%s.json", recurringApplicationChargesBasePath, chargeID, usageChargesPath)
	resource := &UsageChargesResource{}
	err := r.client.Get(path, resource, options)
	return resource.Charges, err
}
