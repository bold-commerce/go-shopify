// Package goshopify provides methods for making requests to Shopify's admin API.
package goshopify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
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

// Client manages communication with the Shopify API.
// A Client must not be copied after first use.
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
	retries int

	rateLimitsMu sync.RWMutex
	rateLimits   RateLimitInfo

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
	ApplicationCredit          ApplicationCreditService
	Balance                    BalanceService
	CancellationRequest        CancellationRequestService
	CollectionListing          CollectionListingService
	Comment                    CommentService
	Country                    CountryService
	Currency                   CurrencyService
	DeprecatedAPICall          DeprecatedAPICallService
	Dispute                    DisputeService
	DisputeEvidence            DisputeEvidenceService
	DisputeFileUpload          DisputeFileUploadService
	Event                      EventService
	LocationsForMove           LocationsForMoveService
	MarketingEvent             MarketingEventService
	MobilePlatformApplication  MobilePlatformApplicationService
	Policy                     PolicyService
	ProductResourceFeedback    ProductResourceFeedbackService
	Province                   ProvinceService
	Refund                     RefundService
	ResourceFeedback           ResourceFeedbackService
	TenderTransaction          TenderTransactionService
	User                       UserService
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

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

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

	c.ApplicationCredit = &ApplicationCreditServiceOp{client: c}
	c.Balance = &BalanceServiceOp{client: c}
	c.CancellationRequest = &CancellationRequestServiceOp{client: c}
	c.CollectionListing = &CollectionListingServiceOp{client: c}
	c.Comment = &CommentServiceOp{client: c}
	c.Country = &CountryServiceOp{client: c}
	c.Currency = &CurrencyServiceOp{client: c}
	c.DeprecatedAPICall = &DeprecatedAPICallServiceOp{client: c}
	c.Dispute = &DisputeServiceOp{client: c}
	c.DisputeEvidence = &DisputeEvidenceServiceOp{client: c}
	c.DisputeFileUpload = &DisputeFileUploadServiceOp{client: c}
	c.Event = &EventServiceOp{client: c}
	c.LocationsForMove = &LocationsForMoveServiceOp{client: c}
	c.MarketingEvent = &MarketingEventServiceOp{client: c}
	c.MobilePlatformApplication = &MobilePlatformApplicationServiceOp{client: c}
	c.Policy = &PolicyServiceOp{client: c}
	c.ProductResourceFeedback = &ProductResourceFeedbackServiceOp{client: c}
	c.Province = &ProvinceServiceOp{client: c}
	c.Refund = &RefundServiceOp{client: c}
	c.ResourceFeedback = &ResourceFeedbackServiceOp{client: c}
	c.TenderTransaction = &TenderTransactionServiceOp{client: c}
	c.User = &UserServiceOp{client: c}

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
	_, err := c.doGetHeaders(req, v, c.retryPolicy(readOnlyMethod(req.Method)), c.recordRateLimitHeaders)
	if err != nil {
		return err
	}

	return nil
}

func readOnlyMethod(method string) bool {
	return method == "" || method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

// A policy belongs to one operation, including both HTTP and GraphQL retries.
type retryPolicy struct {
	remaining int
	retry503  bool
}

func (c *Client) retryPolicy(retry503 bool) *retryPolicy {
	return &retryPolicy{remaining: c.retries, retry503: retry503}
}

// doGetHeaders executes a request, decoding the response into `v` and also returns any response headers.
func (c *Client) doGetHeaders(req *http.Request, v interface{}, retry *retryPolicy, recordRateLimits func(http.Header)) (http.Header, error) {
	resp, responseBody, err := c.doRequest(req, retry, recordRateLimits)
	if err != nil {
		return nil, err
	}
	if v != nil {
		if err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&v); err != nil {
			return nil, err
		}
	}
	return resp.Header, nil
}

