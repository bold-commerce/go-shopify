package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	ordersBasePath     = "orders"
	ordersResourceName = "orders"
)

// OrderService is an interface for interfacing with the orders endpoints of
// the Shopify API.
// See: https://help.shopify.com/api/reference/order
type OrderService interface {
	List(context.Context, interface{}) ([]Order, error)
	ListWithPagination(context.Context, interface{}) ([]Order, *Pagination, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, int64, interface{}) (*Order, error)
	Create(context.Context, Order) (*Order, error)
	Update(context.Context, Order) (*Order, error)
	Cancel(context.Context, int64, interface{}) (*Order, error)
	Close(context.Context, int64) (*Order, error)
	Open(context.Context, int64) (*Order, error)
	Delete(context.Context, int64) error

	// MetafieldsService used for Order resource to communicate with Metafields resource
	MetafieldsService

	// FulfillmentsService used for Order resource to communicate with Fulfillments resource
	FulfillmentsService
}

// OrderServiceOp handles communication with the order related methods of the
// Shopify API.
type OrderServiceOp struct {
	client *Client
}

type orderStatus string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/order#get-orders?status=any
const (
	// Show only open orders.
	OrderStatusOpen orderStatus = "open"

	// Show only closed orders.
	OrderStatusClosed orderStatus = "closed"

	// Show only cancelled orders.
	OrderStatusCancelled orderStatus = "cancelled"

	// Show orders of any status, open, closed, cancellerd, or archived.
	OrderStatusAny orderStatus = "any"
)

type orderFulfillmentStatus string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/order#get-orders?status=any
const (
	// Show orders that have been shipped.
	OrderFulfillmentStatusShipped orderFulfillmentStatus = "shipped"

	// Show partially shipped orders.
	OrderFulfillmentStatusPartial orderFulfillmentStatus = "partial"

	// Show orders that have not yet been shipped.
	OrderFulfillmentStatusUnshipped orderFulfillmentStatus = "unshipped"

	// Show orders of any fulfillment status.
	OrderFulfillmentStatusAny orderFulfillmentStatus = "any"

	// Returns orders with fulfillment_status of null or partial.
	OrderFulfillmentStatusUnfulfilled orderFulfillmentStatus = "unfulfilled"

	//"fulfilled" used to be an acceptable value? Was it deprecated? It isn't noted
	//in the Shopify docs at the provided URL, but it was used in tests and still
	//seems to function.
	OrderFulfillmentStatusFulfilled orderFulfillmentStatus = "fulfilled"
)

type orderFinancialStatus string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/order#get-orders?status=any
const (
	// Show only authorized orders.
	OrderFinancialStatusAuthorized orderFinancialStatus = "authorized"

	// Show only pending orders.
	OrderFinancialStatusPending orderFinancialStatus = "pending"

	// Show only paid orders.
	OrderFinancialStatusPaid orderFinancialStatus = "paid"

	// Show only partially paid orders.
	OrderFinancialStatusPartiallyPaid orderFinancialStatus = "partially_paid"

	// Show only refunded orders.
	OrderFinancialStatusRefunded orderFinancialStatus = "refunded"

	// Show only voided orders.
	OrderFinancialStatusVoided orderFinancialStatus = "voided"

	// Show only partially refunded orders.
	OrderFinancialStatusPartiallyRefunded orderFinancialStatus = "partially_refunded"

	// Show orders of any financial status.
	OrderFinancialStatusAny orderFinancialStatus = "any"

	// Show authorized and partially paid orders.
	OrderFinancialStatusUnpaid orderFinancialStatus = "unpaid"
)

type orderCancelReason string

const (
	// The customer canceled the order.
	OrderCancelReasonCustomer orderCancelReason = "customer"

	// The order was fraudulent.
	OrderCancelReasonFraud orderCancelReason = "fraud"

	// Items in the order were not in inventory.
	OrderCancelReasonInventory orderCancelReason = "inventory"

	// The payment was declined.
	OrderCancelReasonDeclined orderCancelReason = "declined"

	// Cancelled for some other reason.
	OrderCancelReasonOther orderCancelReason = "other"
)

