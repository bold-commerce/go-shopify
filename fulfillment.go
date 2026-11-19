package goshopify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ErrFulfillmentOperationUnsupported indicates that Shopify removed the legacy
// complete/open fulfillment operations. Fulfillment orders have a different lifecycle.
var ErrFulfillmentOperationUnsupported = errors.New("legacy fulfillment complete/open operations are unsupported; use the fulfillment-order workflow")

// ErrFulfillmentRequiresOrder indicates that a fulfillment read needs an order ID.
var ErrFulfillmentRequiresOrder = errors.New("fulfillment reads require an order ID; use Order.ListFulfillments, Order.CountFulfillments or Order.GetFulfillment")

// ErrFulfillmentRequiresFulfillmentOrders indicates that creation needs explicit
// fulfillment-order IDs, rather than order line-item IDs or a location alone.
var ErrFulfillmentRequiresFulfillmentOrders = errors.New("fulfillment creation requires LineItemsByFulfillmentOrder; obtain IDs with FulfillmentOrder.List")

// FulfillmentService is an interface for interfacing with the fulfillment endpoints
// of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/latest/resources/fulfillment
type FulfillmentService interface {
	List(context.Context, interface{}) ([]Fulfillment, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Fulfillment, error)
	Create(context.Context, Fulfillment) (*Fulfillment, error)
	Update(context.Context, Fulfillment) (*Fulfillment, error)
	Complete(context.Context, uint64) (*Fulfillment, error)
	Transition(context.Context, uint64) (*Fulfillment, error)
	Cancel(context.Context, uint64) (*Fulfillment, error)
}

// FulfillmentsService is an interface for other Shopify resources
// to interface with the fulfillment endpoints of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/latest/resources/fulfillment
type FulfillmentsService interface {
	ListFulfillments(context.Context, uint64, interface{}) ([]Fulfillment, error)
	ListFulfillmentsWithPagination(context.Context, uint64, interface{}) ([]Fulfillment, *Pagination, error)
	ListFulfillmentsAll(context.Context, uint64, interface{}) ([]Fulfillment, error)
	CountFulfillments(context.Context, uint64, interface{}) (int, error)
	GetFulfillment(context.Context, uint64, uint64, interface{}) (*Fulfillment, error)
	CreateFulfillment(context.Context, uint64, Fulfillment) (*Fulfillment, error)
	UpdateFulfillment(context.Context, uint64, Fulfillment) (*Fulfillment, error)
	CompleteFulfillment(context.Context, uint64, uint64) (*Fulfillment, error)
	TransitionFulfillment(context.Context, uint64, uint64) (*Fulfillment, error)
	CancelFulfillment(context.Context, uint64, uint64) (*Fulfillment, error)
}

// FulfillmentServiceOp handles communication with the fulfillment
// related methods of the Shopify API.
type FulfillmentServiceOp struct {
	client     *Client
	resource   string
	resourceId uint64
}

// Fulfillment represents a Shopify fulfillment.
type Fulfillment struct {
	Id                          uint64                       `json:"id,omitempty"`
	OrderId                     uint64                       `json:"order_id,omitempty"`
	LocationId                  uint64                       `json:"location_id,omitempty"`
	Status                      string                       `json:"status,omitempty"`
	CreatedAt                   *time.Time                   `json:"created_at,omitempty"`
	Service                     string                       `json:"service,omitempty"`
	UpdatedAt                   *time.Time                   `json:"updated_at,omitempty"`
	TrackingCompany             string                       `json:"tracking_company,omitempty"`
	ShipmentStatus              string                       `json:"shipment_status,omitempty"`
	TrackingInfo                FulfillmentTrackingInfo      `json:"tracking_info,omitempty"`
	TrackingNumber              string                       `json:"tracking_number,omitempty"`
	TrackingNumbers             []string                     `json:"tracking_numbers,omitempty"`
	TrackingUrl                 string                       `json:"tracking_url,omitempty"`
	TrackingUrls                []string                     `json:"tracking_urls,omitempty"`
	Receipt                     Receipt                      `json:"receipt,omitempty"`
	LineItems                   []LineItem                   `json:"line_items,omitempty"`
	LineItemsByFulfillmentOrder []LineItemByFulfillmentOrder `json:"line_items_by_fulfillment_order,omitempty"`
	NotifyCustomer              bool                         `json:"notify_customer"`
	AdminGraphqlApiId           string                       `json:"admin_graphql_api_id,omitempty"`
	Name                        string                       `json:"name,omitempty"`
	// OriginAddress is the address sent on create, or the first address returned by Shopify.
	OriginAddress *Address `json:"origin_address,omitempty"`
	// AdditionalOriginAddresses holds the remaining entries when Shopify returns an array.
	AdditionalOriginAddresses []Address `json:"-"`
	// Message is sent when creating a fulfillment.
	Message string `json:"message,omitempty"`
}

