package goshopify

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"gopkg.in/jarcoal/httpmock.v1"
)

const (
	testApiVersion = "9999-99"
)

var (
	client *Client
	app    App
)

// errReader can be used to simulate a failed call to response.Body.Read
type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("test-error")
}

func (errReader) Close() error {
	return nil
}

func setup() {
	app = App{
		ApiKey:      "apikey",
		ApiSecret:   "hush",
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_products",
		Password:    "privateapppassword",
	}
	client = NewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))
	httpmock.ActivateNonDefault(client.Client)
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func loadFixture(filename string) []byte {
	f, err := ioutil.ReadFile("fixtures/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}
	return f
}

func TestWithVersion(t *testing.T) {
	_ = WithVersion(testApiVersion)
	expected := fmt.Sprintf("admin/api/%s", testApiVersion)
	if globalApiPathPrefix != expected {
		t.Errorf("WithVersion globalApiPathPrefix = %s, expected %s", globalApiPathPrefix, expected)
	}
}

func TestNewClient(t *testing.T) {
	testClient := NewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))
	expected := "https://fooshop.myshopify.com"
	if testClient.baseURL.String() != expected {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.baseURL.String(), expected)
	}
}

func TestNewClientWithNoToken(t *testing.T) {
	testClient := NewClient(app, "fooshop", "", WithVersion(testApiVersion))
	expected := "https://fooshop.myshopify.com"
	if testClient.baseURL.String() != expected {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.baseURL.String(), expected)
	}
}

func TestAppNewClient(t *testing.T) {
	testClient := app.NewClient("fooshop", "abcd", WithVersion(testApiVersion))
	expected := "https://fooshop.myshopify.com"
	if testClient.baseURL.String() != expected {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.baseURL.String(), expected)
	}
}

func TestAppNewClientWithNoToken(t *testing.T) {
	testClient := app.NewClient("fooshop", "", WithVersion(testApiVersion))
	expected := "https://fooshop.myshopify.com"
	if testClient.baseURL.String() != expected {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.baseURL.String(), expected)
	}
}

func TestNewRequest(t *testing.T) {
	testClient := NewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))

	inURL, outURL := "foo?page=1", "https://fooshop.myshopify.com/foo?limit=10&page=1"
	inBody := struct {
		Hello string `json:"hello"`
	}{Hello: "World"}
	outBody := `{"hello":"World"}`

	type extraOptions struct {
		Limit int `url:"limit"`
	}

	req, err := testClient.NewRequest("GET", inURL, inBody, extraOptions{Limit: 10})
	if err != nil {
		t.Fatalf("NewRequest(%v) err = %v, expected nil", inURL, err)
	}

	// Test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// Test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, expected %v", inBody, string(body), outBody)
	}

	// Test user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if userAgent != UserAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, UserAgent)
	}

	// Test token is attached to the request
	token := req.Header.Get("X-Shopify-Access-Token")
	expected := "abcd"
	if token != expected {
		t.Errorf("NewRequest() X-Shopify-Access-Token = %v, expected %v", token, expected)
	}
}

func TestNewRequestForPrivateApp(t *testing.T) {
	testClient := NewClient(app, "fooshop", "", WithVersion(testApiVersion))

	inURL, outURL := "foo?page=1", "https://fooshop.myshopify.com/foo?limit=10&page=1"
	inBody := struct {
		Hello string `json:"hello"`
	}{Hello: "World"}
	outBody := `{"hello":"World"}`

	type extraOptions struct {
		Limit int `url:"limit"`
	}

	req, err := testClient.NewRequest("GET", inURL, inBody, extraOptions{Limit: 10})
	if err != nil {
		t.Fatalf("NewRequest(%v) err = %v, expected nil", inURL, err)
	}

	// Test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// Test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, expected %v", inBody, string(body), outBody)
	}

	// Test user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if userAgent != UserAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, UserAgent)
	}

	// Test token is not attached to the request
	token := req.Header.Get("X-Shopify-Access-Token")
	expected := ""
	if token != expected {
		t.Errorf("NewRequest() X-Shopify-Access-Token = %v, expected %v", token, expected)
	}

	// Test Basic Auth Set
	username, password, ok := req.BasicAuth()
	if username != app.ApiKey {
		t.Errorf("NewRequestPrivateApp() Username = %v, expected %v", username, app.ApiKey)
	}

	if password != app.Password {
		t.Errorf("NewRequestPrivateApp() Password = %v, expected %v", password, app.Password)
	}

	if ok != true {
		t.Errorf("NewRequestPrivateApp() ok = %v, expected %v", ok, true)
	}
}

func TestNewRequestMissingToken(t *testing.T) {
	testClient := NewClient(app, "fooshop", "", WithVersion(testApiVersion))

	req, _ := testClient.NewRequest("GET", "/foo", nil, nil)

	// Test token is not attached to the request
	token := req.Header["X-Shopify-Access-Token"]
	if token != nil {
		t.Errorf("NewRequest() X-Shopify-Access-Token = %v, expected %v", token, nil)
	}
}

