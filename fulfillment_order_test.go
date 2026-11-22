package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func FulfillmentOrderTests(t *testing.T, fulfillmentOrder FulfillmentOrder) {
	// Check that Id is assigned to the returned fulfillment
	expectedInt := uint64(255858046) // in fulfillment_order.json fixture
	if fulfillmentOrder.Id != expectedInt {
		t.Errorf("FulfillmentOrder.Id returned %+v, expected %+v", fulfillmentOrder.Id, expectedInt)
	}
}

func TestFulfillmentOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillment_orders.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment_orders": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentOrderServiceOp{client: client}

	fulfillmentOrders, err := fulfillmentService.List(context.Background(), 123, nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.List returned error: %v", err)
	}

	expected := []FulfillmentOrder{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(fulfillmentOrders, expected) {
		t.Errorf("FulfillmentOrder.List returned %+v, expected %+v", fulfillmentOrders, expected)
	}
}

func TestFulfillmentOrderGet(t *testing.T) {
	setup()
	defer teardown()

	fixture := loadFixture("fulfillment_order.json")
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, fixture))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	fulfillment, err := fulfillmentOrderService.Get(context.Background(), 255858046, nil)
	if err != nil {
		t.Errorf("FulfillmentOrder.Get returned error: %v", err)
	}

	expected := FulfillmentOrderResource{}
	err = json.Unmarshal(fixture, &expected)
	if err != nil {
		t.Errorf("json.Unmarshall returned error : %v", err)
	}

	if !reflect.DeepEqual(fulfillment, expected.FulfillmentOrder) {
		t.Errorf("FulfillmentOrder.Get returned %+v, expected %+v", fulfillment, expected)
	}
}

func TestFulfillmentOrderFulfillments(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/fulfillments.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillments": [{"id":1},{"id":2}]}`))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	fulfillments, err := fulfillmentOrderService.Fulfillments(context.Background(), 255858046)
	if err != nil {
		t.Errorf("FulfillmentOrder.Fulfillments returned error: %v", err)
	}

	expected := []Fulfillment{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(fulfillments, expected) {
		t.Errorf("FulfillmentOrder.Fulfillments returned %+v, expected %+v", fulfillments, expected)
	}
}

func TestFulfillmentOrderCancel(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/cancel.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Cancel(context.Background(), 1)
	if err != nil {
		t.Errorf("FulfillmentOrder.Cancel returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderClose(t *testing.T) {
	setup()
	defer teardown()
	for _, message := range []string{"Not enough inventory to complete this work.", ""} {
		httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/close.json", client.pathPrefix), func(req *http.Request) (*http.Response, error) {
			var got map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			fields := map[string]interface{}{}
			if message != "" {
				fields["message"] = message
			}
			want := map[string]interface{}{"fulfillment_order": fields}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("request = %#v, want %#v", got, want)
			}
			return httpmock.NewBytesResponse(200, loadFixture("fulfillment_order.json")), nil
		})
		got, err := client.FulfillmentOrder.Close(context.Background(), 1, message)
		if err != nil {
			t.Fatal(err)
		}
		FulfillmentOrderTests(t, *got)
	}
}

func TestFulfillmentOrderHold(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/hold.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	returnedFulfillment, err := fulfillmentOrderService.Hold(context.Background(), 1, false, HoldReasonOutOfStock, "test")
	if err != nil {
		t.Errorf("FulfillmentOrder.Hold returned error: %v", err)
	}

	FulfillmentOrderTests(t, *returnedFulfillment)
}

func TestFulfillmentOrderMove(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1046000818/move.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order_move.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}

	req := FulfillmentOrderMoveRequest{
		NewLocationId: 655441491,
		LineItems: []FulfillmentOrderLineItemQuantity{
			{Id: 1058737594, Quantity: 1},
		},
	}

	result, err := fulfillmentOrderService.Move(context.Background(), 1046000818, req)
	if err != nil {
		t.Errorf("FulfillmentOrder.Move returned error: %v", err)
	}
	if result.MovedFulfillmentOrder.AssignedLocationId != 655441491 {
		t.Errorf("FulfillmentOrder.Move result AssignedLocation is is %d, expected %d",
			result.MovedFulfillmentOrder.AssignedLocationId, 655441491)
	}
}

func TestFulfillmentOrderOpen(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/open.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.Open(context.Background(), fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentOrder.Open returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderReleaseHold(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/release_hold.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_order.json")))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.ReleaseHold(context.Background(), fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentOrder.ReleaseHold returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderReschedule(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/255858046/reschedule.json", client.pathPrefix),
		func(req *http.Request) (*http.Response, error) {
			var got map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			want := map[string]interface{}{"fulfillment_order": map[string]interface{}{"new_fulfill_at": "2027-11-01T14:20:00-04:00"}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("request = %#v, want %#v", got, want)
			}
			return httpmock.NewBytesResponse(200, loadFixture("fulfillment_order.json")), nil
		})

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)

	result, err := fulfillmentOrderService.Reschedule(context.Background(), fulfillmentId, time.Date(2027, 11, 1, 14, 20, 0, 0, time.FixedZone("", -4*60*60)))
	if err != nil {
		t.Errorf("FulfillmentOrder.Reschedule returned error: %v", err)
	}

	if result == nil || result.Id != fulfillmentId {
		t.Errorf("Expected Id: %d   got: %d", result.Id, fulfillmentId)
	}
}

