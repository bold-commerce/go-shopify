package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func productTests(t *testing.T, product Product) {
	// Check that ID is assigned to the returned product
	var expectedInt int64 = 1071559748
	if product.ID != expectedInt {
		t.Errorf("Product.ID returned %+v, expected %+v", product.ID, expectedInt)
	}
}

func TestProductList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"products": [{"id":1},{"id":2}]}`))

	products, err := client.Product.List(nil)
	if err != nil {
		t.Errorf("Product.List returned error: %v", err)
	}

	expected := []Product{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Product.List returned %+v, expected %+v", products, expected)
	}
}

func TestProductListFilterByIds(t *testing.T) {
	setup()
	defer teardown()

	params := map[string]string{"ids": "1,2,3"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"products": [{"id":1},{"id":2},{"id":3}]}`))

	listOptions := ListOptions{IDs: []int64{1, 2, 3}}

	products, err := client.Product.List(listOptions)
	if err != nil {
		t.Errorf("Product.List returned error: %v", err)
	}

	expected := []Product{{ID: 1}, {ID: 2}, {ID: 3}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Product.List returned %+v, expected %+v", products, expected)
	}
}

func TestProductCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Product.Count(nil)
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Product.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}
}

func TestProductGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"product": {"id":1}}`))

	product, err := client.Product.Get(1, nil)
	if err != nil {
		t.Errorf("Product.Get returned error: %v", err)
	}

	expected := &Product{ID: 1}
	if !reflect.DeepEqual(product, expected) {
		t.Errorf("Product.Get returned %+v, expected %+v", product, expected)
	}
}

func TestProductCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("product.json")))

	product := Product{
		Title:       "Burton Custom Freestyle 151",
		BodyHTML:    "<strong>Good snowboard!<\\/strong>",
		Vendor:      "Burton",
		ProductType: "Snowboard",
	}

	returnedProduct, err := client.Product.Create(product)
	if err != nil {
		t.Errorf("Product.Create returned error: %v", err)
	}

	productTests(t, *returnedProduct)
}

func TestProductUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("product.json")))

	product := Product{
		ID:          1,
		ProductType: "Skateboard",
	}

	returnedProduct, err := client.Product.Update(product)
	if err != nil {
		t.Errorf("Product.Update returned error: %v", err)
	}

	productTests(t, *returnedProduct)
}

func TestProductDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Product.Delete(1)
	if err != nil {
		t.Errorf("Product.Delete returned error: %v", err)
	}
}

func TestProductListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Product.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("Product.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Product.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestProductCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Product.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("Product.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Product.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Product.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Product.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Product.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestProductGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Product.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("Product.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Product.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestProductCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Product.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Product.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestProductUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Product.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Product.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestProductDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Product.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Product.DeleteMetafield() returned error: %v", err)
	}
}
