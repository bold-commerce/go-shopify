package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	draftOrdersBasePath     = "draft_orders"
	draftOrdersResourceName = "draft_orders"
)

// DraftOrderService is an interface for interfacing with the draft orders endpoints of
// the Shopify API.
// See: https://help.shopify.com/api/reference/orders/draftorder
type DraftOrderService interface {
	List(context.Context, interface{}) ([]DraftOrder, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, int64, interface{}) (*DraftOrder, error)
	Create(context.Context, DraftOrder) (*DraftOrder, error)
	Update(context.Context, DraftOrder) (*DraftOrder, error)
	Delete(context.Context, int64) error
	Invoice(context.Context, int64, DraftOrderInvoice) (*DraftOrderInvoice, error)
	Complete(context.Context, int64, bool) (*DraftOrder, error)

	// MetafieldsService used for DrafT Order resource to communicate with Metafields resource
	MetafieldsService
}

// DraftOrderServiceOp handles communication with the draft order related methods of the
// Shopify API.
type DraftOrderServiceOp struct {
	client *Client
}

// DraftOrder represents a shopify draft order
type DraftOrder struct {
	Id              int64            `json:"id,omitempty"`
	OrderId         int64            `json:"order_id,omitempty"`
	Name            string           `json:"name,omitempty"`
	Customer        *Customer        `json:"customer,omitempty"`
	ShippingAddress *Address         `json:"shipping_address,omitempty"`
	BillingAddress  *Address         `json:"billing_address,omitempty"`
	Note            string           `json:"note,omitempty"`
	NoteAttributes  []NoteAttribute  `json:"note_attributes,omitempty"`
	Email           string           `json:"email,omitempty"`
	Currency        string           `json:"currency,omitempty"`
	InvoiceSentAt   *time.Time       `json:"invoice_sent_at,omitempty"`
	InvoiceURL      string           `json:"invoice_url,omitempty"`
	LineItems       []LineItem       `json:"line_items,omitempty"`
	ShippingLine    *ShippingLines   `json:"shipping_line,omitempty"`
	Tags            string           `json:"tags,omitempty"`
	TaxLines        []TaxLine        `json:"tax_lines,omitempty"`
	AppliedDiscount *AppliedDiscount `json:"applied_discount,omitempty"`
	TaxesIncluded   bool             `json:"taxes_included,omitempty"`
	TotalTax        string           `json:"total_tax,omitempty"`
	TaxExempt       *bool            `json:"tax_exempt,omitempty"`
	TotalPrice      string           `json:"total_price,omitempty"`
	SubtotalPrice   *decimal.Decimal `json:"subtotal_price,omitempty"`
	CompletedAt     *time.Time       `json:"completed_at,omitempty"`
	CreatedAt       *time.Time       `json:"created_at,omitempty"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty"`
	Status          string           `json:"status,omitempty"`
	// only in request to flag using the customer's default address
	UseCustomerDefaultAddress bool `json:"use_customer_default_address,omitempty"`
}

// AppliedDiscount is the discount applied to the line item or the draft order object.
type AppliedDiscount struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
	ValueType   string `json:"value_type,omitempty"`
	Amount      string `json:"amount,omitempty"`
}

