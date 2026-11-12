package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestConcurrentRateLimits(t *testing.T) {
	setup()
	defer teardown()
	baseURL := fmt.Sprintf("https://fooshop.myshopify.com/%s", client.pathPrefix)
	for i := 1; i <= 2; i++ {
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/rate-test/%d.json", baseURL, i),
			createResponderWithHeaders(200, `{}`, map[string]string{
				"X-Shopify-Shop-Api-Call-Limit": fmt.Sprintf("%d/%d", i*10, i*20),
			}))
	}
	httpmock.RegisterResponder("POST", baseURL+"/graphql.json", httpmock.NewStringResponder(200,
		`{"data":{},"extensions":{"cost":{"requestedQueryCost":10,"actualQueryCost":5,"throttleStatus":{"maximumAvailable":1000,"currentlyAvailable":995,"restoreRate":50}}}}`))
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for j := 0; j < 25; j++ {
				var err error
				switch worker % 3 {
				case 0:
					err = client.Get(context.Background(), fmt.Sprintf("rate-test/%d.json", j%2+1), nil, nil)
				case 1:
					err = client.GraphQL.Query(context.Background(), "query {}", nil, nil)
				}
				if err != nil {
					t.Errorf("request failed: %v", err)
					return
				}
				limits := client.RateLimits()
				if limits.BucketSize != limits.RequestCount*2 {
					t.Errorf("inconsistent REST snapshot: %#v", limits)
				}
				if cost := limits.GraphQLCost; cost != nil {
					if cost.ActualQueryCost == nil || cost.RequestedQueryCost != *cost.ActualQueryCost*2 {
						t.Errorf("inconsistent GraphQL snapshot: %#v", cost)
					}
					// A caller may modify its copy without changing shared state.
					cost.RequestedQueryCost = -1
					if cost.ActualQueryCost != nil {
						*cost.ActualQueryCost = -1
					}
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestRateLimitsIndependentSnapshots(t *testing.T) {
	setup()
	defer teardown()
	if got := client.RateLimits(); got != (RateLimitInfo{}) {
		t.Fatalf("initial snapshot = %#v", got)
	}
	baseURL := fmt.Sprintf("https://fooshop.myshopify.com/%s", client.pathPrefix)
	for _, actual := range []string{"null", "5"} {
		httpmock.RegisterResponder("POST", baseURL+"/graphql.json", createResponderWithHeaders(200,
			fmt.Sprintf(`{"data":{},"extensions":{"cost":{"requestedQueryCost":10,"actualQueryCost":%s,"throttleStatus":{"maximumAvailable":1000,"currentlyAvailable":995,"restoreRate":50}}}}`, actual),
			map[string]string{"X-Shopify-Shop-Api-Call-Limit": "10/20"}))
		if err := client.GraphQL.Query(context.Background(), "query {}", nil, nil); err != nil {
			t.Fatal(err)
		}
		original := client.RateLimits()
		if original.RequestCount != 10 || original.BucketSize != 20 || original.GraphQLCost == nil || original.GraphQLCost.RequestedQueryCost != 10 {
			t.Fatalf("snapshot = %#v", original)
		}
		if (original.GraphQLCost.ActualQueryCost == nil) != (actual == "null") {
			t.Fatalf("actual cost for %s = %v", actual, original.GraphQLCost.ActualQueryCost)
		}
		if actual == "5" && *original.GraphQLCost.ActualQueryCost != 5 {
			t.Fatalf("actual cost = %d", *original.GraphQLCost.ActualQueryCost)
		}
		snapshot := client.RateLimits()
		snapshot.RequestCount = 99
		snapshot.GraphQLCost.RequestedQueryCost = 99
		snapshot.GraphQLCost.ThrottleStatus.CurrentlyAvailable = 99
		if snapshot.GraphQLCost.ActualQueryCost != nil {
			*snapshot.GraphQLCost.ActualQueryCost = 99
		}
		if got := client.RateLimits(); !reflect.DeepEqual(got, original) {
			t.Fatalf("changing a snapshot changed the client: got %#v, want %#v", got, original)
		}
		httpmock.RegisterResponder("GET", baseURL+"/shop.json", createResponderWithHeaders(200, `{"shop":{}}`,
			map[string]string{"X-Shopify-Shop-Api-Call-Limit": "15/30"}))
		if _, err := client.Shop.Get(context.Background(), nil); err != nil {
			t.Fatal(err)
		}
		if original.RequestCount != 10 || client.RateLimits().RequestCount != 15 {
			t.Fatal("a later response must update the client without changing an earlier snapshot")
		}
	}
}