// A struct for all available order count options
type OrderCountOptions struct {
	Page              int                    `url:"page,omitempty"`
	Limit             int                    `url:"limit,omitempty"`
	SinceID           int64                  `url:"since_id,omitempty"`
	CreatedAtMin      time.Time              `url:"created_at_min,omitempty"`
	CreatedAtMax      time.Time              `url:"created_at_max,omitempty"`
	UpdatedAtMin      time.Time              `url:"updated_at_min,omitempty"`
	UpdatedAtMax      time.Time              `url:"updated_at_max,omitempty"`
	Order             string                 `url:"order,omitempty"`
	Fields            string                 `url:"fields,omitempty"`
	Status            orderStatus            `url:"status,omitempty"`
	FinancialStatus   orderFinancialStatus   `url:"financial_status,omitempty"`
	FulfillmentStatus orderFulfillmentStatus `url:"fulfillment_status,omitempty"`
}

// A struct for all available order list options.
// See: https://help.shopify.com/api/reference/order#index
type OrderListOptions struct {
	ListOptions
	Status            orderStatus            `url:"status,omitempty"`
	FinancialStatus   orderFinancialStatus   `url:"financial_status,omitempty"`
	FulfillmentStatus orderFulfillmentStatus `url:"fulfillment_status,omitempty"`
	ProcessedAtMin    time.Time              `url:"processed_at_min,omitempty"`
	ProcessedAtMax    time.Time              `url:"processed_at_max,omitempty"`
	Order             string                 `url:"order,omitempty"`
}

// A struct of all available order cancel options.
// See: https://help.shopify.com/api/reference/order#index
type OrderCancelOptions struct {
	Amount   *decimal.Decimal `json:"amount,omitempty"`
	Currency string           `json:"currency,omitempty"`
	Restock  bool             `json:"restock,omitempty"`
	Reason   string           `json:"reason,omitempty"`
	Email    bool             `json:"email,omitempty"`
	Refund   *Refund          `json:"refund,omitempty"`
}

// Order represents a Shopify order
type Order struct {
	ID                     int64                  `json:"id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	CreatedAt              *time.Time             `json:"created_at,omitempty"`
	UpdatedAt              *time.Time             `json:"updated_at,omitempty"`
	CancelledAt            *time.Time             `json:"cancelled_at,omitempty"`
	ClosedAt               *time.Time             `json:"closed_at,omitempty"`
	ProcessedAt            *time.Time             `json:"processed_at,omitempty"`
	Customer               *Customer              `json:"customer,omitempty"`
	BillingAddress         *Address               `json:"billing_address,omitempty"`
	ShippingAddress        *Address               `json:"shipping_address,omitempty"`
	Currency               string                 `json:"currency,omitempty"`
	TotalPrice             *decimal.Decimal       `json:"total_price,omitempty"`
	TotalPriceSet          *AmountSet             `json:"total_price_set,omitempty"`
	TotalShippingPriceSet  *AmountSet             `json:"total_shipping_price_set,omitempty"`
	CurrentTotalPrice      *decimal.Decimal       `json:"current_total_price,omitempty"`
	SubtotalPrice          *decimal.Decimal       `json:"subtotal_price,omitempty"`
	CurrentSubtotalPrice   *decimal.Decimal       `json:"current_subtotal_price,omitempty"`
	TotalDiscounts         *decimal.Decimal       `json:"total_discounts,omitempty"`
	CurrentTotalDiscounts  *decimal.Decimal       `json:"current_total_discounts,omitempty"`
	TotalLineItemsPrice    *decimal.Decimal       `json:"total_line_items_price,omitempty"`
	TaxesIncluded          bool                   `json:"taxes_included,omitempty"`
	TotalTax               *decimal.Decimal       `json:"total_tax,omitempty"`
	TotalTaxSet            *AmountSet             `json:"total_tax_set,omitempty"`
	CurrentTotalTax        *decimal.Decimal       `json:"current_total_tax,omitempty"`
	TaxLines               []TaxLine              `json:"tax_lines,omitempty"`
	TotalWeight            int                    `json:"total_weight,omitempty"`
	FinancialStatus        orderFinancialStatus   `json:"financial_status,omitempty"`
	Fulfillments           []Fulfillment          `json:"fulfillments,omitempty"`
	FulfillmentStatus      orderFulfillmentStatus `json:"fulfillment_status,omitempty"`
	Token                  string                 `json:"token,omitempty"`
	CartToken              string                 `json:"cart_token,omitempty"`
	Number                 int                    `json:"number,omitempty"`
	OrderNumber            int                    `json:"order_number,omitempty"`
	Note                   string                 `json:"note,omitempty"`
	Test                   bool                   `json:"test,omitempty"`
	BrowserIp              string                 `json:"browser_ip,omitempty"`
	BuyerAcceptsMarketing  bool                   `json:"buyer_accepts_marketing,omitempty"`
	CancelReason           orderCancelReason      `json:"cancel_reason,omitempty"`
	NoteAttributes         []NoteAttribute        `json:"note_attributes,omitempty"`
	DiscountCodes          []DiscountCode         `json:"discount_codes,omitempty"`
	LineItems              []LineItem             `json:"line_items,omitempty"`
	ShippingLines          []ShippingLines        `json:"shipping_lines,omitempty"`
	Transactions           []Transaction          `json:"transactions,omitempty"`
	AppID                  int                    `json:"app_id,omitempty"`
	CustomerLocale         string                 `json:"customer_locale,omitempty"`
	LandingSite            string                 `json:"landing_site,omitempty"`
	ReferringSite          string                 `json:"referring_site,omitempty"`
	SourceName             string                 `json:"source_name,omitempty"`
	ClientDetails          *ClientDetails         `json:"client_details,omitempty"`
	Tags                   string                 `json:"tags,omitempty"`
	LocationId             int64                  `json:"location_id,omitempty"`
	PaymentGatewayNames    []string               `json:"payment_gateway_names,omitempty"`
	ProcessingMethod       string                 `json:"processing_method,omitempty"`
	Refunds                []Refund               `json:"refunds,omitempty"`
	UserId                 int64                  `json:"user_id,omitempty"`
	OrderStatusUrl         string                 `json:"order_status_url,omitempty"`
	Gateway                string                 `json:"gateway,omitempty"`
	Confirmed              bool                   `json:"confirmed,omitempty"`
	CheckoutToken          string                 `json:"checkout_token,omitempty"`
	Reference              string                 `json:"reference,omitempty"`
	SourceIdentifier       string                 `json:"source_identifier,omitempty"`
	SourceURL              string                 `json:"source_url,omitempty"`
	DeviceID               int64                  `json:"device_id,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	LandingSiteRef         string                 `json:"landing_site_ref,omitempty"`
	CheckoutID             int64                  `json:"checkout_id,omitempty"`
	ContactEmail           string                 `json:"contact_email,omitempty"`
	Metafields             []Metafield            `json:"metafields,omitempty"`
	SendReceipt            bool                   `json:"send_receipt,omitempty"`
	SendFulfillmentReceipt bool                   `json:"send_fulfillment_receipt,omitempty"`
}

