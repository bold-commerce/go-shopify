package goshopify

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

func TestRefundAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "orders/11/refunds.json",
			response: `{"refunds":[{"id":7}]}`,
			want:     []Refund{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Refund.List(ctx, 11, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "orders/11/refunds/7.json",
			response: `{"refund":{"id":7}}`,
			want:     &Refund{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Refund.Get(ctx, 11, 7, nil)
			},
		},
		{
			name:     "Create",
			method:   "POST",
			path:     "orders/11/refunds.json",
			request:  `{"refund":{"currency":"USD","note":"wrong size","refund_line_items":[{"line_item_id":7,"quantity":1,"restock_type":"return","location_id":9}]}}`,
			response: `{"refund":{"id":7}}`,
			want:     &Refund{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Refund.Create(ctx, 11, Refund{Currency: "USD", Note: "wrong size", RefundLineItems: []RefundLineItem{{LineItemId: 7, Quantity: 1, RestockType: "return", LocationId: 9}}})
			},
		},
		{
			name:     "Calculate",
			method:   "POST",
			path:     "orders/11/refunds/calculate.json",
			request:  `{"refund":{"currency":"USD","refund_line_items":[{"line_item_id":7,"quantity":1,"restock_type":"no_restock"}]}}`,
			response: `{"refund":{"currency":"USD"}}`,
			want:     &Refund{Currency: "USD"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Refund.Calculate(ctx, 11, Refund{Currency: "USD", RefundLineItems: []RefundLineItem{{LineItemId: 7, Quantity: 1, RestockType: "no_restock"}}})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestRefundListWithPagination(t *testing.T) {
	testFeaturePagination(t, "orders/11/refunds.json", `{"refunds":[{"id":7}]}`, []Refund{{Id: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.Refund.ListWithPagination(ctx, 11, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/refund
func TestRefundDocumentedResponse(t *testing.T) {
	var value RefundResource
	testFeatureResponse(t, "fixtures/refund_documented.json", &value, false)
	if !(value.Refund.Shipping.MaximumRefundable.String() == "5" && value.Refund.Transactions[0].MaximumRefundable.String() == "41.94" && value.Refund.Transactions[0].Kind == "suggested_refund") {
		t.Fatalf("incomplete response: %#v", value)
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/refund
func TestRefundCreatedDocumentedResponse(t *testing.T) {
	var value RefundResource
	testFeatureResponse(t, "fixtures/refund_created_documented.json", &value, false)
	if !(value.Refund.Id == 1048576607 && value.Refund.RefundLineItems[0].RestockType == "no_restock" && value.Refund.RefundLineItems[0].Subtotal.IsZero()) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestRefundRequestFields(t *testing.T) {
	amount := decimal.RequireFromString("41.94")
	shipping := decimal.Zero
	notify := false
	parent := int64(801038806)
	refund := Refund{Currency: "USD", Notify: &notify, Shipping: &RefundShipping{Amount: &shipping}, RefundDuties: []RefundDuty{{DutyId: 9, RefundType: "FULL"}}, Transactions: []Transaction{{ParentId: &parent, Amount: &amount, Kind: "refund", Gateway: "bogus"}}}
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Create",
		method:   "POST",
		path:     "orders/11/refunds.json",
		request:  `{"refund":{"currency":"USD","notify":false,"shipping":{"amount":"0"},"refund_duties":[{"duty_id":9,"refund_type":"FULL"}],"transactions":[{"parent_id":801038806,"amount":"41.94","kind":"refund","gateway":"bogus"}]}}`,
		response: `{"refund":{"id":7}}`,
		want:     &Refund{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Refund.Create(ctx, 11, refund)
		},
	})
}