// UnmarshalJSON accepts the object and array forms of origin_address.
func (f *Fulfillment) UnmarshalJSON(data []byte) error {
	type plain Fulfillment
	var wire struct {
		*plain
		Origin json.RawMessage `json:"origin_address"`
	}
	wire.plain = (*plain)(f)
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if len(wire.Origin) == 0 {
		return nil
	}
	if wire.Origin[0] != '[' {
		f.AdditionalOriginAddresses = nil
		return json.Unmarshal(wire.Origin, &f.OriginAddress)
	}
	var origins []Address
	if err := json.Unmarshal(wire.Origin, &origins); err != nil {
		return err
	}
	f.OriginAddress = nil
	if len(origins) > 0 {
		f.OriginAddress = &origins[0]
		origins = origins[1:]
	}
	f.AdditionalOriginAddresses = origins
	return nil
}

// MarshalJSON preserves an origin_address array when re-encoding a response.
func (f Fulfillment) MarshalJSON() ([]byte, error) {
	type plain Fulfillment
	origins := f.AdditionalOriginAddresses
	if origins == nil {
		return json.Marshal(plain(f))
	}
	if f.OriginAddress != nil {
		origins = append([]Address{*f.OriginAddress}, origins...)
	}
	return json.Marshal(struct {
		plain
		Origin []Address `json:"origin_address"`
	}{plain(f), origins})
}

// FulfillmentTrackingInfo represents the tracking information used to create a Fulfillment.
// https://shopify.dev/docs/api/admin-rest/2023-01/resources/fulfillment#post-fulfillments
type FulfillmentTrackingInfo struct {
	Company string `json:"company,omitempty"`
	Number  string `json:"number,omitempty"`
	Url     string `json:"url,omitempty"`
}

// LineItemByFulfillmentOrder represents the FulfillmentOrders (and optionally the items) used to create a Fulfillment.
// https://shopify.dev/docs/api/admin-rest/2023-01/resources/fulfillment#post-fulfillments
type LineItemByFulfillmentOrder struct {
	FulfillmentOrderId        uint64                                   `json:"fulfillment_order_id,omitempty"`
	FulfillmentOrderLineItems []LineItemByFulfillmentOrderItemQuantity `json:"fulfillment_order_line_items,omitempty"`
}

// LineItemByFulfillmentOrderItemQuantity represents the quantity to fulfill for one item.
type LineItemByFulfillmentOrderItemQuantity struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
}

// Receipt represents a Shopify receipt.
type Receipt struct {
	TestCase      bool   `json:"testcase,omitempty"`
	Authorization string `json:"authorization,omitempty"`
}

// FulfillmentResource represents the result from the fulfillments/X.json endpoint
type FulfillmentResource struct {
	Fulfillment *Fulfillment `json:"fulfillment"`
}

// FulfillmentsResource represents the result from the fullfilments.json endpoint
type FulfillmentsResource struct {
	Fulfillments []Fulfillment `json:"fulfillments"`
}

