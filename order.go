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
	ListAll(context.Context, interface{}) ([]Order, error)
	ListWithPagination(context.Context, interface{}) ([]Order, *Pagination, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Order, error)
	Create(context.Context, Order) (*Order, error)
	Update(context.Context, Order) (*Order, error)
	Cancel(context.Context, uint64, interface{}) (*Order, error)
	Close(context.Context, uint64) (*Order, error)
	Open(context.Context, uint64) (*Order, error)
	Delete(context.Context, uint64) error

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

// OrderStatus Filter orders by their status.
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#get-orders?status=any
type OrderStatus string

// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#get-orders?status=any
const (
	// OrderStatusOpen Show only open orders.
	OrderStatusOpen OrderStatus = "open"

	// OrderStatusClosed Show only closed orders.
	OrderStatusClosed OrderStatus = "closed"

	// OrderStatusCancelled Show only cancelled orders.
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusAny Show orders of any status, including archived orders.
	OrderStatusAny OrderStatus = "any"
)

// OrderFulfillmentStatus Filter orders by their fulfillment status.
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#get-orders?status=any
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
type OrderFulfillmentStatus string

const (
	// OrderFulfillmentStatusShipped Show orders that have been shipped. Returns orders with `fulfillment_status` of `fulfilled`.
	OrderFulfillmentStatusShipped OrderFulfillmentStatus = "shipped"

	// OrderFulfillmentStatusPartial Show partially shipped orders.
	OrderFulfillmentStatusPartial OrderFulfillmentStatus = "partial"

	// OrderFulfillmentStatusUnshipped Show orders that have not yet been shipped. Returns orders with `fulfillment_status` of `null`.
	OrderFulfillmentStatusUnshipped OrderFulfillmentStatus = "unshipped"

	// OrderFulfillmentStatusAny Show orders of any fulfillment status.
	OrderFulfillmentStatusAny OrderFulfillmentStatus = "any"

	// OrderFulfillmentStatusUnfulfilled Returns orders with `fulfillment_status` of `null` or `partial`.
	OrderFulfillmentStatusUnfulfilled OrderFulfillmentStatus = "unfulfilled"

	// OrderFulfillmentStatusFulfilled `fulfilled` used to be an acceptable value? Was it deprecated? It isn't noted
	// in the Shopify docs at the provided URL, but it was used in tests and still
	// seems to function.
	OrderFulfillmentStatusFulfilled OrderFulfillmentStatus = "fulfilled"
)

// OrderFinancialStatus Filter orders by their financial status.
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#get-orders?status=any
type OrderFinancialStatus string

const (
	// OrderFinancialStatusAuthorized Show only authorized orders.
	OrderFinancialStatusAuthorized OrderFinancialStatus = "authorized"

	// OrderFinancialStatusPending Show only pending orders.
	OrderFinancialStatusPending OrderFinancialStatus = "pending"

	// OrderFinancialStatusPaid Show only paid orders.
	OrderFinancialStatusPaid OrderFinancialStatus = "paid"

	// OrderFinancialStatusPartiallyPaid Show only partially paid orders.
	OrderFinancialStatusPartiallyPaid OrderFinancialStatus = "partially_paid"

	// OrderFinancialStatusRefunded Show only refunded orders.
	OrderFinancialStatusRefunded OrderFinancialStatus = "refunded"

	// OrderFinancialStatusVoided Show only voided orders.
	OrderFinancialStatusVoided OrderFinancialStatus = "voided"

	// OrderFinancialStatusPartiallyRefunded Show only partially refunded orders.
	OrderFinancialStatusPartiallyRefunded OrderFinancialStatus = "partially_refunded"

	// OrderFinancialStatusAny Show orders of any financial status.
	OrderFinancialStatusAny OrderFinancialStatus = "any"

	// OrderFinancialStatusUnpaid Show authorized and partially paid orders.
	OrderFinancialStatusUnpaid OrderFinancialStatus = "unpaid"
)

// OrderCancelReason The reason why the order was canceled.
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#post-orders-order-id-cancel
type OrderCancelReason string

