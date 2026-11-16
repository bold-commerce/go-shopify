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
)

func TestFulfillmentRequestServiceOp_Send(t *testing.T) {
	setup()
	defer teardown()

	fulfillmentOrderId := uint64(1046000829)
	message := "Fulfill this ASAP please."
	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request.json", client.pathPrefix, fulfillmentOrderId),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_send.json")),
	)

	result, err := client.FulfillmentRequest.Send(context.Background(), fulfillmentOrderId, FulfillmentRequest{Message: message})
	if err != nil {
		t.Errorf("FulfillmentRequest.Send returned error: %v", err)
	}

	sentAt, err := time.Parse(time.RFC3339, "2023-10-03T13:22:49-04:00")
	if err != nil {
		t.Fatal(err)
	}
	expected := &FulfillmentOrder{
		Origin: &FulfillmentOrderAssignedLocation{CountryCode: "DE", LocationId: 24826418, Name: "Apple Api Shipwire"},
		OutgoingRequests: []FulfillmentOrderOutgoingRequest{{
			Kind: "fulfillment_request", Message: message, SentAt: &sentAt,
			RequestOptions: json.RawMessage("{\n          \"notify_customer\": false\n        }"),
		}},
		Id:                 1046000829,
		ShopId:             548380009,
		OrderId:            450789469,
		AssignedLocationId: 24826418,
		RequestStatus:      "submitted",
		Status:             "open",
		SupportedActions:   []string{"cancel_fulfillment_order"},
		Destination: FulfillmentOrderDestination{
			Id:        1046000816,
			Address1:  "Chestnut Street 92",
			City:      "Louisville",
			Company:   "",
			Country:   "United States",
			Email:     "bob.norman@mail.example.com",
			FirstName: "Bob",
			LastName:  "Norman",
			Phone:     "+1(502)-459-2181",
			Province:  "Kentucky",
			Zip:       "40202",
		},
		LineItems: []FulfillmentOrderLineItem{
			{
				Id:                  1058737567,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000829,
				Quantity:            1,
				LineItemId:          466157049,
				InventoryItemId:     39072856,
				FulfillableQuantity: 1,
				VariantId:           39072856,
			},
			{
				Id:                  1058737568,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000829,
				Quantity:            1,
				LineItemId:          518995019,
				InventoryItemId:     49148385,
				FulfillableQuantity: 1,
				VariantId:           49148385,
			},
			{
				Id:                  1058737569,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000829,
				Quantity:            1,
				LineItemId:          703073504,
				InventoryItemId:     457924702,
				FulfillableQuantity: 1,
				VariantId:           457924702,
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FulfillmentRequest.Send returned %+v, expected %+v", result, expected)
	}
}

func TestFulfillmentRequestServiceOp_Accept(t *testing.T) {
	setup()
	defer teardown()

	fulfillmentOrderId := uint64(1046000828)
	message := "We will start processing your fulfillment on the next business day."

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request/accept.json", client.pathPrefix, fulfillmentOrderId),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_accept.json")),
	)

	result, err := client.FulfillmentRequest.Accept(context.Background(), fulfillmentOrderId, FulfillmentRequest{Message: message})
	if err != nil {
		t.Errorf("FulfillmentRequest.Accept returned error: %v", err)
	}

	expected := &FulfillmentOrder{
		Origin:             &FulfillmentOrderAssignedLocation{CountryCode: "DE", LocationId: 24826418, Name: "Apple Api Shipwire"},
		OutgoingRequests:   []FulfillmentOrderOutgoingRequest{},
		Id:                 1046000828,
		ShopId:             548380009,
		OrderId:            450789469,
		AssignedLocationId: 24826418,
		RequestStatus:      "accepted",
		Status:             "in_progress",
		SupportedActions:   []string{"request_cancellation", "create_fulfillment"},
		Destination: FulfillmentOrderDestination{
			Id:        1046000815,
			Address1:  "Chestnut Street 92",
			City:      "Louisville",
			Company:   "",
			Country:   "United States",
			Email:     "bob.norman@mail.example.com",
			FirstName: "Bob",
			LastName:  "Norman",
			Phone:     "+1(502)-459-2181",
			Province:  "Kentucky",
			Zip:       "40202",
		},
		LineItems: []FulfillmentOrderLineItem{
			{
				Id:                  1058737564,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000828,
				Quantity:            1,
				LineItemId:          466157049,
				InventoryItemId:     39072856,
				FulfillableQuantity: 1,
				VariantId:           39072856,
			},
			{
				Id:                  1058737565,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000828,
				Quantity:            1,
				LineItemId:          518995019,
				InventoryItemId:     49148385,
				FulfillableQuantity: 1,
				VariantId:           49148385,
			},
			{
				Id:                  1058737566,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000828,
				Quantity:            1,
				LineItemId:          703073504,
				InventoryItemId:     457924702,
				FulfillableQuantity: 1,
				VariantId:           457924702,
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FulfillmentRequest.Accept returned %+v, expected %+v", result, expected)
	}
}

func TestFulfillmentRequestServiceOp_Reject(t *testing.T) {
	setup()
	defer teardown()

	fulfillmentOrderId := uint64(1046000830)
	rejectionMessage := "Not enough inventory on hand to complete the work."

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request/reject.json", client.pathPrefix, fulfillmentOrderId),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_reject.json")),
	)

	result, err := client.FulfillmentRequest.Reject(context.Background(), fulfillmentOrderId, FulfillmentRequest{Message: rejectionMessage})
	if err != nil {
		t.Errorf("FulfillmentRequest.Reject returned error: %v", err)
	}

	expected := &FulfillmentOrder{
		Origin:             &FulfillmentOrderAssignedLocation{CountryCode: "DE", LocationId: 24826418, Name: "Apple Api Shipwire"},
		OutgoingRequests:   []FulfillmentOrderOutgoingRequest{},
		Id:                 1046000830,
		ShopId:             548380009,
		OrderId:            450789469,
		AssignedLocationId: 24826418,
		RequestStatus:      "rejected",
		Status:             "open",
		SupportedActions:   []string{"request_fulfillment", "create_fulfillment"},
		Destination: FulfillmentOrderDestination{
			Id:        1046000817,
			Address1:  "Chestnut Street 92",
			City:      "Louisville",
			Company:   "",
			Country:   "United States",
			Email:     "bob.norman@mail.example.com",
			FirstName: "Bob",
			LastName:  "Norman",
			Phone:     "+1(502)-459-2181",
			Province:  "Kentucky",
			Zip:       "40202",
		},
		LineItems: []FulfillmentOrderLineItem{
			{
				Id:                  1058737570,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000830,
				Quantity:            1,
				LineItemId:          466157049,
				InventoryItemId:     39072856,
				FulfillableQuantity: 1,
				VariantId:           39072856,
			},
			{
				Id:                  1058737571,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000830,
				Quantity:            1,
				LineItemId:          518995019,
				InventoryItemId:     49148385,
				FulfillableQuantity: 1,
				VariantId:           49148385,
			},
			{
				Id:                  1058737572,
				ShopId:              548380009,
				FulfillmentOrderId:  1046000830,
				Quantity:            1,
				LineItemId:          703073504,
				InventoryItemId:     457924702,
				FulfillableQuantity: 1,
				VariantId:           457924702,
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FulfillmentRequest.Reject returned %+v, expected %+v", result, expected)
	}
}

func TestFulfillmentRequestSendWithResponse(t *testing.T) {
	setup()
	defer teardown()
	for _, tc := range []struct {
		name, body             string
		submitted, unsubmitted uint64
	}{
		{"split", `{"original_fulfillment_order":{"id":1},"submitted_fulfillment_order":{"id":2},"unsubmitted_fulfillment_order":{"id":3}}`, 2, 3},
		{"whole", `{"original_fulfillment_order":{"id":1},"submitted_fulfillment_order":{"id":1},"unsubmitted_fulfillment_order":null}`, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/fulfillment_request.json", client.pathPrefix), func(req *http.Request) (*http.Response, error) {
				var body map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				want := map[string]interface{}{"fulfillment_request": map[string]interface{}{"message": "please fulfill"}}
				if !reflect.DeepEqual(body, want) {
					t.Errorf("body = %#v, want %#v", body, want)
				}
				return httpmock.NewStringResponse(200, tc.body), nil
			})
			got, err := client.FulfillmentRequest.SendWithResponse(context.Background(), 1, FulfillmentRequest{Message: "please fulfill"})
			if err != nil {
				t.Fatal(err)
			}
			if got.OriginalFulfillmentOrder == nil || got.OriginalFulfillmentOrder.Id != 1 {
				t.Fatalf("original = %+v", got.OriginalFulfillmentOrder)
			}
			if got.SubmittedFulfillmentOrder == nil || got.SubmittedFulfillmentOrder.Id != tc.submitted {
				t.Fatalf("submitted = %+v", got.SubmittedFulfillmentOrder)
			}
			if tc.unsubmitted == 0 {
				if got.UnsubmittedFulfillmentOrder != nil {
					t.Fatal("expected no unsubmitted order")
				}
			} else if got.UnsubmittedFulfillmentOrder == nil || got.UnsubmittedFulfillmentOrder.Id != tc.unsubmitted {
				t.Fatalf("unsubmitted = %+v", got.UnsubmittedFulfillmentOrder)
			}
		})
	}
}

func TestFulfillmentRequestSendErrors(t *testing.T) {
	setup()
	defer teardown()
	for _, body := range []string{`{"errors":"cannot submit"}`, `invalid json`} {
		status := 422
		if body == "invalid json" {
			status = 200
		}
		httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/1/fulfillment_request.json", client.pathPrefix), httpmock.NewStringResponder(status, body))
		if _, err := client.FulfillmentRequest.SendWithResponse(context.Background(), 1, FulfillmentRequest{}); err == nil {
			t.Fatal("expected full response error")
		}
		if got, err := client.FulfillmentRequest.Send(context.Background(), 1, FulfillmentRequest{}); err == nil || got != nil {
			t.Fatalf("Send = %+v, %v; want nil and error", got, err)
		}
	}
}
