package goshopify

import (
	"fmt"
	"net/http"
)

// Option is used to configure client with options
type Option func(c *Client)

// WithVersion optionally sets the api-version if the passed string is valid
func WithVersion(apiVersion string) Option {
	return func(c *Client) {
		pathPrefix := defaultApiPathPrefix
		if len(apiVersion) > 0 && (apiVersionRegex.MatchString(apiVersion) || apiVersion == UnstableApiVersion) {
			pathPrefix = fmt.Sprintf("admin/api/%s", apiVersion)
		}
		c.apiVersion = apiVersion
		c.pathPrefix = pathPrefix
	}
}

// WithRetry sets the number of times a request will be retried if a rate limit or service unavailable error is returned.
// Rate limiting can be either REST API limits or GraphQL Cost limits
func WithRetry(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

func WithLogger(logger LeveledLoggerInterface) Option {
	return func(c *Client) {
		c.log = logger
	}
}

// WithHTTPClient is used to set a custom http client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.Client = client
	}
}
