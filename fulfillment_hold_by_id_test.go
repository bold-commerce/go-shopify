package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strings"
	"testing"
)

func TestReleaseHoldByID(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()
	if client.FulfillmentOrder.ReleaseHoldByID(ctx, 42, "") == nil {
		t.Fatal("accepted empty hold ID")
	}
	const hold = "gid://shopify/FulfillmentHold/3"
	const order = "gid://shopify/FulfillmentOrder/42"
	for _, tc := range []struct {
		body  string
		fails bool
	}{
		{`{"data":{"fulfillmentOrderReleaseHold":{"fulfillmentOrder":{"id":"gid://shopify/FulfillmentOrder/42"},"userErrors":[]}}}`, false},
		{`{"data":{"fulfillmentOrderReleaseHold":{"fulfillmentOrder":null,"userErrors":[{"message":"unknown hold"}]}}}`, true},
	} {
		httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/graphql.json", client.pathPrefix), func(req *http.Request) (*http.Response, error) {
			var input struct {
				Query     string
				Variables map[string]string
			}
			if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				return nil, err
			}
			if input.Variables["id"] != order || input.Variables["hold"] != hold || !strings.Contains(input.Query, "holdIds: [$hold]") {
				t.Errorf("unexpected mutation %v", input)
			}
			return httpmock.NewStringResponse(200, tc.body), nil
		})
		if err := client.FulfillmentOrder.ReleaseHoldByID(ctx, 42, hold); (err != nil) != tc.fails {
			t.Fatal(err)
		}
	}
}
