package goshopify

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"time"
)

// MarketingEventService provides access to MarketingEvent endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/marketingevent
type MarketingEventService interface {
	ListWithPagination(ctx context.Context, options interface{}) ([]MarketingEvent, *Pagination, error)
	List(ctx context.Context, options interface{}) ([]MarketingEvent, error)
	Get(ctx context.Context, id uint64, options interface{}) (*MarketingEvent, error)
	Count(ctx context.Context, options interface{}) (int, error)
	Create(ctx context.Context, marketingEvent MarketingEvent) (*MarketingEvent, error)
	Update(ctx context.Context, marketingEvent MarketingEvent) (*MarketingEvent, error)
	Delete(ctx context.Context, id uint64) error
	CreateEngagements(ctx context.Context, eventId uint64, engagements []MarketingEngagement) ([]MarketingEngagement, error)
}

// MarketingEventServiceOp implements MarketingEventService.
type MarketingEventServiceOp struct {
	client *Client
}

// MarketingEvent represents Shopify marketing event.
type MarketingEvent struct {
	Id                  uint64              `json:"id,omitempty"`
	EventType           string              `json:"event_type,omitempty"`
	RemoteId            string              `json:"remote_id,omitempty"`
	StartedAt           *time.Time          `json:"started_at,omitempty"`
	EndedAt             *time.Time          `json:"ended_at,omitempty"`
	ScheduledToEndAt    *time.Time          `json:"scheduled_to_end_at,omitempty"`
	Budget              *decimal.Decimal    `json:"budget,omitempty"`
	Currency            string              `json:"currency,omitempty"`
	ManageURL           string              `json:"manage_url,omitempty"`
	PreviewURL          string              `json:"preview_url,omitempty"`
	UTMCampaign         string              `json:"utm_campaign,omitempty"`
	UTMSource           string              `json:"utm_source,omitempty"`
	UTMMedium           string              `json:"utm_medium,omitempty"`
	BudgetType          string              `json:"budget_type,omitempty"`
	Description         string              `json:"description,omitempty"`
	MarketingChannel    string              `json:"marketing_channel,omitempty"`
	DeliveryChannel     string              `json:"delivery_channel,omitempty"`
	Paid                *bool               `json:"paid,omitempty"`
	ReferringDomain     string              `json:"referring_domain,omitempty"`
	BreadcrumbId        *uint64             `json:"breadcrumb_id,omitempty"`
	MarketingActivityId *uint64             `json:"marketing_activity_id,omitempty"`
	AdminGraphqlApiId   string              `json:"admin_graphql_api_id,omitempty"`
	MarketedResources   *[]MarketedResource `json:"marketed_resources,omitempty"`
}

// MarketedResource represents Shopify marketed resource.
type MarketedResource struct {
	Type string `json:"type"`
	Id   uint64 `json:"id"`
}

// MarketingEngagement represents Shopify marketing engagement.
type MarketingEngagement struct {
	OccurredOn        *OnlyDate        `json:"occurred_on,omitempty"`
	FetchedAt         *time.Time       `json:"fetched_at,omitempty"`
	ViewsCount        *int64           `json:"views_count,omitempty"`
	ImpressionsCount  *int64           `json:"impressions_count,omitempty"`
	ClicksCount       *int64           `json:"clicks_count,omitempty"`
	FavoritesCount    *int64           `json:"favorites_count,omitempty"`
	CommentsCount     *int64           `json:"comments_count,omitempty"`
	SharesCount       *int64           `json:"shares_count,omitempty"`
	AdSpend           *decimal.Decimal `json:"ad_spend,omitempty"`
	CurrencyCode      string           `json:"currency_code,omitempty"`
	IsCumulative      *bool            `json:"is_cumulative,omitempty"`
	UnsubscribesCount *int64           `json:"unsubscribes_count,omitempty"`
	ComplaintsCount   *int64           `json:"complaints_count,omitempty"`
	FailsCount        *int64           `json:"fails_count,omitempty"`
	SendsCount        *int64           `json:"sends_count,omitempty"`
	UniqueViewsCount  *int64           `json:"unique_views_count,omitempty"`
	UniqueClicksCount *int64           `json:"unique_clicks_count,omitempty"`
	UTCOffset         *int             `json:"utc_offset,omitempty"`
}

// MarketingEventsResource represents the marketing_events response.
type MarketingEventsResource struct {
	MarketingEvents []MarketingEvent `json:"marketing_events"`
}

// MarketingEventResource represents the marketing_event response.
type MarketingEventResource struct {
	MarketingEvent *MarketingEvent `json:"marketing_event"`
}

// MarketingEngagementsResource represents the engagements response.
type MarketingEngagementsResource struct {
	MarketingEngagements []MarketingEngagement `json:"engagements"`
}

// MarketingEventListOptions controls filtering and pagination.
type MarketingEventListOptions struct {
	ListOptions
	Offset int `url:"offset,omitempty"`
}

// List retrieves marketing event records.
func (s *MarketingEventServiceOp) List(ctx context.Context, options interface{}) ([]MarketingEvent, error) {
	items, _, err := s.ListWithPagination(ctx, options)
	return items, err
}

// Get retrieves a marketing event record.
func (s *MarketingEventServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*MarketingEvent, error) {
	resource := new(MarketingEventResource)
	err := s.client.Get(ctx, fmt.Sprintf("marketing_events/%d.json", id), resource, options)
	return resource.MarketingEvent, err
}

// Count counts marketing event records.
func (s *MarketingEventServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	return s.client.Count(ctx, "marketing_events/count.json", options)
}

// Create creates a marketing event record.
func (s *MarketingEventServiceOp) Create(ctx context.Context, marketingEvent MarketingEvent) (*MarketingEvent, error) {
	resource := new(MarketingEventResource)
	err := s.client.Post(ctx, "marketing_events.json", MarketingEventResource{MarketingEvent: &marketingEvent}, resource)
	return resource.MarketingEvent, err
}

// Update updates a marketing event record.
func (s *MarketingEventServiceOp) Update(ctx context.Context, marketingEvent MarketingEvent) (*MarketingEvent, error) {
	resource := new(MarketingEventResource)
	err := s.client.Put(ctx, fmt.Sprintf("marketing_events/%d.json", marketingEvent.Id), MarketingEventResource{MarketingEvent: &marketingEvent}, resource)
	return resource.MarketingEvent, err
}

// Delete deletes a marketing event record.
func (s *MarketingEventServiceOp) Delete(ctx context.Context, id uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("marketing_events/%d.json", id))
}

// CreateEngagements records daily engagement metrics for a marketing event.
func (s *MarketingEventServiceOp) CreateEngagements(ctx context.Context, eventId uint64, engagements []MarketingEngagement) ([]MarketingEngagement, error) {
	resource := new(MarketingEngagementsResource)
	err := s.client.Post(ctx, fmt.Sprintf("marketing_events/%d/engagements.json", eventId), MarketingEngagementsResource{MarketingEngagements: engagements}, resource)
	return resource.MarketingEngagements, err
}

// ListWithPagination returns a page and its next/previous page options.
func (s *MarketingEventServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]MarketingEvent, *Pagination, error) {
	resource := new(MarketingEventsResource)
	pagination, err := s.client.ListWithPagination(ctx, "marketing_events.json", resource, options)
	return resource.MarketingEvents, pagination, err
}
