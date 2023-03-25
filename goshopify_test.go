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
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

const (
	testApiVersion = "9999-99"
	maxRetries     = 3
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
	client = NewClient(app, "fooshop", "abcd",
		WithVersion(testApiVersion),
		WithRetry(maxRetries))
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

func TestBadShopNamePanic(t *testing.T) {
	func() {
		var tried string
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TestBadShopNamePanic should have panicked tried %s", tried)
			}
		}()

		for _, shopName := range []string{
			"foo shop",
			"foo, shop, stuff commas",
		} {
			tried = shopName
			_ = NewClient(app, shopName, "abcd", WithVersion(testApiVersion))
		}
	}()
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
				Message: "Not Acceptable",
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
		req, err := client.NewRequest("GET", c.url, nil, nil)
		if err != nil {
			t.Error("error creating request: ", err)
		}

		err = client.Do(req, body)
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

func TestRetry(t *testing.T) {
	setup()
	defer teardown()

	type MyStruct struct {
		Foo string `json:"foo"`
	}

	var retries int
	urlFormat := "https://fooshop.myshopify.com/%s"

	cases := []struct {
		relPath   string
		responder httpmock.Responder
		expected  interface{}
		retries   int
	}{
		{ // no retries
			relPath:  "foo/1",
			retries:  1,
			expected: &MyStruct{Foo: "bar"},
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, `{"foo": "bar"}`), nil
			},
		},
		{ // 2 retries rate limited, 3 succeeds
			relPath:  "foo/2",
			retries:  maxRetries,
			expected: &MyStruct{Foo: "bar"},
			responder: func(req *http.Request) (*http.Response, error) {
				if retries > 1 {
					resp := httpmock.NewStringResponse(http.StatusTooManyRequests, `{"errors":"Exceeded 2 calls per second for api client. Reduce request rates to resume uninterrupted service."}`)
					resp.Header.Add("Retry-After", "2.0")
					retries--
					return resp, nil
				}

				return httpmock.NewStringResponse(http.StatusOK, `{"foo": "bar"}`), nil
			},
		},
		{ // all retries rate limited
			relPath: "foo/3",
			retries: maxRetries,
			expected: RateLimitError{
				RetryAfter: 2,
				ResponseError: ResponseError{
					Status:  429,
					Message: "Exceeded 2 calls per second for api client. Reduce request rates to resume uninterrupted service.",
				},
			},
			responder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(http.StatusTooManyRequests, `{"errors":"Exceeded 2 calls per second for api client. Reduce request rates to resume uninterrupted service."}`)
				resp.Header.Add("Retry-After", "2.0")
				return resp, nil
			},
		},
		{ // 2 retries 503, 3 succeeds
			relPath:  "foo/4",
			retries:  maxRetries,
			expected: &MyStruct{Foo: "bar"},
			responder: func(req *http.Request) (*http.Response, error) {
				if retries > 1 {
					retries--
					return httpmock.NewStringResponse(http.StatusServiceUnavailable, "<html></html>"), nil
				}

				return httpmock.NewStringResponse(http.StatusOK, `{"foo": "bar"}`), nil
			},
		},
		{ // all retries 503
			relPath: "foo/5",
			retries: maxRetries,
			expected: ResponseError{
				Status: http.StatusServiceUnavailable,
			},
			responder: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusServiceUnavailable, ""), nil
			},
		},
	}

	for _, c := range cases {
		retries = c.retries
		httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat, c.relPath), c.responder)
		body := new(MyStruct)
		req, err := client.NewRequest("GET", c.relPath, nil, nil)
		if err != nil {
			t.Error("error creating request: ", err)
		}

		err = client.Do(req, body)

		if client.attempts != c.retries {
			t.Errorf("Do(): attempts do not match retries %#v, actual %#v", client.attempts, c.retries)
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

func TestClientDoAutoApiVersion(t *testing.T) {
	u := "foo/1"
	responder := func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, ``)
		resp.Header.Add("X-Shopify-API-Version", testApiVersion)
		return resp, nil
	}
	expected := testApiVersion

	testClient := NewClient(app, "fooshop", "abcd")
	httpmock.ActivateNonDefault(testClient.Client)
	shopUrl := fmt.Sprintf("https://fooshop.myshopify.com/%v", u)
	httpmock.RegisterResponder("GET", shopUrl, responder)

	req, err := testClient.NewRequest("GET", u, nil, nil)
	if err != nil {
		t.Errorf("TestClientDoApiVersion(): errored %s", err)
	}

	err = testClient.Do(req, nil)
	if err != nil {
		t.Errorf("TestClientDoApiVersion(): errored %s", err)
	}

	if expected != testClient.apiVersion {
		t.Errorf(
			"TestClientDoApiVersion(): client unable to get API Version from X-Shopify-API-Version: expected %s received %s",
			expected, testClient.apiVersion)
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

	mockPrefix := fmt.Sprintf("https://fooshop.myshopify.com/%s/", client.pathPrefix)

	cases := []struct {
		url       string
		responder httpmock.Responder
		options   interface{}
		expected  interface{}
	}{
		{
			"foo/1",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			nil,
			&MyStruct{Foo: "bar"},
		},
		{
			"foo/2",
			httpmock.NewStringResponder(404, `{"error": "does not exist"}`),
			nil,
			ResponseError{Status: 404, Message: "does not exist"},
		},

		// non relPath will get auto fixed by CreateAndDo but the httpmock endpoints above will respond for them
		{
			"/foo/1",
			httpmock.NewStringResponder(200, `{"foo": "bar"}`),
			nil,
			&MyStruct{Foo: "bar"},
		},
		{
			"/foo/2",
			httpmock.NewStringResponder(404, `{"error": "does not exist"}`),
			nil,
			ResponseError{Status: 404, Message: "does not exist"},
		},
		// Problem with options to test c.NewRequest() returning error in createAndDoGetHeaders()
		{
			"foo/1",
			httpmock.NewStringResponder(500, ""),
			123,
			errors.New("query: Values() expects struct input. Got int"),
		},
	}

	for _, c := range cases {
		httpmock.RegisterResponder("GET", mockPrefix+c.url, c.responder)
		body := new(MyStruct)
		err := client.CreateAndDo("GET", c.url, nil, c.options, body)

		if err != nil && fmt.Sprint(err) != fmt.Sprint(c.expected) {
			t.Errorf("CreateAndDo(): expected error %v, actual %v", c.expected, err)
		} else if err == nil && !reflect.DeepEqual(body, c.expected) {
			t.Errorf("CreateAndDo(): expected %#v, actual %#v", c.expected, body)
		}
	}

	// test invalid invalid protocol
	httpmock.RegisterResponder("GET", "://fooshop.myshopify.com/foo/2", httpmock.NewStringResponder(200, ""))
	body := new(MyStruct)
	_, err := client.NewRequest("GET", "://fooshop.myshopify.com/foo/2", body, nil)

	expected := errors.New("parse \"://fooshop.myshopify.com/foo/2\": missing protocol scheme")

	if err != nil && !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Errorf("CreateAndDo(): expected error %v, actual %v", expected, err)
	} else if err == nil && !reflect.DeepEqual(body, expected) {
		t.Errorf("CreateAndDo(): expected %#v, actual %#v", expected, body)
	}
}

