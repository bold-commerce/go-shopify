package goshopify

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RateLimitInfo struct {
	RequestCount      int
	BucketSize        int
	GraphQLCost       *GraphQLCost
	RetryAfterSeconds float64
}

func (c *Client) recordRateLimitHeaders(header http.Header) {
	c.rateLimitsMu.Lock()
	defer c.rateLimitsMu.Unlock()
	if s := strings.Split(header.Get("X-Shopify-Shop-Api-Call-Limit"), "/"); len(s) == 2 {
		c.rateLimits.RequestCount, _ = strconv.Atoi(s[0])
		c.rateLimits.BucketSize, _ = strconv.Atoi(s[1])
	}

	c.rateLimits.RetryAfterSeconds, _ = strconv.ParseFloat(header.Get("Retry-After"), 64)
}

// RateLimits returns an independent snapshot of the latest recorded rate-limit
// information. It is safe to call while requests are in progress.
func (c *Client) RateLimits() RateLimitInfo {
	c.rateLimitsMu.RLock()
	defer c.rateLimitsMu.RUnlock()

	limits := c.rateLimits
	if limits.GraphQLCost != nil {
		cost := *limits.GraphQLCost
		if cost.ActualQueryCost != nil {
			actual := *cost.ActualQueryCost
			cost.ActualQueryCost = &actual
		}
		limits.GraphQLCost = &cost
	}
	return limits
}

// waitForRetry stops the backoff when the caller cancels the request.
func waitForRetry(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