const (
	// OrderCancelReasonCustomer The customer canceled the order.
	OrderCancelReasonCustomer OrderCancelReason = "customer"

	// OrderCancelReasonFraud The order was fraudulent.
	OrderCancelReasonFraud OrderCancelReason = "fraud"

	// OrderCancelReasonInventory Items in the order were not in inventory.
	OrderCancelReasonInventory OrderCancelReason = "inventory"

	// OrderCancelReasonDeclined The payment was declined.
	OrderCancelReasonDeclined OrderCancelReason = "declined"

	// OrderCancelReasonOther Cancelled for some other reason.
	OrderCancelReasonOther OrderCancelReason = "other"
)

// DiscountAllocationMethod
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
type DiscountAllocationMethod string

const (
	// DiscountAllocationMethodAcross The value is spread across all entitled lines.
	DiscountAllocationMethodAcross DiscountAllocationMethod = "across"

	// DiscountAllocationMethodEach The value is applied onto every entitled line.
	DiscountAllocationMethodEach DiscountAllocationMethod = "each"

	// DiscountAllocationMethodOne The value is applied onto a single line.
	DiscountAllocationMethodOne DiscountAllocationMethod = "one"
)

// DiscountTargetSelection The lines on the order, of the type defined by `target_type`, that the discount is allocated over
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
type DiscountTargetSelection string

const (
	// DiscountTargetSelectionAll The discount is allocated onto all lines
	DiscountTargetSelectionAll DiscountTargetSelection = "all"

	// DiscountTargetSelectionEntitled The discount is allocated only onto lines it is entitled for.
	DiscountTargetSelectionEntitled DiscountTargetSelection = "entitled"

	// DiscountTargetSelectionExplicit The discount is allocated onto explicitly selected lines.
	DiscountTargetSelectionExplicit DiscountTargetSelection = "explicit"
)

// DiscountTargetType The type of line on the order that the discount is applicable on
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
type DiscountTargetType string

const (
	// DiscountTargetTypeLineItem The discount applies to line items.
	DiscountTargetTypeLineItem DiscountTargetType = "line_item"

	// DiscountTargetTypeShippingLine The discount applies to shipping lines.
	DiscountTargetTypeShippingLine DiscountTargetType = "shipping_line"
)

// DiscountType The discount application type
type DiscountType string

const (
	// DiscountTypeAutomatic The discount was applied automatically, such as by a Buy X Get Y automatic discount.
	DiscountTypeAutomatic DiscountType = "automatic"

	// DiscountTypeDiscountCode The discount was applied by a discount code.
	DiscountTypeDiscountCode DiscountType = "discount_code"

	// DiscountTypeManual The discount was manually applied by the merchant (for example, by using an app or creating a draft order).
	DiscountTypeManual DiscountType = "manual"

	// DiscountTypeScript The discount was applied by a Shopify Script.
	DiscountTypeScript DiscountType = "script"
)

// DiscountValueType The type of value of the discount
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#resource-object
type DiscountValueType string

const (
	// DiscountValueTypeFixedAmount A fixed amount discount value in the currency of the order.
	DiscountValueTypeFixedAmount DiscountValueType = "fixed_amount"

	// DiscountValueTypePercentage A percentage discount value.
	DiscountValueTypePercentage DiscountValueType = "percentage"
)

// OrderCountOptions A struct for all available order count options
type OrderCountOptions struct {
	Page              int                    `url:"page,omitempty"`
	Limit             int                    `url:"limit,omitempty"`
	SinceId           uint64                 `url:"since_id,omitempty"`
	CreatedAtMin      time.Time              `url:"created_at_min,omitempty"`
	CreatedAtMax      time.Time              `url:"created_at_max,omitempty"`
	UpdatedAtMin      time.Time              `url:"updated_at_min,omitempty"`
	UpdatedAtMax      time.Time              `url:"updated_at_max,omitempty"`
	Order             string                 `url:"order,omitempty"`
	Fields            string                 `url:"fields,omitempty"`
	Status            OrderStatus            `url:"status,omitempty"`
	FinancialStatus   OrderFinancialStatus   `url:"financial_status,omitempty"`
	FulfillmentStatus OrderFulfillmentStatus `url:"fulfillment_status,omitempty"`
}

// OrderListOptions A struct for all available order list options.
// See: https://help.shopify.com/api/reference/order#index
type OrderListOptions struct {
	ListOptions
	Status            OrderStatus            `url:"status,omitempty"`
	FinancialStatus   OrderFinancialStatus   `url:"financial_status,omitempty"`
	FulfillmentStatus OrderFulfillmentStatus `url:"fulfillment_status,omitempty"`
	ProcessedAtMin    time.Time              `url:"processed_at_min,omitempty"`
	ProcessedAtMax    time.Time              `url:"processed_at_max,omitempty"`
	Order             string                 `url:"order,omitempty"`
}

