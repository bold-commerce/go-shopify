package goshopify

import (
	"fmt"
	"time"
)

// FulfillmentOrderService is an interface for interfacing with the fulfillment
// order endpoints  of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/2023-01/resources/fulfillmentorder#resource-object
type FulfillmentOrderService interface {
	List(int64, interface{}) ([]FulfillmentOrder, error)
	Get(int64, interface{}) (*FulfillmentOrder, error)
	Cancel(int64) (*FulfillmentOrder, error)
	Close(int64, string) (*FulfillmentOrder, error)
	Hold(int64, bool, FulfillmentOrderHoldReason, string) (*FulfillmentOrder, error)
	Open(int64) (*FulfillmentOrder, error)
	ReleaseHold(int64) (*FulfillmentOrder, error)
	Reschedule(int64) (*FulfillmentOrder, error)
	SetDeadline([]int64, time.Time) error
	Move(int64, FulfillmentOrderMoveRequest, interface{}) (*FulfillmentOrderMoveResource, error)
}

// FulfillmentOrderHoldReason represents the reason for a fulfillment hold
type FulfillmentOrderHoldReason string

const (
	HoldReasonAwaitingPayment  FulfillmentOrderHoldReason = "awaiting_payment"
	HoldReasonHighRiskOfFraud                             = "high_risk_of_fraud"
	HoldReasonIncorrectAddress                            = "incorrect_address"
	HoldReasonOutOfStock                                  = "inventory_out_of_stock"
	HoldReasonOther                                       = "other"
)

// FulfillmentOrderServiceOp handles communication with the fulfillment order
// related methods of the Shopify API.
type FulfillmentOrderServiceOp struct {
	client *Client
}

type FulfillmentOrderLineItemQuantity struct {
	Id       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
}

type FulfillmentOrderMoveRequest struct {
	NewLocationId int64                              `json:"new_location_id"`
	LineItems     []FulfillmentOrderLineItemQuantity `json:"fulfillment_order_line_items,omitempty"`
}

// FulfillmentOrderDeliveryMethod represents a delivery method for a FulfillmentOrder
type FulfillmentOrderDeliveryMethod struct {
	Id                  int64     `json:"id,omitempty"`
	MethodType          string    `json:"method_type,omitempty"`
	MinDeliveryDateTime time.Time `json:"min_delivery_date_time,omitempty"`
	MaxDeliveryDateTime time.Time `json:"max_delivery_date_time,omitempty"`
}

// FulfillmentOrderDestination represents a destination for a FulfillmentOrder
type FulfillmentOrderDestination struct {
	Id        int64  `json:"id,omitempty"`
	Address1  string `json:"address1,omitempty"`
	Address2  string `json:"address2,omitempty"`
	City      string `json:"city,omitempty"`
	Company   string `json:"company,omitempty"`
	Country   string `json:"country,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Province  string `json:"province,omitempty"`
	Zip       string `json:"zip,omitempty"`
}

// FulfillmentOrderHold represents a fulfillment hold for a FulfillmentOrder
type FulfillmentOrderHold struct {
	Reason      FulfillmentOrderHoldReason `json:"reason,omitempty"`
	ReasonNotes string                     `json:"reason_notes,omitempty"`
}

// FulfillmentOrderInternationalDuties represents an InternationalDuty for a FulfillmentOrder
type FulfillmentOrderInternationalDuties struct {
	IncoTerm string `json:"incoterm,omitempty"`
}

// FulfillmentOrderLineItem represents a line item for a FulfillmentOrder
type FulfillmentOrderLineItem struct {
	Id                  int64 `json:"id,omitempty"`
	ShopId              int64 `json:"shop_id,omitempty"`
	FulfillmentOrderId  int64 `json:"fulfillment_order_id,omitempty"`
	LineItemId          int64 `json:"line_item_id,omitempty"`
	InventoryItemId     int64 `json:"inventory_item_id,omitempty"`
	Quantity            int64 `json:"quantity,omitempty"`
	FulfillableQuantity int64 `json:"fulfillable_quantity,omitempty"`
	VariantId           int64 `json:"variant_id,omitempty"`
}

// FulfillmentOrderMerchantRequest represents a MerchantRequest for a FulfillmentOrder
type FulfillmentOrderMerchantRequest struct {
	Message        string `json:"message,omitempty"`
	RequestOptions struct {
		ShippingMethod string    `json:"shipping_method,omitempty"`
		Note           string    `json:"note,omitempty"`
		Date           time.Time `json:"date,omitempty"`
	} `json:"request_options"`
	Kind string `json:"kind,omitempty"`
}

// FulfillmentOrderAssignedLocation represents an AssignedLocation for a FulfillmentOrder
type FulfillmentOrderAssignedLocation struct {
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	LocationId  int64  `json:"location_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Province    string `json:"province,omitempty"`
	Zip         string `json:"zip,omitempty"`
}