type Address struct {
	ID           int64   `json:"id,omitempty"`
	Address1     string  `json:"address1,omitempty"`
	Address2     string  `json:"address2,omitempty"`
	City         string  `json:"city,omitempty"`
	Company      string  `json:"company,omitempty"`
	Country      string  `json:"country,omitempty"`
	CountryCode  string  `json:"country_code,omitempty"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	Name         string  `json:"name,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Province     string  `json:"province,omitempty"`
	ProvinceCode string  `json:"province_code,omitempty"`
	Zip          string  `json:"zip,omitempty"`
}

type DiscountCode struct {
	Amount *decimal.Decimal `json:"amount,omitempty"`
	Code   string           `json:"code,omitempty"`
	Type   string           `json:"type,omitempty"`
}

type LineItem struct {
	ID                         int64                  `json:"id,omitempty"`
	ProductID                  int64                  `json:"product_id,omitempty"`
	VariantID                  int64                  `json:"variant_id,omitempty"`
	Quantity                   int                    `json:"quantity,omitempty"`
	Price                      *decimal.Decimal       `json:"price,omitempty"`
	TotalDiscount              *decimal.Decimal       `json:"total_discount,omitempty"`
	Title                      string                 `json:"title,omitempty"`
	VariantTitle               string                 `json:"variant_title,omitempty"`
	Name                       string                 `json:"name,omitempty"`
	SKU                        string                 `json:"sku,omitempty"`
	Vendor                     string                 `json:"vendor,omitempty"`
	GiftCard                   bool                   `json:"gift_card,omitempty"`
	Taxable                    bool                   `json:"taxable,omitempty"`
	FulfillmentService         string                 `json:"fulfillment_service,omitempty"`
	RequiresShipping           bool                   `json:"requires_shipping,omitempty"`
	VariantInventoryManagement string                 `json:"variant_inventory_management,omitempty"`
	PreTaxPrice                *decimal.Decimal       `json:"pre_tax_price,omitempty"`
	Properties                 []NoteAttribute        `json:"properties,omitempty"`
	ProductExists              bool                   `json:"product_exists,omitempty"`
	FulfillableQuantity        int                    `json:"fulfillable_quantity,omitempty"`
	Grams                      int                    `json:"grams,omitempty"`
	FulfillmentStatus          orderFulfillmentStatus `json:"fulfillment_status,omitempty"`
	TaxLines                   []TaxLine              `json:"tax_lines,omitempty"`

	// Deprecated: See 2022-10 release notes: https://shopify.dev/docs/api/release-notes/2022-10
	OriginLocation *Address `json:"origin_location,omitempty"`

	// Deprecated: See 2022-10 release notes: https://shopify.dev/docs/api/release-notes/2022-10
	DestinationLocation *Address `json:"destination_location,omitempty"`

	AppliedDiscount     *AppliedDiscount      `json:"applied_discount,omitempty"`
	DiscountAllocations []DiscountAllocations `json:"discount_allocations,omitempty"`
}

