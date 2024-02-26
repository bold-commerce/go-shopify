package goshopify

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestFulfillmentEventServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	orderId := uint64(1234567890)
	fulfillmentId := uint64(987654321)
	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/%d/fulfillments/%d/events.json", client.pathPrefix, orderId, fulfillmentId),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_events.json")),
	)

	fulfillmentEvents, err := client.FulfillmentEvent.List(context.Background(), orderId, fulfillmentId)
	if err != nil {
		t.Errorf("FulfillmentEvent.List returned error: %v", err)
	}

	expected := []FulfillmentEvent{
		{
			Id:                  944956391,
			FulfillmentId:       255858046,
			Status:              "in_transit",
			Message:             "",
			HappenedAt:          "2023-10-20T23:39:23-04:00",
			City:                "",
			Province:            "",
			Country:             "",
			Zip:                 "",
			Address1:            "",
			Latitude:            0,
			Longitude:           0,
			ShopId:              548380009,
			CreatedAt:           "2023-10-20T23:39:23-04:00",
			UpdatedAt:           "2023-10-20T23:39:23-04:00",
			EstimatedDeliveryAt: "",
			OrderId:             450789469,
		},
	}
	if !reflect.DeepEqual(fulfillmentEvents, expected) {
		t.Errorf("FulfillmentEvent.List returned %+v, expected %+v", fulfillmentEvents, expected)
	}
}

func TestFulfillmentEventServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	orderId := uint64(1234567890)
	fulfillmentId := uint64(987654321)
	eventId := uint64(123123123)

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/%d/fulfillments/%d/events/%d.json", client.pathPrefix, orderId, fulfillmentId, eventId),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_event.json")),
	)

	fulfillmentEvent, err := client.FulfillmentEvent.Get(context.Background(), orderId, fulfillmentId, eventId)
	if err != nil {
		t.Errorf("FulfillmentEvent.Get returned error: %v", err)
	}

	expected := &FulfillmentEvent{
		Id:                  944956393,
		FulfillmentId:       255858046,
		Status:              "in_transit",
		Message:             "",
		HappenedAt:          "2023-10-20T23:39:27-04:00",
		City:                "",
		Province:            "",
		Country:             "",
		Zip:                 "",
		Address1:            "",
		Latitude:            0,
		Longitude:           0,
		ShopId:              548380009,
		CreatedAt:           "2023-10-20T23:39:27-04:00",
		UpdatedAt:           "2023-10-20T23:39:27-04:00",
		EstimatedDeliveryAt: "",
		OrderId:             450789469,
	}
	if !reflect.DeepEqual(fulfillmentEvent, expected) {
		t.Errorf("FulfillmentEvent.Get returned %+v, expected %+v", fulfillmentEvent, expected)
	}
}

func TestFulfillmentEventServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	orderId := uint64(1234567890)
	fulfillmentId := uint64(987654321)

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/%d/fulfillments/%d/events.json", client.pathPrefix, orderId, fulfillmentId),
		httpmock.NewBytesResponder(201, loadFixture("fulfillment_event.json")),
	)

	event := FulfillmentEvent{
		Id:                  944956393,
		FulfillmentId:       255858046,
		Status:              "in_transit",
		Message:             "",
		HappenedAt:          "2023-10-20T23:39:27-04:00",
		City:                "",
		Province:            "",
		Country:             "",
		Zip:                 "",
		Address1:            "",
		Latitude:            0,
		Longitude:           0,
		ShopId:              548380009,
		CreatedAt:           "2023-10-20T23:39:27-04:00",
		UpdatedAt:           "2023-10-20T23:39:27-04:00",
		EstimatedDeliveryAt: "",
		OrderId:             450789469,
	}
	createdEvent, err := client.FulfillmentEvent.Create(context.Background(), orderId, fulfillmentId, event)
	if err != nil {
		t.Errorf("FulfillmentEvent.Create returned error: %v", err)
	}

	expected := &FulfillmentEvent{
		Id:                  944956393,
		FulfillmentId:       255858046,
		Status:              "in_transit",
		Message:             "",
		HappenedAt:          "2023-10-20T23:39:27-04:00",
		City:                "",
		Province:            "",
		Country:             "",
		Zip:                 "",
		Address1:            "",
		Latitude:            0,
		Longitude:           0,
		ShopId:              548380009,
		CreatedAt:           "2023-10-20T23:39:27-04:00",
		UpdatedAt:           "2023-10-20T23:39:27-04:00",
		EstimatedDeliveryAt: "",
		OrderId:             450789469,
	}
	if !reflect.DeepEqual(createdEvent, expected) {
		t.Errorf("FulfillmentEvent.Create returned %+v, expected %+v", createdEvent, expected)
	}
}

func TestFulfillmentEventServiceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	orderId := uint64(1234567890)
	fulfillmentId := uint64(987654321)
	eventId := uint64(123123123)

	httpmock.RegisterResponder(
		http.MethodDelete,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/%d/fulfillments/%d/events/%d.json", client.pathPrefix, orderId, fulfillmentId, eventId),
		httpmock.NewStringResponder(200, ""),
	)

	err := client.FulfillmentEvent.Delete(context.Background(), orderId, fulfillmentId, eventId)
	if err != nil {
		t.Errorf("FulfillmentEvent.Delete returned error: %v", err)
	}
}