// FulfillmentOrder represents a Shopify Fulfillment Order
type FulfillmentOrder struct {
	Id                  int64                               `json:"id,omitempty"`
	AssignedLocation    FulfillmentOrderAssignedLocation    `json:"assigned_location,omitempty"`
	AssignedLocationId  int64                               `json:"assigned_location_id,omitempty"`
	CreatedAt           time.Time                           `json:"created_at,omitempty"`
	DeliveryMethod      FulfillmentOrderDeliveryMethod      `json:"delivery_method,omitempty"`
	Destination         FulfillmentOrderDestination         `json:"destination,omitempty"`
	FulfillAt           OnlyDate                            `json:"fulfill_at,omitempty"`
	FulfillBy           OnlyDate                            `json:"fulfill_by,omitempty"`
	FulfillmentHolds    []FulfillmentOrderHold              `json:"fulfillment_holds,omitempty"`
	InternationalDuties FulfillmentOrderInternationalDuties `json:"international_duties,omitempty"`
	LineItems           []FulfillmentOrderLineItem          `json:"line_items,omitempty"`
	MerchantRequests    []FulfillmentOrderMerchantRequest   `json:"merchant_requests,omitempty"`
	OrderId             int64                               `json:"order_id,omitempty"`
	RequestStatus       string                              `json:"request_status,omitempty"`
	ShopId              int64                               `json:"shop_id,omitempty"`
	Status              string                              `json:"status,omitempty"`
	SupportedActions    []string                            `json:"supported_actions,omitempty"`
	UpdatedAt           time.Time                           `json:"updated_at,omitempty"`
}

// FulfillmentOrdersResource represents the result from the fulfilment_orders.json endpoint
type FulfillmentOrdersResource struct {
	FulfillmentOrders []FulfillmentOrder `json:"fulfillment_orders"`
}

// FulfillmentOrderResource represents the result from the fulfilment_orders/<id>.json endpoint
type FulfillmentOrderResource struct {
	FulfillmentOrder *FulfillmentOrder `json:"fulfillment_order"`
}

// FulfillmentOrderMoveResource represents the result from the move.json endpoint
type FulfillmentOrderMoveResource struct {
	OriginalFulfillmentOrder FulfillmentOrder `json:"original_fulfillment_order"`
	MovedFulfillmentOrder    FulfillmentOrder `json:"moved_fulfillment_order"`
}

// FulfillmentOrderPathPrefix returns the prefix for a fulfillmentorder path
func FulfillmentOrderPathPrefix(resource string, resourceID int64) string {
	return fmt.Sprintf("%s/%d", resource, resourceID)
}

// List gets FulfillmentOrder items for an order
func (s *FulfillmentOrderServiceOp) List(orderId int64, options interface{}) ([]FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("orders", orderId)
	path := fmt.Sprintf("%s/fulfillment_orders.json", prefix)
	resource := new(FulfillmentOrdersResource)
	err := s.client.Get(path, resource, options)
	return resource.FulfillmentOrders, err
}

