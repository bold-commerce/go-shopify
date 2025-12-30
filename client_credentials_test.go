package goshopify

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestAppGetClientCredentialsToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{
			"access_token": "client_credentials_token",
			"token_type": "Bearer",
			"expires_in": 86400,
			"scope": "read_products"
		}`))

	app.Client = client
	response, err := app.GetClientCredentialsToken(context.Background(), "fooshop")
	if err != nil {
		t.Fatalf("App.GetClientCredentialsToken(): %v", err)
	}

	if response.AccessToken != "client_credentials_token" {
		t.Errorf("AccessToken = %v, expected %v", response.AccessToken, "client_credentials_token")
	}
	if response.TokenType != "Bearer" {
		t.Errorf("TokenType = %v, expected %v", response.TokenType, "Bearer")
	}
	if response.ExpiresIn != 86400 {
		t.Errorf("ExpiresIn = %v, expected %v", response.ExpiresIn, 86400)
	}
	if response.Scope != "read_products" {
		t.Errorf("Scope = %v, expected %v", response.Scope, "read_products")
	}
}

func TestAppGetClientCredentialsTokenError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(401, `{"error": "invalid_client"}`))

	app.Client = client
	response, err := app.GetClientCredentialsToken(context.Background(), "fooshop")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if response != nil {
		t.Errorf("Expected nil response, got %v", response)
	}
}

func TestNewTokenManager(t *testing.T) {
	testApp := App{
		ApiKey:    "test-key",
		ApiSecret: "test-secret",
	}

	tm := NewTokenManager(testApp, "testshop")

	if tm.app.ApiKey != "test-key" {
		t.Errorf("ApiKey = %v, expected %v", tm.app.ApiKey, "test-key")
	}
	if tm.shopName != "testshop" {
		t.Errorf("shopName = %v, expected %v", tm.shopName, "testshop")
	}
	if tm.refreshBuffer != defaultRefreshBuffer {
		t.Errorf("refreshBuffer = %v, expected %v", tm.refreshBuffer, defaultRefreshBuffer)
	}
}

func TestNewTokenManagerWithOptions(t *testing.T) {
	testApp := App{
		ApiKey:    "test-key",
		ApiSecret: "test-secret",
	}

	customBuffer := 30 * time.Minute
	tm := NewTokenManager(testApp, "testshop", WithRefreshBuffer(customBuffer))

	if tm.refreshBuffer != customBuffer {
		t.Errorf("refreshBuffer = %v, expected %v", tm.refreshBuffer, customBuffer)
	}
}

func TestTokenManagerGetAccessToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{
			"access_token": "test_token_123",
			"token_type": "Bearer",
			"expires_in": 86400,
			"scope": "read_products"
		}`))

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	token, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("TokenManager.GetAccessToken(): %v", err)
	}

	if token != "test_token_123" {
		t.Errorf("Token = %v, expected %v", token, "test_token_123")
	}
}

func TestTokenManagerGetAccessTokenCached(t *testing.T) {
	setup()
	defer teardown()

	var callCount int
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			callCount++
			return httpmock.NewStringResponse(200, `{
				"access_token": "cached_token",
				"token_type": "Bearer",
				"expires_in": 86400,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	// First call should fetch token
	token1, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("First GetAccessToken(): %v", err)
	}

	// Second call should use cached token
	token2, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("Second GetAccessToken(): %v", err)
	}

	if token1 != token2 {
		t.Errorf("Tokens should be equal: %v != %v", token1, token2)
	}

	if callCount != 1 {
		t.Errorf("API should be called once, was called %d times", callCount)
	}
}

func TestTokenManagerForceRefresh(t *testing.T) {
	setup()
	defer teardown()

	var callCount int
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			callCount++
			return httpmock.NewStringResponse(200, `{
				"access_token": "refreshed_token",
				"token_type": "Bearer",
				"expires_in": 86400,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	// Get initial token
	_, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken(): %v", err)
	}

	// Force refresh
	token, err := tm.ForceRefresh(context.Background())
	if err != nil {
		t.Fatalf("ForceRefresh(): %v", err)
	}

	if token != "refreshed_token" {
		t.Errorf("Token = %v, expected %v", token, "refreshed_token")
	}

	if callCount != 2 {
		t.Errorf("API should be called twice, was called %d times", callCount)
	}
}

func TestTokenManagerGetTokenInfo(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{
			"access_token": "info_token",
			"token_type": "Bearer",
			"expires_in": 86400,
			"scope": "read_products"
		}`))

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	// Before fetching token
	info := tm.GetTokenInfo()
	if info.HasToken {
		t.Error("HasToken should be false before fetching")
	}
	if info.IsValid {
		t.Error("IsValid should be false before fetching")
	}

	// Fetch token
	_, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken(): %v", err)
	}

	// After fetching token
	info = tm.GetTokenInfo()
	if !info.HasToken {
		t.Error("HasToken should be true after fetching")
	}
	if !info.IsValid {
		t.Error("IsValid should be true after fetching")
	}
	if info.ExpiresInSeconds <= 0 {
		t.Errorf("ExpiresInSeconds should be positive, got %d", info.ExpiresInSeconds)
	}
}

func TestTokenManagerClearToken(t *testing.T) {
	setup()
	defer teardown()

	var callCount int
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			callCount++
			return httpmock.NewStringResponse(200, `{
				"access_token": "clear_test_token",
				"token_type": "Bearer",
				"expires_in": 86400,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	// Fetch token
	_, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken(): %v", err)
	}

	// Clear token
	tm.ClearToken()

	info := tm.GetTokenInfo()
	if info.HasToken {
		t.Error("HasToken should be false after clearing")
	}

	// Next GetAccessToken should fetch new token
	_, err = tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken() after clear: %v", err)
	}

	if callCount != 2 {
		t.Errorf("API should be called twice after clear, was called %d times", callCount)
	}
}

