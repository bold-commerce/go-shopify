package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	draftOrdersBasePath     = "draft_orders"
	draftOrdersResourceName = "draft_orders"
)

// DraftOrderService is an interface for interfacing with the draft orders endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/draftorder
type DraftOrderService interface {
	List(context.Context, interface{}) ([]DraftOrder, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]DraftOrder, *Pagination, error)
	ListAll(ctx context.Context, options interface{}) ([]DraftOrder, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*DraftOrder, error)
	Create(context.Context, DraftOrder) (*DraftOrder, error)
	Update(context.Context, DraftOrder) (*DraftOrder, error)
	Delete(context.Context, uint64) error
	Invoice(context.Context, uint64, DraftOrderInvoice) (*DraftOrderInvoice, error)
	Complete(context.Context, uint64, bool) (*DraftOrder, error)

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
	Id              uint64           `json:"id,omitempty"`
	OrderId         uint64           `json:"order_id,omitempty"`
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
	TaxesIncluded   *bool            `json:"taxes_included,omitempty"`
	TotalTax        string           `json:"total_tax,omitempty"`
	TaxExempt       *bool            `json:"tax_exempt,omitempty"`
	TotalPrice      string           `json:"total_price,omitempty"`
	SubtotalPrice   *decimal.Decimal `json:"subtotal_price,omitempty"`
	CompletedAt     *time.Time       `json:"completed_at,omitempty"`
	CreatedAt       *time.Time       `json:"created_at,omitempty"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty"`
	Status          string           `json:"status,omitempty"`
	// only in request to flag using the customer's default address
	UseCustomerDefaultAddress *bool `json:"use_customer_default_address,omitempty"`

	PaymentTerms                 *PaymentTerms `json:"payment_terms,omitempty"`
	SourceName                   string        `json:"source_name,omitempty"`
	AdminGraphqlApiId            string        `json:"admin_graphql_api_id,omitempty"`
	AllowDiscountCodesInCheckout *bool         `json:"allow_discount_codes_in_checkout?,omitempty"`
	AmountDueLaterSet            *AmountSet    `json:"amount_due_later_set,omitempty"`
	AmountDueNowSet              *AmountSet    `json:"amount_due_now_set,omitempty"`
	ApiClientId                  uint64        `json:"api_client_id,omitempty"`
	B2B                          *bool         `json:"b2b?,omitempty"`
	CreatedOnApiVersionHandle    string        `json:"created_on_api_version_handle,omitempty"`
	PresentmentCurrency          string        `json:"presentment_currency,omitempty"`
	SubtotalPriceSet             *AmountSet    `json:"subtotal_price_set,omitempty"`
	TotalAdditionalFeesSet       *AmountSet    `json:"total_additional_fees_set,omitempty"`
	TotalDiscountsSet            *AmountSet    `json:"total_discounts_set,omitempty"`
	TotalDutiesSet               *AmountSet    `json:"total_duties_set,omitempty"`
	TotalLineItemsPriceSet       *AmountSet    `json:"total_line_items_price_set,omitempty"`
	TotalPriceSet                *AmountSet    `json:"total_price_set,omitempty"`
	TotalShippingPriceSet        *AmountSet    `json:"total_shipping_price_set,omitempty"`
	TotalTaxSet                  *AmountSet    `json:"total_tax_set,omitempty"`
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

// DraftOrderStatus is a status used to filter draft orders.
type DraftOrderStatus string

const (
	DraftOrderStatusOpen        DraftOrderStatus = "open"
	DraftOrderStatusInvoiceSent DraftOrderStatus = "invoice_sent"
	DraftOrderStatusCompleted   DraftOrderStatus = "completed"
)

// DraftOrderListOptions represents the possible options that can be used
// to further query the list draft orders endpoint
type DraftOrderListOptions struct {
	Fields       string           `url:"fields,omitempty"`
	Limit        int              `url:"limit,omitempty"`
	SinceId      uint64           `url:"since_id,omitempty"`
	UpdatedAtMin *time.Time       `url:"updated_at_min,omitempty"`
	UpdatedAtMax *time.Time       `url:"updated_at_max,omitempty"`
	Ids          string           `url:"ids,omitempty"`
	Status       DraftOrderStatus `url:"status,omitempty"`
}

// DraftOrderCountOptions represents the possible options to the count draft orders endpoint
type DraftOrderCountOptions struct {
	SinceId      uint64           `url:"since_id,omitempty"`
	Status       DraftOrderStatus `url:"status,omitempty"`
	UpdatedAtMin *time.Time       `url:"updated_at_min,omitempty"`
	UpdatedAtMax *time.Time       `url:"updated_at_max,omitempty"`
	// Deprecated: Shopify does not document fields for draft order counts.
	Fields string `url:"fields,omitempty"`
	// Deprecated: Shopify does not document limit for draft order counts.
	Limit int `url:"limit,omitempty"`
	// Deprecated: Shopify does not document ids for draft order counts.
	Ids string `url:"ids,omitempty"`
}

// Create draft order
func (s *DraftOrderServiceOp) Create(ctx context.Context, draftOrder DraftOrder) (*DraftOrder, error) {
	path := fmt.Sprintf("%s.json", draftOrdersBasePath)
	return s.write(ctx, http.MethodPost, path, draftOrder)
}

