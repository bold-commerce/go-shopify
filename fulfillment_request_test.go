package goshopify

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestFulfillmentRequestServiceOp_Send(t *testing.T) {
	setup()
	defer teardown()

	fulfillmentOrderID := int64(1046000829)
	message := "Fulfill this ASAP please."
	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request.json", client.pathPrefix, fulfillmentOrderID),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_send.json")),
	)

	result, err := client.FulfillmentRequest.Send(context.Background(), fulfillmentOrderID, FulfillmentRequest{Message: message})
	if err != nil {
		t.Errorf("FulfillmentRequest.Send returned error: %v", err)
	}

	expected := &FulfillmentOrder{
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

	fulfillmentOrderID := int64(1046000828)
	message := "We will start processing your fulfillment on the next business day."

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request/accept.json", client.pathPrefix, fulfillmentOrderID),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_accept.json")),
	)

	result, err := client.FulfillmentRequest.Accept(context.Background(), fulfillmentOrderID, FulfillmentRequest{Message: message})
	if err != nil {
		t.Errorf("FulfillmentRequest.Accept returned error: %v", err)
	}

	expected := &FulfillmentOrder{
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

	fulfillmentOrderID := int64(1046000830)
	rejectionMessage := "Not enough inventory on hand to complete the work."

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_orders/%d/fulfillment_request/reject.json", client.pathPrefix, fulfillmentOrderID),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_reject.json")),
	)

	result, err := client.FulfillmentRequest.Reject(context.Background(), fulfillmentOrderID, FulfillmentRequest{Message: rejectionMessage})
	if err != nil {
		t.Errorf("FulfillmentRequest.Reject returned error: %v", err)
	}

	expected := &FulfillmentOrder{
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