// OrderCancelOptions A struct of all available order cancel options.
// See: https://help.shopify.com/api/reference/order#index
type OrderCancelOptions struct {
	Amount   *decimal.Decimal `json:"amount,omitempty"`
	Currency string           `json:"currency,omitempty"`
	Restock  bool             `json:"restock,omitempty"`
	Reason   string           `json:"reason,omitempty"`
	Email    bool             `json:"email,omitempty"`
	Refund   *Refund          `json:"refund,omitempty"`
}

// OrderInventoryBehaviour The behaviour to use when updating inventory.
//
// https://shopify.dev/docs/api/admin-rest/2024-01/resources/order#post-orders
type OrderInventoryBehaviour string

const (
	// OrderInventoryBehaviourBypass Do not claim inventory.
	OrderInventoryBehaviourBypass OrderInventoryBehaviour = "bypass"

	// OrderInventoryBehaviourDecrementIgnoringPolicy Ignore the product's inventory policy and claim inventory.
	OrderInventoryBehaviourDecrementIgnoringPolicy OrderInventoryBehaviour = "decrement_ignoring_policy"

	// OrderInventoryBehaviourDecrementObeyingPolicy Follow the product's inventory policy and claim inventory, if possible.
	OrderInventoryBehaviourDecrementObeyingPolicy OrderInventoryBehaviour = "decrement_obeying_policy"
)

