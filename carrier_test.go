package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCarrierList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/carrier_services.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("carrier_services.json")))

	carriers, err := client.CarrierService.List(context.Background())
	if err != nil {
		t.Errorf("Carrier.List returned error: %v", err)
	}

	trueVar := true

	expected := []CarrierService{
		{
			Id:                 1,
			Name:               "Shipping Rate Provider",
			Active:             &trueVar,
			ServiceDiscovery:   true,
			CarrierServiceType: "api",
			AdminGraphqlApiId:  "gid://shopify/DeliveryCarrierService/1",
			Format:             "json",
			CallbackUrl:        "https://fooshop.example.com/shipping",
		},
	}
	if !reflect.DeepEqual(carriers, expected) {
		t.Errorf("Carrier.List returned %+v, expected %+v", carriers, expected)
	}
}

func TestCarrierGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/carrier_services/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("carrier_service.json")))

	carrier, err := client.CarrierService.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("Carrier.Get returned error: %v", err)
	}

	trueVar := true

	expected := &CarrierService{
		Id:                 1,
		Name:               "Shipping Rate Provider",
		Active:             &trueVar,
		ServiceDiscovery:   true,
		CarrierServiceType: "api",
		AdminGraphqlApiId:  "gid://shopify/DeliveryCarrierService/1",
		Format:             "json",
		CallbackUrl:        "https://fooshop.example.com/shipping",
	}
	if !reflect.DeepEqual(carrier, expected) {
		t.Errorf("Carrier.Get returned %+v, expected %+v", carrier, expected)
	}
}

func TestCarrierCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/carrier_services.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("carrier_service.json")))

	carrier, err := client.CarrierService.Create(context.Background(), CarrierService{})
	if err != nil {
		t.Errorf("Carrier.Create returned error: %v", err)
	}

	trueVar := true

	expected := &CarrierService{
		Id:                 1,
		Name:               "Shipping Rate Provider",
		Active:             &trueVar,
		ServiceDiscovery:   true,
		CarrierServiceType: "api",
		AdminGraphqlApiId:  "gid://shopify/DeliveryCarrierService/1",
		Format:             "json",
		CallbackUrl:        "https://fooshop.example.com/shipping",
	}
	if !reflect.DeepEqual(carrier, expected) {
		t.Errorf("Carrier.Create returned %+v, expected %+v", carrier, expected)
	}
}

func TestCarrierUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/carrier_services/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("carrier_service.json")))

	carrier, err := client.CarrierService.Update(context.Background(), CarrierService{Id: 1})
	if err != nil {
		t.Errorf("Carrier.Update returned error: %v", err)
	}

	trueVar := true

	expected := &CarrierService{
		Id:                 1,
		Name:               "Shipping Rate Provider",
		Active:             &trueVar,
		ServiceDiscovery:   true,
		CarrierServiceType: "api",
		AdminGraphqlApiId:  "gid://shopify/DeliveryCarrierService/1",
		Format:             "json",
		CallbackUrl:        "https://fooshop.example.com/shipping",
	}
	if !reflect.DeepEqual(carrier, expected) {
		t.Errorf("Carrier.Update returned %+v, expected %+v", carrier, expected)
	}
}

func TestCarrierDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/carrier_services/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{}`))

	err := client.CarrierService.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Carrier.Delete returned error: %v", err)
	}
}
