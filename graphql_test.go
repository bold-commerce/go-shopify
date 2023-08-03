package goshopify

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGraphQLQuery(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"data":{"foo":"bar"}}`),
	)

	resp := struct {
		Foo string `json:"foo"`
	}{}
	_, err := client.GraphQL.Query("query {}", nil, &resp)

	if err != nil {
		t.Errorf("GraphQL.Query returned error: %v", err)
	}

	expectedFoo := "bar"
	if resp.Foo != expectedFoo {
		t.Errorf("resp.Foo returned %s expected %s", resp.Foo, expectedFoo)
	}
}

func TestGraphQLQueryWithError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"errors":[{"message":"oops"}]}`),
	)

	resp := struct {
		Foo string `json:"foo"`
	}{}
	_, err := client.GraphQL.Query("query {}", nil, &resp)

	if err == nil {
		t.Error("GraphQL.Query should return error!")
	}

	expectedError := "oops"
	if err.Error() != expectedError {
		t.Errorf("GraphQL.Query returned error message %s but expected %s", err.Error(), expectedError)
	}
}

func TestGraphQLQueryWithRetries(t *testing.T) {
	setup()
	defer teardown()

	type MyStruct struct {
		Foo string `json:"foo"`
	}

	var retries int

	cases := []struct {
		description string
		responder   httpmock.Responder
		expected    interface{}
		retries     int
	}{
		{
			description: "no retries",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, `{"data":{"foo":"bar"}}`), nil
			},
			expected: MyStruct{Foo: "bar"},
			retries:  1,
		},
		{
			description: "3 throttled retries",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, `
					{
						"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}],
						"extensions":{
							"cost":{
								"requestedQueryCost":400,
								"throttleStatus":{
									"maximumAvailable":1000.0,
									"currentlyAvailable":300,
									"restoreRate":50.0
								}
							}
						}
					}`), nil
			},
			expected: RateLimitError{
				ResponseError: ResponseError{
					Status:  200,
					Message: "Throttled",
				},
				RetryAfter: 2,
			},
			retries: maxRetries,
		},
		{
			description: "2 throttled then success",
			responder: func(req *http.Request) (*http.Response, error) {
				if retries > 1 {
					retries--
					return httpmock.NewStringResponse(200, `
					{
						"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}],
						"extensions":{
							"cost":{
								"requestedQueryCost":400,
								"throttleStatus":{
									"maximumAvailable":1000.0,
									"currentlyAvailable":300,
									"restoreRate":50.0
								}
							}
						}
					}`), nil
				}

				return httpmock.NewStringResponse(200, `{"data":{"foo":"bar"}}`), nil
			},
			expected: MyStruct{Foo: "bar"},
			retries:  maxRetries,
		},
		{
			description: "1 503, 1 throttled then success",
			responder: func(req *http.Request) (*http.Response, error) {
				if retries > 2 {
					retries--
					return httpmock.NewStringResponse(http.StatusServiceUnavailable, "<html></html>"), nil
				}

				if retries > 1 {
					retries--
					return httpmock.NewStringResponse(200, `
					{
						"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}],
						"extensions":{
							"cost":{
								"requestedQueryCost":400,
								"throttleStatus":{
									"maximumAvailable":1000.0,
									"currentlyAvailable":300,
									"restoreRate":50.0
								}
							}
						}
					}`), nil
				}

				return httpmock.NewStringResponse(200, `{"data":{"foo":"bar"}}`), nil
			},
			expected: MyStruct{Foo: "bar"},
			retries:  maxRetries,
		},

		{
			description: "3 503s",
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusServiceUnavailable, ""), nil
			},
			expected: ResponseError{
				Status: http.StatusServiceUnavailable,
			},
			retries: maxRetries,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			// used to track retries in case clojure
			retries = c.retries

			httpmock.RegisterResponder(
				"POST",
				fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix),
				c.responder,
			)

			resp := MyStruct{}
			attempts, err := client.GraphQL.Query("query {}", nil, &resp)

			if attempts != c.retries {
				t.Errorf("GraphQL.Query attempts equal %d but expected %d", attempts, c.retries)
			}

			if err != nil {
				if !reflect.DeepEqual(err, c.expected) {
					t.Errorf("GraphQL.Query got error %#v but expected %#v", err, c.expected)
				}
			} else if !reflect.DeepEqual(resp, c.expected) {
				t.Errorf("GraphQL.Query responsed %#v but expected %#v", resp, c.expected)
			}
		})
	}
}

