// Package goshopify provides methods for making requests to Shopify's admin API.
package goshopify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	UserAgent = "goshopify/1.0.0"
	// UnstableApiVersion Shopify API version for accessing unstable API features
	UnstableApiVersion = "unstable"

	// Shopify API version YYYY-MM - defaults to admin which uses the oldest stable version of the api
	defaultApiPathPrefix = "admin"
	defaultApiVersion    = "stable"
	defaultHttpTimeout   = 10
)

// version regex match
var apiVersionRegex = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}$`)

// App represents basic app settings such as Api key, secret, scope, and redirect url.
// See oauth.go for OAuth related helper functions.
type App struct {
	ApiKey      string
	ApiSecret   string
	RedirectUrl string
	Scope       string
	Password    string
	Client      *Client // see GetAccessToken
}

type RateLimitInfo struct {
	RequestCount      int
	BucketSize        int
	GraphQLCost       *GraphQLCost
	RetryAfterSeconds float64
}

// Client manages communication with the Shopify API.
type Client struct {
	// HTTP client used to communicate with the Shopify API.
	Client *http.Client
	log    LeveledLoggerInterface

	// App settings
	app App

	// Base URL for API requests.
	// This is set on a per-store basis which means that each store must have
	// its own client.
	baseURL *url.URL

	// URL Prefix, defaults to "admin" see WithVersion
	pathPrefix string

	// version you're currently using of the api, defaults to "stable"
	apiVersion string

	// A permanent access token
	token string

	// max number of retries, defaults to 0 for no retries see WithRetry option
	retries  int
	attempts int

	RateLimits RateLimitInfo

	// Services used for communicating with the API
	Product                    ProductService
	CustomCollection           CustomCollectionService
	SmartCollection            SmartCollectionService
	Customer                   CustomerService
	CustomerAddress            CustomerAddressService
	Order                      OrderService
	Fulfillment                FulfillmentService
	DraftOrder                 DraftOrderService
	AbandonedCheckout          AbandonedCheckoutService
	Shop                       ShopService
	Webhook                    WebhookService
	Variant                    VariantService
	Image                      ImageService
	Transaction                TransactionService
	Theme                      ThemeService
	Asset                      AssetService
	ScriptTag                  ScriptTagService
	RecurringApplicationCharge RecurringApplicationChargeService
	UsageCharge                UsageChargeService
	Metafield                  MetafieldService
	Blog                       BlogService
	ApplicationCharge          ApplicationChargeService
	Redirect                   RedirectService
	Page                       PageService
	StorefrontAccessToken      StorefrontAccessTokenService
	Collect                    CollectService
	Collection                 CollectionService
	Location                   LocationService
	DiscountCode               DiscountCodeService
	PriceRule                  PriceRuleService
	InventoryItem              InventoryItemService
	ShippingZone               ShippingZoneService
	ProductListing             ProductListingService
	InventoryLevel             InventoryLevelService
	AccessScopes               AccessScopesService
	FulfillmentService         FulfillmentServiceService
	CarrierService             CarrierServiceService
	Payouts                    PayoutsService
	GiftCard                   GiftCardService
	FulfillmentOrder           FulfillmentOrderService
	GraphQL                    GraphQLService
	AssignedFulfillmentOrder   AssignedFulfillmentOrderService
	FulfillmentEvent           FulfillmentEventService
	FulfillmentRequest         FulfillmentRequestService
	PaymentsTransactions       PaymentsTransactionsService
	OrderRisk                  OrderRiskService
	ApiPermissions             ApiPermissionsService
	Article                    ArticlesService
}

// A general response error that follows a similar layout to Shopify's response
// errors, i.e. either a single message or a list of messages.
type ResponseError struct {
	Status  int
	Message string
	Errors  []string
}

// GetStatus returns http  response status
func (e ResponseError) GetStatus() int {
	return e.Status
}

// GetMessage returns response error message
func (e ResponseError) GetMessage() string {
	return e.Message
}

// GetErrors returns response errors list
func (e ResponseError) GetErrors() []string {
	return e.Errors
}

func (e ResponseError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	sort.Strings(e.Errors)
	s := strings.Join(e.Errors, ", ")

	if s != "" {
		return s
	}

	return "Unknown Error"
}

// ResponseDecodingError occurs when the response body from Shopify could
// not be parsed.
type ResponseDecodingError struct {
	Body    []byte
	Message string
	Status  int
}

func (e ResponseDecodingError) Error() string {
	return e.Message
}

// An error specific to a rate-limiting response. Embeds the ResponseError to
// allow consumers to handle it the same was a normal ResponseError.
type RateLimitError struct {
	ResponseError
	RetryAfter int
}

// Creates an API request. A relative URL can be provided in urlStr, which will
// be resolved to the BaseURL of the Client. Relative URLS should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *Client) NewRequest(ctx context.Context, method, relPath string, body, options interface{}) (*http.Request, error) {
	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	// Make the full url based on the relative path
	u := c.baseURL.ResolveReference(rel)

	// Add custom options
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}

		for k, values := range u.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
		u.RawQuery = optionsQuery.Encode()
	}

	// A bit of JSON ceremony
	var js []byte = nil

	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)

	if c.token != "" {
		req.Header.Add("X-Shopify-Access-Token", c.token)
	} else if c.app.Password != "" {
		req.SetBasicAuth(c.app.ApiKey, c.app.Password)
	}

	return req, nil
}

// MustNewClient returns a new Shopify API client with an already authenticated shopname and
// token. The shopName parameter is the shop's myshopify domain,
// e.g. "theshop.myshopify.com", or simply "theshop"
// a.NewClient(shopName, token, opts) is equivalent to NewClient(a, shopName, token, opts)
func (app App) NewClient(shopName, token string, opts ...Option) (*Client, error) {
	return NewClient(app, shopName, token, opts...)
}

// Returns a new Shopify API client with an already authenticated shopname and
// token. The shopName parameter is the shop's myshopify domain,
// e.g. "theshop.myshopify.com", or simply "theshop"
// panics if an error occurs
func MustNewClient(app App, shopName, token string, opts ...Option) *Client {
	c, err := NewClient(app, shopName, token, opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// Returns a new Shopify API client with an already authenticated shopname and
// token. The shopName parameter is the shop's myshopify domain,
// e.g. "theshop.myshopify.com", or simply "theshop"
func NewClient(app App, shopName, token string, opts ...Option) (*Client, error) {
	baseURL, err := url.Parse(ShopBaseUrl(shopName))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client: &http.Client{
			Timeout: time.Second * defaultHttpTimeout,
		},
		log:        &LeveledLogger{},
		app:        app,
		baseURL:    baseURL,
		token:      token,
		apiVersion: defaultApiVersion,
		pathPrefix: defaultApiPathPrefix,
	}

	c.Product = &ProductServiceOp{client: c}
	c.CustomCollection = &CustomCollectionServiceOp{client: c}
	c.SmartCollection = &SmartCollectionServiceOp{client: c}
	c.Customer = &CustomerServiceOp{client: c}
	c.CustomerAddress = &CustomerAddressServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Fulfillment = &FulfillmentServiceOp{client: c}
	c.DraftOrder = &DraftOrderServiceOp{client: c}
	c.AbandonedCheckout = &AbandonedCheckoutServiceOp{client: c}
	c.Shop = &ShopServiceOp{client: c}
	c.Webhook = &WebhookServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Image = &ImageServiceOp{client: c}
	c.Transaction = &TransactionServiceOp{client: c}
	c.Theme = &ThemeServiceOp{client: c}
	c.Asset = &AssetServiceOp{client: c}
	c.ScriptTag = &ScriptTagServiceOp{client: c}
	c.RecurringApplicationCharge = &RecurringApplicationChargeServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.Blog = &BlogServiceOp{client: c}
	c.ApplicationCharge = &ApplicationChargeServiceOp{client: c}
	c.Redirect = &RedirectServiceOp{client: c}
	c.Page = &PageServiceOp{client: c}
	c.StorefrontAccessToken = &StorefrontAccessTokenServiceOp{client: c}
	c.UsageCharge = &UsageChargeServiceOp{client: c}
	c.Collect = &CollectServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Location = &LocationServiceOp{client: c}
	c.DiscountCode = &DiscountCodeServiceOp{client: c}
	c.PriceRule = &PriceRuleServiceOp{client: c}
	c.InventoryItem = &InventoryItemServiceOp{client: c}
	c.ShippingZone = &ShippingZoneServiceOp{client: c}
	c.ProductListing = &ProductListingServiceOp{client: c}
	c.InventoryLevel = &InventoryLevelServiceOp{client: c}
	c.AccessScopes = &AccessScopesServiceOp{client: c}
	c.FulfillmentService = &FulfillmentServiceServiceOp{client: c}
	c.CarrierService = &CarrierServiceOp{client: c}
	c.Payouts = &PayoutsServiceOp{client: c}
	c.GiftCard = &GiftCardServiceOp{client: c}
	c.FulfillmentOrder = &FulfillmentOrderServiceOp{client: c}
	c.GraphQL = &GraphQLServiceOp{client: c}
	c.AssignedFulfillmentOrder = &AssignedFulfillmentOrderServiceOp{client: c}
	c.FulfillmentEvent = &FulfillmentEventServiceOp{client: c}
	c.FulfillmentRequest = &FulfillmentRequestServiceOp{client: c}
	c.PaymentsTransactions = &PaymentsTransactionsServiceOp{client: c}
	c.OrderRisk = &OrderRiskServiceOp{client: c}
	c.ApiPermissions = &ApiPermissionsServiceOp{client: c}
	c.Article = &ArticlesServiceOp{client: c}

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Do sends an API request and populates the given interface with the parsed
// response. It does not make much sense to call Do without a prepared
// interface instance.
func (c *Client) Do(req *http.Request, v interface{}) error {
	_, err := c.doGetHeaders(req, v)
	if err != nil {
		return err
	}

	return nil
}

// doGetHeaders executes a request, decoding the response into `v` and also returns any response headers.
func (c *Client) doGetHeaders(req *http.Request, v interface{}) (http.Header, error) {
	var resp *http.Response
	var err error
	retries := c.retries
	c.attempts = 0
	c.logRequest(req)

	// copy request body so it can be re-used
	var body []byte
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	for {
		c.attempts++
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		resp, err = c.Client.Do(req)
		c.logResponse(resp)
		if err != nil {
			return nil, err // http client errors, not api responses
		}

		respErr := CheckResponseError(resp)
		if respErr == nil {
			break // no errors, break out of the retry loop
		}

		// retry scenario, close resp and any continue will retry
		resp.Body.Close()

		if retries <= 1 {
			return nil, respErr
		}

		if rateLimitErr, isRetryErr := respErr.(RateLimitError); isRetryErr {
			// back off and retry

			wait := time.Duration(rateLimitErr.RetryAfter) * time.Second
			c.log.Debugf("rate limited waiting %s", wait.String())
			time.Sleep(wait)
			retries--
			continue
		}

		var doRetry bool
		switch resp.StatusCode {
		case http.StatusServiceUnavailable:
			c.log.Debugf("service unavailable, retrying")
			doRetry = true
			retries--
		}

		if doRetry {
			continue
		}

		// no retry attempts, just return the err
		return nil, respErr
	}

	defer resp.Body.Close()

	if c.apiVersion == defaultApiVersion && resp.Header.Get("X-Shopify-API-Version") != "" {
		// if using stable on first request set the api version
		c.apiVersion = resp.Header.Get("X-Shopify-API-Version")
		c.log.Infof("api version not set, now using %s", c.apiVersion)
	}

	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&v)
		if err != nil {
			return nil, err
		}
	}

	if s := strings.Split(resp.Header.Get("X-Shopify-Shop-Api-Call-Limit"), "/"); len(s) == 2 {
		c.RateLimits.RequestCount, _ = strconv.Atoi(s[0])
		c.RateLimits.BucketSize, _ = strconv.Atoi(s[1])
	}

	c.RateLimits.RetryAfterSeconds, _ = strconv.ParseFloat(resp.Header.Get("Retry-After"), 64)

	return resp.Header, nil
}

func (c *Client) logRequest(req *http.Request) {
	if req == nil {
		return
	}
	if req.URL != nil {
		c.log.Debugf("%s: %s", req.Method, req.URL.String())
	}
	c.logBody(&req.Body, "SENT: %s")
}

func (c *Client) logResponse(res *http.Response) {
	if res == nil {
		return
	}

	c.log.Debugf("Shopify X-Request-Id: %s", res.Header.Get("X-Request-Id"))
	c.log.Debugf("RECV %d: %s", res.StatusCode, res.Status)
	c.logBody(&res.Body, "RESP: %s")
}

func (c *Client) logBody(body *io.ReadCloser, format string) {
	if body == nil {
		return
	}
	b, err := ioutil.ReadAll(*body)
	if err != nil && !errors.Is(err, io.EOF) {
		return
	}
	if len(b) > 0 {
		c.log.Debugf(format, string(b))
	}
	*body = ioutil.NopCloser(bytes.NewBuffer(b))
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	// see https://www.shopify.dev/concepts/about-apis/response-codes
	if err.Status == http.StatusTooManyRequests {
		f, _ := strconv.ParseFloat(r.Header.Get("Retry-After"), 64)
		return RateLimitError{
			ResponseError: err,
			RetryAfter:    int(f),
		}
	}

	// if err.Status == http.StatusSeeOther {
	// todo
	// The response to the request can be found under a different URL in the
	// Location header and can be retrieved using a GET method on that resource.
	// }

	if err.Status == http.StatusNotAcceptable {
		err.Message = http.StatusText(err.Status)
	}

	return err
}

func CheckResponseError(r *http.Response) error {
	if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	// Create an anonoymous struct to parse the JSON data into.
	shopifyError := struct {
		Error  string      `json:"error"`
		Errors interface{} `json:"errors"`
	}{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// empty body, this probably means shopify returned an error with no body
	// we'll handle that error in wrapSpecificError()
	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &shopifyError)
		if err != nil {
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	}

	// Create the response error from the Shopify error.
	responseError := ResponseError{
		Status:  r.StatusCode,
		Message: shopifyError.Error,
	}

	// If the errors field is not filled out, we can return here.
	if shopifyError.Errors == nil {
		return wrapSpecificError(r, responseError)
	}

	// Shopify errors usually have the form:
	// {
	//   "errors": {
	//     "title": [
	//       "something is wrong"
	//     ]
	//   }
	// }
	// This structure is flattened to a single array:
	// [ "title: something is wrong" ]
	//
	// Unfortunately, "errors" can also be a single string so we have to deal
	// with that. Lots of reflection :-(
	switch reflect.TypeOf(shopifyError.Errors).Kind() {
	case reflect.String:
		// Single string, use as message
		responseError.Message = shopifyError.Errors.(string)
	case reflect.Slice:
		// An array, parse each entry as a string and join them on the message
		// json always serializes JSON arrays into []interface{}
		for _, elem := range shopifyError.Errors.([]interface{}) {
			responseError.Errors = append(responseError.Errors, fmt.Sprint(elem))
		}
		responseError.Message = strings.Join(responseError.Errors, ", ")
	case reflect.Map:
		// A map, parse each error for each key in the map.
		// json always serializes into map[string]interface{} for objects
		for k, v := range shopifyError.Errors.(map[string]interface{}) {
			switch reflect.TypeOf(v).Kind() {
			// Check to make sure the interface is a slice
			// json always serializes JSON arrays into []interface{}
			case reflect.Slice:
				for _, elem := range v.([]interface{}) {
					// If the primary message of the response error is not set, use
					// any message.
					if responseError.Message == "" {
						responseError.Message = fmt.Sprintf("%v: %v", k, elem)
					}
					topicAndElem := fmt.Sprintf("%v: %v", k, elem)
					responseError.Errors = append(responseError.Errors, topicAndElem)
				}
			case reflect.String:
				elem := v.(string)
				if responseError.Message == "" {
					responseError.Message = fmt.Sprintf("%v: %v", k, elem)
				}
				topicAndElem := fmt.Sprintf("%v: %v", k, elem)
				responseError.Errors = append(responseError.Errors, topicAndElem)
			}
		}
	}

	return wrapSpecificError(r, responseError)
}

// General list options that can be used for most collections of entities.
type ListOptions struct {
	// PageInfo is used with new pagination search.
	PageInfo string `url:"page_info,omitempty"`

	// Page is used to specify a specific page to load.
	// It is the deprecated way to do pagination.
	Page         int       `url:"page,omitempty"`
	Limit        int       `url:"limit,omitempty"`
	SinceId      *uint64   `url:"since_id,omitempty"`
	CreatedAtMin time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax time.Time `url:"updated_at_max,omitempty"`
	Order        string    `url:"order,omitempty"`
	Fields       string    `url:"fields,omitempty"`
	Vendor       string    `url:"vendor,omitempty"`
	Ids          []uint64  `url:"ids,omitempty,comma"`
}

// General count options that can be used for most collection counts.
type CountOptions struct {
	CreatedAtMin time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax time.Time `url:"updated_at_max,omitempty"`
}

func (c *Client) Count(ctx context.Context, path string, options interface{}) (int, error) {
	resource := struct {
		Count int `json:"count"`
	}{}
	err := c.Get(ctx, path, &resource, options)
	return resource.Count, err
}

// CreateAndDo performs a web request to Shopify with the given method (GET,
// POST, PUT, DELETE) and relative path (e.g. "/admin/orders.json").
// The data, options and resource arguments are optional and only relevant in
// certain situations.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
// The options argument is used for specifying request options such as search
// parameters like created_at_min
// Any data returned from Shopify will be marshalled into resource argument.
func (c *Client) CreateAndDo(ctx context.Context, method, relPath string, data, options, resource interface{}) error {
	_, err := c.createAndDoGetHeaders(ctx, method, relPath, data, options, resource)
	if err != nil {
		return err
	}
	return nil
}

// createAndDoGetHeaders creates an executes a request while returning the response headers.
func (c *Client) createAndDoGetHeaders(ctx context.Context, method, relPath string, data, options, resource interface{}) (http.Header, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	relPath = path.Join(c.pathPrefix, relPath)
	req, err := c.NewRequest(ctx, method, relPath, data, options)
	if err != nil {
		return nil, err
	}

	return c.doGetHeaders(req, resource)
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *Client) Get(ctx context.Context, path string, resource, options interface{}) error {
	return c.CreateAndDo(ctx, "GET", path, nil, options, resource)
}

// ListWithPagination performs a GET request for the given path and saves the result in the
// given resource and returns the pagination.
func (c *Client) ListWithPagination(ctx context.Context, path string, resource, options interface{}) (*Pagination, error) {
	headers, err := c.createAndDoGetHeaders(ctx, "GET", path, nil, options, resource)
	if err != nil {
		return nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}

// extractPagination extracts pagination info from linkHeader.
// Details on the format are here:
// https://shopify.dev/docs/api/admin-rest/latest/resources/paginated-rest-results
func extractPagination(linkHeader string) (*Pagination, error) {
	pagination := new(Pagination)

	if linkHeader == "" {
		return pagination, nil
	}

	for _, link := range strings.Split(linkHeader, ",") {
		match := linkRegex.FindStringSubmatch(link)
		// Make sure the link is not empty or invalid
		if len(match) != 3 {
			// We expect 3 values:
			// match[0] = full match
			// match[1] is the URL and match[2] is either 'previous' or 'next'
			err := ResponseDecodingError{
				Message: "could not extract pagination link header",
			}
			return nil, err
		}

		rel, err := url.Parse(match[1])
		if err != nil {
			err = ResponseDecodingError{
				Message: "pagination does not contain a valid URL",
			}
			return nil, err
		}

		params, err := url.ParseQuery(rel.RawQuery)
		if err != nil {
			return nil, err
		}

		paginationListOptions := ListOptions{}

		paginationListOptions.PageInfo = params.Get("page_info")
		if paginationListOptions.PageInfo == "" {
			err = ResponseDecodingError{
				Message: "page_info is missing",
			}
			return nil, err
		}

		limit := params.Get("limit")
		if limit != "" {
			paginationListOptions.Limit, err = strconv.Atoi(params.Get("limit"))
			if err != nil {
				return nil, err
			}
		}

		// 'rel' is either next or previous
		if match[2] == "next" {
			pagination.NextPageOptions = &paginationListOptions
		} else {
			pagination.PreviousPageOptions = &paginationListOptions
		}
	}

	return pagination, nil
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *Client) Post(ctx context.Context, path string, data, resource interface{}) error {
	return c.CreateAndDo(ctx, "POST", path, data, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *Client) Put(ctx context.Context, path string, data, resource interface{}) error {
	return c.CreateAndDo(ctx, "PUT", path, data, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.DeleteWithOptions(ctx, path, nil)
}

// DeleteWithOptions performs a DELETE request for the given path WithOptions
func (c *Client) DeleteWithOptions(ctx context.Context, path string, options interface{}) error {
	return c.CreateAndDo(ctx, "DELETE", path, nil, options, nil)
}
