package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func variantTests(t *testing.T, variant Variant) {
	// Check that the ID is assigned to the returned variant
	expectedInt := int64(1)
	if variant.ID != expectedInt {
		t.Errorf("Variant.ID returned %+v, expected %+v", variant.ID, expectedInt)
	}

	// Check that the Title is assigned to the returned variant
	expectedTitle := "Yellow"
	if variant.Title != expectedTitle {
		t.Errorf("Variant.Title returned %+v, expected %+v", variant.Title, expectedTitle)
	}

	expectedInventoryItemId := int64(1)
	if variant.InventoryItemId != expectedInventoryItemId {
		t.Errorf("Variant.InventoryItemId returned %+v, expected %+v", variant.InventoryItemId, expectedInventoryItemId)
	}

	expectedMetafieldCount := 0
	if len(variant.Metafields) != expectedMetafieldCount {
		t.Errorf("Variant.Metafield returned %+v, expected %+v", variant.Metafields, expectedMetafieldCount)
	}
}

func variantWithMetafieldsTests(t *testing.T, variant Variant) {
	// Check that the ID is assigned to the returned variant
	expectedInt := int64(2)
	if variant.ID != expectedInt {
		t.Errorf("Variant.ID returned %+v, expected %+v", variant.ID, expectedInt)
	}

	// Check that the Title is assigned to the returned variant
	expectedTitle := "Blue"
	if variant.Title != expectedTitle {
		t.Errorf("Variant.Title returned %+v, expected %+v", variant.Title, expectedTitle)
	}

	expectedInventoryItemId := int64(1)
	if variant.InventoryItemId != expectedInventoryItemId {
		t.Errorf("Variant.InventoryItemId returned %+v, expected %+v", variant.InventoryItemId, expectedInventoryItemId)
	}

	expectedMetafieldCount := 1
	if len(variant.Metafields) != expectedMetafieldCount {
		t.Errorf("Variant.Metafield returned %+v, expected %+v", variant.Metafields, expectedMetafieldCount)
	}

	expectedMetafieldDescription := "description"
	if variant.Metafields[0].Description != expectedMetafieldDescription {
		t.Errorf("Variant.Metafield.Description returned %+v, expected %+v", variant.Metafields[0].Description, expectedMetafieldDescription)
	}
}

func TestVariantList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"variants": [{"id":1},{"id":2}]}`))

	variants, err := client.Variant.List(1, nil)
	if err != nil {
		t.Errorf("Variant.List returned error: %v", err)
	}

	expected := []Variant{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(variants, expected) {
		t.Errorf("Variant.List returned %+v, expected %+v", variants, expected)
	}
}

func TestVariantCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Variant.Count(1, nil)
	if err != nil {
		t.Errorf("Variant.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Variant.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Variant.Count(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Variant.Count returned %d, expected %d", cnt, expected)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Variant.Count returned %d, expected %d", cnt, expected)
	}
}

func TestVariantGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"variant": {"id":1}}`))

	variant, err := client.Variant.Get(1, nil)
	if err != nil {
		t.Errorf("Variant.Get returned error: %v", err)
	}

	expected := &Variant{ID: 1}
	if !reflect.DeepEqual(variant, expected) {
		t.Errorf("Variant.Get returned %+v, expected %+v", variant, expected)
	}
}

func TestVariantCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("variant.json")))

	price := decimal.NewFromFloat(1)

	variant := Variant{
		Option1: "Yellow",
		Price:   &price,
	}
	result, err := client.Variant.Create(1, variant)
	if err != nil {
		t.Errorf("Variant.Create returned error: %v", err)
	}
	variantTests(t, *result)
}

func TestVariantCreateWithMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("variant_with_metafields.json")))

	price := decimal.NewFromFloat(2)

	variant := Variant{
		Option1: "Blue",
		Price:   &price,
	}
	result, err := client.Variant.Create(1, variant)
	if err != nil {
		t.Errorf("Variant.Create returned error: %v", err)
	}
	variantWithMetafieldsTests(t, *result)
}

func TestVariantUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("variant.json")))

	variant := Variant{
		ID:      1,
		Option1: "Green",
	}

	variant.Option1 = "Yellow"

	returnedVariant, err := client.Variant.Update(variant)
	if err != nil {
		t.Errorf("Variant.Update returned error: %v", err)
	}
	variantTests(t, *returnedVariant)
}

func TestVariantWithMetafieldsUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("variant_with_metafields.json")))

	variant := Variant{
		ID:      2,
		Option1: "Green",
		Metafields: []Metafield{
			{
				ID:          123,
				Description: "Original",
			},
		},
	}

	variant.Option1 = "Blue"
	variant.Metafields[0].Description = "description"

	returnedVariant, err := client.Variant.Update(variant)
	if err != nil {
		t.Errorf("Variant.Update returned error: %v", err)
	}

	variantWithMetafieldsTests(t, *returnedVariant)
}

func TestVariantDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Variant.Delete(1, 1)
	if err != nil {
		t.Errorf("Variant.Delete returned error: %v", err)
	}
}

func TestVariantListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Variant.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("Variant.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Variant.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestVariantCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Variant.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("Variant.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Variant.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Variant.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Variant.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Variant.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestVariantGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Variant.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("Variant.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Variant.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestVariantCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Variant.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Variant.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestVariantUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Variant.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Variant.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestVariantDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Variant.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Variant.DeleteMetafield() returned error: %v", err)
	}
}

func TestVariantListWithTaxCode(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"variants": [{"id":1, "tax_code":"P0000000"},{"id":2, "tax_code":"P0000000"}]}`))

	variants, err := client.Variant.List(1, nil)
	if err != nil {
		t.Errorf("Variant.List returned error: %v", err)
	}

	expected := []Variant{{ID: 1, TaxCode: "P0000000"}, {ID: 2, TaxCode: "P0000000"}}
	if !reflect.DeepEqual(variants, expected) {
		t.Errorf("Variant.List returned %+v, expected %+v", variants, expected)
	}
}

func TestVariantGetWithTaxCode(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/variants/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"variant": {"id":1, "tax_code":"P0000000"}}`))

	variant, err := client.Variant.Get(1, nil)
	if err != nil {
		t.Errorf("Variant.Get returned error: %v", err)
	}

	expected := &Variant{ID: 1, TaxCode: "P0000000"}
	if !reflect.DeepEqual(variant, expected) {
		t.Errorf("Variant.Get returned %+v, expected %+v", variant, expected)
	}
}

func TestVariantCreateWithTaxCode(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/variants.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("variant_with_taxcode.json")))

	price := decimal.NewFromFloat(1)

	variant := Variant{
		Option1: "Yellow",
		Price:   &price,
		TaxCode: "P0000000",
	}
	result, err := client.Variant.Create(1, variant)
	if err != nil {
		t.Errorf("Variant.Create returned error: %v", err)
	}
	variantTestsWithTaxCode(t, *result)
}

func variantTestsWithTaxCode(t *testing.T, variant Variant) {
	// Check that the ID is assigned to the returned variant
	expectedInt := int64(1)
	if variant.ID != expectedInt {
		t.Errorf("Variant.ID returned %+v, expected %+v", variant.ID, expectedInt)
	}

	// Check that the Title is assigned to the returned variant
	expectedTitle := "Green"
	if variant.Title != expectedTitle {
		t.Errorf("Variant.Title returned %+v, expected %+v", variant.Title, expectedTitle)
	}

	expectedInventoryItemId := int64(1)
	if variant.InventoryItemId != expectedInventoryItemId {
		t.Errorf("Variant.InventoryItemId returned %+v, expected %+v", variant.InventoryItemId, expectedInventoryItemId)
	}

	expectedMetafieldCount := 0
	if len(variant.Metafields) != expectedMetafieldCount {
		t.Errorf("Variant.Metafield returned %+v, expected %+v", variant.Metafields, expectedMetafieldCount)
	}

	// Check that the Tax_code is assigned to the returned variant
	expectedTacCode := "P0000000"
	if variant.TaxCode != expectedTacCode {
		t.Errorf("Variant.TaxCode returned %+v, expected %+v", variant.TaxCode, expectedTacCode)
	}

}
