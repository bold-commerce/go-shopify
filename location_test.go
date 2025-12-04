package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestLocationServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("locations.json")))

	products, err := client.Location.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Location.List returned error: %v", err)
	}

	created, _ := time.Parse(time.RFC3339, "2018-02-19T16:18:59-05:00")
	updated, _ := time.Parse(time.RFC3339, "2018-02-19T16:19:00-05:00")

	expected := []Location{{
		Id:                4688969785,
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
		AdminGraphqlApiId: "gid://shopify/Location/4688969785",
	}}

	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Location.List returned %+v, expected %+v", products, expected)
	}
}

func TestLocationServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/4688969785.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("location.json")))

	product, err := client.Location.Get(context.Background(), 4688969785, nil)
	if err != nil {
		t.Errorf("Location.Get returned error: %v", err)
	}

	created, _ := time.Parse(time.RFC3339, "2018-02-19T16:18:59-05:00")
	updated, _ := time.Parse(time.RFC3339, "2018-02-19T16:19:00-05:00")

	expected := &Location{
		Id:                4688969785,
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
		AdminGraphqlApiId: "gid://shopify/Location/4688969785",
	}

	if !reflect.DeepEqual(product, expected) {
		t.Errorf("Location.Get returned %+v, expected %+v", product, expected)
	}
}

func TestLocationServiceOp_Count(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	cnt, err := client.Location.Count(context.Background(), nil)
	if err != nil {
		t.Errorf("Location.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Location.Count returned %d, expected %d", cnt, expected)
	}
}

func TestLocationListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Location.ListMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Location.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Location.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestLocationCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Location.CountMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Location.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Location.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Location.CountMetafields(context.Background(), 1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Location.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Location.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestLocationGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Location.GetMetafield(context.Background(), 1, 2, nil)
	if err != nil {
		t.Errorf("Location.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{Id: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Location.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestLocationCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Location.CreateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Location.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestLocationUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Id:        2,
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Location.UpdateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Location.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestLocationDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/locations/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Location.DeleteMetafield(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("Location.DeleteMetafield() returned error: %v", err)
	}
}
