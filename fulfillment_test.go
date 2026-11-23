package goshopify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func FulfillmentTests(t *testing.T, fulfillment Fulfillment) {
	// Check that Id is assigned to the returned fulfillment
	expectedInt := uint64(1022782888)
	if fulfillment.Id != expectedInt {
		t.Errorf("Fulfillment.Id returned %+v, expected %+v", fulfillment.Id, expectedInt)
	}
}

func TestFulfillmentList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillments": [{"id":1},{"id":2}]}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillments, err := fulfillmentService.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Fulfillment.List returned error: %v", err)
	}

	expected := []Fulfillment{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(fulfillments, expected) {
		t.Errorf("Fulfillment.List returned %+v, expected %+v", fulfillments, expected)
	}
}

func TestFulfillmentCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	cnt, err := fulfillmentService.Count(context.Background(), nil)
	if err != nil {
		t.Errorf("Fulfillment.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Fulfillment.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = fulfillmentService.Count(context.Background(), CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Fulfillment.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Fulfillment.Count returned %d, expected %d", cnt, expected)
	}
}

func TestFulfillmentGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123/fulfillments/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment": {"id":1}}`))

	fulfillmentService := &FulfillmentServiceOp{client: client, resource: ordersResourceName, resourceId: 123}

	fulfillment, err := fulfillmentService.Get(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Fulfillment.Get returned error: %v", err)
	}

	expected := &Fulfillment{Id: 1}
	if !reflect.DeepEqual(fulfillment, expected) {
		t.Errorf("Fulfillment.Get returned %+v, expected %+v", fulfillment, expected)
	}
}

// Check the exact wire payload, including false values and omitted legacy fields.
func fulfillmentWriteResponder(t *testing.T, want string) httpmock.Responder {
	t.Helper()
	return func(req *http.Request) (*http.Response, error) {
		var gotBody, wantBody interface{}
		if err := json.NewDecoder(req.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(want), &wantBody); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(gotBody, wantBody) {
			t.Errorf("request body = %#v; want %#v", gotBody, wantBody)
		}
		return httpmock.NewBytesResponse(200, loadFixture("fulfillment.json")), nil
	}
}

