package goshopify

import (
	"time"
)

// GraphQLService is an interface to interact with the graphql endpoint
// of the Shopify API
// See https://shopify.dev/docs/admin-api/graphql/reference
type GraphQLService interface {
	Query(string, interface{}, interface{}) error
}

// GraphQLServiceOp handles communication with the graphql endpoint of
// the Shopify API.
type GraphQLServiceOp struct {
	client *Client
}

type graphQLResponse struct {
	Data       interface{}       `json:"data"`
	Errors     []graphQLError    `json:"errors"`
	Extensions *graphQLExtension `json:"extensions"`
}

type graphQLExtension struct {
	Cost GraphQLCost `json:"cost"`
}

// GraphQLCost represents the cost of the graphql query
type GraphQLCost struct {
	RequestedQueryCost int                   `json:"requestedQueryCost"`
	ActualQueryCost    *int                  `json:"actualQueryCost"`
	ThrottleStatus     GraphQLThrottleStatus `json:"throttleStatus"`
}

// GraphQLThrottleStatus represents the status of the shop's rate limit points
type GraphQLThrottleStatus struct {
	MaximumAvailable   float64 `json:"maximumAvailable"`
	CurrentlyAvailable float64 `json:"currentlyAvailable"`
	RestoreRate        float64 `json:"restoreRate"`
}

type graphQLError struct {
	Message    string                 `json:"message"`
	Extensions *graphQLErrorExtension `json:"extensions"`
	Locations  []graphQLErrorLocation `json:"locations"`
}

type graphQLErrorExtension struct {
	Code          string
	Documentation string
}

const (
	graphQLErrorCodeThrottled = "THROTTLED"
)

type graphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Query creates a graphql query against the Shopify API
// the "data" portion of the response is unmarshalled into resp
func (s *GraphQLServiceOp) Query(q string, vars, resp interface{}) error {
	data := struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables"`
	}{
		Query:     q,
		Variables: vars,
	}

	attempts := 0

	for {
		gr := graphQLResponse{
			Data: resp,
		}

		err := s.client.Post("graphql.json", data, &gr)
		// internal attempts count towards outer total
		attempts += s.client.attempts
		s.client.attempts = attempts

		var ra float64

		if gr.Extensions != nil {
			ra = gr.Extensions.Cost.RetryAfterSeconds()
			s.client.RateLimits.GraphQLCost = &gr.Extensions.Cost
			s.client.RateLimits.RetryAfterSeconds = ra
		}

		if len(gr.Errors) > 0 {
			re := ResponseError{Status: 200}
			var doRetry bool

			for _, err := range gr.Errors {
				if err.Extensions != nil && err.Extensions.Code == graphQLErrorCodeThrottled {
					if attempts >= s.client.retries {
						return RateLimitError{
							RetryAfter: int(ra),
							ResponseError: ResponseError{
								Status:  200,
								Message: err.Message,
							},
						}
					}

					// only need to retry graphql throttled retries
					doRetry = true
				}

				re.Errors = append(re.Errors, err.Message)
			}

			if doRetry {
				wait := time.Duration(ra) * time.Second
				s.client.log.Debugf("rate limited waiting %s", wait.String())
				time.Sleep(wait)
				continue
			}

			err = re
		}

		return err
	}
}

// RetryAfterSeconds returns the estimated retry after seconds based on
// the requested query cost and throttle status
func (c GraphQLCost) RetryAfterSeconds() float64 {
	var diff float64

	if c.ActualQueryCost != nil {
		diff = c.ThrottleStatus.CurrentlyAvailable - float64(*c.ActualQueryCost)
	} else {
		diff = c.ThrottleStatus.CurrentlyAvailable - float64(c.RequestedQueryCost)
	}

	if diff < 0 {
		return -diff / c.ThrottleStatus.RestoreRate
	}

	return 0
}
