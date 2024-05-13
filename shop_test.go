package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestShopGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shop.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("shop.json")))

	shop, err := client.Shop.Get(context.Background(), nil)
	if err != nil {
		t.Errorf("Shop.Get returned error: %v", err)
	}

	// Check that dates are parsed
	d := time.Date(2007, time.December, 31, 19, 0o0, 0o0, 0, time.UTC)
	if !d.Equal(*shop.CreatedAt) {
		t.Errorf("Shop.CreatedAt returned %+v, expected %+v", shop.CreatedAt, d)
	}

	// Test a few fields
	cases := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"Id", uint64(690933842), shop.Id},
		{"ShopOwner", "Steve Jobs", shop.ShopOwner},
		{"Address1", "1 Infinite Loop", shop.Address1},
		{"Name", "Apple Computers", shop.Name},
		{"Email", "steve@apple.com", shop.Email},
		{"HasStorefront", true, shop.HasStorefront},
		{"Source", "", shop.Source},
		{"GoogleAppsDomain", "", shop.GoogleAppsDomain},
		{"GoogleAppsLoginEnabled", false, shop.GoogleAppsLoginEnabled},
		{"MoneyInEmailsFormat", "${{amount}}", shop.MoneyInEmailsFormat},
		{"MoneyWithCurrencyInEmailsFormat", "${{amount}} USD", shop.MoneyWithCurrencyInEmailsFormat},
		{"EligibleForPayments", true, shop.EligibleForPayments},
		{"RequiresExtraPaymentsAgreement", false, shop.RequiresExtraPaymentsAgreement},
		{"PreLaunchEnabled", false, shop.PreLaunchEnabled},
	}

	for _, c := range cases {
		if c.expected != c.actual {
			t.Errorf("Shop.%v returned %v, expected %v", c.field, c.actual, c.expected)
		}
	}
}

func TestShopListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Shop.ListMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Shop.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Shop.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestShopCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Shop.CountMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Shop.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Shop.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Shop.CountMetafields(context.Background(), 1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Shop.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Shop.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestShopGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Shop.GetMetafield(context.Background(), 1, 2, nil)
	if err != nil {
		t.Errorf("Shop.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{Id: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Shop.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestShopCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Shop.CreateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Shop.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestShopUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Id:        2,
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Shop.UpdateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Shop.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestShopDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Shop.DeleteMetafield(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("Shop.DeleteMetafield() returned error: %v", err)
	}
}