func TestFulfillmentOrderSetDeadline(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/set_fulfillment_orders_deadline.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	fulfillmentOrderService := &FulfillmentOrderServiceOp{client: client}
	fulfillmentId := uint64(255858046)
	newDeadline := time.Now().Add(time.Hour * 24 * 7)
	err := fulfillmentOrderService.SetDeadline(context.Background(), []uint64{fulfillmentId}, newDeadline)
	if err != nil {
		t.Errorf("FulfillmentOrder.SetDeadline returned error: %v", err)
	}
}

func TestFulfillmentOrderCancelWithResponse(t *testing.T) {
	setup()
	defer teardown()
	for _, replacement := range []string{`{"id":2}`, `null`} {
		httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/cancel.json", client.pathPrefix), httpmock.NewStringResponder(200, `{"fulfillment_order":{"id":1},"replacement_fulfillment_order":`+replacement+`}`))
		got, err := client.FulfillmentOrder.CancelWithResponse(context.Background(), 1)
		if err != nil {
			t.Fatal(err)
		}
		if got.FulfillmentOrder == nil || got.FulfillmentOrder.Id != 1 {
			t.Fatalf("cancelled = %+v", got.FulfillmentOrder)
		}
		if replacement == "null" {
			if got.ReplacementFulfillmentOrder != nil {
				t.Fatal("expected no replacement")
			}
		} else if got.ReplacementFulfillmentOrder == nil || got.ReplacementFulfillmentOrder.Id != 2 {
			t.Fatalf("replacement = %+v", got.ReplacementFulfillmentOrder)
		}
	}
}

func TestFulfillmentOrderCancelErrors(t *testing.T) {
	setup()
	defer teardown()
	for _, body := range []string{`{"errors":"cannot cancel"}`, `invalid json`} {
		status := 422
		if body == "invalid json" {
			status = 200
		}
		httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/cancel.json", client.pathPrefix), httpmock.NewStringResponder(status, body))
		if _, err := client.FulfillmentOrder.CancelWithResponse(context.Background(), 1); err == nil {
			t.Fatal("expected full response error")
		}
		if got, err := client.FulfillmentOrder.Cancel(context.Background(), 1); err == nil || got != nil {
			t.Fatalf("Cancel = %+v, %v; want nil and error", got, err)
		}
	}
}

func TestFulfillmentOrderOptionalReadFields(t *testing.T) {
	options := FulfillmentOrderOptions{IncludeFinancialSummaries: true, IncludeOrderReferenceFields: true}
	const query = "include_financial_summaries=true&include_order_reference_fields=true"
	const order = `{"order_name":"#1001","order_processed_at":"2026-01-10T11:00:00Z","channel_id":9007199254740993,"line_items":[{"id":11,"financial_summaries":[{"quantity":1,"original_unit_price_set":"199.00","approximate_discounted_unit_price_set":"195.67","discount_allocations":[{"amount":"3.33","discount_application":{"allocation_method":"across","target_selection":"all","target_type":"line_item"}}]}]}]}`
	processed := time.Date(2026, 1, 10, 11, 0, 0, 0, time.UTC)
	channel := uint64(9007199254740993)
	price, discounted, amount := decimal.RequireFromString("199.00"), decimal.RequireFromString("195.67"), decimal.RequireFromString("3.33")
	want := FulfillmentOrder{OrderName: "#1001", OrderProcessedAt: &processed, ChannelId: &channel, LineItems: []FulfillmentOrderLineItem{{Id: 11}},
		FinancialSummaries: map[uint64][]FulfillmentOrderFinancialSummary{11: {{
			Quantity: 1, OriginalUnitPriceSet: &price, ApproximateDiscountedUnitPriceSet: &discounted,
			DiscountAllocations: []DiscountAllocations{{Amount: &amount, DiscountApplication: &DiscountApplication{AllocationMethod: "across", TargetSelection: "all", TargetType: "line_item"}}},
		}}},
	}
	for _, tc := range []featureEndpoint{
		{name: "Get", method: "GET", path: "fulfillment_orders/7.json", query: query, response: `{"fulfillment_order":` + order + "}", want: &want,
			call: func(ctx context.Context) (interface{}, error) { return client.FulfillmentOrder.Get(ctx, 7, options) }},
		{name: "List", method: "GET", path: "orders/7/fulfillment_orders.json", query: query, response: `{"fulfillment_orders":[` + order + "]}", want: []FulfillmentOrder{want},
			call: func(ctx context.Context) (interface{}, error) { return client.FulfillmentOrder.List(ctx, 7, options) }},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestFulfillmentOrderLineItemCompatibility(t *testing.T) {
	const input = `{"line_items":[{"id":7,"financial_summaries":[{"quantity":2}]},{"id":8,"financial_summaries":[]}]}`
	var order FulfillmentOrder
	if err := json.Unmarshal([]byte(input), &order); err != nil {
		t.Fatal(err)
	}
	keys := map[FulfillmentOrderLineItem]int{{Id: 7}: 1, {Id: 8}: 2}
	if keys[order.LineItems[0]] != 1 || keys[order.LineItems[1]] != 2 {
		t.Fatal("financial summaries changed line item equality")
	}
	encoded, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}
	var got, want struct {
		LineItems []map[string]interface{} `json:"line_items"`
	}
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(input), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lost nested summaries: %s", encoded)
	}
	if err := json.Unmarshal([]byte(`{"line_items":[{"id":7}]}`), &order); err != nil || order.FinancialSummaries != nil {
		t.Fatal("summaries from the previous response survived", err)
	}
}