// Get gets an individual fulfillment order
func (s *FulfillmentOrderServiceOp) Get(fulfillmentID int64, options interface{}) (*FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Get(path, resource, options)
	return resource.FulfillmentOrder, err
}

// Cancel cancels a fulfillment order
func (s *FulfillmentOrderServiceOp) Cancel(fulfillmentID int64) (*FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/cancel.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, nil, resource)
	return resource.FulfillmentOrder, err
}

// Close marks a fulfillment order as incomplete with an optional message
func (s *FulfillmentOrderServiceOp) Close(fulfillmentID int64, message string) (*FulfillmentOrder, error) {
	type closeRequest struct {
		Message string `json:"message,omitempty"`
	}
	req := closeRequest{
		Message: message,
	}
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/close.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, req, resource)
	return resource.FulfillmentOrder, err
}

// Hold applies a fulfillment hold on an open fulfillment order
func (s *FulfillmentOrderServiceOp) Hold(fulfillmentID int64, notify bool, reason FulfillmentOrderHoldReason, notes string) (*FulfillmentOrder, error) {
	type holdRequest struct {
		Reason         FulfillmentOrderHoldReason `json:"reason"`
		ReasonNotes    string                     `json:"reason_notes,omitempty"`
		NotifyMerchant bool                       `json:"notify_merchant"`
	}
	type wrappedRequest struct {
		FulfillmentHold holdRequest `json:"fulfillment_hold"`
	}
	req := wrappedRequest{
		FulfillmentHold: holdRequest{
			Reason:         reason,
			ReasonNotes:    notes,
			NotifyMerchant: notify,
		},
	}
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/hold.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, req, resource)
	return resource.FulfillmentOrder, err
}

// Open marks the fulfillment order as open
func (s *FulfillmentOrderServiceOp) Open(fulfillmentID int64) (*FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/open.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, nil, resource)
	return resource.FulfillmentOrder, err
}

// ReleaseHold releases the fulfillment hold on a fulfillment order
func (s *FulfillmentOrderServiceOp) ReleaseHold(fulfillmentID int64) (*FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/release_hold.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, nil, resource)
	return resource.FulfillmentOrder, err
}

// Reschedule reschedules the fulfill_at time of a scheduled fulfillment order
func (s *FulfillmentOrderServiceOp) Reschedule(fulfillmentID int64) (*FulfillmentOrder, error) {
	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/reschedule.json", prefix)
	resource := new(FulfillmentOrderResource)
	err := s.client.Post(path, nil, resource)
	return resource.FulfillmentOrder, err
}

// SetDeadline sets deadline for fulfillment orders
func (s *FulfillmentOrderServiceOp) SetDeadline(fulfillmentIDs []int64, deadline time.Time) error {
	type deadlineRequest struct {
		FulfillmentOrderIds []int64   `json:"fulfillment_order_ids"`
		FulfillmentDeadline time.Time `json:"fulfillment_deadline"`
	}

	req := deadlineRequest{
		FulfillmentOrderIds: fulfillmentIDs,
		FulfillmentDeadline: deadline,
	}
	path := "fulfillment_orders/set_fulfillment_orders_deadline.json"
	err := s.client.Post(path, req, nil)
	return err
}

// Move moves a fulfillment order to a new location
func (s *FulfillmentOrderServiceOp) Move(fulfillmentID int64, moveRequest FulfillmentOrderMoveRequest, options interface{}) (*FulfillmentOrderMoveResource, error) {
	type request struct {
		FulfillmentOrder FulfillmentOrderMoveRequest `json:"fulfillment_order"`
	}
	wrappedRequest := request{
		FulfillmentOrder: moveRequest,
	}

	prefix := FulfillmentOrderPathPrefix("fulfillment_orders", fulfillmentID)
	path := fmt.Sprintf("%s/move.json", prefix)
	resource := new(FulfillmentOrderMoveResource)
	err := s.client.Post(path, wrappedRequest, resource)
	return resource, err
}
