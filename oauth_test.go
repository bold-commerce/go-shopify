package goshopify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAppAuthorizeUrl(t *testing.T) {
	setup()
	defer teardown()

	cases := []struct {
		shopName    string
		nonce       string
		expected    string
		errExpected string
	}{
		{
			"fooshop",
			"thenonce",
			"https://fooshop.myshopify.com/admin/oauth/authorize?client_id=apikey&redirect_uri=https%3A%2F%2Fexample.com%2Fcallback&scope=read_products&state=thenonce",
			"",
		},
		{
			"foo^^shop",
			"thenonce",
			"",
			`parse "https://foo^^shop.myshopify.com": invalid character "^" in host name`,
		},
	}

	for _, c := range cases {
		actual, err := app.AuthorizeUrl(c.shopName, c.nonce)
		if (c.errExpected == "" && err != nil) || (c.errExpected != "" && err.Error() != c.errExpected) {
			t.Fatalf("App.AuthorizeUrl(): err expected %s, err actual %s", c.errExpected, err)
		}
		if actual != c.expected {
			t.Fatalf("App.AuthorizeUrl(): expected %s, actual %s", c.expected, actual)
		}
	}
}

func TestAppGetAccessToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://requested-shop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"footoken"}`))

	app.Client = client
	token, err := app.GetAccessToken(context.Background(), "requested-shop", "foocode")
	if err != nil {
		t.Fatalf("App.GetAccessToken(): %v", err)
	}

	expected := "footoken"
	if token != expected {
		t.Errorf("Token = %v, expected %v", token, expected)
	}
}

func TestAppGetAccessTokenWithClientCredentialsGrant(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://requested-shop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"footoken","scope":"read_products","expires_in":86399}`))

	app.Client = client
	token, err := app.GetAccessTokenWithClientCredentialsGrant(context.Background(), "requested-shop")
	if err != nil {
		t.Fatalf("App.GetAccessTokenWithClientCredentialsGrant(): %v", err)
	}

	if token.AccessToken != "footoken" {
		t.Errorf("Token = %v, expected footoken", token.AccessToken)
	}
	if token.Scope != "read_products" {
		t.Errorf("Scope = %v, expected read_products", token.Scope)
	}
	if token.ExpiresIn != 86399 {
		t.Errorf("ExpiresIn = %v, expected 86399", token.ExpiresIn)
	}
}

func TestAppGetAccessTokenWithClientCredentialsGrantError(t *testing.T) {
	setup()
	defer teardown()

	// Case 1: client.Do returns error (e.g. 404)
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(404, `{"errors": "application_cannot_be_found"}`))

	app.Client = client
	expectedError := errors.New("application_cannot_be_found")

	token, err := app.GetAccessTokenWithClientCredentialsGrant(context.Background(), "fooshop")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != nil && token.AccessToken != "" {
		t.Errorf("Expected nil token received %v", token.AccessToken)
	}

	expectedError = errors.New("parse ://example.com: missing protocol scheme")
	oldPath := accessTokenRelPath
	accessTokenRelPath = "://example.com"
	defer func() { accessTokenRelPath = oldPath }()

	token, err = app.GetAccessTokenWithClientCredentialsGrant(context.Background(), "fooshop")
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != nil {
		t.Errorf("Expected nil token received %v", token)
	}
}

func TestAppGetAccessTokenWithClientCredentialsGrantNilClient(t *testing.T) {
	setup()
	defer teardown()

	oldPath := accessTokenRelPath
	accessTokenRelPath = "://example.com"
	defer func() { accessTokenRelPath = oldPath }()

	token, err := app.GetAccessTokenWithClientCredentialsGrant(context.Background(), "fooshop")
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Errorf("Expected missing protocol scheme error, got: %v", err)
	}
	if token != nil {
		t.Errorf("Expected nil token received %v", token)
	}
}

func TestAppGetAccessTokenError(t *testing.T) {
	setup()
	defer teardown()

	app.Client = client
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token", httpmock.NewStringResponder(404, `{"errors":"application_cannot_be_found"}`))
	expectedError := errors.New("application_cannot_be_found")

	token, err := app.GetAccessToken(context.Background(), "fooshop", "")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token received %s", token)
	}

	expectedError = errors.New("parse ://example.com: missing protocol scheme")
	oldPath := accessTokenRelPath
	accessTokenRelPath = "://example.com" // cause NewRequest to trip a parse error
	defer func() { accessTokenRelPath = oldPath }()
	token, err = app.GetAccessToken(context.Background(), "fooshop", "")
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token received %s", token)
	}
}

