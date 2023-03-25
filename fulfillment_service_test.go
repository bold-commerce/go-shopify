package goshopify

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestFulfillmentServiceServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_services.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_services.json")),
	)

	options := FulfillmentServiceOptions{Scope: "all"}

	fulfillmentServices, err := client.FulfillmentService.List(options)
	if err != nil {
		t.Errorf("fulfillmentService.List returned error: %v", err)
	}

	expected := []FulfillmentServiceData{
		{
			Id:                     1061774487,
			Name:                   "Jupiter Fulfillment",
			Email:                  "aaa@gmail.com",
			ServiceName:            "Jupiter Fulfillment",
			Handle:                 "jupiter-fulfillment",
			FulfillmentOrdersOptIn: false,
			IncludePendingStock:    false,
			ProviderId:             1234,
			LocationId:             1072404542,
			CallbackURL:            "https://google.com/",
			TrackingSupport:        false,
			InventoryManagement:    false,
			AdminGraphqlApiId:      "gid://shopify/ApiFulfillmentService/1061774487",
			PermitsSkuSharing:      false,
		},
	}
	if !reflect.DeepEqual(fulfillmentServices, expected) {
		t.Errorf("fulfillmentService.List returned %+v, expected %+v", fulfillmentServices, expected)
	}
}

func TestFulfillmentServiceServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_services/1061774487.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_service.json")),
	)

	fulfillmentService, err := client.FulfillmentService.Get(1061774487, nil)
	if err != nil {
		t.Errorf("FulfillmentService.Get returned error: %v", err)
	}

	expected := &FulfillmentServiceData{
		Id:                     1061774487,
		Name:                   "Jupiter Fulfillment",
		Email:                  "aaa@gmail.com",
		ServiceName:            "Jupiter Fulfillment",
		Handle:                 "jupiter-fulfillment",
		FulfillmentOrdersOptIn: false,
		IncludePendingStock:    false,
		ProviderId:             1234,
		LocationId:             1072404542,
		CallbackURL:            "https://google.com/",
		TrackingSupport:        false,
		InventoryManagement:    false,
		AdminGraphqlApiId:      "gid://shopify/ApiFulfillmentService/1061774487",
		PermitsSkuSharing:      false,
	}
	if !reflect.DeepEqual(fulfillmentService, expected) {
		t.Errorf("FulfillmentService.Get returned %+v, expected %+v", fulfillmentService, expected)
	}
}

func TestFulfillmentServiceServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_services.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_service.json")),
	)

	fulfillmentService, err := client.FulfillmentService.Create(FulfillmentServiceData{
		Name: "jupiter-fulfillment",
	})
	if err != nil {
		t.Errorf("FulfillmentService.Get returned error: %v", err)
	}

	expectedFulfillmentServiceID := int64(1061774487)
	if fulfillmentService.Id != expectedFulfillmentServiceID {
		t.Errorf("FulfillmentService.Id returned %+v, expected %+v", fulfillmentService.Id, expectedFulfillmentServiceID)
	}
}

func TestFulfillmentServiceServiceOp_Update(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		http.MethodPut,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_services/1061774487.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment_service.json")),
	)

	fulfillmentService, err := client.FulfillmentService.Update(FulfillmentServiceData{
		Id:     1061774487,
		Handle: "jupiter-fulfillment",
	})
	if err != nil {
		t.Errorf("FulfillmentService.Update returned error: %v", err)
	}

	expected := &FulfillmentServiceData{
		Id:                     1061774487,
		Name:                   "Jupiter Fulfillment",
		Email:                  "aaa@gmail.com",
		ServiceName:            "Jupiter Fulfillment",
		Handle:                 "jupiter-fulfillment",
		FulfillmentOrdersOptIn: false,
		IncludePendingStock:    false,
		ProviderId:             1234,
		LocationId:             1072404542,
		CallbackURL:            "https://google.com/",
		TrackingSupport:        false,
		InventoryManagement:    false,
		AdminGraphqlApiId:      "gid://shopify/ApiFulfillmentService/1061774487",
		PermitsSkuSharing:      false,
	}
	if !reflect.DeepEqual(fulfillmentService, expected) {
		t.Errorf("FulfillmentService.Update returned %+v, expected %+v", fulfillmentService, expected)
	}
}

func TestFulfillmentServiceServiceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		http.MethodDelete,
		fmt.Sprintf("https://fooshop.myshopify.com/%s/fulfillment_services/1061774487.json", client.pathPrefix),
		httpmock.NewStringResponder(200, ""),
	)

	if err := client.FulfillmentService.Delete(1061774487); err != nil {
		t.Errorf("FulfillmentService.Delete returned error: %v", err)
	}
}
