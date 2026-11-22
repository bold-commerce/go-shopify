package goshopify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestConcurrentRetryBudgets(t *testing.T) {
	setup()
	defer teardown()
	client.retries = 2
	const requests = 16
	paths := make([]string, requests)
	urls := make([]string, requests)
	for i := range paths {
		paths[i] = fmt.Sprintf("retry-test/%d.json", i)
		urls[i] = fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, paths[i])
		// Exhaust retries to isolate request counters from response rate-limit state.
		httpmock.RegisterResponder("GET", urls[i], httpmock.NewStringResponder(503, `{"errors":"unavailable"}`))
	}
	start := make(chan struct{})
	results := make(chan error, requests)
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			<-start
			results <- client.Get(context.Background(), path, nil, nil)
		}(path)
	}
	close(start)
	wg.Wait()
	close(results)
	for err := range results {
		var responseErr ResponseError
		if !errors.As(err, &responseErr) || responseErr.Status != 503 {
			t.Errorf("expected exhausted 503 retries, got %v", err)
		}
	}
	counts := httpmock.GetCallCountInfo()
	for _, url := range urls {
		if got := counts["GET "+url]; got != 3 {
			t.Errorf("%s: got %d attempts, want 3", url, got)
		}
	}
}

func TestRetryBudget(t *testing.T) {
	setup()
	defer teardown()
	for _, graphql := range []bool{false, true} {
		for _, retries := range []int{0, 1, 2} {
			for _, succeed := range []bool{false, true} {
				calls := 0
				client.retries = retries
				method, path := "GET", "shop.json"
				if graphql {
					method, path = "POST", "graphql.json"
				}
				httpmock.RegisterResponder(method, fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, path), func(req *http.Request) (*http.Response, error) {
					calls++
					if succeed && calls == retries+1 {
						if graphql {
							return httpmock.NewStringResponse(200, `{"data":{}}`), nil
						}
						return httpmock.NewStringResponse(200, `{"shop":{}}`), nil
					}
					if graphql {
						return httpmock.NewStringResponse(200, `{"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}]}`), nil
					}
					return httpmock.NewStringResponse(503, `{"errors":"unavailable"}`), nil
				})
				var err error
				if graphql {
					err = client.GraphQL.Query(context.Background(), "query {}", nil, nil)
				} else {
					_, err = client.Shop.Get(context.Background(), nil)
				}
				if (err == nil) != succeed {
					t.Fatalf("graphql=%v retries=%d succeed=%v: %v", graphql, retries, succeed, err)
				}
				if calls != retries+1 {
					t.Fatalf("graphql=%v retries=%d: got %d calls", graphql, retries, calls)
				}
			}
		}
	}
}

func TestFractionalRetryAfter(t *testing.T) {
	for _, tc := range []struct {
		header  string
		seconds int
	}{
		{"0.5", 1}, {"1.1", 2}, {"2.0", 2}, {"", 0}, {"invalid", 0},
	} {
		response := httpmock.NewStringResponse(429, `{"errors":"Throttled"}`)
		response.Header.Set("Retry-After", tc.header)
		err := CheckResponseError(response)
		limit, ok := err.(RateLimitError)
		if !ok || limit.RetryAfter != tc.seconds {
			t.Errorf("Retry-After %q: got %#v, want %d seconds", tc.header, err, tc.seconds)
		}
	}
}

func TestRetryCancellation(t *testing.T) {
	setup()
	defer teardown()
	client.retries = 1
	for _, graphql := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0
		method, path := "GET", "shop.json"
		if graphql {
			method, path = "POST", "graphql.json"
		}
		httpmock.RegisterResponder(method, fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, path), func(req *http.Request) (*http.Response, error) {
			calls++
			time.AfterFunc(20*time.Millisecond, cancel)
			if graphql {
				return httpmock.NewStringResponse(200, `{"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}],"extensions":{"cost":{"requestedQueryCost":60,"throttleStatus":{"currentlyAvailable":0,"restoreRate":1}}}}`), nil
			}
			response := httpmock.NewStringResponse(429, `{"errors":"Throttled"}`)
			response.Header.Set("Retry-After", "60")
			return response, nil
		})
		start := time.Now()
		var err error
		if graphql {
			err = client.GraphQL.Query(ctx, "query {}", nil, nil)
		} else {
			_, err = client.Shop.Get(ctx, nil)
		}
		cancel()
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("graphql=%v: got %v", graphql, err)
		}
		if calls != 1 {
			t.Fatalf("graphql=%v: made %d requests after cancellation", graphql, calls)
		}
		if time.Since(start) > time.Second {
			t.Fatalf("graphql=%v: backoff ignored cancellation", graphql)
		}
	}
}

func TestWaitForRetryCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := waitForRetry(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestGraphQLMixedRetryBudget(t *testing.T) {
	const throttled = `{"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}]}`
	for _, tc := range []struct {
		name       string
		statuses   []int
		retries    int
		wantStatus int
	}{
		{"HTTP then GraphQL", []int{429, 200}, 1, 200},
		{"GraphQL then HTTP", []int{200, 503}, 1, 503},
		{"mixed success", []int{503, 200, 429, 200}, 3, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()
			client.retries = tc.retries
			calls := 0
			endpoint := fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix)
			httpmock.RegisterResponder("POST", endpoint, func(req *http.Request) (*http.Response, error) {
				if calls >= len(tc.statuses) {
					t.Fatal("sent an extra request")
				}
				status := tc.statuses[calls]
				calls++
				body := `{"errors":"unavailable"}`
				if status == 200 {
					body = throttled
					if tc.wantStatus == 0 && calls == len(tc.statuses) {
						body = `{"data":{}}`
					}
				}
				return httpmock.NewStringResponse(status, body), nil
			})
			err := client.GraphQL.Query(context.Background(), "query {shop {id}}", nil, nil)
			switch tc.wantStatus {
			case 0:
				if err != nil {
					t.Fatal(err)
				}
			case 200:
				var limit RateLimitError
				if !errors.As(err, &limit) || limit.Status != 200 {
					t.Fatalf("got %v, want GraphQL throttle", err)
				}
			case 503:
				var response ResponseError
				if !errors.As(err, &response) || response.Status != 503 {
					t.Fatalf("got %v, want HTTP 503", err)
				}
			}
			if calls != tc.retries+1 {
				t.Fatalf("got %d calls, want %d", calls, tc.retries+1)
			}
		})
	}
}
