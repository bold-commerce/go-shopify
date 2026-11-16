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
			want: []MarketingEvent{
				{Id: 7, EventType: "post", UTMCampaign: "summer", UTMSource: "newsletter", UTMMedium: "email", MarketingChannel: "email"},
			},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "marketing_events/7.json",
			response: `{"marketing_event":{"id":7,"event_type":"post","utm_campaign":"summer","utm_source":"newsletter","utm_medium":"email","marketing_channel":"email"}}`,
			want: &MarketingEvent{
				Id:               7,
				EventType:        "post",
				UTMCampaign:      "summer",
				UTMSource:        "newsletter",
				UTMMedium:        "email",
				MarketingChannel: "email",
			},
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
	testFeaturePagination(t, "marketing_events.json", `{"marketing_events":[{"id":7}]}`, []MarketingEvent{{Id: 7}}, `{"marketing_events":[{"id":8}]}`, []MarketingEvent{{Id: 8}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
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
	amount := decimal.RequireFromString("0")
	for _, tc := range []featureEndpoint{
		{
			name:     "Create",
			status:   201,
			method:   "POST",
			path:     "marketing_events.json",
			request:  `{"marketing_event":{"started_at":"2026-01-15T00:00:00Z","event_type":"post","utm_campaign":"summer","utm_source":"newsletter","utm_medium":"email","marketing_channel":"email","paid":false,"marketed_resources":[]}}`,
			response: `{"marketing_event":{"id":7}}`,
			want:     &MarketingEvent{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Create(ctx, MarketingEvent{StartedAt: &at, EventType: "post", UTMCampaign: "summer", UTMSource: "newsletter", UTMMedium: "email", MarketingChannel: "email", Paid: &disabled, MarketedResources: &resources})
			},
		},
		{
			name:     "UpdateBudget",
			method:   "PUT",
			path:     "marketing_events/7.json",
			request:  `{"marketing_event":{"id":7,"budget":"0","budget_type":"daily","currency":"USD"}}`,
			response: `{"marketing_event":{"id":7,"budget":"0","budget_type":"daily","currency":"USD"}}`,
			want:     &MarketingEvent{Id: 7, Budget: &amount, BudgetType: "daily", Currency: "USD"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MarketingEvent.Update(ctx, MarketingEvent{Id: 7, Budget: &amount, BudgetType: "daily", Currency: "USD"})
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

func TestMarketingEventSupplementalResponse(t *testing.T) {
	paid := true
	budget := decimal.RequireFromString("100.25")
	at := featureTime(t, "2026-11-01T09:00:00-07:00")
	done := featureTime(t, "2026-11-02T09:00:00-07:00")
	resources := []MarketedResource{{Type: "product", Id: 9007199254740993}, {Type: "collection", Id: 11}}
	testFeatureEndpoint(t, featureEndpoint{
		name:     "PopulatedOptionalFields",
		method:   "GET",
		path:     "marketing_events/7.json",
		response: string(loadFixture("marketing_event_expanded.json")),
		want: &MarketingEvent{
			Id:                7,
			RemoteId:          "remote-1",
			EventType:         "ad",
			MarketingChannel:  "social",
			Paid:              &paid,
			ReferringDomain:   "example.com",
			Budget:            &budget,
			Currency:          "USD",
			BudgetType:        "lifetime",
			StartedAt:         at,
			ScheduledToEndAt:  done,
			EndedAt:           done,
			UTMCampaign:       "Summer & autumn",
			UTMSource:         "social",
			UTMMedium:         "cpc",
			Description:       "Campaign details",
			ManageURL:         "https://example.com/manage/7",
			PreviewURL:        "https://example.com/preview/7",
			MarketedResources: &resources,
		},
		call: func(ctx context.Context) (interface{}, error) { return client.MarketingEvent.Get(ctx, 7, nil) },
	})
}

func TestMarketingEngagementPopulatedFields(t *testing.T) {
	day := OnlyDate{Time: *featureTime(t, "2024-02-29T00:00:00Z")}
	at := featureTime(t, "2024-03-01T10:00:00+05:30")
	counts := []int64{100, 200, 30, 4, 5, 6, 7, 8, 9, 1000, 80, 20}
	amount := decimal.RequireFromString("12.34")
	cumulative, offset := true, 330
	engagement := MarketingEngagement{
		OccurredOn: &day, FetchedAt: at, AdSpend: &amount, CurrencyCode: "USD", IsCumulative: &cumulative, UTCOffset: &offset,
		ViewsCount: &counts[0], ImpressionsCount: &counts[1], ClicksCount: &counts[2], FavoritesCount: &counts[3],
		CommentsCount: &counts[4], SharesCount: &counts[5], UnsubscribesCount: &counts[6], ComplaintsCount: &counts[7],
		FailsCount: &counts[8], SendsCount: &counts[9], UniqueViewsCount: &counts[10], UniqueClicksCount: &counts[11],
	}
	body := `{"engagements":[{"occurred_on":"2024-02-29","fetched_at":"2024-03-01T10:00:00+05:30","views_count":100,"impressions_count":200,"clicks_count":30,"favorites_count":4,"comments_count":5,"shares_count":6,"unsubscribes_count":7,"complaints_count":8,"fails_count":9,"sends_count":1000,"unique_views_count":80,"unique_clicks_count":20,"ad_spend":"12.34","currency_code":"USD","is_cumulative":true,"utc_offset":330}]}`
	testFeatureEndpoint(t, featureEndpoint{
		method:           "POST",
		status:           201,
		path:             "marketing_events/7/engagements.json",
		request:          body,
		response:         body,
		want:             []MarketingEngagement{engagement},
		invalidResponses: []string{`{"engagements":[{"occurred_on":"2025-02-29"}]}`, `{"engagements":[{"clicks_count":1.5}]}`},
		call: func(ctx context.Context) (interface{}, error) {
			return client.MarketingEvent.CreateEngagements(ctx, 7, []MarketingEngagement{engagement})
		},
	})
}