func TestNewRequestError(t *testing.T) {
	testClient := NewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))

	cases := []struct {
		method  string
		inURL   string
		body    interface{}
		options interface{}
	}{
		{"GET", "://example.com", nil, nil}, // Error for malformed url
		{"bad method", "/foo", nil, nil},    // Error for invalid method
		{"GET", "/foo", func() {}, nil},     // Error for invalid body
		{"GET", "/foo", nil, 123},           // Error for invalid options
	}

	for _, c := range cases {
		_, err := testClient.NewRequest(c.method, c.inURL, c.body, c.options)

		if err == nil {
			t.Errorf("NewRequest(%v, %v, %v, %v) err = %v, expected error", c.method, c.inURL, c.body, c.options, err)
		}
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type MyStruct struct {
		Foo string `json:"foo"`
	}

	cases := []struct {
		url       string
		responder httpmock.Responder
		expected  interface{}
	}{
		{
			"foo/1",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			&MyStruct{Foo: "bar"},
		},
		{
			"foo/2",
			httpmock.NewStringResponder(404, `{"error": "does not exist"}`),
			ResponseError{Status: 404, Message: "does not exist"},
		},
		{
			"foo/3",
			httpmock.NewStringResponder(400, `{"errors": {"title": ["wrong"]}}`),
			ResponseError{Status: 400, Message: "title: wrong", Errors: []string{"title: wrong"}},
		},
		{
			"foo/4",
			httpmock.NewErrorResponder(errors.New("something something")),
			errors.New("something something"),
		},
		{
			"foo/5",
			httpmock.NewStringResponder(200, `{foo:bar}`),
			errors.New("invalid character 'f' looking for beginning of object key string"),
		},
		{
			"foo/6",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(429, `{"errors":"Exceeded 2 calls per second for api client. Reduce request rates to resume uninterrupted service."}`)
				resp.Header.Add("Retry-After", "2.0")
				return resp, nil
			},
			RateLimitError{
				RetryAfter: 2,
				ResponseError: ResponseError{
					Status:  429,
					Message: "Exceeded 2 calls per second for api client. Reduce request rates to resume uninterrupted service.",
				},
			},
		},
		{
			"foo/7",
			httpmock.NewStringResponder(406, ``),
			ResponseError{
				Status:  406,
				Message: "Not acceptable",
			},
		},
		{
			"foo/8",
			httpmock.NewStringResponder(500, "<html></html>"),
			ResponseDecodingError{
				Body:    []byte("<html></html>"),
				Message: "invalid character '<' looking for beginning of value",
				Status:  500,
			},
		},
	}

	for _, c := range cases {
		shopUrl := fmt.Sprintf("https://fooshop.myshopify.com/%v", c.url)
		httpmock.RegisterResponder("GET", shopUrl, c.responder)

		body := new(MyStruct)
		req, _ := client.NewRequest("GET", c.url, nil, nil)
		err := client.Do(req, body)

		if err != nil {
			if e, ok := err.(*url.Error); ok {
				err = e.Err
			} else if e, ok := err.(*json.SyntaxError); ok {
				err = errors.New(e.Error())
			}

			if !reflect.DeepEqual(err, c.expected) {
				t.Errorf("Do(): expected error %#v, actual %#v", c.expected, err)
			}
		} else if err == nil && !reflect.DeepEqual(body, c.expected) {
			t.Errorf("Do(): expected %#v, actual %#v", c.expected, body)
		}
	}
}

func TestCustomHTTPClientDo(t *testing.T) {
	setup()
	defer teardown()

	type MyStruct struct {
		Foo string `json:"foo"`
	}

	cases := []struct {
		url       string
		responder httpmock.Responder
		expected  interface{}
		client    *http.Client
	}{
		{
			"foo/1",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			&MyStruct{Foo: "bar"},
			http.DefaultClient,
		},
		{
			"foo/2",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			&MyStruct{Foo: "bar"},
			&http.Client{
				Timeout: time.Second * 1,
			},
		},
		{
			"foo/3",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			&MyStruct{Foo: "bar"},
			&http.Client{
				Timeout: time.Second * 1,
				Transport: &http.Transport{
					Dial: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
					}).Dial,
					TLSHandshakeTimeout:   10 * time.Second,
					ResponseHeaderTimeout: 10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
					TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				},
			},
		},
	}

	for _, c := range cases {

		client.Client = c.client
		httpmock.ActivateNonDefault(client.Client)

		shopUrl := fmt.Sprintf("https://fooshop.myshopify.com/%v", c.url)
		httpmock.RegisterResponder("GET", shopUrl, c.responder)

		body := new(MyStruct)
		req, err := client.NewRequest("GET", c.url, nil, nil)
		if err != nil {
			t.Fatal(c.url, err)
		}
		err = client.Do(req, body)
		if err != nil {
			t.Fatal(c.url, err)
		}

		if err != nil {
			if e, ok := err.(*url.Error); ok {
				err = e.Err
			} else if e, ok := err.(*json.SyntaxError); ok {
				err = errors.New(e.Error())
			}

			if !reflect.DeepEqual(err, c.expected) {
				t.Errorf("Do(): expected error %#v, actual %#v", c.expected, err)
			}
		} else if err == nil && !reflect.DeepEqual(body, c.expected) {
			t.Errorf("Do(): expected %#v, actual %#v", c.expected, body)
		}
	}
}