// List draft orders
func (s *DraftOrderServiceOp) List(ctx context.Context, options interface{}) ([]DraftOrder, error) {
	path := fmt.Sprintf("%s.json", draftOrdersBasePath)
	resource := new(DraftOrdersResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.DraftOrders, err
}

// ListWithPagination returns one page and its next/previous page options.
func (s *DraftOrderServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]DraftOrder, *Pagination, error) {
	path := fmt.Sprintf("%s.json", draftOrdersBasePath)
	resource := new(DraftOrdersResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.DraftOrders, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *DraftOrderServiceOp) ListAll(ctx context.Context, options interface{}) ([]DraftOrder, error) {
	items := []DraftOrder{}
	err := forEachPage(options, func(options interface{}) (*Pagination, error) {
		page, pagination, err := s.ListWithPagination(ctx, options)
		if err == nil {
			items = append(items, page...)
		}
		return pagination, err
	})
	return items, err
}

// Count draft orders
func (s *DraftOrderServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", draftOrdersBasePath)
	return s.client.Count(ctx, path, options)
}

// Delete draft orders
func (s *DraftOrderServiceOp) Delete(ctx context.Context, draftOrderId uint64) error {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrderId)
	return s.client.Delete(ctx, path)
}

// Invoice a draft order
func (s *DraftOrderServiceOp) Invoice(ctx context.Context, draftOrderId uint64, draftOrderInvoice DraftOrderInvoice) (*DraftOrderInvoice, error) {
	path := fmt.Sprintf("%s/%d/send_invoice.json", draftOrdersBasePath, draftOrderId)
	wrappedData := DraftOrderInvoiceResource{DraftOrderInvoice: &draftOrderInvoice}
	resource := new(DraftOrderInvoiceResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.DraftOrderInvoice, err
}

// Get individual draft order
func (s *DraftOrderServiceOp) Get(ctx context.Context, draftOrderId uint64, options interface{}) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrderId)
	resource := new(DraftOrderResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.DraftOrder, err
}

// Update draft order
func (s *DraftOrderServiceOp) Update(ctx context.Context, draftOrder DraftOrder) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d.json", draftOrdersBasePath, draftOrder.Id)
	return s.write(ctx, http.MethodPut, path, draftOrder)
}

// Complete draft order
func (s *DraftOrderServiceOp) Complete(ctx context.Context, draftOrderId uint64, paymentPending bool) (*DraftOrder, error) {
	path := fmt.Sprintf("%s/%d/complete.json?payment_pending=%t", draftOrdersBasePath, draftOrderId, paymentPending)
	resource := new(DraftOrderResource)
	err := s.client.Put(ctx, path, nil, resource)
	return resource.DraftOrder, err
}

// List metafields for an order
func (s *DraftOrderServiceOp) ListMetafields(ctx context.Context, draftOrderId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.List(ctx, options)
}

// Count metafields for an order
func (s *DraftOrderServiceOp) CountMetafields(ctx context.Context, draftOrderId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for an order
func (s *DraftOrderServiceOp) GetMetafield(ctx context.Context, draftOrderId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for an order
func (s *DraftOrderServiceOp) CreateMetafield(ctx context.Context, draftOrderId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for an order
func (s *DraftOrderServiceOp) UpdateMetafield(ctx context.Context, draftOrderId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Update(ctx, metafield)
}

// Delete an existing metafield for an order
func (s *DraftOrderServiceOp) DeleteMetafield(ctx context.Context, draftOrderId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.Delete(ctx, metafieldId)
}

// ListMetafieldsWithPagination retrieves one page and its continuation options.
func (s *DraftOrderServiceOp) ListMetafieldsWithPagination(ctx context.Context, draftOrderId uint64, options interface{}) ([]Metafield, *Pagination, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.ListWithPagination(ctx, options)
}

// ListMetafieldsAll retrieves all pages, returning completed pages on error.
func (s *DraftOrderServiceOp) ListMetafieldsAll(ctx context.Context, draftOrderId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: draftOrdersResourceName, resourceId: draftOrderId}
	return metafieldService.ListAll(ctx, options)
}

// Shopify can accept the write before shipping and taxes are ready. Only poll
// with GET after acceptance, since repeating the write can create a second draft.
func (s *DraftOrderServiceOp) write(ctx context.Context, method, endpoint string, draft DraftOrder) (*DraftOrder, error) {
	c := s.client
	req, err := c.NewRequest(ctx, method, path.Join(c.pathPrefix, endpoint), DraftOrderResource{DraftOrder: &draft}, nil)
	if err != nil {
		return nil, err
	}
	for {
		resp, body, err := c.doRequest(req, c.retryPolicy(req.Method == http.MethodGet), c.recordRateLimitHeaders)
		if err != nil {
			if req.Method == http.MethodGet {
				return nil, fmt.Errorf("draft order write accepted, polling failed: %w", err)
			}
			return nil, err
		}
		if resp.StatusCode != http.StatusAccepted {
			var result DraftOrderResource
			if err := json.Unmarshal(body, &result); err != nil {
				return nil, err
			}
			return result.DraftOrder, nil
		}
		location := resp.Header.Get("Location")
		if location == "" && req.Method != http.MethodGet {
			return nil, fmt.Errorf("draft order write accepted without a polling location")
		}
		target, err := req.URL.Parse(location)
		if err != nil || target.User != nil || target.Scheme != c.baseURL.Scheme || !strings.EqualFold(target.Host, c.baseURL.Host) {
			return nil, fmt.Errorf("draft order write accepted with invalid polling location")
		}
		delay := time.Second
		if h := resp.Header.Get("Retry-After"); h != "" {
			seconds, err := strconv.ParseUint(h, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("draft order write accepted with invalid Retry-After: %w", err)
			}
			delay = time.Duration(seconds) * time.Second
		}
		if err := waitForRetry(ctx, delay); err != nil {
			return nil, fmt.Errorf("draft order write accepted, polling interrupted: %w", err)
		}
		req, err = c.NewRequest(ctx, http.MethodGet, target.String(), nil, nil)
		if err != nil {
			return nil, err
		}
	}
}