func TestTokenManagerConcurrentAccess(t *testing.T) {
	setup()
	defer teardown()

	var callCount int32
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&callCount, 1)
			// Simulate slow response
			time.Sleep(10 * time.Millisecond)
			return httpmock.NewStringResponse(200, `{
				"access_token": "concurrent_token",
				"token_type": "Bearer",
				"expires_in": 86400,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	var wg sync.WaitGroup
	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := tm.GetAccessToken(context.Background())
			if err != nil {
				errors <- err
				return
			}
			if token != "concurrent_token" {
				errors <- fmt.Errorf("unexpected token: %s", token)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}

	// Due to double-check locking, the API should only be called once
	// (or at most a few times if goroutines race before the first token is cached)
	if atomic.LoadInt32(&callCount) > 3 {
		t.Errorf("API called too many times in concurrent scenario: %d", callCount)
	}
}

func TestTokenManagerTokenExpiry(t *testing.T) {
	setup()
	defer teardown()

	var callCount int
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			callCount++
			return httpmock.NewStringResponse(200, `{
				"access_token": "expiry_token",
				"token_type": "Bearer",
				"expires_in": 2,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	// Set very small refresh buffer for testing
	tm := NewTokenManager(app, "fooshop", WithRefreshBuffer(1*time.Second))

	// First call
	_, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("First GetAccessToken(): %v", err)
	}

	// Token should be valid
	if !tm.GetTokenInfo().IsValid {
		t.Error("Token should be valid immediately after fetch")
	}

	// Wait for token to enter refresh buffer period
	time.Sleep(1500 * time.Millisecond)

	// Token should be considered invalid now (within refresh buffer)
	if tm.GetTokenInfo().IsValid {
		t.Error("Token should be invalid when within refresh buffer")
	}

	// Next call should refresh
	_, err = tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("Second GetAccessToken(): %v", err)
	}

	if callCount != 2 {
		t.Errorf("API should be called twice due to expiry, was called %d times", callCount)
	}
}

func TestTokenManagerDefaultExpiresIn(t *testing.T) {
	setup()
	defer teardown()

	// Response without expires_in field
	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		httpmock.NewStringResponder(200, `{
			"access_token": "no_expiry_token",
			"token_type": "Bearer",
			"scope": "read_products"
		}`))

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	_, err := tm.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken(): %v", err)
	}

	info := tm.GetTokenInfo()
	// Should default to ~24 hours
	if info.ExpiresInSeconds < 23*60*60 || info.ExpiresInSeconds > 24*60*60+10 {
		t.Errorf("ExpiresInSeconds should be ~24 hours, got %d seconds", info.ExpiresInSeconds)
	}
}

func TestTokenManagerContextCancellation(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/oauth/access_token",
		func(req *http.Request) (*http.Response, error) {
			// Simulate slow response
			time.Sleep(100 * time.Millisecond)
			return httpmock.NewStringResponse(200, `{
				"access_token": "cancelled_token",
				"token_type": "Bearer",
				"expires_in": 86400,
				"scope": "read_products"
			}`), nil
		})

	app.Client = client
	tm := NewTokenManager(app, "fooshop")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := tm.GetAccessToken(ctx)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
}