type DiscountAllocations struct {
	Amount                   *decimal.Decimal `json:"amount,omitempty"`
	DiscountApplicationIndex int              `json:"discount_application_index,omitempty"`
	AmountSet                *AmountSet       `json:"amount_set,omitempty"`
}

type AmountSet struct {
	ShopMoney        AmountSetEntry `json:"shop_money,omitempty"`
	PresentmentMoney AmountSetEntry `json:"presentment_money,omitempty"`
}

type AmountSetEntry struct {
	Amount       *decimal.Decimal `json:"amount,omitempty"`
	CurrencyCode string           `json:"currency_code,omitempty"`
}

// UnmarshalJSON custom unmarsaller for LineItem required to mitigate some older orders having LineItem.Properies
// which are empty JSON objects rather than the expected array.
func (li *LineItem) UnmarshalJSON(data []byte) error {
	type alias LineItem
	aux := &struct {
		Properties json.RawMessage `json:"properties"`
		*alias
	}{alias: (*alias)(li)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if len(aux.Properties) == 0 {
		return nil
	} else if aux.Properties[0] == '[' { // if the first character is a '[' we unmarshal into an array
		var p []NoteAttribute
		err = json.Unmarshal(aux.Properties, &p)
		if err != nil {
			return err
		}
		li.Properties = p
	} else { // else we unmarshal it into a struct
		var p NoteAttribute
		err = json.Unmarshal(aux.Properties, &p)
		if err != nil {
			return err
		}
		if p.Name == "" && p.Value == nil { // if the struct is empty we set properties to nil
			li.Properties = nil
		} else {
			li.Properties = []NoteAttribute{p} // else we set them to an array with the property nested
		}
	}

	return nil
}

type LineItemProperty struct {
	Message string `json:"message"`
}

type NoteAttribute struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Represents the result from the orders/X.json endpoint
type OrderResource struct {
	Order *Order `json:"order"`
}

// Represents the result from the orders.json endpoint
type OrdersResource struct {
	Orders []Order `json:"orders"`
}

type PaymentDetails struct {
	AVSResultCode     string `json:"avs_result_code,omitempty"`
	CreditCardBin     string `json:"credit_card_bin,omitempty"`
	CVVResultCode     string `json:"cvv_result_code,omitempty"`
	CreditCardNumber  string `json:"credit_card_number,omitempty"`
	CreditCardCompany string `json:"credit_card_company,omitempty"`
}

type ShippingLines struct {
	ID                            int64            `json:"id,omitempty"`
	Title                         string           `json:"title,omitempty"`
	Price                         *decimal.Decimal `json:"price,omitempty"`
	PriceSet                      *AmountSet       `json:"price_set,omitempty"`
	DiscountedPrice               *decimal.Decimal `json:"discounted_price,omitempty"`
	DiscountedPriceSet            *AmountSet       `json:"discounted_price_set,omitempty"`
	Code                          string           `json:"code,omitempty"`
	Source                        string           `json:"source,omitempty"`
	Phone                         string           `json:"phone,omitempty"`
	RequestedFulfillmentServiceID string           `json:"requested_fulfillment_service_id,omitempty"`
	DeliveryCategory              string           `json:"delivery_category,omitempty"`
	CarrierIdentifier             string           `json:"carrier_identifier,omitempty"`
	TaxLines                      []TaxLine        `json:"tax_lines,omitempty"`
	Handle                        string           `json:"handle,omitempty"`
}

// UnmarshalJSON custom unmarshaller for ShippingLines implemented to handle requested_fulfillment_service_id being
// returned as json numbers or json nulls instead of json strings
func (sl *ShippingLines) UnmarshalJSON(data []byte) error {
	type alias ShippingLines
	aux := &struct {
		*alias
		RequestedFulfillmentServiceID interface{} `json:"requested_fulfillment_service_id"`
	}{alias: (*alias)(sl)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	switch aux.RequestedFulfillmentServiceID.(type) {
	case nil:
		sl.RequestedFulfillmentServiceID = ""
	default:
		sl.RequestedFulfillmentServiceID = fmt.Sprintf("%v", aux.RequestedFulfillmentServiceID)
	}

	return nil
}

type TaxLine struct {
	Title string           `json:"title,omitempty"`
	Price *decimal.Decimal `json:"price,omitempty"`
	Rate  *decimal.Decimal `json:"rate,omitempty"`
}

type Transaction struct {
	ID             int64            `json:"id,omitempty"`
	OrderID        int64            `json:"order_id,omitempty"`
	Amount         *decimal.Decimal `json:"amount,omitempty"`
	Kind           string           `json:"kind,omitempty"`
	Gateway        string           `json:"gateway,omitempty"`
	Status         string           `json:"status,omitempty"`
	Message        string           `json:"message,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	Test           bool             `json:"test,omitempty"`
	Authorization  string           `json:"authorization,omitempty"`
	Currency       string           `json:"currency,omitempty"`
	LocationID     *int64           `json:"location_id,omitempty"`
	UserID         *int64           `json:"user_id,omitempty"`
	ParentID       *int64           `json:"parent_id,omitempty"`
	DeviceID       *int64           `json:"device_id,omitempty"`
	ErrorCode      string           `json:"error_code,omitempty"`
	SourceName     string           `json:"source_name,omitempty"`
	Source         string           `json:"source,omitempty"`
	PaymentDetails *PaymentDetails  `json:"payment_details,omitempty"`
}

type ClientDetails struct {
	AcceptLanguage string `json:"accept_language,omitempty"`
	BrowserHeight  int    `json:"browser_height,omitempty"`
	BrowserIp      string `json:"browser_ip,omitempty"`
	BrowserWidth   int    `json:"browser_width,omitempty"`
	SessionHash    string `json:"session_hash,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
}

type Refund struct {
	Id               int64             `json:"id,omitempty"`
	OrderId          int64             `json:"order_id,omitempty"`
	CreatedAt        *time.Time        `json:"created_at,omitempty"`
	Note             string            `json:"note,omitempty"`
	Restock          bool              `json:"restock,omitempty"`
	UserId           int64             `json:"user_id,omitempty"`
	RefundLineItems  []RefundLineItem  `json:"refund_line_items,omitempty"`
	Transactions     []Transaction     `json:"transactions,omitempty"`
	OrderAdjustments []OrderAdjustment `json:"order_adjustments,omitempty"`
}

type OrderAdjustment struct {
	Id           int64               `json:"id,omitempty"`
	OrderId      int64               `json:"order_id,omitempty"`
	RefundId     int64               `json:"refund_id,omitempty"`
	Amount       *decimal.Decimal    `json:"amount,omitempty"`
	TaxAmount    *decimal.Decimal    `json:"tax_amount,omitempty"`
	Kind         OrderAdjustmentType `json:"kind,omitempty"`
	Reason       string              `json:"reason,omitempty"`
	AmountSet    *AmountSet          `json:"amount_set,omitempty"`
	TaxAmountSet *AmountSet          `json:"tax_amount_set,omitempty"`
}

type OrderAdjustmentType string

const (
	OrderAdjustmentTypeShippingRefund    OrderAdjustmentType = "shipping_refund"
	OrderAdjustmentTypeRefundDiscrepancy OrderAdjustmentType = "refund_discrepancy"
)

type RefundLineItem struct {
	Id         int64            `json:"id,omitempty"`
	Quantity   int              `json:"quantity,omitempty"`
	LineItemId int64            `json:"line_item_id,omitempty"`
	LineItem   *LineItem        `json:"line_item,omitempty"`
	Subtotal   *decimal.Decimal `json:"subtotal,omitempty"`
	TotalTax   *decimal.Decimal `json:"total_tax,omitempty"`
}

// List orders
func (s *OrderServiceOp) List(ctx context.Context, options interface{}) ([]Order, error) {
	orders, _, err := s.ListWithPagination(ctx, options)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Order, *Pagination, error) {
	path := fmt.Sprintf("%s.json", ordersBasePath)
	resource := new(OrdersResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Orders, pagination, nil
}

// Count orders
func (s *OrderServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", ordersBasePath)
	return s.client.Count(ctx, path, options)
}

// Get individual order
func (s *OrderServiceOp) Get(ctx context.Context, orderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, orderID)
	resource := new(OrderResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Order, err
}

// Create order
func (s *OrderServiceOp) Create(ctx context.Context, order Order) (*Order, error) {
	path := fmt.Sprintf("%s.json", ordersBasePath)
	wrappedData := OrderResource{Order: &order}
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Order, err
}

// Update order
func (s *OrderServiceOp) Update(ctx context.Context, order Order) (*Order, error) {
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, order.ID)
	wrappedData := OrderResource{Order: &order}
	resource := new(OrderResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Order, err
}

// Cancel order
func (s *OrderServiceOp) Cancel(ctx context.Context, orderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d/cancel.json", ordersBasePath, orderID)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, options, resource)
	return resource.Order, err
}

// Close order
func (s *OrderServiceOp) Close(ctx context.Context, orderID int64) (*Order, error) {
	path := fmt.Sprintf("%s/%d/close.json", ordersBasePath, orderID)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, nil, resource)
	return resource.Order, err
}

// Open order
func (s *OrderServiceOp) Open(ctx context.Context, orderID int64) (*Order, error) {
	path := fmt.Sprintf("%s/%d/open.json", ordersBasePath, orderID)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, nil, resource)
	return resource.Order, err
}

// Delete order
func (s *OrderServiceOp) Delete(ctx context.Context, orderID int64) error {
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, orderID)
	err := s.client.Delete(ctx, path)
	return err
}