func TestGraphQLQueryWithMultipleErrors(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"errors":[{"message":"oops"},{"message":"I did it again"}]}`),
	)

	resp := struct {
		Foo string `json:"foo"`
	}{}
	_, err := client.GraphQL.Query("query {}", nil, &resp)

	if err == nil {
		t.Error("GraphQL.Query should return error!")
	}

	expectedError := "I did it again, oops"
	if err.Error() != expectedError {
		t.Errorf("GraphQL.Query returned error message %s but expected %s", err.Error(), expectedError)
	}
}

func TestGraphQLQueryWithThrottledError(t *testing.T) {
	setup()
	defer teardown()
	client.retries = 1

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `
			{
				"errors":[{"message":"Throttled","extensions":{"code":"THROTTLED"}}],
				"extensions":{
					"cost":{
						"requestedQueryCost":400,
						"throttleStatus":{
							"maximumAvailable":1000.0,
							"currentlyAvailable":300,
							"restoreRate":50.0
						}
					}
				}
			}`),
	)

	resp := struct {
		Foo string `json:"foo"`
	}{}
	_, err := client.GraphQL.Query("query {}", nil, &resp)

	if err == nil {
		t.Error("GraphQL.Query should return error!")
	}

	expectedError := "Throttled"
	if err.Error() != expectedError {
		t.Errorf("GraphQL.Query returned error message %s but expected %s", err.Error(), expectedError)
	}

	rle, ok := err.(RateLimitError)
	if !ok {
		t.Errorf("GraphQL.Query returned error not of type RateLimitError")
	}

	expectedRetryAfterSeconds := 2.0
	if rle.RetryAfter != int(expectedRetryAfterSeconds) {
		t.Errorf("GraphQL.Query rle.RetryAfter is %d but expected %d", rle.RetryAfter, int(expectedRetryAfterSeconds))
	}

	if client.RateLimits.GraphQLCost == nil {
		t.Errorf("GraphQL.Query should have assigned client.RateLimits.GraphQLCost")
	}

	if client.RateLimits.RetryAfterSeconds != expectedRetryAfterSeconds {
		t.Errorf("GraphQL.Query client.RateLimits.RetryAfterSeconds is %f but expected %f", client.RateLimits.RetryAfterSeconds, expectedRetryAfterSeconds)
	}
}

func TestGraphQLCostRetryAfterSeconds(t *testing.T) {
	cases := []struct {
		description string
		GraphQLCost GraphQLCost
		expected    float64
	}{
		{
			"last query passed, does not need to be throttled",
			GraphQLCost{
				RequestedQueryCost: 300,
				ActualQueryCost:    makeIntPointer(50),
				ThrottleStatus: GraphQLThrottleStatus{
					MaximumAvailable:   1000,
					CurrentlyAvailable: 400,
					RestoreRate:        50,
				},
			},
			0,
		},
		{
			"last query failed, needs to be throttled",
			GraphQLCost{
				RequestedQueryCost: 300,
				ActualQueryCost:    nil,
				ThrottleStatus: GraphQLThrottleStatus{
					MaximumAvailable:   1000,
					CurrentlyAvailable: 200,
					RestoreRate:        50,
				},
			},
			2,
		},
		{
			"last query passed, does not need to be throttled",
			GraphQLCost{
				RequestedQueryCost: 300,
				ActualQueryCost:    makeIntPointer(50),
				ThrottleStatus: GraphQLThrottleStatus{
					MaximumAvailable:   1000,
					CurrentlyAvailable: 200,
					RestoreRate:        50,
				},
			},
			0,
		},
		{
			"last query passed, needs to be throttled",
			GraphQLCost{
				RequestedQueryCost: 300,
				ActualQueryCost:    makeIntPointer(100),
				ThrottleStatus: GraphQLThrottleStatus{
					MaximumAvailable:   1000,
					CurrentlyAvailable: 50,
					RestoreRate:        50,
				},
			},
			1,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			s := c.GraphQLCost.RetryAfterSeconds()

			if s != c.expected {
				t.Errorf("GraphQLCost.RetryAfterSeconds returned %f expected %f (%s)", s, c.expected, c.description)
			}
		})
	}
}

func makeIntPointer(v int) *int {
	return &v
}
