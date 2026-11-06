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
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const shopifyChecksumHeader = "X-Shopify-Hmac-Sha256"

var accessTokenRelPath = "admin/oauth/access_token"

// Returns a Shopify oauth authorization url for the given shopname and state.
//
// State is a unique value that can be used to check the authenticity during a
// callback from Shopify.
func (app App) AuthorizeUrl(shopName string, state string) (string, error) {
	shopUrl, err := url.Parse(ShopBaseUrl(shopName))
	if err != nil {
		return "", err
	}
	shopUrl.Path = "/admin/oauth/authorize"
	query := shopUrl.Query()
	query.Set("client_id", app.ApiKey)
	query.Set("redirect_uri", app.RedirectUrl)
	query.Set("scope", app.Scope)
	query.Set("state", state)
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
	data := struct {
		GrantType    string `json:"grant_type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		GrantType:    "client_credentials",
		ClientId:     app.ApiKey,
		ClientSecret: app.ApiSecret,
	}

	client := app.Client
	if client == nil {
		client = MustNewClient(app, shopName, "")
	}

	req, err := client.NewRequest(ctx, "POST", accessTokenRelPath, data, nil)
	if err != nil {
		return nil, err
	}
	clientCredentialsToken := &ClientCredentialsToken{}
	err = client.Do(req, clientCredentialsToken)
	return clientCredentialsToken, err
}

func (app App) GetAccessToken(ctx context.Context, shopName string, code string) (string, error) {
	type Token struct {
		Token string `json:"access_token"`
	}

	data := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}{
		ClientId:     app.ApiKey,
		ClientSecret: app.ApiSecret,
		Code:         code,
	}

	client := app.Client
	if client == nil {
		client = MustNewClient(app, shopName, "")
	}

	req, err := client.NewRequest(ctx, "POST", accessTokenRelPath, data, nil)
	if err != nil {
		return "", err
	}

	token := new(Token)
	err = client.Do(req, token)
	return token.Token, err
}

const (
	tokenExchangeGrantType = "urn:ietf:params:oauth:grant-type:token-exchange"
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

// GetExpiringAccessToken exchanges an authorization code for an expiring
// offline access token. It is the expiring counterpart of GetAccessToken.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/offline-access-tokens
func (app App) GetExpiringAccessToken(ctx context.Context, shopName string, code string) (*ExpiringAccessToken, error) {
	data := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		Expiring     string `json:"expiring"`
	}{
		ClientId:     app.ApiKey,
		ClientSecret: app.ApiSecret,
		Code:         code,
		Expiring:     "1",
	}

	return app.requestExpiringAccessToken(ctx, shopName, data)
}

// RefreshAccessToken gets a new expiring offline access token using a refresh
// token. Shopify rotates the refresh token on every refresh, so the returned
// RefreshToken replaces the one that was passed in.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/offline-access-tokens
func (app App) RefreshAccessToken(ctx context.Context, shopName string, refreshToken string) (*ExpiringAccessToken, error) {
	data := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}{
		ClientId:     app.ApiKey,
		ClientSecret: app.ApiSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	return app.requestExpiringAccessToken(ctx, shopName, data)
}

// MigrateToExpiringAccessToken exchanges a non-expiring offline access token
// for an expiring one. Shopify invalidates the non-expiring token when the
// exchange succeeds.
// Docs: https://shopify.dev/docs/apps/build/authentication-authorization/migrate-to-expiring-offline-access-tokens
func (app App) MigrateToExpiringAccessToken(ctx context.Context, shopName string, accessToken string) (*ExpiringAccessToken, error) {
	data := struct {
		ClientId           string `json:"client_id"`
		ClientSecret       string `json:"client_secret"`
		GrantType          string `json:"grant_type"`
		SubjectToken       string `json:"subject_token"`
		SubjectTokenType   string `json:"subject_token_type"`
		RequestedTokenType string `json:"requested_token_type"`
		Expiring           string `json:"expiring"`
	}{
		ClientId:           app.ApiKey,
		ClientSecret:       app.ApiSecret,
		GrantType:          tokenExchangeGrantType,
		SubjectToken:       accessToken,
		SubjectTokenType:   offlineAccessTokenType,
		RequestedTokenType: offlineAccessTokenType,
		Expiring:           "1",
	}

	return app.requestExpiringAccessToken(ctx, shopName, data)
}

func (app App) requestExpiringAccessToken(ctx context.Context, shopName string, data interface{}) (*ExpiringAccessToken, error) {
	client := app.Client
	if client == nil {
		client = MustNewClient(app, shopName, "")
	}

	req, err := client.NewRequest(ctx, "POST", accessTokenRelPath, data, nil)
	if err != nil {
		return nil, err
	}

	token := new(ExpiringAccessToken)
	if err := client.Do(req, token); err != nil {
		return nil, err
	}
	return token, nil
}

// Verify a message against a message HMAC
func (app App) VerifyMessage(message, messageMAC string) bool {
	mac := hmac.New(sha256.New, []byte(app.ApiSecret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	// shopify HMAC is in hex so it needs to be decoded
	actualMac, _ := hex.DecodeString(messageMAC)

	return hmac.Equal(actualMac, expectedMAC)
}

// Verifying URL callback parameters.
func (app App) VerifyAuthorizationURL(u *url.URL) (bool, error) {
	q := u.Query()
	messageMAC := q.Get("hmac")

	// Remove hmac and signature and leave the rest of the parameters alone.
	q.Del("hmac")
	q.Del("signature")

	message, err := url.QueryUnescape(q.Encode())

	return app.VerifyMessage(message, messageMAC), err
}

// Verifies a webhook http request, sent by Shopify.
// The body of the request is still readable after invoking the method.
func (app App) VerifyWebhookRequest(httpRequest *http.Request) bool {
	shopifySha256 := httpRequest.Header.Get(shopifyChecksumHeader)
	actualMac := []byte(shopifySha256)

	mac := hmac.New(sha256.New, []byte(app.ApiSecret))
	requestBody, _ := ioutil.ReadAll(httpRequest.Body)
	httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	mac.Write(requestBody)
	macSum := mac.Sum(nil)
	expectedMac := []byte(base64.StdEncoding.EncodeToString(macSum))

	return hmac.Equal(actualMac, expectedMac)
}

// Verifies a webhook http request, sent by Shopify.
// The body of the request is still readable after invoking the method.
// This method has more verbose error output which is useful for debugging.
func (app App) VerifyWebhookRequestVerbose(httpRequest *http.Request) (bool, error) {
	if app.ApiSecret == "" {
		return false, errors.New("ApiSecret is empty")
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
	requestBody, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		return false, err
	}

	httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	if len(requestBody) == 0 {
		return false, errors.New("request body is empty")
	}

	// Sha256 write doesn't actually return an error
	mac.Write(requestBody)

	computedHMAC := mac.Sum(nil)
	HMACSame := hmac.Equal(decodedReceivedHMAC, computedHMAC)
	if !HMACSame {
		return HMACSame, fmt.Errorf("expected hash %x does not equal %x", computedHMAC, decodedReceivedHMAC)
	}

	return HMACSame, nil
}

// Verifies an app proxy request, sent by Shopify.
// When Shopify proxies HTTP requests to the proxy URL,
// Shopify adds a signature paramter that is used to verify that the request was sent by Shopify.
// https://shopify.dev/tutorials/display-dynamic-store-data-with-app-proxies
func (app App) VerifySignature(u *url.URL) bool {
	val := u.Query()
	sig := val.Get("signature")
	val.Del("signature")

	keys := []string{}
	for k, v := range val {
		keys = append(keys, fmt.Sprintf("%s=%s", k, strings.Join(v, ",")))
	}
	sort.Strings(keys)

	joined := strings.Join(keys, "")

	return hmacSHA256([]byte(app.ApiSecret), []byte(joined), []byte(sig))
}

func hmacSHA256(key, body, expected []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(body)
	result := mac.Sum(nil)

	dst := make([]byte, hex.EncodedLen(len(result)))
	hex.Encode(dst, result)

	return hmac.Equal(dst, expected)
}