const expiringAccessTokenResponse = `{"access_token":"footoken","scope":"read_products","expires_in":3600,"refresh_token":"foorefresh","refresh_token_expires_in":7776000}`

// expiringAccessTokenResponder responds with an expiring access token and
// records the JSON body of the request in body.
func expiringAccessTokenResponder(body *map[string]string) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		if err := json.NewDecoder(req.Body).Decode(body); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(200, expiringAccessTokenResponse), nil
	}
}

func expiringAccessTokenTest(t *testing.T, token *ExpiringAccessToken, body, expectedBody map[string]string) {
	t.Helper()

	expectedToken := ExpiringAccessToken{
		AccessToken:           "footoken",
		Scope:                 "read_products",
		ExpiresIn:             3600,
		RefreshToken:          "foorefresh",
		RefreshTokenExpiresIn: 7776000,
	}
	if token == nil || *token != expectedToken {
		t.Errorf("Token = %+v, expected %+v", token, expectedToken)
	}
	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Request body = %v, expected %v", body, expectedBody)
	}
}

func TestAppGetExpiringAccessToken(t *testing.T) {
	setup()
	defer teardown()

	var body map[string]string
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		expiringAccessTokenResponder(&body))

	app.Client = client
	token, err := app.GetExpiringAccessToken(context.Background(), "fooshop", "foocode")
	if err != nil {
		t.Fatalf("App.GetExpiringAccessToken(): %v", err)
	}

	expiringAccessTokenTest(t, token, body, map[string]string{
		"client_id":     "apikey",
		"client_secret": "hush",
		"code":          "foocode",
		"expiring":      "1",
	})
}

func TestAppRefreshAccessToken(t *testing.T) {
	setup()
	defer teardown()

	var body map[string]string
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		expiringAccessTokenResponder(&body))

	app.Client = client
	token, err := app.RefreshAccessToken(context.Background(), "fooshop", "oldrefresh")
	if err != nil {
		t.Fatalf("App.RefreshAccessToken(): %v", err)
	}

	expiringAccessTokenTest(t, token, body, map[string]string{
		"client_id":     "apikey",
		"client_secret": "hush",
		"grant_type":    "refresh_token",
		"refresh_token": "oldrefresh",
	})
}

func TestAppRefreshAccessTokenError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(401, `{"error":"invalid_request","error_description":"This request requires an active refresh_token"}`))

	app.Client = client
	token, err := app.RefreshAccessToken(context.Background(), "fooshop", "expiredrefresh")

	var responseError ResponseError
	if !errors.As(err, &responseError) {
		t.Fatalf("Expected a ResponseError, got: %v", err)
	}
	if responseError.Status != http.StatusUnauthorized {
		t.Errorf("Status = %v, expected %v", responseError.Status, http.StatusUnauthorized)
	}
	if responseError.Message != "invalid_request" {
		t.Errorf("Message = %v, expected invalid_request", responseError.Message)
	}
	if token != nil {
		t.Errorf("Expected nil token received %v", token)
	}
}

func TestAppMigrateToExpiringAccessToken(t *testing.T) {
	setup()
	defer teardown()

	var body map[string]string
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		expiringAccessTokenResponder(&body))

	app.Client = client
	token, err := app.MigrateToExpiringAccessToken(context.Background(), "fooshop", "nonexpiringtoken")
	if err != nil {
		t.Fatalf("App.MigrateToExpiringAccessToken(): %v", err)
	}

	expiringAccessTokenTest(t, token, body, map[string]string{
		"client_id":            "apikey",
		"client_secret":        "hush",
		"grant_type":           "urn:ietf:params:oauth:grant-type:token-exchange",
		"subject_token":        "nonexpiringtoken",
		"subject_token_type":   "urn:shopify:params:oauth:token-type:offline-access-token",
		"requested_token_type": "urn:shopify:params:oauth:token-type:offline-access-token",
		"expiring":             "1",
	})
}

func TestAppExpiringAccessTokenNilClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, expiringAccessTokenResponse))

	app := App{ApiKey: "apikey", ApiSecret: "hush"}
	token, err := app.GetExpiringAccessToken(context.Background(), "fooshop", "foocode")
	if err != nil {
		t.Fatal(err)
	}
	if token == nil || token.AccessToken != "footoken" || token.RefreshToken != "foorefresh" {
		t.Fatal("expected the access and refresh tokens from the requested shop")
	}
}

func TestAppVerifyAuthorizationURL(t *testing.T) {
	app := App{ApiSecret: "hush"}
	// These credentials are from the Shopify example page:
	// https://shopify.dev/docs/api/admin-rest/latest/resources/oauth#verification
	urlOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&timestamp=1337178173")
	urlOkWithState, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=7db6973c2aff68295ebcf354c2ce528a6b09aef1146baafccc2e0b369fff5f6d&shop=some-shop.myshopify.com&timestamp=1337178173&state=abcd")
	urlNotOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&timestamp=133717817")

	cases := []struct {
		u        *url.URL
		expected bool
	}{
		{urlOk, true},
		{urlOkWithState, true},
		{urlNotOk, false},
	}

	for _, c := range cases {
		actual, err := app.VerifyAuthorizationURL(c.u)
		if err != nil {
			t.Errorf("App.VerifyAuthorizationURL(..., %s) returned an error: %v", c.u, err)
		}
		if actual != c.expected {
			t.Errorf("App.VerifyAuthorizationURL(..., %s): expected %v, actual %v", c.u, c.expected, actual)
		}
	}
}

func TestSignature(t *testing.T) {
	setup()
	defer teardown()

	// https://shopify.dev/tutorials/display-data-on-an-online-store-with-an-application-proxy-app-extension
	queryString := "extra=1&extra=2&shop=shop-name.myshopify.com&path_prefix=%2Fapps%2Fawesome_reviews&timestamp=1317327555&signature=a9718877bea71c2484f91608a7eaea1532bdf71f5c56825065fa4ccabe549ef3"

	urlOk, _ := url.Parse(fmt.Sprintf("http://example.com/proxied?%s", queryString))
	urlNotOk, _ := url.Parse(fmt.Sprintf("http://example.com/proxied?%s&notok=true", queryString))

	cases := []struct {
		u        *url.URL
		expected bool
	}{
		{urlOk, true},
		{urlNotOk, false},
	}

	for _, c := range cases {
		ok := app.VerifySignature(c.u)
		if ok != c.expected {
			t.Errorf("VerifySignature expected: |%v| but got: |%v|", c.expected, ok)
		}
	}
}

func TestVerifyWebhookRequest(t *testing.T) {
	setup()
	defer teardown()

	cases := []struct {
		hmac     string
		message  string
		expected bool
	}{
		{"hMTq0K2x7oyOjoBwGYeTj5oxfnaVYXzbanUG9aajpKI=", "my secret message", true},
		{"wronghash", "my secret message", false},
		{"wronghash", "", false},
		{"hMTq0K2x7oyOjoBwGYeTj5oxfnaVYXzbanUG9aajpKI=", "", false},
		{"", "", false},
		{"hMTq0K2x7oyOjoBwGYeTj5oxfnaVYXzbanUG9aajpKIthissignatureisnowwaytoolongohmanicantbelievehowlongthisis=", "my secret message", false},
	}

	for _, c := range cases {

		testClient := MustNewClient(App{}, "", "")
		req, err := testClient.NewRequest(context.Background(), "GET", "", c.message, nil)
		if err != nil {
			t.Fatalf("Webhook.verify err = %v, expected true", err)
		}
		if c.hmac != "" {
			req.Header.Add("X-Shopify-Hmac-Sha256", c.hmac)
		}

		isValid := app.VerifyWebhookRequest(req)

		if isValid != c.expected {
			t.Errorf("Webhook.verify was expecting %t got %t", c.expected, isValid)
		}
	}
}