// Order represents a Shopify order
//
// Docs:
//
//   - [The Order resource]
//   - [Retrieve a specific order]
//
// [The Order resource]: https://shopify.dev/docs/api/admin-rest/2024-04/resources/order#resource-object
// [Retrieve a specific order]: https://shopify.dev/docs/api/admin-rest/2024-04/resources/order#get-orders-order-id
type Order struct {
	Id                       uint64                  `json:"id,omitempty"`
	Name                     string                  `json:"name,omitempty"`
	Email                    string                  `json:"email,omitempty"`
	CreatedAt                *time.Time              `json:"created_at,omitempty"`
	UpdatedAt                *time.Time              `json:"updated_at,omitempty"`
	CancelledAt              *time.Time              `json:"cancelled_at,omitempty"`
	ClosedAt                 *time.Time              `json:"closed_at,omitempty"`
	ProcessedAt              *time.Time              `json:"processed_at,omitempty"`
	Customer                 *Customer               `json:"customer,omitempty"`
	BillingAddress           *Address                `json:"billing_address,omitempty"`
	ShippingAddress          *Address                `json:"shipping_address,omitempty"`
	Currency                 string                  `json:"currency,omitempty"`
	TotalPrice               *decimal.Decimal        `json:"total_price,omitempty"`
	TotalPriceSet            *AmountSet              `json:"total_price_set,omitempty"`
	TotalShippingPriceSet    *AmountSet              `json:"total_shipping_price_set,omitempty"`
	CurrentTotalPrice        *decimal.Decimal        `json:"current_total_price,omitempty"`
	SubtotalPrice            *decimal.Decimal        `json:"subtotal_price,omitempty"`
	CurrentSubtotalPrice     *decimal.Decimal        `json:"current_subtotal_price,omitempty"`
	TotalDiscounts           *decimal.Decimal        `json:"total_discounts,omitempty"`
	TotalDiscountSet         *AmountSet              `json:"total_discount_set,omitempty"`
	CurrentTotalDiscounts    *decimal.Decimal        `json:"current_total_discounts,omitempty"`
	CurrentTotalDiscountsSet *AmountSet              `json:"current_total_discounts_set,omitempty"`
	TotalLineItemsPrice      *decimal.Decimal        `json:"total_line_items_price,omitempty"`
	TaxesIncluded            bool                    `json:"taxes_included,omitempty"`
	TotalTax                 *decimal.Decimal        `json:"total_tax,omitempty"`
	TotalTaxSet              *AmountSet              `json:"total_tax_set,omitempty"`
	CurrentTotalTax          *decimal.Decimal        `json:"current_total_tax,omitempty"`
	CurrentTotalTaxSet       *AmountSet              `json:"current_total_tax_set,omitempty"`
	TaxLines                 []TaxLine               `json:"tax_lines,omitempty"`
	TotalWeight              int                     `json:"total_weight,omitempty"`
	TotalTipReceived         string                  `json:"total_tip_received,omitempty"`
	FinancialStatus          OrderFinancialStatus    `json:"financial_status,omitempty"`
	Fulfillments             []Fulfillment           `json:"fulfillments,omitempty"`
	FulfillmentStatus        OrderFulfillmentStatus  `json:"fulfillment_status,omitempty"`
	Token                    string                  `json:"token,omitempty"`
	CartToken                string                  `json:"cart_token,omitempty"`
	Number                   int                     `json:"number,omitempty"`
	OrderNumber              int                     `json:"order_number,omitempty"`
	Note                     string                  `json:"note,omitempty"`
	Test                     bool                    `json:"test,omitempty"`
	BrowserIp                string                  `json:"browser_ip,omitempty"`
	BuyerAcceptsMarketing    bool                    `json:"buyer_accepts_marketing,omitempty"`
	CancelReason             OrderCancelReason       `json:"cancel_reason,omitempty"`
	NoteAttributes           []NoteAttribute         `json:"note_attributes,omitempty"`
	DiscountCodes            []DiscountCode          `json:"discount_codes,omitempty"`
	DiscountApplications     []DiscountApplication   `json:"discount_applications,omitempty"`
	LineItems                []LineItem              `json:"line_items,omitempty"`
	ShippingLines            []ShippingLines         `json:"shipping_lines,omitempty"`
	Transactions             []Transaction           `json:"transactions,omitempty"`
	AppId                    int                     `json:"app_id,omitempty"`
	CustomerLocale           string                  `json:"customer_locale,omitempty"`
	LandingSite              string                  `json:"landing_site,omitempty"`
	ReferringSite            string                  `json:"referring_site,omitempty"`
	SourceName               string                  `json:"source_name,omitempty"`
	ClientDetails            *ClientDetails          `json:"client_details,omitempty"`
	Tags                     string                  `json:"tags,omitempty"`
	LocationId               uint64                  `json:"location_id,omitempty"`
	PaymentGatewayNames      []string                `json:"payment_gateway_names,omitempty"`
	ProcessingMethod         string                  `json:"processing_method,omitempty"`
	Refunds                  []Refund                `json:"refunds,omitempty"`
	UserId                   uint64                  `json:"user_id,omitempty"`
	OrderStatusUrl           string                  `json:"order_status_url,omitempty"`
	Gateway                  string                  `json:"gateway,omitempty"`
	Confirmed                bool                    `json:"confirmed,omitempty"`
	CheckoutToken            string                  `json:"checkout_token,omitempty"`
	Reference                string                  `json:"reference,omitempty"`
	SourceIdentifier         string                  `json:"source_identifier,omitempty"`
	SourceURL                string                  `json:"source_url,omitempty"`
	DeviceId                 uint64                  `json:"device_id,omitempty"`
	Phone                    string                  `json:"phone,omitempty"`
	LandingSiteRef           string                  `json:"landing_site_ref,omitempty"`
	CheckoutId               uint64                  `json:"checkout_id,omitempty"`
	ContactEmail             string                  `json:"contact_email,omitempty"`
	Metafields               []Metafield             `json:"metafields,omitempty"`
	SendReceipt              bool                    `json:"send_receipt,omitempty"`
	SendFulfillmentReceipt   bool                    `json:"send_fulfillment_receipt,omitempty"`
	PresentmentCurrency      string                  `json:"presentment_currency,omitempty"`
	InventoryBehaviour       OrderInventoryBehaviour `json:"inventory_behaviour,omitempty"`
}