func TestCreateAndDo(t *testing.T) {
	setup()
	defer teardown()

	type MyStruct struct {
		Foo string `json:"foo"`
	}

	cases := []struct {
		url       string
		responder httpmock.Responder
		expected  interface{}
	}{
		{
			"https://fooshop.myshopify.com/foo/1",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			&MyStruct{Foo: "bar"},
		},
		{
			"https://fooshop.myshopify.com/foo/2",
			httpmock.NewStringResponder(404, `{"error": "does not exist"}`),
			ResponseError{Status: 404, Message: "does not exist"},
		},
		{
			"://fooshop.myshopify.com/foo/2",
			httpmock.NewStringResponder(200, ""),
			errors.New("parse ://fooshop.myshopify.com/foo/2: missing protocol scheme"),
		},
	}

	for _, c := range cases {
		httpmock.RegisterResponder("GET", c.url, c.responder)
		body := new(MyStruct)
		err := client.CreateAndDo("GET", c.url, nil, nil, body)

		if err != nil && fmt.Sprint(err) != fmt.Sprint(c.expected) {
			t.Errorf("CreateAndDo(): expected error %v, actual %v", c.expected, err)
		} else if err == nil && !reflect.DeepEqual(body, c.expected) {
			t.Errorf("CreateAndDo(): expected %#v, actual %#v", c.expected, body)
		}
	}
}

func TestResponseErrorError(t *testing.T) {
	cases := []struct {
		err      ResponseError
		expected string
	}{
		{
			ResponseError{Message: "oh no"},
			"oh no",
		},
		{
			ResponseError{},
			"Unknown Error",
		},
		{
			ResponseError{Errors: []string{"title: not a valid title"}},
			"title: not a valid title",
		},
		{
			ResponseError{Errors: []string{
				"not a valid title",
				"not a valid description",
			}},
			// The strings are sorted description comes first
			"not a valid description, not a valid title",
		},
	}

	for _, c := range cases {
		actual := fmt.Sprint(c.err)
		if actual != c.expected {
			t.Errorf("ResponseError.Error(): expected %s, actual %s", c.expected, actual)
		}
	}
}

func TestCheckResponseError(t *testing.T) {
	cases := []struct {
		resp     *http.Response
		expected error
	}{
		{
			httpmock.NewStringResponse(200, `{"foo": "bar"}`),
			nil,
		},
		{
			httpmock.NewStringResponse(299, `{"foo": "bar"}`),
			nil,
		},
		{
			httpmock.NewStringResponse(400, `{"error": "bad request"}`),
			ResponseError{Status: 400, Message: "bad request"},
		},
		{
			httpmock.NewStringResponse(500, `{"error": "terrible error"}`),
			ResponseError{Status: 500, Message: "terrible error"},
		},
		{
			httpmock.NewStringResponse(500, `{"errors": "This action requires read_customers scope"}`),
			ResponseError{Status: 500, Message: "This action requires read_customers scope"},
		},
		{
			httpmock.NewStringResponse(500, `{"errors": ["not", "very good"]}`),
			ResponseError{Status: 500, Message: "not, very good", Errors: []string{"not", "very good"}},
		},
		{
			httpmock.NewStringResponse(400, `{"errors": { "order": ["order is wrong"] }}`),
			ResponseError{Status: 400, Message: "order: order is wrong", Errors: []string{"order: order is wrong"}},
		},
		{
			httpmock.NewStringResponse(400, `{error:bad request}`),
			errors.New("invalid character 'e' looking for beginning of object key string"),
		},
		{
			&http.Response{StatusCode: 400, Body: errReader{}},
			errors.New("test-error"),
		},
	}

	for _, c := range cases {
		actual := CheckResponseError(c.resp)
		if fmt.Sprint(actual) != fmt.Sprint(c.expected) {
			t.Errorf("CheckResponseError(): expected %v, actual %v", c.expected, actual)
		}
	}
}

func TestCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/foocount",
		httpmock.NewStringResponder(200, `{"count": 5}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://fooshop.myshopify.com/foocount",
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	// Test without options
	cnt, err := client.Count("foocount", nil)
	if err != nil {
		t.Errorf("Client.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("Client.Count returned %d, expected %d", cnt, expected)
	}

	// Test with options
	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Count("foocount", CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Client.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Client.Count returned %d, expected %d", cnt, expected)
	}
}