func TestFulfillmentCreate(t *testing.T) {
	setup()
	defer teardown()
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillments.json", client.pathPrefix),
		fulfillmentWriteResponder(t, `{"fulfillment":{"line_items_by_fulfillment_order":[{"fulfillment_order_id":10,"fulfillment_order_line_items":[{"id":20,"quantity":2}]}],"tracking_info":{"company":"UPS","number":"123","url":"https://tracking.example/123"},"notify_customer":false}}`))
	got, err := client.Fulfillment.Create(context.Background(), Fulfillment{
		Id: 99, Status: "success", OrderId: 1, LocationId: 2,
		TrackingInfo:   FulfillmentTrackingInfo{Company: "UPS", Number: "123", Url: "https://tracking.example/123"},
		TrackingNumber: "old number",
		LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{
			{FulfillmentOrderId: 10, FulfillmentOrderLineItems: []LineItemByFulfillmentOrderItemQuantity{{Id: 20, Quantity: 2}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	FulfillmentTests(t, *got)
	if httpmock.GetTotalCallCount() != 1 {
		t.Fatal("unexpected additional request")
	}
}

func TestFulfillmentUpdateTracking(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    Fulfillment
		tracking string
	}{
		{"modern", Fulfillment{TrackingInfo: FulfillmentTrackingInfo{Company: "UPS", Number: "new", Url: "https://tracking.example/new"}, TrackingNumber: "old", TrackingUrls: []string{"old1", "old2"}}, `{"company":"UPS","number":"new","url":"https://tracking.example/new"}`},
		{"legacy singular", Fulfillment{TrackingCompany: "UPS", TrackingNumber: "123", TrackingUrl: "https://tracking.example/123"}, `{"company":"UPS","number":"123","url":"https://tracking.example/123"}`},
		{"legacy arrays", Fulfillment{TrackingNumbers: []string{"123"}, TrackingUrls: []string{"https://tracking.example/123"}}, `{"number":"123","url":"https://tracking.example/123"}`},
		{"singular takes precedence", Fulfillment{TrackingNumber: "new", TrackingUrl: "https://tracking.example/new", TrackingNumbers: []string{"old"}, TrackingUrls: []string{"old"}}, `{"number":"new","url":"https://tracking.example/new"}`},
		{"modern is not merged with stale legacy fields", Fulfillment{TrackingInfo: FulfillmentTrackingInfo{Number: "new"}, TrackingUrl: "old"}, `{"number":"new"}`},
		{"empty tracking", Fulfillment{}, `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()
			tc.input.Id = 1022782888
			tc.input.NotifyCustomer = true
			tc.input.Status = "success"
			tc.input.LineItems = []LineItem{{Id: 9}}
			httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillments/1022782888/update_tracking.json", client.pathPrefix),
				fulfillmentWriteResponder(t, `{"fulfillment":{"tracking_info":`+tc.tracking+`,"notify_customer":true}}`))
			got, err := client.Fulfillment.Update(context.Background(), tc.input)
			if err != nil {
				t.Fatal(err)
			}
			FulfillmentTests(t, *got)
			if httpmock.GetTotalCallCount() != 1 {
				t.Fatal("unexpected additional request")
			}
		})
	}
}

func TestFulfillmentCancel(t *testing.T) {
	setup()
	defer teardown()
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillments/1/cancel.json", client.pathPrefix), httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))
	got, err := client.Fulfillment.Cancel(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	FulfillmentTests(t, *got)
	if httpmock.GetTotalCallCount() != 1 {
		t.Fatal("unexpected additional request")
	}
}

func TestFulfillmentRemovedOperations(t *testing.T) {
	setup()
	defer teardown()
	for _, call := range []func() (*Fulfillment, error){
		func() (*Fulfillment, error) { return client.Fulfillment.Complete(context.Background(), 1) },
		func() (*Fulfillment, error) { return client.Fulfillment.Transition(context.Background(), 1) },
		func() (*Fulfillment, error) { return client.Order.CompleteFulfillment(context.Background(), 1, 2) },
		func() (*Fulfillment, error) { return client.Order.TransitionFulfillment(context.Background(), 1, 2) },
	} {
		got, err := call()
		if got != nil || !errors.Is(err, ErrFulfillmentOperationUnsupported) {
			t.Fatalf("got %v, %v", got, err)
		}
	}
	if httpmock.GetTotalCallCount() != 0 {
		t.Fatal("removed operations must not send requests")
	}
}

func TestFulfillmentUnscopedReads(t *testing.T) {
	setup()
	defer teardown()
	got, err := client.Fulfillment.List(context.Background(), nil)
	if got != nil || !errors.Is(err, ErrFulfillmentRequiresOrder) {
		t.Fatalf("got %v, %v", got, err)
	}
	count, err := client.Fulfillment.Count(context.Background(), nil)
	if count != 0 || !errors.Is(err, ErrFulfillmentRequiresOrder) {
		t.Fatalf("got %v, %v", count, err)
	}
	item, err := client.Fulfillment.Get(context.Background(), 1, nil)
	if item != nil || !errors.Is(err, ErrFulfillmentRequiresOrder) {
		t.Fatalf("got %v, %v", item, err)
	}
	if httpmock.GetTotalCallCount() != 0 {
		t.Fatal("unscoped reads must not send requests")
	}
}

func TestFulfillmentInvalidCreation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Fulfillment
		want  error
	}{
		{"missing fulfillment orders", Fulfillment{}, ErrFulfillmentRequiresFulfillmentOrders},
		{"explicit empty selection", Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10, FulfillmentOrderLineItems: []LineItemByFulfillmentOrderItemQuantity{}}}}, nil},
		{"zero fulfillment order ID", Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{}}}, ErrFulfillmentRequiresFulfillmentOrders},
		{"legacy line items", Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10}}, LineItems: []LineItem{{Id: 20, Quantity: 1}}}, nil},
		{"multiple tracking numbers", Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10}}, TrackingNumbers: []string{"1", "2"}}, nil},
		{"multiple tracking URLs", Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10}}, TrackingUrls: []string{"a", "b"}}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()
			got, err := client.Order.CreateFulfillment(context.Background(), 1, tc.input)
			if got != nil || err == nil || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("got %v, %v", got, err)
			}
			if httpmock.GetTotalCallCount() != 0 {
				t.Fatal("invalid creation must not send requests")
			}
		})
	}
}

func TestFulfillmentInvalidTrackingUpdate(t *testing.T) {
	setup()
	defer teardown()
	got, err := client.Fulfillment.Update(context.Background(), Fulfillment{Id: 1, TrackingNumbers: []string{"1", "2"}})
	if got != nil || err == nil {
		t.Fatalf("got %v, %v", got, err)
	}
	if httpmock.GetTotalCallCount() != 0 {
		t.Fatal("invalid tracking must not send requests")
	}
}

func TestFulfillmentOrderMembership(t *testing.T) {
	for _, action := range []string{"create", "update", "cancel"} {
		for _, tc := range []struct {
			name   string
			status int
			body   string
		}{
			{"not found", 404, `{"errors":"Not Found"}`},
			{"empty", 200, `{}`},
			{"wrong order", 200, `{"fulfillment":{"id":2,"order_id":99},"fulfillment_orders":[{"id":99}]}`},
			{"wrong fulfillment", 200, `{"fulfillment":{"id":99,"order_id":1},"fulfillment_orders":[{"id":99}]}`},
		} {
			t.Run(action+"/"+tc.name, func(t *testing.T) {
				setup()
				defer teardown()
				path := "orders/1/fulfillments/2.json"
				if action == "create" {
					path = "orders/1/fulfillment_orders.json"
				}
				httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, path), httpmock.NewStringResponder(tc.status, tc.body))
				var got *Fulfillment
				var err error
				switch action {
				case "create":
					got, err = client.Order.CreateFulfillment(context.Background(), 1, Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10}}})
				case "update":
					got, err = client.Order.UpdateFulfillment(context.Background(), 1, Fulfillment{Id: 2})
				case "cancel":
					got, err = client.Order.CancelFulfillment(context.Background(), 1, 2)
				}
				if got != nil || err == nil {
					t.Fatalf("got %v, %v", got, err)
				}
				if tc.status == 404 {
					var responseErr ResponseError
					if !errors.As(err, &responseErr) || responseErr.Status != 404 {
						t.Fatalf("lost API error: %v", err)
					}
				}
				if httpmock.GetTotalCallCount() != 1 {
					t.Fatal("failed membership check must not send a write")
				}
			})
		}
	}
}

func TestFulfillmentWriteErrors(t *testing.T) {
	for _, action := range []string{"create", "update", "cancel"} {
		t.Run(action, func(t *testing.T) {
			setup()
			defer teardown()
			paths := map[string]string{"create": "fulfillments.json", "update": "fulfillments/1/update_tracking.json", "cancel": "fulfillments/1/cancel.json"}
			httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, paths[action]), httpmock.NewStringResponder(422, `{"errors":"invalid fulfillment"}`))
			var got *Fulfillment
			var err error
			switch action {
			case "create":
				got, err = client.Fulfillment.Create(context.Background(), Fulfillment{LineItemsByFulfillmentOrder: []LineItemByFulfillmentOrder{{FulfillmentOrderId: 10}}})
			case "update":
				got, err = client.Fulfillment.Update(context.Background(), Fulfillment{Id: 1})
			case "cancel":
				got, err = client.Fulfillment.Cancel(context.Background(), 1)
			}
			var responseErr ResponseError
			if got != nil || !errors.As(err, &responseErr) || responseErr.Status != 422 {
				t.Fatalf("got %v, %v", got, err)
			}
		})
	}
}
