package goshopify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAppAuthorizeUrl(t *testing.T) {
	setup()
	defer teardown()

	cases := []struct {
		shopName string
		nonce    string
		expected string
	}{
		{"fooshop", "thenonce", "https://fooshop.myshopify.com/admin/oauth/authorize?client_id=apikey&redirect_uri=https%3A%2F%2Fexample.com%2Fcallback&scope=read_products&state=thenonce"},
	}

	for _, c := range cases {
		actual := app.AuthorizeUrl(c.shopName, c.nonce)
		if actual != c.expected {
			t.Errorf("App.AuthorizeUrl(): expected %s, actual %s", c.expected, actual)
		}
	}
}

func TestAppGetAccessToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"footoken"}`))

	app.Client = client
	token, err := app.GetAccessToken("fooshop", "foocode")

	if err != nil {
		t.Fatalf("App.GetAccessToken(): %v", err)
	}

	expected := "footoken"
	if token != expected {
		t.Errorf("Token = %v, expected %v", token, expected)
	}
}

func TestAppGetAccessTokenError(t *testing.T) {
	setup()
	defer teardown()

	// app.Client isn't specified so NewClient called
	expectedError := errors.New("invalid_request")

	token, err := app.GetAccessToken("fooshop", "")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token received %s", token)
	}

	expectedError = errors.New("parse ://example.com: missing protocol scheme")
	accessTokenRelPath = "://example.com" // cause NewRequest to trip a parse error
	token, err = app.GetAccessToken("fooshop", "")
	if err == nil || !strings.Contains(err.Error(), "missing protocol scheme") {
		t.Errorf("Expected error %s got error %s", expectedError.Error(), err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token received %s", token)
	}
}

func TestAppVerifyAuthorizationURL(t *testing.T) {
	// These credentials are from the Shopify example page:
	// https://help.shopify.com/api/guides/authentication/oauth#verification
	urlOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&signature=11813d1e7bbf4629edcda0628a3f7a20&timestamp=1337178173")
	urlOkWithState, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=7db6973c2aff68295ebcf354c2ce528a6b09aef1146baafccc2e0b369fff5f6d&shop=some-shop.myshopify.com&signature=11813d1e7bbf4629edcda0628a3f7a20&timestamp=1337178173&state=abcd")
	urlNotOk, _ := url.Parse("http://example.com/callback?code=0907a61c0c8d55e99db179b68161bc00&hmac=4712bf92ffc2917d15a2f5a273e39f0116667419aa4b6ac0b3baaf26fa3c4d20&shop=some-shop.myshopify.com&signature=11813d1e7bbf4629edcda0628a3f7a20&timestamp=133717817")

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

		testClient := NewClient(App{}, "", "")
		req, err := testClient.NewRequest("GET", "", c.message, nil)
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
	validHMACBytes, _ := base64.StdEncoding.DecodeString(validHMACSignature)
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
		{validHMACSignature, "my invalid message", false, fmt.Errorf("expected hash %s does not equal %x", "ac3560a67dd9ed3f46cf1807856b78f27c9543b5ae98f5292d8eccf6252254f0", validHMACBytes)},
		{"", "", false, fmt.Errorf("header %s not set", shopifyChecksumHeader)},
		{validHMACSignatureEmptyBody, "", false, errors.New("request body is empty")},
	}

	for _, c := range cases {

		testClient := NewClient(App{}, "", "")

		// We actually want to test nil body's, not ""
		if c.message == "" {
			req, err = testClient.NewRequest("GET", "", nil, nil)
		} else {
			req, err = testClient.NewRequest("GET", "", c.message, nil)
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
	oldSecret := app.ApiSecret
	app.ApiSecret = ""
	isValid, err := app.VerifyWebhookRequestVerbose(req)
	if err == nil || isValid == true || err.Error() != errors.New("ApiSecret is empty").Error() {
		t.Errorf("Expected error %s got nil or true", errors.New("ApiSecret is empty"))
	}

	app.ApiSecret = oldSecret

	req.Body = errReader{}
	isValid, err = app.VerifyWebhookRequestVerbose(req)
	if err == nil || isValid == true || err.Error() != errors.New("test-error").Error() {
		t.Errorf("Expected error %s got %s", errors.New("test-error"), err)
	}

}
