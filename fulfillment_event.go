package goshopify

import (
	"fmt"
)

const (
	fulfillmentEventBasePath = "orders"
)

// FulfillmentEventService is an interface for interfacing with the fulfillment event service
// of the Shopify API.
// https://help.shopify.com/api/reference/fulfillmentevent
type FulfillmentEventService interface {
	List(orderID int64, fulfillmentID int64) ([]FulfillmentEvent, error)
	Get(orderID int64, fulfillmentID int64, eventID int64) (*FulfillmentEvent, error)
	Create(orderID int64, fulfillmentID int64, event FulfillmentEvent) (*FulfillmentEvent, error)
	Delete(orderID int64, fulfillmentID int64, eventID int64) error
}

// FulfillmentEvent represents a Shopify fulfillment event.
type FulfillmentEvent struct {
	ID                  int64   `json:"id"`
	Address1            string  `json:"address1"`
	City                string  `json:"city"`
	Country             string  `json:"country"`
	CreatedAt           string  `json:"created_at"`
	EstimatedDeliveryAt string  `json:"estimated_delivery_at"`
	FulfillmentID       int64   `json:"fulfillment_id"`
	HappenedAt          string  `json:"happened_at"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Message             string  `json:"message"`
	OrderID             int64   `json:"order_id"`
	Province            string  `json:"province"`
	ShopID              int64   `json:"shop_id"`
	Status              string  `json:"status"`
	UpdatedAt           string  `json:"updated_at"`
	Zip                 string  `json:"zip"`
}

type fulfillmentEventResource struct {
	Event *FulfillmentEvent `json:"fulfillment_event"`
}

type fulfillmentEventsResource struct {
	Events []FulfillmentEvent `json:"fulfillment_events"`
}

// FulfillmentEventServiceOp handles communication with the fulfillment event related methods of the Shopify API.
type FulfillmentEventServiceOp struct {
	client *Client
}

// List of all FulfillmentEvents for an order's fulfillment
func (s *FulfillmentEventServiceOp) List(orderID int64, fulfillmentID int64) ([]FulfillmentEvent, error) {
	path := fmt.Sprintf("%s/%d/fulfillments/%d/events.json", fulfillmentEventBasePath, orderID, fulfillmentID)
	resource := new(fulfillmentEventsResource)
	err := s.client.Get(path, resource, nil)
	return resource.Events, err
}

// Get a single FulfillmentEvent
func (s *FulfillmentEventServiceOp) Get(orderID int64, fulfillmentID int64, eventID int64) (*FulfillmentEvent, error) {
	path := fmt.Sprintf("%s/%d/fulfillments/%d/events/%d.json", fulfillmentEventBasePath, orderID, fulfillmentID, eventID)
	resource := new(fulfillmentEventResource)
	err := s.client.Get(path, resource, nil)
	return resource.Event, err
}

// Create a new FulfillmentEvent
func (s *FulfillmentEventServiceOp) Create(orderID int64, fulfillmentID int64, event FulfillmentEvent) (*FulfillmentEvent, error) {
	path := fmt.Sprintf("%s/%d/fulfillments/%d/events.json", fulfillmentEventBasePath, orderID, fulfillmentID)
	wrappedData := fulfillmentEventResource{Event: &event}
	resource := new(fulfillmentEventResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Event, err
}

// Delete an existing FulfillmentEvent
func (s *FulfillmentEventServiceOp) Delete(orderID int64, fulfillmentID int64, eventID int64) error {
	path := fmt.Sprintf("%s/%d/fulfillments/%d/events/%d.json", fulfillmentEventBasePath, orderID, fulfillmentID, eventID)
	return s.client.Delete(path)
}