func (c *Client) doRequest(req *http.Request, retry *retryPolicy, recordRateLimits func(http.Header)) (*http.Response, []byte, error) {
	var resp *http.Response
	var err error

	// copy request body so it can be re-used
	var body, responseBody []byte
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, nil, err
		}
	}

	c.logRequest(req, body)

	for {
		if req.Body != nil && req.Body != http.NoBody {
			req.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		resp, err = c.Client.Do(req)
		if err != nil {
			return nil, nil, err // http client errors, not api responses
		}

		recordRateLimits(resp.Header)
		responseBody, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		c.logResponse(req, resp, responseBody)

		if http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices {
			break
		}
		respErr := decodeResponseError(resp, responseBody)
		if retry.remaining <= 0 {
			return nil, nil, respErr
		}

		if rateLimitErr, isRetryErr := respErr.(RateLimitError); isRetryErr {
			// back off and retry

			wait := time.Duration(rateLimitErr.RetryAfter) * time.Second
			c.log.Debugf("rate limited waiting %s", wait.String())
			if err := waitForRetry(req.Context(), wait); err != nil {
				return nil, nil, err
			}
			retry.remaining--
			continue
		}

		// A write may have succeeded before Shopify returned the 503.
		if retry.retry503 && resp.StatusCode == http.StatusServiceUnavailable {
			c.log.Debugf("service unavailable, retrying")
			retry.remaining--
			continue
		}
		return nil, nil, respErr
	}

	if c.apiVersion == defaultApiVersion && resp.Header.Get("X-Shopify-API-Version") != "" {
		// if using stable on first request set the api version
		c.apiVersion = resp.Header.Get("X-Shopify-API-Version")
		c.log.Infof("api version not set, now using %s", c.apiVersion)
	}

	return resp, responseBody, nil
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	// see https://www.shopify.dev/concepts/about-apis/response-codes
	if err.Status == http.StatusTooManyRequests {
		f, _ := strconv.ParseFloat(r.Header.Get("Retry-After"), 64)
		return RateLimitError{
			ResponseError: err,
			RetryAfter:    int(math.Ceil(f)),
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return decodeResponseError(r, body)
}

func decodeResponseError(r *http.Response, bodyBytes []byte) error {
	var shopifyError struct {
		Error  string      `json:"error"`
		Errors interface{} `json:"errors"`
	}

	// empty body, this probably means shopify returned an error with no body
	// we'll handle that error in wrapSpecificError()
	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &shopifyError)
		if err != nil {
			if r.StatusCode == http.StatusTooManyRequests {
				return wrapSpecificError(r, ResponseError{Status: r.StatusCode, Message: strings.TrimSpace(string(bodyBytes))})
			}
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
	// Shopify can return a message, a list of messages, or errors keyed by field.
	switch details := shopifyError.Errors.(type) {
	case string:
		responseError.Message = details
	case []interface{}:
		for _, elem := range details {
			responseError.Errors = append(responseError.Errors, fmt.Sprint(elem))
		}
		responseError.Message = strings.Join(responseError.Errors, ", ")
	case map[string]interface{}:
		for field, value := range details {
			var messages []interface{}
			switch v := value.(type) {
			case []interface{}:
				messages = v
			case string:
				messages = []interface{}{v}
			}
			for _, elem := range messages {
				message := fmt.Sprintf("%s: %v", field, elem)
				if responseError.Message == "" {
					responseError.Message = message
				}
				responseError.Errors = append(responseError.Errors, message)
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
	_, err := c.createAndDoGetHeaders(ctx, method, relPath, data, options, resource, c.retryPolicy(readOnlyMethod(method)))
	if err != nil {
		return err
	}
	return nil
}

// createAndDoGetHeaders creates an executes a request while returning the response headers.
func (c *Client) createAndDoGetHeaders(ctx context.Context, method, relPath string, data, options, resource interface{}, retry *retryPolicy) (http.Header, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	relPath = path.Join(c.pathPrefix, relPath)
	req, err := c.NewRequest(ctx, method, relPath, data, options)
	if err != nil {
		return nil, err
	}

	return c.doGetHeaders(req, resource, retry, c.recordRateLimitHeaders)
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *Client) Get(ctx context.Context, path string, resource, options interface{}) error {
	return c.CreateAndDo(ctx, "GET", path, nil, options, resource)
}

// forEachPage calls fetchPage with the initial options, then each next-page cursor.
// The callback collects the results of each successful page.
func forEachPage(options interface{}, fetchPage func(interface{}) (*Pagination, error)) error {
	for {
		pagination, err := fetchPage(options)
		if err != nil {
			return err
		}
		if pagination.NextPageOptions == nil {
			return nil
		}
		options = pagination.NextPageOptions
	}
}

// ListWithPagination performs a GET request for the given path and saves the result in the
// given resource and returns the pagination.
func (c *Client) ListWithPagination(ctx context.Context, path string, resource, options interface{}) (*Pagination, error) {
	headers, err := c.createAndDoGetHeaders(ctx, "GET", path, nil, options, resource, c.retryPolicy(true))
	if err != nil {
		return nil, err
	}

	// Extract pagination info from header
	linkHeader := strings.Join(headers["Link"], ", ")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}

// extractPagination extracts pagination info from linkHeader.
// Details on the format are here:
// https://shopify.dev/docs/api/admin-rest/usage/pagination
func extractPagination(linkHeader string) (*Pagination, error) {
	pagination := new(Pagination)

	if linkHeader == "" {
		return pagination, nil
	}

	// Only commas outside <...> separate links. URLs can contain commas in fields.
	var links []string
	start := 0
	inURL := false
	for i, char := range linkHeader {
		switch char {
		case '<':
			inURL = true
		case '>':
			inURL = false
		case ',':
			if !inURL {
				links = append(links, linkHeader[start:i])
				start = i + 1
			}
		}
	}
	links = append(links, linkHeader[start:])

	for _, link := range links {
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
		paginationListOptions.Fields = params.Get("fields")

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