func TestResponseErrorStructError(t *testing.T) {
	res := ResponseError{
		Status:  400,
		Message: "invalid email",
		Errors:  []string{"invalid email"},
	}

	expected := ResponseError{
		Status:  res.GetStatus(),
		Message: res.GetMessage(),
		Errors:  res.GetErrors(),
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("ResponseError returned  %+v, expected %+v", res, expected)
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
			httpmock.NewStringResponse(400, `{"errors": { "order": ["order is wrong"] }}`),
			ResponseError{Status: 400, Message: "order: order is wrong", Errors: []string{"order: order is wrong"}},
		},
		{
			httpmock.NewStringResponse(400, `{"errors": { "collection_id": "collection_id is wrong" }}`),
			ResponseError{Status: 400, Message: "collection_id: collection_id is wrong", Errors: []string{"collection_id: collection_id is wrong"}},
		},
		{
			httpmock.NewStringResponse(400, `{error:bad request}`),
			errors.New("invalid character 'e' looking for beginning of object key string"),
		},
		{
			&http.Response{StatusCode: 400, Body: errReader{}},
			errors.New("test-error"),
		},
		{
			httpmock.NewStringResponse(422, `{"error": "Unprocessable Entity - ok"}`),
			ResponseError{Status: 422, Message: "Unprocessable Entity - ok"},
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
	}

	for _, c := range cases {
		actual := CheckResponseError(c.resp)
		if fmt.Sprint(actual) != fmt.Sprint(c.expected) {
			t.Errorf("CheckResponseError(): expected [%v], actual [%v]", c.expected, actual)
		}
	}
}

func TestCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/foocount", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 5}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/foocount", client.pathPrefix),
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

func createResponderWithHeaders(status int, body string, headers map[string]string) httpmock.Responder {
	header := http.Header{}
	resp := httpmock.NewStringResponse(status, body)
	for k, v := range headers {
		header.Add(k, v)
	}
	resp.Header = header

	return httpmock.ResponderFromResponse(resp)
}

func TestDoRateLimit(t *testing.T) {
	setup()
	defer teardown()

	cases := []struct {
		description string
		url         string
		responder   httpmock.Responder
		expected    RateLimitInfo
	}{
		{
			"valid request count and bucket size should set values properly",
			"foo/1",
			createResponderWithHeaders(200, `{"foo": "bar"}`, map[string]string{
				"X-Shopify-Shop-Api-Call-Limit": "15/30",
			}),
			RateLimitInfo{
				RequestCount:      15,
				BucketSize:        30,
				RetryAfterSeconds: 0,
			},
		},
		{
			"valid retry should set RetryAfterSeconds properly",
			"foo/1",
			createResponderWithHeaders(200, `{"foo": "bar"}`, map[string]string{
				"X-Shopify-Shop-Api-Call-Limit": "0/30",
				"Retry-after":                   "30",
			}),
			RateLimitInfo{
				RequestCount:      0,
				BucketSize:        30,
				RetryAfterSeconds: 30,
			},
		},
		{
			"invalid headers should set all values to 0",
			"foo/1",
			createResponderWithHeaders(200, `{"foo": "bar"}`, map[string]string{
				"X-Shopify-Shop-Api-Call-Limit": "invalid/invalid",
				"Retry-after":                   "invalid",
			}),
			RateLimitInfo{
				RequestCount:      0,
				BucketSize:        0,
				RetryAfterSeconds: 0,
			},
		},
		{
			"missing headers should set all values to 0",
			"foo/1",
			createResponderWithHeaders(200, `{"foo": "bar"}`, map[string]string{}),
			RateLimitInfo{
				RequestCount:      0,
				BucketSize:        0,
				RetryAfterSeconds: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {

			shopUrl := fmt.Sprintf("https://fooshop.myshopify.com/%v", c.url)
			httpmock.RegisterResponder("GET", shopUrl, c.responder)

			var reqBody struct {
				Foo string `json:"foo"`
			}
			req, _ := client.NewRequest("GET", c.url, nil, nil)
			err := client.Do(req, reqBody)

			if err != nil {
				if e, ok := err.(*url.Error); ok {
					err = e.Err
				} else if e, ok := err.(*json.SyntaxError); ok {
					err = errors.New(e.Error())
				}

				if !reflect.DeepEqual(err, c.expected) {
					t.Errorf("Do(): expected error %#v, actual %#v", c.expected, err)
				}
			} else if err == nil && !reflect.DeepEqual(client.RateLimits, c.expected) {
				t.Errorf("%s: expected %#v, actual %#v", c.description, c.expected, client.RateLimits)
			}
		})
	}
}

func TestListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/locations", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("locations.json")).
			HeaderSet(http.Header{
				"Link": {
					fmt.Sprintf(
						`<https://fooshop.myshopify.com/%s/locations.json?page_info=abc&limit=10>; rel="next", <https://fooshop.myshopify.com/%s/locations.json?page_info=123&limit=10>; rel="previous"`,
						client.pathPrefix,
						client.pathPrefix,
					),
				},
			}))

	var locations LocationsResource
	pagination, err := client.ListWithPagination("locations", &locations, nil)
	if err != nil {
		t.Fatalf("Client.ListWithPagination returned error: %v", err)
	}

	if pagination == nil || pagination.NextPageOptions == nil || pagination.PreviousPageOptions == nil {
		t.Fatalf("Expected pagination options but found at least one of them nil")
	}

	t.Logf("b: %#v \n", *pagination.NextPageOptions)

	if pagination.NextPageOptions.PageInfo != "abc" {
		t.Fatalf("Expected next page: %s   got: %s", "abc", pagination.NextPageOptions.PageInfo)
	}
	if pagination.PreviousPageOptions.PageInfo != "123" {
		t.Fatalf("Expected prev page: %s   got: %s", "123", pagination.PreviousPageOptions.PageInfo)
	}
}
