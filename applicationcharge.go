package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const applicationChargesBasePath = "application_charges"

// ApplicationChargeService is an interface for interacting with the
// ApplicationCharge endpoints of the Shopify API.
// See https://help.shopify.com/api/reference/billing/applicationcharge
type ApplicationChargeService interface {
	Create(context.Context, ApplicationCharge) (*ApplicationCharge, error)
	Get(context.Context, int64, interface{}) (*ApplicationCharge, error)
	List(context.Context, interface{}) ([]ApplicationCharge, error)
	Activate(context.Context, ApplicationCharge) (*ApplicationCharge, error)
}

type ApplicationChargeServiceOp struct {
	client *Client
}

type ApplicationCharge struct {
	Id                 int64            `json:"id"`
	Name               string           `json:"name"`
	APIClientId        int64            `json:"api_client_id"`
	Price              *decimal.Decimal `json:"price"`
	Status             string           `json:"status"`
	ReturnURL          string           `json:"return_url"`
	Test               *bool            `json:"test"`
	CreatedAt          *time.Time       `json:"created_at"`
	UpdatedAt          *time.Time       `json:"updated_at"`
	ChargeType         *string          `json:"charge_type"`
	DecoratedReturnURL string           `json:"decorated_return_url"`
	ConfirmationURL    string           `json:"confirmation_url"`
}

// ApplicationChargeResource represents the result from the
// admin/application_charges{/X{/activate.json}.json}.json endpoints.
type ApplicationChargeResource struct {
	Charge *ApplicationCharge `json:"application_charge"`
}

// ApplicationChargesResource represents the result from the
// admin/application_charges.json endpoint.
type ApplicationChargesResource struct {
	Charges []ApplicationCharge `json:"application_charges"`
}

// Create creates new application charge.
func (a ApplicationChargeServiceOp) Create(ctx context.Context, charge ApplicationCharge) (*ApplicationCharge, error) {
	path := fmt.Sprintf("%s.json", applicationChargesBasePath)
	resource := &ApplicationChargeResource{}
	return resource.Charge, a.client.Post(ctx, path, ApplicationChargeResource{Charge: &charge}, resource)
}

// Get gets individual application charge.
func (a ApplicationChargeServiceOp) Get(ctx context.Context, chargeId int64, options interface{}) (*ApplicationCharge, error) {
	path := fmt.Sprintf("%s/%d.json", applicationChargesBasePath, chargeId)
	resource := &ApplicationChargeResource{}
	return resource.Charge, a.client.Get(ctx, path, resource, options)
}

// List gets all application charges.
func (a ApplicationChargeServiceOp) List(ctx context.Context, options interface{}) ([]ApplicationCharge, error) {
	path := fmt.Sprintf("%s.json", applicationChargesBasePath)
	resource := &ApplicationChargesResource{}
	return resource.Charges, a.client.Get(ctx, path, resource, options)
}

// Activate activates application charge.
func (a ApplicationChargeServiceOp) Activate(ctx context.Context, charge ApplicationCharge) (*ApplicationCharge, error) {
	path := fmt.Sprintf("%s/%d/activate.json", applicationChargesBasePath, charge.Id)
	resource := &ApplicationChargeResource{}
	return resource.Charge, a.client.Post(ctx, path, ApplicationChargeResource{Charge: &charge}, resource)
}
