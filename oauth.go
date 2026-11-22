package goshopify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Limit webhook request bodies to 10 MiB.
const maxWebhookBodyBytes = 10 * 1024 * 1024

const shopifyChecksumHeader = "X-Shopify-Hmac-Sha256"

var accessTokenRelPath = "admin/oauth/access_token"

var shopDomain = regexp.MustCompile(`(?i)^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.myshopify\.com$`)

// Returns a Shopify oauth authorization url for the given shopname and state.
//
// State is a unique value that can be used to check the authenticity during a
// callback from Shopify.
func (app App) AuthorizeUrl(shopName string, state string) (string, error) {
	return app.authorizeURL(shopName, state, false)
}

// AuthorizeOnlineURL builds the authorization URL for a per-user online access token.
func (app App) AuthorizeOnlineURL(shopName, state string) (string, error) {
	return app.authorizeURL(shopName, state, true)
}

func (app App) authorizeURL(shopName, state string, online bool) (string, error) {
	shopUrl, err := url.Parse(ShopBaseUrl(shopName))
	if err != nil {
		return "", err
	}
	if !shopDomain.MatchString(ShopFullName(shopName)) {
		return "", errors.New("invalid Shopify shop domain")
	}
	shopUrl.Path = "/admin/oauth/authorize"
	query := shopUrl.Query()
	query.Set("client_id", app.ApiKey)
	query.Set("redirect_uri", app.RedirectUrl)
	query.Set("scope", app.Scope)
	query.Set("state", state)
	if online {
		query.Set("grant_options[]", "per-user")
	}
	shopUrl.RawQuery = query.Encode()
	return shopUrl.String(), nil
}

type ClientCredentialsToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessTokenWithClientCredentialsGrant gets an access token using the client credentials grant type.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/client-credentials-grant
func (app App) GetAccessTokenWithClientCredentialsGrant(ctx context.Context, shopName string) (*ClientCredentialsToken, error) {
	var token ClientCredentialsToken
	if err := app.requestToken(ctx, shopName, map[string]string{"grant_type": "client_credentials"}, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (app App) GetAccessToken(ctx context.Context, shopName string, code string) (string, error) {
	var token struct {
		AccessToken string `json:"access_token"`
	}
	err := app.requestToken(ctx, shopName, map[string]string{"code": code}, &token)
	return token.AccessToken, err
}

const (
	tokenExchangeGrantType = "urn:ietf:params:oauth:grant-type:token-exchange"
	idTokenType            = "urn:ietf:params:oauth:token-type:id_token"
	onlineAccessTokenType  = "urn:shopify:params:oauth:token-type:online-access-token"
	offlineAccessTokenType = "urn:shopify:params:oauth:token-type:offline-access-token"
)

// ExpiringAccessToken is an expiring offline access token, together with the
// refresh token that renews it.
type ExpiringAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	// ExpiresIn is the lifetime of AccessToken in seconds.
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	// RefreshTokenExpiresIn is the lifetime of RefreshToken in seconds.
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in"`
}

// OnlineAccessToken authorizes Admin API calls for the associated user.
type OnlineAccessToken struct {
	AccessToken         string `json:"access_token"`
	Scope               string `json:"scope"`
	ExpiresIn           int    `json:"expires_in"`
	AssociatedUserScope string `json:"associated_user_scope"`
	AssociatedUser      struct {
		Id            uint64 `json:"id"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		AccountOwner  bool   `json:"account_owner"`
		Locale        string `json:"locale"`
		Collaborator  bool   `json:"collaborator"`
	} `json:"associated_user"`
}

// GetOnlineAccessToken exchanges the authorization code returned by AuthorizeOnlineURL.
func (app App) GetOnlineAccessToken(ctx context.Context, shopName, code string) (*OnlineAccessToken, error) {
	return app.requestOnlineAccessToken(ctx, shopName, map[string]string{"code": code})
}

// ExchangeIDTokenForOnlineAccessToken gets an online token for an embedded app user.
func (app App) ExchangeIDTokenForOnlineAccessToken(ctx context.Context, shopName, idToken string) (*OnlineAccessToken, error) {
	return app.requestOnlineAccessToken(ctx, shopName, tokenExchange(idToken, idTokenType, onlineAccessTokenType))
}

func (app App) requestOnlineAccessToken(ctx context.Context, shopName string, data map[string]string) (*OnlineAccessToken, error) {
	var token OnlineAccessToken
	if err := app.requestToken(ctx, shopName, data, &token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" || token.AssociatedUser.Id == 0 {
		return nil, errors.New("online access token response requires an access token and an associated user ID")
	}
	return &token, nil
}

// ExchangeIDTokenForExpiringAccessToken gets an expiring offline token for an embedded app.
func (app App) ExchangeIDTokenForExpiringAccessToken(ctx context.Context, shopName, idToken string) (*ExpiringAccessToken, error) {
	return app.requestExpiringAccessToken(ctx, shopName, tokenExchange(idToken, idTokenType, offlineAccessTokenType))
}

// GetExpiringAccessToken exchanges an authorization code for an expiring
// offline access token. It is the expiring counterpart of GetAccessToken.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/offline-access-tokens
func (app App) GetExpiringAccessToken(ctx context.Context, shopName string, code string) (*ExpiringAccessToken, error) {
	return app.requestExpiringAccessToken(ctx, shopName, map[string]string{"code": code, "expiring": "1"})
}

// RefreshAccessToken gets a new expiring offline access token using a refresh
// token. Shopify rotates the refresh token on every refresh, so the returned
// RefreshToken replaces the one that was passed in.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/offline-access-tokens
func (app App) RefreshAccessToken(ctx context.Context, shopName string, refreshToken string) (*ExpiringAccessToken, error) {
	return app.requestExpiringAccessToken(ctx, shopName, map[string]string{"grant_type": "refresh_token", "refresh_token": refreshToken})
}

// MigrateToExpiringAccessToken exchanges a non-expiring offline access token
// for an expiring one. Shopify invalidates the non-expiring token when the
// exchange succeeds.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/migrate-to-expiring-offline-access-tokens
func (app App) MigrateToExpiringAccessToken(ctx context.Context, shopName string, accessToken string) (*ExpiringAccessToken, error) {
	return app.requestExpiringAccessToken(ctx, shopName, tokenExchange(accessToken, offlineAccessTokenType, offlineAccessTokenType))
}

func tokenExchange(subject, subjectType, requestedType string) map[string]string {
	data := map[string]string{
		"grant_type":           tokenExchangeGrantType,
		"subject_token":        subject,
		"subject_token_type":   subjectType,
		"requested_token_type": requestedType,
	}
	if requestedType == offlineAccessTokenType {
		data["expiring"] = "1"
	}
	return data
}

func (app App) requestExpiringAccessToken(ctx context.Context, shopName string, data map[string]string) (*ExpiringAccessToken, error) {
	var token ExpiringAccessToken
	if err := app.requestToken(ctx, shopName, data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (app App) requestToken(ctx context.Context, shopName string, data map[string]string, result interface{}) error {
	if !shopDomain.MatchString(ShopFullName(shopName)) {
		return errors.New("invalid Shopify shop domain")
	}
	// Use the requested shop and keep existing API credentials out of the exchange.
	client, err := NewClient(App{}, shopName, "")
	if err != nil {
		return err
	}
	recordRateLimits := client.recordRateLimitHeaders
	if app.Client != nil {
		client.Client = app.Client.Client
		client.retries = app.Client.retries
		// A transport can serve several shops, but each shop has its own quota.
		if app.Client.baseURL != nil && strings.EqualFold(app.Client.baseURL.Host, client.baseURL.Host) {
			recordRateLimits = app.Client.recordRateLimitHeaders
		}
	}
	// A redirect can forward the grant and the client secret to another endpoint.
	httpClient := *client.Client
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	client.Client = &httpClient

	data["client_id"] = app.ApiKey
	data["client_secret"] = app.ApiSecret
	req, err := client.NewRequest(ctx, http.MethodPost, accessTokenRelPath, data, nil)
	if err != nil {
		return err
	}
	_, err = client.doGetHeaders(req, result, client.retryPolicy(false), recordRateLimits)
	return err
}

// Verify a message against a message HMAC
func (app App) VerifyMessage(message, messageMAC string) bool {
	if app.ApiSecret == "" || len(messageMAC) != hex.EncodedLen(sha256.Size) {
		return false
	}
	actualMac, err := hex.DecodeString(messageMAC)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(app.ApiSecret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(actualMac, expectedMAC)
}

// Verifying URL callback parameters.
func (app App) VerifyAuthorizationURL(u *url.URL) (bool, error) {
	if u == nil {
		return false, errors.New("missing callback URL")
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return false, err
	}
	if len(q["hmac"]) != 1 {
		return false, nil
	}
	messageMAC := q.Get("hmac")

	// Shopify signs every callback parameter except hmac.
	q.Del("hmac")

	message, err := url.QueryUnescape(q.Encode())

	return app.VerifyMessage(message, messageMAC), err
}

// Verifies a webhook http request, sent by Shopify.
// The body of the request is still readable after invoking the method.
func (app App) VerifyWebhookRequest(httpRequest *http.Request) bool {
	valid, _ := app.VerifyWebhookRequestVerbose(httpRequest)
	return valid
}

// Verifies a webhook http request, sent by Shopify.
// The body of the request is still readable after invoking the method.
// Bodies larger than 10 MiB are rejected.
func (app App) VerifyWebhookRequestVerbose(httpRequest *http.Request) (bool, error) {
	if app.ApiSecret == "" {
		return false, errors.New("ApiSecret is empty")
	}

	if httpRequest == nil || httpRequest.Body == nil {
		return false, errors.New("request body is empty")
	}

	shopifySha256 := httpRequest.Header.Get(shopifyChecksumHeader)
	if shopifySha256 == "" {
		return false, fmt.Errorf("header %s not set", shopifyChecksumHeader)
	}

	decodedReceivedHMAC, err := base64.StdEncoding.DecodeString(shopifySha256)
	if err != nil {
		return false, err
	}
	if len(decodedReceivedHMAC) != 32 {
		return false, fmt.Errorf("received HMAC is not of length 32, it is of length %d", len(decodedReceivedHMAC))
	}

	mac := hmac.New(sha256.New, []byte(app.ApiSecret))
	if httpRequest.ContentLength > maxWebhookBodyBytes {
		return false, errors.New("webhook body exceeds 10 MiB")
	}
	body := httpRequest.Body
	requestBody, err := ioutil.ReadAll(io.LimitReader(body, maxWebhookBodyBytes+1))
	var replay io.Reader = bytes.NewReader(requestBody)
	if err != nil || len(requestBody) > maxWebhookBodyBytes {
		// A partial read can leave unread bytes in the original body.
		replay = io.MultiReader(replay, body)
	}
	httpRequest.Body = struct {
		io.Reader
		io.Closer
	}{replay, body}
	if err != nil {
		return false, err
	}
	if len(requestBody) > maxWebhookBodyBytes {
		return false, errors.New("webhook body exceeds 10 MiB")
	}
	if len(requestBody) == 0 {
		return false, errors.New("request body is empty")
	}

	// Sha256 write doesn't actually return an error
	mac.Write(requestBody)

	if !hmac.Equal(decodedReceivedHMAC, mac.Sum(nil)) {
		return false, errors.New("webhook HMAC mismatch")
	}
	return true, nil
}

// Verifies an app proxy request, sent by Shopify.
// When Shopify proxies HTTP requests to the proxy URL,
// Shopify adds a signature paramter that is used to verify that the request was sent by Shopify.
// https://shopify.dev/tutorials/display-dynamic-store-data-with-app-proxies
func (app App) VerifySignature(u *url.URL) bool {
	return app.verifySignatureAt(u, time.Now().Unix())
}

func (app App) verifySignatureAt(u *url.URL, now int64) bool {
	if u == nil {
		return false
	}
	val, err := url.ParseQuery(u.RawQuery)
	if err != nil || len(val["signature"]) != 1 {
		return false
	}
	if len(val["timestamp"]) != 1 || len(val["shop"]) > 1 || len(val["hmac"]) > 1 {
		return false
	}
	timestamp, err := strconv.ParseInt(val.Get("timestamp"), 10, 64)
	// Match the 90-second clock tolerance of Shopify's app proxy validator.
	if err != nil || timestamp < now-90 || timestamp > now+90 {
		return false
	}
	sig := val.Get("signature")
	val.Del("signature")

	keys := []string{}
	for k, v := range val {
		keys = append(keys, fmt.Sprintf("%s=%s", k, strings.Join(v, ",")))
	}
	sort.Strings(keys)

	joined := strings.Join(keys, "")

	return app.VerifyMessage(joined, sig)
}
