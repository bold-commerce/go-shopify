package goshopify

import (
	"context"
	"sync"
	"time"
)

const (
	// defaultRefreshBuffer is the time before token expiry when we should refresh.
	// Shopify client credentials tokens expire after 24 hours.
	defaultRefreshBuffer = 1 * time.Hour
)

// ClientCredentialsResponse represents the response from Shopify's token endpoint
// when using the client credentials grant type.
type ClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// GetClientCredentialsToken exchanges client credentials for an access token.
// This implements the OAuth 2.0 client credentials grant flow for Shopify.
//
// The client credentials grant is used for server-to-server authentication
// where the app needs to access shop data without user interaction.
// The app must be installed in the shop to use this flow.
//
// Tokens obtained via this method expire after 24 hours and should be refreshed.
// For automatic token management, consider using TokenManager instead.
//
// See: https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/client-credentials-grant
func (app App) GetClientCredentialsToken(ctx context.Context, shopName string) (*ClientCredentialsResponse, error) {
	data := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}{
		ClientId:     app.ApiKey,
		ClientSecret: app.ApiSecret,
		GrantType:    "client_credentials",
	}

	client := app.Client
	if client == nil {
		client = MustNewClient(app, shopName, "")
	}

	req, err := client.NewRequest(ctx, "POST", accessTokenRelPath, data, nil)
	if err != nil {
		return nil, err
	}

	response := new(ClientCredentialsResponse)
	err = client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// TokenManager handles automatic token refresh for Shopify client credentials.
// It provides thread-safe access to tokens and automatically refreshes them
// before they expire.
//
// TokenManager is safe for concurrent use by multiple goroutines.
type TokenManager struct {
	app      App
	shopName string

	mutex       sync.RWMutex
	accessToken string
	expiresAt   time.Time

	refreshBuffer time.Duration
	log           LeveledLoggerInterface
}

// TokenManagerOption is used to configure TokenManager with options.
type TokenManagerOption func(*TokenManager)

// TokenInfo contains information about the current token state.
type TokenInfo struct {
	HasToken         bool          `json:"has_token"`
	ExpiresAt        time.Time     `json:"expires_at"`
	ExpiresInSeconds int           `json:"expires_in_seconds"`
	IsValid          bool          `json:"is_valid"`
	RefreshBuffer    time.Duration `json:"refresh_buffer"`
}

// NewTokenManager creates a new TokenManager for managing client credentials tokens.
//
// The TokenManager automatically handles token caching and refresh. By default,
// tokens are refreshed 1 hour before expiry (Shopify tokens last 24 hours).
//
// Example:
//
//	app := goshopify.App{
//	    ApiKey:    "your-client-id",
//	    ApiSecret: "your-client-secret",
//	}
//	tm := goshopify.NewTokenManager(app, "shop-name")
//	token, err := tm.GetAccessToken(ctx)
func NewTokenManager(app App, shopName string, opts ...TokenManagerOption) *TokenManager {
	tm := &TokenManager{
		app:           app,
		shopName:      shopName,
		refreshBuffer: defaultRefreshBuffer,
	}

	for _, opt := range opts {
		opt(tm)
	}

	return tm
}

// WithRefreshBuffer sets the duration before token expiry when a refresh should occur.
// Default is 1 hour. For example, WithRefreshBuffer(30*time.Minute) will refresh
// the token 30 minutes before it expires.
func WithRefreshBuffer(duration time.Duration) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.refreshBuffer = duration
	}
}

// WithTokenManagerLogger sets a custom logger for the TokenManager.
func WithTokenManagerLogger(logger LeveledLoggerInterface) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.log = logger
	}
}

// GetAccessToken returns a valid access token, refreshing if necessary.
// This method is safe for concurrent use.
//
// If the cached token is still valid (not expired and not within the refresh buffer),
// it is returned immediately. Otherwise, a new token is fetched from Shopify.
func (tm *TokenManager) GetAccessToken(ctx context.Context) (string, error) {
	tm.mutex.RLock()
	if tm.isTokenValid() {
		token := tm.accessToken
		tm.mutex.RUnlock()
		return token, nil
	}
	tm.mutex.RUnlock()

	// Need to refresh - acquire write lock
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Double-check after acquiring write lock (another goroutine might have refreshed)
	if tm.isTokenValid() {
		return tm.accessToken, nil
	}

	return tm.refreshToken(ctx)
}

// ForceRefresh forces a token refresh regardless of current token validity.
// This is useful when you receive an authentication error and need to get a fresh token.
func (tm *TokenManager) ForceRefresh(ctx context.Context) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	return tm.refreshToken(ctx)
}

// GetTokenInfo returns information about the current token state.
// This is useful for monitoring and debugging.
func (tm *TokenManager) GetTokenInfo() TokenInfo {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	hasToken := tm.accessToken != ""
	var expiresInSeconds int
	if hasToken && !tm.expiresAt.IsZero() {
		expiresInSeconds = int(time.Until(tm.expiresAt).Seconds())
		if expiresInSeconds < 0 {
			expiresInSeconds = 0
		}
	}

	return TokenInfo{
		HasToken:         hasToken,
		ExpiresAt:        tm.expiresAt,
		ExpiresInSeconds: expiresInSeconds,
		IsValid:          tm.isTokenValid(),
		RefreshBuffer:    tm.refreshBuffer,
	}
}

// isTokenValid checks if the current token is valid.
// Must be called with at least a read lock held.
func (tm *TokenManager) isTokenValid() bool {
	if tm.accessToken == "" {
		return false
	}
	// Check if token will expire within the buffer period
	return time.Now().Add(tm.refreshBuffer).Before(tm.expiresAt)
}

// refreshToken fetches a new access token using client credentials grant.
// Must be called with write lock held.
func (tm *TokenManager) refreshToken(ctx context.Context) (string, error) {
	response, err := tm.app.GetClientCredentialsToken(ctx, tm.shopName)
	if err != nil {
		return "", err
	}

	tm.accessToken = response.AccessToken

	// Calculate expiry time
	if response.ExpiresIn > 0 {
		tm.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	} else {
		// Default to 24 hours if not specified (Shopify's default)
		tm.expiresAt = time.Now().Add(24 * time.Hour)
	}

	if tm.log != nil {
		tm.log.Infof("Shopify access token refreshed, expires at: %s", tm.expiresAt.Format(time.RFC3339))
	}

	return tm.accessToken, nil
}

// ClearToken clears the cached token, forcing a refresh on the next GetAccessToken call.
func (tm *TokenManager) ClearToken() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.accessToken = ""
	tm.expiresAt = time.Time{}
}

