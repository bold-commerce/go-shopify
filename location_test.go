package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestLocationServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("locations.json")))

	products, err := client.Location.List(nil)
	if err != nil {
		t.Errorf("Location.List returned error: %v", err)
	}

	created, _ := time.Parse(time.RFC3339, "2018-02-19T16:18:59-05:00")
	updated, _ := time.Parse(time.RFC3339, "2018-02-19T16:19:00-05:00")

	expected := []Location{{
		ID:                4688969785,
		Name:              "Bajkowa",
		Address1:          "Bajkowa",
		Address2:          "",
		City:              "Olsztyn",
		Zip:               "10-001",
		Country:           "PL",
		Phone:             "12312312",
		CreatedAt:         created,
		UpdatedAt:         updated,
		CountryCode:       "PL",
		CountryName:       "Poland",
		Legacy:            false,
		Active:            true,
		AdminGraphqlAPIID: "gid://shopify/Location/4688969785",
	}}

	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Location.List returned %+v, expected %+v", products, expected)
	}
}

func TestLocationServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/4688969785.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("location.json")))

	product, err := client.Location.Get(4688969785, nil)
	if err != nil {
		t.Errorf("Location.Get returned error: %v", err)
	}

	created, _ := time.Parse(time.RFC3339, "2018-02-19T16:18:59-05:00")
	updated, _ := time.Parse(time.RFC3339, "2018-02-19T16:19:00-05:00")

	expected := &Location{
		ID:                4688969785,
		Name:              "Bajkowa",
		Address1:          "Bajkowa",
		Address2:          "",
		City:              "Olsztyn",
		Zip:               "10-001",
		Country:           "PL",
		Phone:             "12312312",
		CreatedAt:         created,
		UpdatedAt:         updated,
		CountryCode:       "PL",
		CountryName:       "Poland",
		Legacy:            false,
		Active:            true,
		AdminGraphqlAPIID: "gid://shopify/Location/4688969785",
	}

	if !reflect.DeepEqual(product, expected) {
		t.Errorf("Location.Get returned %+v, expected %+v", product, expected)
	}
}

func TestLocationServiceOp_Count(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	cnt, err := client.Location.Count(nil)
	if err != nil {
		t.Errorf("Location.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Location.Count returned %d, expected %d", cnt, expected)
	}
}
