package goshopify

import (
	"context"
	"time"
)

type RateLimitInfo struct {
	RequestCount      int
	BucketSize        int
	GraphQLCost       *GraphQLCost
	RetryAfterSeconds float64
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
