package goshopify

import (
	"context"
	"testing"
)

// Request shapes follow the REST fulfillment and fulfillmentorder documentation:
// https://shopify.dev/docs/api/admin-rest/2026-10/resources/fulfillment
// https://shopify.dev/docs/api/admin-rest/2026-10/resources/fulfillmentorder
func TestFulfillmentRequestOptions(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{
			name: "create with message and origin", method: "POST", path: "fulfillments.json",
			request:  `{"fulfillment":{"message":"The package was shipped this morning.","origin_address":{"address1":"150 Elgin St","address2":"Suite 1","city":"Ottawa","country_code":"CA","province_code":"ON","zip":"K2P 1L4"},"notify_customer":false,"tracking_info":{},"line_items_by_fulfillment_order":[{"fulfillment_order_id":7}]}}`,
			response: `{"fulfillment":{"id":9,"origin_address":{"country_code":"CA"}}}`,
			want:     &Fulfillment{Id: 9, OriginAddress: &Address{CountryCode: "CA"}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Fulfillment.Create(ctx, Fulfillment{
					Message:                     "The package was shipped this morning.",
					OriginAddress:               &Address{Address1: "150 Elgin St", Address2: "Suite 1", City: "Ottawa", CountryCode: "CA", ProvinceCode: "ON", Zip: "K2P 1L4"},
					LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 7}},
				})
			},
		},
		{
			name: "create without optional fields", method: "POST", path: "fulfillments.json",
			request:  `{"fulfillment":{"notify_customer":false,"tracking_info":{},"line_items_by_fulfillment_order":[{"fulfillment_order_id":7}]}}`,
			response: `{"fulfillment":{"id":9}}`, want: &Fulfillment{Id: 9},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Fulfillment.Create(ctx, Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 7}}})
			},
		},
		{
			name: "tracking update excludes create-only fields", method: "POST", path: "fulfillments/9/update_tracking.json",
			request:  `{"fulfillment":{"notify_customer":false,"tracking_info":{"number":"123"}}}`,
			response: `{"fulfillment":{"id":9}}`, want: &Fulfillment{Id: 9},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Fulfillment.Update(ctx, Fulfillment{Id: 9, Message: "create only", OriginAddress: &Address{CountryCode: "CA"}, TrackingInfo: FulfillmentTrackingInfo{Number: "123"}})
			},
		},
		{
			name: "partial hold", method: "POST", path: "fulfillment_orders/7/hold.json",
			request:  `{"fulfillment_hold":{"reason":"inventory_out_of_stock","reason_notes":"Not enough inventory","notify_merchant":true,"fulfillment_order_line_items":[{"id":11,"quantity":1},{"id":12,"quantity":2}]}}`,
			response: `{"fulfillment_order":{"id":7,"status":"on_hold"}}`, want: &FulfillmentOrder{Id: 7, Status: "on_hold"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.FulfillmentOrder.HoldWithOptions(ctx, 7, FulfillmentOrderHoldOptions{Reason: HoldReasonOutOfStock, ReasonNotes: "Not enough inventory", NotifyMerchant: true, LineItems: []FulfillmentOrderLineItemQuantity{{Id: 11, Quantity: 1}, {Id: 12, Quantity: 2}}})
			},
		},
		{
			name: "whole order hold", method: "POST", path: "fulfillment_orders/7/hold.json",
			request:  `{"fulfillment_hold":{"reason":"awaiting_payment","notify_merchant":false}}`,
			response: `{"fulfillment_order":{"id":7,"status":"on_hold"}}`, want: &FulfillmentOrder{Id: 7, Status: "on_hold"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.FulfillmentOrder.HoldWithOptions(ctx, 7, FulfillmentOrderHoldOptions{Reason: HoldReasonAwaitingPayment})
			},
		},
		{
			name: "existing hold signature", method: "POST", path: "fulfillment_orders/7/hold.json",
			request:  `{"fulfillment_hold":{"reason":"other","reason_notes":"Check order","notify_merchant":false}}`,
			response: `{"fulfillment_order":{"id":7,"status":"on_hold"}}`, want: &FulfillmentOrder{Id: 7, Status: "on_hold"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.FulfillmentOrder.Hold(ctx, 7, false, HoldReasonOther, "Check order")
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
