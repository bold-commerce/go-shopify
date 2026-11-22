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
	testFeaturePagination(t, "orders/11/refunds.json", `{"refunds":[{"id":7}]}`, []Refund{{Id: 7}}, `{"refunds":[{"id":8}]}`, []Refund{{Id: 8}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
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

func TestRefundFullShippingRequest(t *testing.T) {
	enabled := true
	testFeatureEndpoint(t, featureEndpoint{
		method:   "POST",
		path:     "orders/11/refunds.json",
		status:   201,
		request:  `{"refund":{"notify":true,"discrepancy_reason":"customer","shipping":{"full_refund":true},"refund_duties":[{"duty_id":9007199254740993,"refund_type":"FULL"}],"refund_line_items":[{"line_item_id":9007199254740995,"quantity":2,"restock_type":"return","location_id":9}]}}`,
		response: `{"refund":{"id":7}}`,
		want:     &Refund{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Refund.Create(ctx, 11, Refund{
				Notify: &enabled, DiscrepancyReason: "customer",
				Shipping:        &RefundShipping{FullRefund: &enabled},
				RefundDuties:    []RefundDuty{{DutyId: 9007199254740993, RefundType: "FULL"}},
				RefundLineItems: []RefundLineItem{{LineItemId: 9007199254740995, Quantity: 2, RestockType: "return", LocationId: 9}},
			})
		},
	})
}

func TestRefundCalculateLineItems(t *testing.T) {
	amount := decimal.RequireFromString("10.50")
	tax := decimal.RequireFromString("0.75")
	price := decimal.RequireFromString("5.00")
	discountedPrice := decimal.RequireFromString("4.75")
	discountedTotalPrice := decimal.RequireFromString("9.50")
	totalCartDiscountAmount := decimal.RequireFromString("0.50")
	parent := int64(9007199254740993)
	items := []RefundLineItem{{LineItemId: 7, Quantity: 2, RestockType: "return", LocationId: 9}}
	testFeatureEndpoint(t, featureEndpoint{
		method:   "POST",
		path:     "orders/11/refunds/calculate.json",
		request:  `{"refund":{"currency":"USD","refund_line_items":[{"line_item_id":7,"quantity":2,"restock_type":"return","location_id":9}]}}`,
		response: `{"refund":{"currency":"USD","refund_line_items":[{"line_item_id":7,"quantity":2,"restock_type":"return","location_id":9,"price":"5.00","discounted_price":"4.75","discounted_total_price":"9.50","total_cart_discount_amount":"0.50","subtotal":"10.50","total_tax":"0.75"}],"shipping":{"amount":"10.50","tax":"0.75","maximum_refundable":"10.50"},"transactions":[{"parent_id":9007199254740993,"amount":"10.50","maximum_refundable":"10.50","kind":"suggested_refund","currency":"USD"}]}}`,
		want: &Refund{
			Currency:        "USD",
			RefundLineItems: []RefundLineItem{{LineItemId: 7, Quantity: 2, RestockType: "return", LocationId: 9, Price: &price, DiscountedPrice: &discountedPrice, DiscountedTotalPrice: &discountedTotalPrice, TotalCartDiscountAmount: &totalCartDiscountAmount, Subtotal: &amount, TotalTax: &tax}},
			Shipping:        &RefundShipping{Amount: &amount, Tax: &tax, MaximumRefundable: &amount},
			Transactions:    []Transaction{{ParentId: &parent, Amount: &amount, MaximumRefundable: &amount, Kind: "suggested_refund", Currency: "USD"}},
		},
		call: func(ctx context.Context) (interface{}, error) {
			return client.Refund.Calculate(ctx, 11, Refund{Currency: "USD", RefundLineItems: items})
		},
	})
}

func TestRefundCalculateAlreadyStocked(t *testing.T) {
	stocked, unstocked := true, false
	for _, tc := range []struct {
		name    string
		stocked *bool
		field   string
	}{
		{name: "Omitted"},
		{name: "False", stocked: &unstocked, field: `,"already_stocked":false`},
		{name: "True", stocked: &stocked, field: `,"already_stocked":true`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calculated *Refund
			request := `{"refund":{"refund_line_items":[{"line_item_id":7,"quantity":1,"restock_type":"return","location_id":9}]}}`
			response := `{"refund":{"refund_line_items":[{"line_item_id":7,"quantity":1,"restock_type":"return","location_id":9` + tc.field + `}]}}`
			want := &Refund{RefundLineItems: []RefundLineItem{{LineItemId: 7, Quantity: 1, RestockType: "return", LocationId: 9, AlreadyStocked: tc.stocked}}}
			testFeatureEndpoint(t, featureEndpoint{
				method:   "POST",
				path:     "orders/11/refunds/calculate.json",
				request:  request,
				response: response,
				want:     want,
				call: func(ctx context.Context) (interface{}, error) {
					refund, err := client.Refund.Calculate(ctx, 11, Refund{RefundLineItems: []RefundLineItem{{LineItemId: 7, Quantity: 1, RestockType: "return", LocationId: 9}}})
					if err == nil {
						calculated = refund
					}
					return refund, err
				},
			})
			testFeatureEndpoint(t, featureEndpoint{
				method:   "POST",
				path:     "orders/11/refunds.json",
				status:   201,
				request:  response,
				response: response,
				want:     want,
				call: func(ctx context.Context) (interface{}, error) {
					return client.Refund.Create(ctx, 11, *calculated)
				},
			})
		})
	}
}