// List lists fulfillments for an order. Use Order.ListFulfillments on a Client.
func (s *FulfillmentServiceOp) List(ctx context.Context, options interface{}) ([]Fulfillment, error) {
	if s.resource != ordersResourceName || s.resourceId == 0 {
		return nil, ErrFulfillmentRequiresOrder
	}
	prefix := FulfillmentPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s.json", prefix)
	resource := new(FulfillmentsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Fulfillments, err
}

// Count counts fulfillments for an order. Use Order.CountFulfillments on a Client.
func (s *FulfillmentServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	if s.resource != ordersResourceName || s.resourceId == 0 {
		return 0, ErrFulfillmentRequiresOrder
	}
	prefix := FulfillmentPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s/count.json", prefix)
	return s.client.Count(ctx, path, options)
}

// Get retrieves a fulfillment from an order. Use Order.GetFulfillment on a Client.
func (s *FulfillmentServiceOp) Get(ctx context.Context, fulfillmentId uint64, options interface{}) (*Fulfillment, error) {
	if s.resource != ordersResourceName || s.resourceId == 0 {
		return nil, ErrFulfillmentRequiresOrder
	}
	prefix := FulfillmentPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s/%d.json", prefix, fulfillmentId)
	resource := new(FulfillmentResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Fulfillment, err
}

// Create creates a fulfillment using LineItemsByFulfillmentOrder. For a partial
// fulfillment, use fulfillment-order line-item IDs, not order line-item IDs.
// LocationId is not sent: Shopify uses the fulfillment orders' assigned location.
// Message and OriginAddress are optional. OriginAddress requires CountryCode.
func (s *FulfillmentServiceOp) Create(ctx context.Context, fulfillment Fulfillment) (*Fulfillment, error) {
	if len(fulfillment.LineItemsByFulfillmentOrder) == 0 {
		return nil, ErrFulfillmentRequiresFulfillmentOrders
	}
	if len(fulfillment.LineItems) != 0 {
		return nil, errors.New("use LineItemsByFulfillmentOrder for partial fulfillment; legacy LineItems cannot be sent")
	}
	for _, order := range fulfillment.LineItemsByFulfillmentOrder {
		if order.FulfillmentOrderId == 0 {
			return nil, ErrFulfillmentRequiresFulfillmentOrders
		}
	}
	tracking, err := fulfillmentTrackingInfo(fulfillment)
	if err != nil {
		return nil, err
	}
	if s.resource == ordersResourceName {
		orders, err := s.client.FulfillmentOrder.List(ctx, s.resourceId, nil)
		if err != nil {
			return nil, err
		}
		orderIDs := make(map[uint64]bool, len(orders))
		for _, order := range orders {
			orderIDs[order.Id] = true
		}
		for _, order := range fulfillment.LineItemsByFulfillmentOrder {
			if !orderIDs[order.FulfillmentOrderId] {
				return nil, fmt.Errorf("fulfillment order %d does not belong to order %d or is not accessible", order.FulfillmentOrderId, s.resourceId)
			}
		}
	}
	wrappedData := struct {
		Fulfillment struct {
			LineItemsByFulfillmentOrder []LineItemByFulfillmentOrder `json:"line_items_by_fulfillment_order"`
			TrackingInfo                FulfillmentTrackingInfo      `json:"tracking_info"`
			NotifyCustomer              bool                         `json:"notify_customer"`
			Message                     string                       `json:"message,omitempty"`
			OriginAddress               *Address                     `json:"origin_address,omitempty"`
		} `json:"fulfillment"`
	}{}
	wrappedData.Fulfillment.LineItemsByFulfillmentOrder = fulfillment.LineItemsByFulfillmentOrder
	wrappedData.Fulfillment.TrackingInfo = tracking
	wrappedData.Fulfillment.NotifyCustomer = fulfillment.NotifyCustomer
	wrappedData.Fulfillment.Message = fulfillment.Message
	wrappedData.Fulfillment.OriginAddress = fulfillment.OriginAddress
	resource := new(FulfillmentResource)
	err = s.client.Post(ctx, "fulfillments.json", wrappedData, resource)
	return resource.Fulfillment, err
}

// Update updates tracking information and the customer-notification flag.
// Other fulfillment fields are not sent. TrackingInfo takes precedence over the
// legacy TrackingCompany, TrackingNumber(s) and TrackingUrl(s) fields.
func (s *FulfillmentServiceOp) Update(ctx context.Context, fulfillment Fulfillment) (*Fulfillment, error) {
	tracking, err := fulfillmentTrackingInfo(fulfillment)
	if err != nil {
		return nil, err
	}
	if err := s.checkOrder(ctx, fulfillment.Id); err != nil {
		return nil, err
	}
	wrappedData := struct {
		Fulfillment struct {
			TrackingInfo   FulfillmentTrackingInfo `json:"tracking_info"`
			NotifyCustomer bool                    `json:"notify_customer"`
		} `json:"fulfillment"`
	}{}
	wrappedData.Fulfillment.TrackingInfo = tracking
	wrappedData.Fulfillment.NotifyCustomer = fulfillment.NotifyCustomer
	path := fmt.Sprintf("fulfillments/%d/update_tracking.json", fulfillment.Id)
	resource := new(FulfillmentResource)
	err = s.client.Post(ctx, path, wrappedData, resource)
	return resource.Fulfillment, err
}

// Complete returns ErrFulfillmentOperationUnsupported without sending a request.
//
// Deprecated: Shopify removed this operation. Create a fulfillment against
// fulfillment orders when items are shipped; there is no equivalent completion call.
func (s *FulfillmentServiceOp) Complete(ctx context.Context, fulfillmentId uint64) (*Fulfillment, error) {
	return nil, ErrFulfillmentOperationUnsupported
}

// Transition returns ErrFulfillmentOperationUnsupported without sending a request.
//
// Deprecated: Shopify removed the pending-to-open fulfillment transition.
// FulfillmentOrder.Open is a different operation and requires a fulfillment-order ID.
func (s *FulfillmentServiceOp) Transition(ctx context.Context, fulfillmentId uint64) (*Fulfillment, error) {
	return nil, ErrFulfillmentOperationUnsupported
}

// Cancel cancels a fulfillment. Order-scoped calls first check order membership.
func (s *FulfillmentServiceOp) Cancel(ctx context.Context, fulfillmentId uint64) (*Fulfillment, error) {
	if err := s.checkOrder(ctx, fulfillmentId); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("fulfillments/%d/cancel.json", fulfillmentId)
	resource := new(FulfillmentResource)
	err := s.client.Post(ctx, path, nil, resource)
	return resource.Fulfillment, err
}

// checkOrder preserves the order boundary when using the global write endpoints.
func (s *FulfillmentServiceOp) checkOrder(ctx context.Context, fulfillmentId uint64) error {
	if s.resource != ordersResourceName {
		return nil
	}
	fulfillment, err := s.Get(ctx, fulfillmentId, nil)
	if err != nil {
		return err
	}
	if fulfillment == nil || fulfillment.Id != fulfillmentId || fulfillment.OrderId != s.resourceId {
		return fmt.Errorf("fulfillment %d does not belong to order %d", fulfillmentId, s.resourceId)
	}
	return nil
}

func fulfillmentTrackingInfo(fulfillment Fulfillment) (FulfillmentTrackingInfo, error) {
	if fulfillment.TrackingInfo != (FulfillmentTrackingInfo{}) {
		return fulfillment.TrackingInfo, nil
	}
	if len(fulfillment.TrackingNumbers) > 1 || len(fulfillment.TrackingUrls) > 1 {
		return FulfillmentTrackingInfo{}, errors.New("REST fulfillments support only one tracking number and URL; use GraphQL for multiple shipments")
	}
	tracking := FulfillmentTrackingInfo{
		Company: fulfillment.TrackingCompany,
		Number:  fulfillment.TrackingNumber,
		Url:     fulfillment.TrackingUrl,
	}
	if tracking.Number == "" && len(fulfillment.TrackingNumbers) == 1 {
		tracking.Number = fulfillment.TrackingNumbers[0]
	}
	if tracking.Url == "" && len(fulfillment.TrackingUrls) == 1 {
		tracking.Url = fulfillment.TrackingUrls[0]
	}
	return tracking, nil
}