// DraftOrderInvoice is the struct used to create an invoice for a draft order
type DraftOrderInvoice struct {
	To            string   `json:"to,omitempty"`
	From          string   `json:"from,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	CustomMessage string   `json:"custom_message,omitempty"`
	Bcc           []string `json:"bcc,omitempty"`
}

type DraftOrdersResource struct {
	DraftOrders []DraftOrder `json:"draft_orders"`
}

type DraftOrderResource struct {
	DraftOrder *DraftOrder `json:"draft_order"`
}

type DraftOrderInvoiceResource struct {
	DraftOrderInvoice *DraftOrderInvoice `json:"draft_order_invoice,omitempty"`
}

// DraftOrderListOptions represents the possible options that can be used
// to further query the list draft orders endpoint
type DraftOrderListOptions struct {
	Fields       string      `url:"fields,omitempty"`
	Limit        int         `url:"limit,omitempty"`
	SinceId      int64       `url:"since_id,omitempty"`
	UpdatedAtMin *time.Time  `url:"updated_at_min,omitempty"`
	UpdatedAtMax *time.Time  `url:"updated_at_max,omitempty"`
	Ids          string      `url:"ids,omitempty"`
	Status       orderStatus `url:"status,omitempty"`
}

// DraftOrderCountOptions represents the possible options to the count draft orders endpoint
type DraftOrderCountOptions struct {
	Fields  string      `url:"fields,omitempty"`
	Limit   int         `url:"limit,omitempty"`
	SinceId int64       `url:"since_id,omitempty"`
	Ids     string      `url:"ids,omitempty"`
	Status  orderStatus `url:"status,omitempty"`
}

// Create draft order
func (s *DraftOrderServiceOp) Create(ctx context.Context, draftOrder DraftOrder) (*DraftOrder, error) {
	path := fmt.Sprintf("%s.json", draftOrdersBasePath)
	wrappedData := DraftOrderResource{DraftOrder: &draftOrder}
	resource := new(DraftOrderResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.DraftOrder, err
}

// List draft orders
func (s *DraftOrderServiceOp) List(ctx context.Context, options interface{}) ([]DraftOrder, error) {
	path := fmt.Sprintf("%s.json", draftOrdersBasePath)
	resource := new(DraftOrdersResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.DraftOrders, err
}

// Count draft orders
func (s *DraftOrderServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", draftOrdersBasePath)
	return s.client.Count(ctx, path, options)
}

// Delete draft orders
func (s *DraftOrderServiceOp) Delete(ctx context.Context, draftOrderId int64) error {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrderId)
	return s.client.Delete(ctx, path)
}

// Invoice a draft order
func (s *DraftOrderServiceOp) Invoice(ctx context.Context, draftOrderId int64, draftOrderInvoice DraftOrderInvoice) (*DraftOrderInvoice, error) {
	path := fmt.Sprintf("%s/%d/send_invoice.json", draftOrdersBasePath, draftOrderId)
	wrappedData := DraftOrderInvoiceResource{DraftOrderInvoice: &draftOrderInvoice}
	resource := new(DraftOrderInvoiceResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.DraftOrderInvoice, err
}

// Get individual draft order
func (s *DraftOrderServiceOp) Get(ctx context.Context, draftOrderId int64, options interface{}) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrderId)
	resource := new(DraftOrderResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.DraftOrder, err
}

// Update draft order
func (s *DraftOrderServiceOp) Update(ctx context.Context, draftOrder DraftOrder) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrder.Id)
	wrappedData := DraftOrderResource{DraftOrder: &draftOrder}
	resource := new(DraftOrderResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.DraftOrder, err
}

// Complete draft order
func (s *DraftOrderServiceOp) Complete(ctx context.Context, draftOrderId int64, paymentPending bool) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d/complete.json?payment_pending=%t", draftOrdersBasePath, draftOrderId, paymentPending)
	resource := new(DraftOrderResource)
	err := s.client.Put(ctx, path, nil, resource)
	return resource.DraftOrder, err
}

// List metafields for an order
func (s *DraftOrderServiceOp) ListMetafields(ctx context.Context, draftOrderId int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.List(ctx, options)
}

// Count metafields for an order
func (s *DraftOrderServiceOp) CountMetafields(ctx context.Context, draftOrderId int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for an order
func (s *DraftOrderServiceOp) GetMetafield(ctx context.Context, draftOrderId int64, metafieldId int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for an order
func (s *DraftOrderServiceOp) CreateMetafield(ctx context.Context, draftOrderId int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for an order
func (s *DraftOrderServiceOp) UpdateMetafield(ctx context.Context, draftOrderId int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Update(ctx, metafield)
}

// Delete an existing metafield for an order
func (s *DraftOrderServiceOp) DeleteMetafield(ctx context.Context, draftOrderId int64, metafieldId int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Delete(ctx, metafieldId)
}
