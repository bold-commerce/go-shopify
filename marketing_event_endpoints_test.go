package goshopify

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMarketingEventAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "marketing_events.json",
			response: `{"marketing_events":[{"id":7,"event_type":"post","utm_campaign":"summer","utm_source":"newsletter","utm_medium":"email","marketing_channel":"email"}]}`,
			want:     []MarketingEvent{{Id: 7, EventType: "post", UTMCampaign: "summer", UTMSource: "newsletter", UTMMedium: "email", MarketingChannel: "email"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "marketing_events/7.json",
			response: `{"marketing_event":{"id":7,"event_type":"post","utm_campaign":"summer","utm_source":"newsletter","utm_medium":"email","marketing_channel":"email"}}`,
			want:     &MarketingEvent{Id: 7, EventType: "post", UTMCampaign: "summer", UTMSource: "newsletter", UTMMedium: "email", MarketingChannel: "email"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Get(ctx, 7, nil)
			},
		},
		{
			name:     "Count",
			method:   "GET",
			path:     "marketing_events/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Count(ctx, nil)
			},
		},
		{
			name:   "Delete",
			method: "DELETE",
			path:   "marketing_events/7.json",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.MarketingEvent.Delete(ctx, 7)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestMarketingEventListWithPagination(t *testing.T) {
	testFeaturePagination(t, "marketing_events.json", `{"marketing_events":[{"id":7}]}`, []MarketingEvent{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.MarketingEvent.ListWithPagination(ctx, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/marketingevent
func TestMarketingEventDocumentedResponse(t *testing.T) {
	var value MarketingEventResource
	testFeatureResponse(t, "fixtures/marketing_event_documented.json", &value, true)
	if !(value.MarketingEvent.Id == 1069063883 && *value.MarketingEvent.Paid && value.MarketingEvent.StartedAt != nil) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/marketingevent
func TestMarketingEngagementDocumentedResponse(t *testing.T) {
	var value MarketingEngagementsResource
	testFeatureResponse(t, "fixtures/marketing_engagement_documented.json", &value, true)
	if !(len(value.MarketingEngagements) == 3 && value.MarketingEngagements[0].OccurredOn.Format("2006-01-02") == "2026-01-15" && *value.MarketingEngagements[0].ViewsCount == 0 && value.MarketingEngagements[0].AdSpend.String() == "10") {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestMarketingEventWriteFields(t *testing.T) {
	at := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	day := OnlyDate{Time: at}
	disabled := false
	resources := []MarketedResource{}
	zero := int64(0)
	amount := decimal.Zero
	for _, tc := range []featureEndpoint{
		{
			name:     "Create",
			method:   "POST",
			path:     "marketing_events.json",
			request:  `{"marketing_event":{"started_at":"2026-01-15T00:00:00Z","event_type":"post","utm_campaign":"summer","utm_source":"newsletter","utm_medium":"email","marketing_channel":"email","paid":false}}`,
			response: `{"marketing_event":{"id":7}}`,
			want:     &MarketingEvent{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Create(ctx, MarketingEvent{StartedAt: &at, EventType: "post", UTMCampaign: "summer", UTMSource: "newsletter", UTMMedium: "email", MarketingChannel: "email", Paid: &disabled})
			},
		},
		{
			name:     "UpdateFalse",
			method:   "PUT",
			path:     "marketing_events/7.json",
			request:  `{"marketing_event":{"id":7,"paid":false,"marketed_resources":[]}}`,
			response: `{"marketing_event":{"id":7}}`,
			want:     &MarketingEvent{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Update(ctx, MarketingEvent{Id: 7, Paid: &disabled, MarketedResources: &resources})
			},
		},
		{
			name:     "EngagementsZero",
			method:   "POST",
			path:     "marketing_events/7/engagements.json",
			request:  `{"engagements":[{"occurred_on":"2026-01-15","views_count":0,"clicks_count":0,"ad_spend":"0","is_cumulative":false}]}`,
			response: `{"engagements":[]}`,
			want:     []MarketingEngagement{},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.CreateEngagements(ctx, 7, []MarketingEngagement{{OccurredOn: &day, ViewsCount: &zero, ClicksCount: &zero, AdSpend: &amount, IsCumulative: &disabled}})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
func TestMarketingEventQueryFields(t *testing.T) {
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Filters",
		method:   "GET",
		path:     "marketing_events.json",
		query:    "offset=2&limit=1",
		response: `{"marketing_events":[]}`,
		want:     []MarketingEvent{},
		call: func(ctx context.Context) (interface{}, error) {
			return client.MarketingEvent.List(ctx, MarketingEventListOptions{Offset: 2, ListOptions: ListOptions{Limit: 1}})
		},
	})
}
