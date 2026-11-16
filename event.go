package goshopify

import (
	"context"
	"fmt"
	"time"
)

// EventService provides access to Event endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/event
type EventService interface {
	ListForProductWithPagination(ctx context.Context, id uint64, options interface{}) ([]Event, *Pagination, error)
	ListForOrderWithPagination(ctx context.Context, id uint64, options interface{}) ([]Event, *Pagination, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]Event, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]Event, error)
	Get(ctx context.Context, id uint64, options interface{}) (*Event, error)
	Count(ctx context.Context, options interface{}) (int, error)
	ListForOrder(ctx context.Context, id uint64, options interface{}) ([]Event, error)
	ListForProduct(ctx context.Context, id uint64, options interface{}) ([]Event, error)
}

// EventServiceOp implements EventService.
type EventServiceOp struct {
	client *Client
}

// Event represents Shopify event.
type Event struct {
	Id          uint64        `json:"id,omitempty"`
	Arguments   []interface{} `json:"arguments,omitempty"`
	Body        string        `json:"body,omitempty"`
	CreatedAt   *time.Time    `json:"created_at,omitempty"`
	Description string        `json:"description,omitempty"`
	Path        string        `json:"path,omitempty"`
	Message     string        `json:"message,omitempty"`
	SubjectId   uint64        `json:"subject_id,omitempty"`
	SubjectType string        `json:"subject_type,omitempty"`
	Verb        string        `json:"verb,omitempty"`
	Author      string        `json:"author,omitempty"`
}

// EventsResource represents the events response.
type EventsResource struct {
	Events []Event `json:"events"`
}

// EventResource represents the event response.
type EventResource struct {
	Event *Event `json:"event"`
}

// EventListOptions controls filtering and pagination.
type EventListOptions struct {
	ListOptions
	Filter string `url:"filter,omitempty"`
	Verb   string `url:"verb,omitempty"`
}

// List retrieves event records.
func (s *EventServiceOp) List(ctx context.Context, options interface{}) ([]Event, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a event record.
func (s *EventServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*Event, error) {
	resource := new(EventResource)
	err := s.client.Get(ctx, fmt.Sprintf("events/%d.json", id), resource, options)
	return resource.Event, err
}

// Count counts event records.
func (s *EventServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	return s.client.Count(ctx, "events/count.json", options)
}

// ListForOrder retrieves events for an order.
func (s *EventServiceOp) ListForOrder(ctx context.Context, id uint64, options interface{}) ([]Event, error) {
	items, _, err := s.ListForOrderWithPagination(ctx, id, options)
	return items, err
}

// ListForProduct retrieves events for a product.
func (s *EventServiceOp) ListForProduct(ctx context.Context, id uint64, options interface{}) ([]Event, error) {
	items, _, err := s.ListForProductWithPagination(ctx, id, options)
	return items, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *EventServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Event, *Pagination, error) {
	resource := new(EventsResource)
	pagination, err := s.client.ListWithPagination(ctx, "events.json", resource, options)
	return resource.Events, pagination, err
}

// ListForOrderWithPagination returns a page and its next/previous page options.
func (s *EventServiceOp) ListForOrderWithPagination(ctx context.Context, id uint64, options interface{}) ([]Event, *Pagination, error) {
	resource := new(EventsResource)
	pagination, err := s.client.ListWithPagination(ctx, fmt.Sprintf("orders/%d/events.json", id), resource, options)
	return resource.Events, pagination, err
}

// ListForProductWithPagination returns a page and its next/previous page options.
func (s *EventServiceOp) ListForProductWithPagination(ctx context.Context, id uint64, options interface{}) ([]Event, *Pagination, error) {
	resource := new(EventsResource)
	pagination, err := s.client.ListWithPagination(ctx, fmt.Sprintf("products/%d/events.json", id), resource, options)
	return resource.Events, pagination, err
}