// List metafields for an order
func (s *OrderServiceOp) ListMetafields(ctx context.Context, orderID int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.List(ctx, options)
}

// Count metafields for an order
func (s *OrderServiceOp) CountMetafields(ctx context.Context, orderID int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for an order
func (s *OrderServiceOp) GetMetafield(ctx context.Context, orderID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.Get(ctx, metafieldID, options)
}

// Create a new metafield for an order
func (s *OrderServiceOp) CreateMetafield(ctx context.Context, orderID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for an order
func (s *OrderServiceOp) UpdateMetafield(ctx context.Context, orderID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.Update(ctx, metafield)
}

// Delete an existing metafield for an order
func (s *OrderServiceOp) DeleteMetafield(ctx context.Context, orderID int64, metafieldID int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return metafieldService.Delete(ctx, metafieldID)
}

// List fulfillments for an order
func (s *OrderServiceOp) ListFulfillments(ctx context.Context, orderID int64, options interface{}) ([]Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.List(ctx, options)
}

// Count fulfillments for an order
func (s *OrderServiceOp) CountFulfillments(ctx context.Context, orderID int64, options interface{}) (int, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Count(ctx, options)
}

// Get individual fulfillment for an order
func (s *OrderServiceOp) GetFulfillment(ctx context.Context, orderID int64, fulfillmentID int64, options interface{}) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Get(ctx, fulfillmentID, options)
}

// Create a new fulfillment for an order
func (s *OrderServiceOp) CreateFulfillment(ctx context.Context, orderID int64, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Create(ctx, fulfillment)
}

// Update an existing fulfillment for an order
func (s *OrderServiceOp) UpdateFulfillment(ctx context.Context, orderID int64, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Update(ctx, fulfillment)
}

// Complete an existing fulfillment for an order
func (s *OrderServiceOp) CompleteFulfillment(ctx context.Context, orderID int64, fulfillmentID int64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Complete(ctx, fulfillmentID)
}

// Transition an existing fulfillment for an order
func (s *OrderServiceOp) TransitionFulfillment(ctx context.Context, orderID int64, fulfillmentID int64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Transition(ctx, fulfillmentID)
}

// Cancel an existing fulfillment for an order
func (s *OrderServiceOp) CancelFulfillment(ctx context.Context, orderID int64, fulfillmentID int64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Cancel(ctx, fulfillmentID)
}
