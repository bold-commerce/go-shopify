package goshopify

import "fmt"

// Option is used to configure client with options
type Option func(c *Client)

// WithVersion optionally sets the api-version if the passed string is valid
func WithVersion(apiVersion string) Option {
	return func(c *Client) {
		pathPrefix := defaultApiPathPrefix
		if len(apiVersion) > 0 && apiVersionRegex.MatchString(apiVersion) {
			pathPrefix = fmt.Sprintf("admin/api/%s", apiVersion)
		}
		c.apiVersion = apiVersion
		c.pathPrefix = pathPrefix
	}
}

func WithRetry(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}