type Address struct {
	Id           uint64  `json:"id,omitempty"`
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

type DiscountApplication struct {
	AllocationMethod DiscountAllocationMethod `json:"allocation_method,omitempty"`
	Code             string                   `json:"code"`
	Description      string                   `json:"description"`
	TargetSelection  DiscountTargetSelection  `json:"target_selection"`
	TargetType       DiscountTargetType       `json:"target_type"`
	Title            string                   `json:"title"`
	Type             DiscountType             `json:"type"`
	Value            *decimal.Decimal         `json:"value"`
	ValueType        DiscountValueType        `json:"value_type"`
}

type LineItem struct {
	Id                         uint64                 `json:"id,omitempty"`
	ProductId                  uint64                 `json:"product_id,omitempty"`
	VariantId                  uint64                 `json:"variant_id,omitempty"`
	Quantity                   int                    `json:"quantity,omitempty"`
	CurrentQuantity            int                    `json:"current_quantity,omitempty"`
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
	FulfillmentStatus          OrderFulfillmentStatus `json:"fulfillment_status,omitempty"`
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
	Id                            uint64           `json:"id,omitempty"`
	Title                         string           `json:"title,omitempty"`
	Price                         *decimal.Decimal `json:"price,omitempty"`
	PriceSet                      *AmountSet       `json:"price_set,omitempty"`
	DiscountedPrice               *decimal.Decimal `json:"discounted_price,omitempty"`
	DiscountedPriceSet            *AmountSet       `json:"discounted_price_set,omitempty"`
	Code                          string           `json:"code,omitempty"`
	Source                        string           `json:"source,omitempty"`
	Phone                         string           `json:"phone,omitempty"`
	RequestedFulfillmentServiceId string           `json:"requested_fulfillment_service_id,omitempty"`
	CarrierIdentifier             string           `json:"carrier_identifier,omitempty"`
	TaxLines                      []TaxLine        `json:"tax_lines,omitempty"`
	Handle                        string           `json:"handle,omitempty"`
}

// UnmarshalJSON custom unmarshaller for ShippingLines implemented
// to handle requested_fulfillment_service_id being
// returned as json numbers or json nulls instead of json strings
func (sl *ShippingLines) UnmarshalJSON(data []byte) error {
	type alias ShippingLines
	aux := &struct {
		*alias
		RequestedFulfillmentServiceId interface{} `json:"requested_fulfillment_service_id"`
	}{alias: (*alias)(sl)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	switch aux.RequestedFulfillmentServiceId.(type) {
	case nil:
		sl.RequestedFulfillmentServiceId = ""
	default:
		sl.RequestedFulfillmentServiceId = fmt.Sprintf("%v", aux.RequestedFulfillmentServiceId)
	}

	return nil
}

type TaxLine struct {
	Title string           `json:"title,omitempty"`
	Price *decimal.Decimal `json:"price,omitempty"`
	Rate  *decimal.Decimal `json:"rate,omitempty"`
}

type Transaction struct {
	Id             uint64           `json:"id,omitempty"`
	OrderId        uint64           `json:"order_id,omitempty"`
	Amount         *decimal.Decimal `json:"amount,omitempty"`
	Kind           string           `json:"kind,omitempty"`
	Gateway        string           `json:"gateway,omitempty"`
	Status         string           `json:"status,omitempty"`
	Message        string           `json:"message,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	Test           bool             `json:"test,omitempty"`
	Authorization  string           `json:"authorization,omitempty"`
	Currency       string           `json:"currency,omitempty"`
	LocationId     *int64           `json:"location_id,omitempty"`
	UserId         *int64           `json:"user_id,omitempty"`
	ParentId       *int64           `json:"parent_id,omitempty"`
	DeviceId       *int64           `json:"device_id,omitempty"`
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
	Id               uint64            `json:"id,omitempty"`
	OrderId          uint64            `json:"order_id,omitempty"`
	CreatedAt        *time.Time        `json:"created_at,omitempty"`
	Note             string            `json:"note,omitempty"`
	Restock          bool              `json:"restock,omitempty"`
	UserId           uint64            `json:"user_id,omitempty"`
	RefundLineItems  []RefundLineItem  `json:"refund_line_items,omitempty"`
	Transactions     []Transaction     `json:"transactions,omitempty"`
	OrderAdjustments []OrderAdjustment `json:"order_adjustments,omitempty"`
}

type OrderAdjustment struct {
	Id           uint64              `json:"id,omitempty"`
	OrderId      uint64              `json:"order_id,omitempty"`
	RefundId     uint64              `json:"refund_id,omitempty"`
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
	Id          uint64           `json:"id,omitempty"`
	Quantity    int              `json:"quantity,omitempty"`
	LineItemId  uint64           `json:"line_item_id,omitempty"`
	LineItem    *LineItem        `json:"line_item,omitempty"`
	Subtotal    *decimal.Decimal `json:"subtotal,omitempty"`
	TotalTax    *decimal.Decimal `json:"total_tax,omitempty"`
	SubTotalSet *AmountSet       `json:"subtotal_set,omitempty"`
	TotalTaxSet *AmountSet       `json:"total_tax_set,omitempty"`
}

// List orders
func (s *OrderServiceOp) List(ctx context.Context, options interface{}) ([]Order, error) {
	orders, _, err := s.ListWithPagination(ctx, options)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// ListAll Lists all orders, iterating over pages
func (s *OrderServiceOp) ListAll(ctx context.Context, options interface{}) ([]Order, error) {
	collector := []Order{}

	for {
		entities, pagination, err := s.ListWithPagination(ctx, options)

		if err != nil {
			return collector, err
		}

		collector = append(collector, entities...)

		if pagination.NextPageOptions == nil {
			break
		}

		options = pagination.NextPageOptions
	}

	return collector, nil
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
func (s *OrderServiceOp) Get(ctx context.Context, orderId uint64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, orderId)
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
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, order.Id)
	wrappedData := OrderResource{Order: &order}
	resource := new(OrderResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Order, err
}

// Cancel order
func (s *OrderServiceOp) Cancel(ctx context.Context, orderId uint64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d/cancel.json", ordersBasePath, orderId)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, options, resource)
	return resource.Order, err
}

// Close order
func (s *OrderServiceOp) Close(ctx context.Context, orderId uint64) (*Order, error) {
	path := fmt.Sprintf("%s/%d/close.json", ordersBasePath, orderId)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, nil, resource)
	return resource.Order, err
}

// Open order
func (s *OrderServiceOp) Open(ctx context.Context, orderId uint64) (*Order, error) {
	path := fmt.Sprintf("%s/%d/open.json", ordersBasePath, orderId)
	resource := new(OrderResource)
	err := s.client.Post(ctx, path, nil, resource)
	return resource.Order, err
}

// Delete order
func (s *OrderServiceOp) Delete(ctx context.Context, orderId uint64) error {
	path := fmt.Sprintf("%s/%d.json", ordersBasePath, orderId)
	err := s.client.Delete(ctx, path)
	return err
}

// List metafields for an order
func (s *OrderServiceOp) ListMetafields(ctx context.Context, orderId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.List(ctx, options)
}

// Count metafields for an order
func (s *OrderServiceOp) CountMetafields(ctx context.Context, orderId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for an order
func (s *OrderServiceOp) GetMetafield(ctx context.Context, orderId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for an order
func (s *OrderServiceOp) CreateMetafield(ctx context.Context, orderId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for an order
func (s *OrderServiceOp) UpdateMetafield(ctx context.Context, orderId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.Update(ctx, metafield)
}

// Delete an existing metafield for an order
func (s *OrderServiceOp) DeleteMetafield(ctx context.Context, orderId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return metafieldService.Delete(ctx, metafieldId)
}

// List fulfillments for an order
func (s *OrderServiceOp) ListFulfillments(ctx context.Context, orderId uint64, options interface{}) ([]Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.List(ctx, options)
}

// Count fulfillments for an order
func (s *OrderServiceOp) CountFulfillments(ctx context.Context, orderId uint64, options interface{}) (int, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Count(ctx, options)
}

// Get individual fulfillment for an order
func (s *OrderServiceOp) GetFulfillment(ctx context.Context, orderId uint64, fulfillmentId uint64, options interface{}) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Get(ctx, fulfillmentId, options)
}

// Create a new fulfillment for an order
func (s *OrderServiceOp) CreateFulfillment(ctx context.Context, orderId uint64, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Create(ctx, fulfillment)
}

// Update an existing fulfillment for an order
func (s *OrderServiceOp) UpdateFulfillment(ctx context.Context, orderId uint64, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Update(ctx, fulfillment)
}

// Complete an existing fulfillment for an order
func (s *OrderServiceOp) CompleteFulfillment(ctx context.Context, orderId uint64, fulfillmentId uint64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Complete(ctx, fulfillmentId)
}

// Transition an existing fulfillment for an order
func (s *OrderServiceOp) TransitionFulfillment(ctx context.Context, orderId uint64, fulfillmentId uint64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Transition(ctx, fulfillmentId)
}

// Cancel an existing fulfillment for an order
func (s *OrderServiceOp) CancelFulfillment(ctx context.Context, orderId uint64, fulfillmentId uint64) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceId: orderId}
	return fulfillmentService.Cancel(ctx, fulfillmentId)
}