func TestVerifyWebhookRequestVerbose(t *testing.T) {
	setup()
	defer teardown()

	var (
		shortHMAC                   = "YmxhaGJsYWgK"
		invalidBase64               = "XXXXXaGVsbG8="
		validHMACSignature          = "hMTq0K2x7oyOjoBwGYeTj5oxfnaVYXzbanUG9aajpKI="
		validHMACSignatureEmptyBody = "ZAZ6P4c14f6v798OCPYCodtdf9g8Z+GfthdfCgyhUYg="
		longHMAC                    = "VGhpc2lzdGhlc29uZ3RoYXRuZXZlcmVuZHN5ZXNpdGdvZXNvbmFuZG9ubXlmcmllbmRzc29tZXBlb3BsZXN0YXJ0aW5nc2luZ2luZ2l0bm90a25vd2luZ3doYXRpdHdhc2FuZG5vd2NvbnRpbnVlc2luZ2luZ2l0Zm9yZXZlcmp1c3RiZWNhdXNlCg=="
		req                         *http.Request
		err                         error
	)
	shortHMACBytes, _ := base64.StdEncoding.DecodeString(shortHMAC)
	longHMACBytes, _ := base64.StdEncoding.DecodeString(longHMAC)

	cases := []struct {
		hmac          string
		message       string
		expected      bool
		expectedError error
	}{
		{validHMACSignature, "my secret message", true, nil},
		{invalidBase64, "my secret message", false, errors.New("illegal base64 data at input byte 12")},
		{shortHMAC, "my secret message", false, fmt.Errorf("received HMAC is not of length 32, it is of length %d", len(shortHMACBytes))},
		{longHMAC, "my secret message", false, fmt.Errorf("received HMAC is not of length 32, it is of length %d", len(longHMACBytes))},
		{shortHMAC, "", false, fmt.Errorf("received HMAC is not of length 32, it is of length %d", len(shortHMACBytes))},
		{validHMACSignature, "my invalid message", false, errors.New("webhook HMAC mismatch")},
		{"", "", false, fmt.Errorf("header %s not set", shopifyChecksumHeader)},
		{validHMACSignatureEmptyBody, "", false, errors.New("request body is empty")},
	}

	for _, c := range cases {

		testClient := MustNewClient(App{}, "", "")

		// We actually want to test nil body's, not ""
		if c.message == "" {
			req, err = testClient.NewRequest(context.Background(), "GET", "", nil, nil)
		} else {
			req, err = testClient.NewRequest(context.Background(), "GET", "", c.message, nil)
		}

		if err != nil {
			t.Fatalf("Webhook.verify err = %v, expected true", err)
		}

		// We actually want to test not sending the header, not empty headers
		if c.hmac != "" {
			req.Header.Add("X-Shopify-Hmac-Sha256", c.hmac)
		}

		isValid, err := app.VerifyWebhookRequestVerbose(req)
		if err == nil && c.expectedError != nil {
			t.Errorf("Expected error %s got nil", c.expectedError.Error())
		}
		if c.expectedError == nil && err != nil {
			t.Errorf("Expected nil got error %s", err.Error())
		}

		if isValid != c.expected {
			t.Errorf("Webhook.verify was expecting %t got %t", c.expected, isValid)
			if err != nil {
				t.Errorf("Error returned %s header passed is %s", err.Error(), c.hmac)
			} else {
				t.Errorf("Header passed is %s", c.hmac)
			}
		}

		if c.expectedError != nil && err.Error() != c.expectedError.Error() {
			t.Errorf("Expected error %s got error %s", c.expectedError.Error(), err.Error())
		}
	}

	// Other error cases
	if app.VerifyWebhookRequest(nil) {
		t.Fatal("accepted a missing request")
	}
	oldSecret := app.ApiSecret
	app.ApiSecret = ""
	if app.VerifyWebhookRequest(req) {
		t.Fatal("accepted an empty secret")
	}
	isValid, err := app.VerifyWebhookRequestVerbose(req)
	if err == nil || isValid == true || err.Error() != errors.New("ApiSecret is empty").Error() {
		t.Errorf("Expected error %s got nil or true", errors.New("ApiSecret is empty"))
	}

	app.ApiSecret = oldSecret

	req.Body = errReader{}
	if app.VerifyWebhookRequest(req) {
		t.Fatal("accepted an unreadable body")
	}
	isValid, err = app.VerifyWebhookRequestVerbose(req)
	if err == nil || isValid == true || err.Error() != errors.New("test-error").Error() {
		t.Errorf("Expected error %s got %s", errors.New("test-error"), err)
	}
}

