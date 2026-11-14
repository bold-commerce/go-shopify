package goshopify

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"context"
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestDiscountCodeAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Count",
			method:   "GET",
			path:     "discount_codes/count.json",
			response: `{"count":3}`,
			want:     3,
			call: func(ctx context.Context) (interface{}, error) {
				return client.DiscountCode.Count(ctx, nil)
			},
		},
		{
			name:     "Lookup",
			method:   "GET",
			path:     "discount_codes/lookup.json",
			query:    "code=A%26B+%2B",
			response: `{"discount_code":{"id":7,"code":"A&B +"}}`,
			want:     &PriceRuleDiscountCode{Id: 7, Code: "A&B +"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DiscountCode.Lookup(ctx, "A&B +")
			},
		},
		{
			name:     "CreateBatch",
			method:   "POST",
			path:     "price_rules/11/batch.json",
			request:  `{"discount_codes":[{"code":"SUMMER"}]}`,
			response: `{"discount_code_creation":{"id":7,"status":"queued"}}`,
			want:     &DiscountCodeCreation{Id: 7, Status: "queued"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DiscountCode.CreateBatch(ctx, 11, []PriceRuleDiscountCode{{Code: "SUMMER"}})
			},
		},
		{
			name:     "GetBatch",
			method:   "GET",
			path:     "price_rules/11/batch/7.json",
			response: `{"discount_code_creation":{"id":7,"status":"completed"}}`,
			want:     &DiscountCodeCreation{Id: 7, Status: "completed"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DiscountCode.GetBatch(ctx, 11, 7)
			},
		},
		{
			name:     "ListBatchCodes",
			method:   "GET",
			path:     "price_rules/11/batch/7/discount_codes.json",
			response: `{"discount_codes":[{"id":null,"code":"SUMMER","errors":{"code":["already exists"]}}]}`,
			want:     []BatchDiscountCode{{Code: "SUMMER", Errors: map[string][]string{"code": {"already exists"}}}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DiscountCode.ListBatchCodes(ctx, 11, 7, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestDiscountCodeLookupRedirect(t *testing.T) {
	setup()
	defer teardown()
	httpmock.DeactivateAndReset()
	prefix := "/" + client.pathPrefix
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Header.Get("X-Shopify-Access-Token") != client.token {
			t.Error("missing token")
		}
		switch r.URL.Path {
		case prefix + "/discount_codes/lookup.json":
			if r.URL.Query().Get("code") != "A&B +" {
				t.Errorf("incorrect code: %v", r.URL.Query())
			}
			http.Redirect(w, r, prefix+"/price_rules/11/discount_codes/7.json", http.StatusFound)
		case prefix + "/price_rules/11/discount_codes/7.json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"discount_code":{"id":7,"price_rule_id":11,"code":"A&B +"}}`))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client.Client = server.Client()
	client.baseURL, _ = url.Parse(server.URL)
	code, err := client.DiscountCode.Lookup(context.Background(), "A&B +")
	if err != nil {
		t.Fatal(err)
	}
	if code == nil || code.Id != 7 || code.PriceRuleId != 11 || code.Code != "A&B +" || requests != 2 {
		t.Fatalf("code=%#v, requests=%d", code, requests)
	}
}
