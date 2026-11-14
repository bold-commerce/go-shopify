package goshopify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

type featureEndpoint struct {
	name, method, path, query, request, response string
	want                                         interface{}
	call                                         func(context.Context) (interface{}, error)
}

// Exercise the HTTP boundary, including response decoding and API errors.
func testFeatureEndpoint(t *testing.T, tc featureEndpoint) {
	t.Helper()
	setup()
	defer teardown()
	endpoint := "https://fooshop.myshopify.com/" + client.pathPrefix + "/" + tc.path
	type contextKey struct{}
	ctx := context.WithValue(context.Background(), contextKey{}, "request")
	httpmock.RegisterResponder(tc.method, endpoint, func(req *http.Request) (*http.Response, error) {
		if req.Context().Value(contextKey{}) != "request" {
			t.Error("request context was lost")
		}
		if req.Header.Get("X-Shopify-Access-Token") != client.token {
			t.Error("request token was lost")
		}
		wantQuery, err := url.ParseQuery(tc.query)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(req.URL.Query(), wantQuery) {
			t.Errorf("query = %v, want %v", req.URL.Query(), wantQuery)
		}
		var data []byte
		if req.Body != nil {
			data, err = ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
		}
		if tc.request == "" {
			if len(data) != 0 {
				t.Errorf("unexpected request body: %s", data)
			}
		} else {
			var got, want interface{}
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.request), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body = %s, want %s", data, tc.request)
			}
		}
		return httpmock.NewStringResponse(http.StatusOK, tc.response), nil
	})
	got, err := tc.call(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, tc.want) {
		t.Fatalf("result = %#v, want %#v", got, tc.want)
	}
	httpmock.RegisterResponder(tc.method, endpoint, httpmock.NewStringResponder(422, `{"errors":"invalid request"}`))
	_, err = tc.call(ctx)
	var responseError ResponseError
	if !errors.As(err, &responseError) || responseError.Status != 422 {
		t.Fatalf("API error = %v, want status 422", err)
	}
	networkError := errors.New("connection failed")
	httpmock.RegisterResponder(tc.method, endpoint, func(*http.Request) (*http.Response, error) { return nil, networkError })
	if _, err := tc.call(ctx); !errors.Is(err, networkError) {
		t.Fatalf("network error was lost: %v", err)
	}
	if tc.response != "" {
		httpmock.RegisterResponder(tc.method, endpoint, httpmock.NewStringResponder(200, `{`))
		if _, err := tc.call(context.Background()); err == nil {
			t.Fatal("decode error was lost")
		}
	}
}

// Follow the returned cursor and check that each wrapper keeps the route and options.
func testFeaturePagination(t *testing.T, path, response string, want interface{}, call func(context.Context, interface{}) (interface{}, *Pagination, error)) {
	t.Helper()
	setup()
	defer teardown()
	endpoint := "https://fooshop.myshopify.com/" + client.pathPrefix + "/" + path
	requests := 0
	httpmock.RegisterResponder("GET", endpoint, func(req *http.Request) (*http.Response, error) {
		requests++
		expected := url.Values{"limit": {"1"}}
		if requests == 2 {
			expected.Set("page_info", "next")
		}
		if !reflect.DeepEqual(req.URL.Query(), expected) {
			t.Errorf("query = %v, want %v", req.URL.Query(), expected)
		}
		resp := httpmock.NewStringResponse(200, response)
		if requests == 1 {
			resp.Header.Set("Link", "<"+endpoint+"?page_info=next&limit=1>; rel=\"next\"")
		}
		return resp, nil
	})
	got, page, err := call(context.Background(), ListOptions{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("first page = %#v, want %#v", got, want)
	}
	if page == nil || page.NextPageOptions == nil || page.NextPageOptions.PageInfo != "next" {
		t.Fatalf("missing cursor: %#v", page)
	}
	got, page, err = call(context.Background(), page.NextPageOptions)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("second page = %#v, want %#v", got, want)
	}
	if page == nil || page.NextPageOptions != nil || requests != 2 {
		t.Fatalf("unexpected final page: %#v, requests: %d", page, requests)
	}
	for _, bad := range []struct {
		status int
		body   string
	}{{422, `{"errors":"invalid request"}`}, {200, `{`}} {
		httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(bad.status, bad.body))
		if _, _, err := call(context.Background(), nil); err == nil {
			t.Fatal("error was lost")
		}
	}
}

func testFeatureResponse(t *testing.T, fixture string, value interface{}, strict bool) {
	t.Helper()
	data, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if strict {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(value); err != nil {
		t.Fatal(err)
	}
}