func TestVerifySignedQueryInputs(t *testing.T) {
	app := App{ApiSecret: "hush"}
	const digest = "4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20"
	const message = "code=0907a61c0c8d55e99db179b68161bc00&shop=some-shop.myshopify.com&timestamp=1337178173"
	mac := hmac.New(sha256.New, nil)
	mac.Write([]byte(message))
	emptyKeyMAC := hex.EncodeToString(mac.Sum(nil))
	if (App{}).VerifyMessage(message, emptyKeyMAC) || app.VerifyMessage(message, digest+"zz") {
		t.Fatal("accepted invalid authentication")
	}
	if (App{}).VerifySignature(&url.URL{RawQuery: "signature=b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad"}) {
		t.Fatal("accepted empty proxy secret")
	}
	for _, extra := range []string{"&signature=extra", "&field=%zz", "&hmac=" + digest} {
		if ok, _ := app.VerifyAuthorizationURL(&url.URL{RawQuery: message + "&hmac=" + digest + extra}); ok {
			t.Fatal("accepted invalid callback", extra)
		}
	}
}

func TestWebhookBodySizeLimit(t *testing.T) {
	app := App{ApiSecret: "test"}
	for _, tc := range []struct {
		size        int
		knownLength bool
		valid       bool
	}{
		{maxWebhookBodyBytes, true, true},
		{maxWebhookBodyBytes + 1, true, false},
		{maxWebhookBodyBytes + 8, false, false},
	} {
		body := bytes.Repeat([]byte("x"), tc.size)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		if !tc.knownLength {
			req.ContentLength = -1
		}
		mac := hmac.New(sha256.New, []byte(app.ApiSecret))
		mac.Write(body)
		req.Header.Set(shopifyChecksumHeader, base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		valid, err := app.VerifyWebhookRequestVerbose(req)
		if valid != tc.valid || (err == nil) != tc.valid {
			t.Fatalf("size %d: valid=%v err=%v", tc.size, valid, err)
		}
		replayed, err := ioutil.ReadAll(req.Body)
		if err != nil || !bytes.Equal(body, replayed) {
			t.Fatalf("size %d: replay lost data: %v", tc.size, err)
		}
		req.Body.Close()
	}
}

func TestOnlineAuthorizationAndIDTokenExchange(t *testing.T) {
	setup()
	defer teardown()
	app.Client = client
	ctx := context.Background()
	authURL, err := app.AuthorizeOnlineURL("second-shop", "state")
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(authURL)
	if err != nil || u.Host != "second-shop.myshopify.com" || u.Query().Get("grant_options[]") != "per-user" || u.Query().Get("state") != "state" {
		t.Fatalf("online authorization URL = %s", authURL)
	}

	const online = `{"access_token":"online","expires_in":42,"scope":"read_orders,read_customers","associated_user_scope":"read_orders","associated_user":{"id":9007199254740993,"first_name":"A","last_name":"B","email":"a@example.com","email_verified":true,"account_owner":true,"collaborator":false,"locale":"en"}}`
	const endpoint = "https://second-shop.myshopify.com/admin/oauth/access_token"
	for _, tc := range []struct {
		name     string
		params   map[string]string
		response string
	}{
		{"online code", map[string]string{"code": "code"}, online},
		{"online ID token", map[string]string{"grant_type": tokenExchangeGrantType, "subject_token": "id-token", "subject_token_type": idTokenType, "requested_token_type": onlineAccessTokenType}, online},
		{"offline ID token", map[string]string{"grant_type": tokenExchangeGrantType, "subject_token": "id-token", "subject_token_type": idTokenType, "requested_token_type": offlineAccessTokenType, "expiring": "1"}, expiringAccessTokenResponse},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.params["client_id"], tc.params["client_secret"] = app.ApiKey, app.ApiSecret
			httpmock.RegisterResponder("POST", endpoint, func(req *http.Request) (*http.Response, error) {
				var got map[string]string
				if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
					return nil, err
				}
				if !reflect.DeepEqual(got, tc.params) {
					t.Errorf("grant fields = %v, want %v", got, tc.params)
				}
				if req.Header.Get("X-Shopify-Access-Token") != "" || req.Header.Get("Authorization") != "" {
					t.Error("sent API credentials with token request")
				}
				return httpmock.NewStringResponse(200, tc.response), nil
			})
			if tc.name == "offline ID token" {
				token, err := app.ExchangeIDTokenForExpiringAccessToken(ctx, "second-shop", "id-token")
				if err != nil || token.AccessToken != "footoken" || token.RefreshToken != "foorefresh" || token.ExpiresIn != 3600 {
					t.Fatalf("offline token = %+v, %v", token, err)
				}
				return
			}
			var token *OnlineAccessToken
			var err error
			if tc.name == "online code" {
				token, err = app.GetOnlineAccessToken(ctx, "second-shop", "code")
			} else {
				token, err = app.ExchangeIDTokenForOnlineAccessToken(ctx, "second-shop", "id-token")
			}
			if err != nil || token.AccessToken != "online" || token.ExpiresIn != 42 || token.Scope != "read_orders,read_customers" || token.AssociatedUserScope != "read_orders" || token.AssociatedUser.Id != 9007199254740993 || !token.AssociatedUser.EmailVerified || token.AssociatedUser.Email != "a@example.com" {
				t.Fatalf("online token = %+v, %v", token, err)
			}
		})
	}
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(400, `{"errors":"invalid token"}`))
	if token, err := app.ExchangeIDTokenForOnlineAccessToken(ctx, "second-shop", "expired"); err == nil || token != nil {
		t.Fatal("expected invalid token error")
	}
	if token, err := app.ExchangeIDTokenForExpiringAccessToken(ctx, "bad/shop", "id-token"); err == nil || token != nil {
		t.Fatal("expected invalid shop error")
	}
}

func TestOnlineAccessTokenRequiresUserAndToken(t *testing.T) {
	setup()
	defer teardown()
	app.Client = client
	for _, exchange := range []func(context.Context, string, string) (*OnlineAccessToken, error){
		app.GetOnlineAccessToken, app.ExchangeIDTokenForOnlineAccessToken,
	} {
		for _, response := range []string{`{}`, expiringAccessTokenResponse, `{"associated_user":{"id":42}}`, `{"access_token":"online","associated_user":{"id":0}}`} {
			httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token", httpmock.NewStringResponder(200, response))
			if token, err := exchange(context.Background(), "fooshop", "fixture"); err == nil || token != nil {
				t.Fatalf("expected a nil token and an error for incomplete online grant")
			}
		}
	}
}

func TestVerifiedWebhookReplaysOnlyVerifiedBytes(t *testing.T) {
	app := App{ApiSecret: "test"}
	const payload = "webhook body"
	body := &trackedBody{Reader: strings.NewReader(payload)}
	req := httptest.NewRequest("POST", "/", body)
	mac := hmac.New(sha256.New, []byte(app.ApiSecret))
	mac.Write([]byte(payload))
	req.Header.Set(shopifyChecksumHeader, base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	if valid, err := app.VerifyWebhookRequestVerbose(req); !valid || err != nil {
		t.Fatal("verification failed", err)
	}
	// Reading the replay must work even when the underlying reader is unavailable.
	body.Reader = errReader{}
	replayed, err := ioutil.ReadAll(req.Body)
	if err != nil || string(replayed) != payload {
		t.Fatalf("replayed body = %q, err = %v", replayed, err)
	}
	if err := req.Body.Close(); err != nil || body.closes != 1 {
		t.Fatalf("original body closed %d times, err = %v", body.closes, err)
	}
}

func TestTokenExchangeRecordsRateLimitsForItsShop(t *testing.T) {
	setup()
	defer teardown()
	app.Client = client
	client.retries = 0
	for _, tc := range []struct {
		shop   string
		status int
		limit  string
		want   int
	}{
		{"fooshop", 200, "7/40", 7},
		{"fooshop.myshopify.com", 429, "40/40", 40},
		{"different-shop", 200, "1/40", 40},
	} {
		endpoint := ShopBaseUrl(tc.shop) + "/admin/oauth/access_token"
		httpmock.RegisterResponder("POST", endpoint, func(*http.Request) (*http.Response, error) {
			body := expiringAccessTokenResponse
			if tc.status == 429 {
				body = `{"errors":"throttled"}`
			}
			response := httpmock.NewStringResponse(tc.status, body)
			response.Header.Set("X-Shopify-Shop-Api-Call-Limit", tc.limit)
			return response, nil
		})
		_, err := app.RefreshAccessToken(context.Background(), tc.shop, "refresh")
		if (err != nil) != (tc.status == 429) {
			t.Fatalf("%s: unexpected error %v", tc.shop, err)
		}
		limits := client.RateLimits()
		if limits.RequestCount != tc.want || limits.BucketSize != 40 {
			t.Fatalf("%s: got %d/%d, want %d/40", tc.shop, limits.RequestCount, limits.BucketSize, tc.want)
		}
	}
}
