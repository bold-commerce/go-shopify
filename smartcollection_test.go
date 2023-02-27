package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func smartCollectionTests(t *testing.T, collection SmartCollection) {
	// Test a few fields
	cases := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", int64(30497275952), collection.ID},
		{"Handle", "macbooks", collection.Handle},
		{"Title", "Macbooks", collection.Title},
		{"BodyHTML", "Macbook Body", collection.BodyHTML},
		{"SortOrder", "best-selling", collection.SortOrder},
		{"Column", "title", collection.Rules[0].Column},
		{"Relation", "contains", collection.Rules[0].Relation},
		{"Condition", "mac", collection.Rules[0].Condition},
		{"Disjunctive", true, collection.Disjunctive},
	}

	for _, c := range cases {
		if c.expected != c.actual {
			t.Errorf("SmartCollection.%v returned %v, expected %v", c.field, c.actual, c.expected)
		}
	}
}

func TestSmartCollectionList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"smart_collections": [{"id":1},{"id":2}]}`))

	collections, err := client.SmartCollection.List(nil)
	if err != nil {
		t.Errorf("SmartCollection.List returned error: %v", err)
	}

	expected := []SmartCollection{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(collections, expected) {
		t.Errorf("SmartCollection.List returned %+v, expected %+v", collections, expected)
	}
}

func TestSmartCollectionCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 5}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.SmartCollection.Count(nil)
	if err != nil {
		t.Errorf("SmartCollection.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("SmartCollection.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.SmartCollection.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("SmartCollection.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("SmartCollection.Count returned %d, expected %d", cnt, expected)
	}
}

func TestSmartCollectionGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"smart_collection": {"id":1}}`))

	collection, err := client.SmartCollection.Get(1, nil)
	if err != nil {
		t.Errorf("SmartCollection.Get returned error: %v", err)
	}

	expected := &SmartCollection{ID: 1}
	if !reflect.DeepEqual(collection, expected) {
		t.Errorf("SmartCollection.Get returned %+v, expected %+v", collection, expected)
	}
}

func TestSmartCollectionCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("smartcollection.json")))

	collection := SmartCollection{
		Title: "Macbooks",
	}

	returnedCollection, err := client.SmartCollection.Create(collection)
	if err != nil {
		t.Errorf("SmartCollection.Create returned error: %v", err)
	}

	smartCollectionTests(t, *returnedCollection)
}

func TestSmartCollectionUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("smartcollection.json")))

	collection := SmartCollection{
		ID:    1,
		Title: "Macbooks",
	}

	returnedCollection, err := client.SmartCollection.Update(collection)
	if err != nil {
		t.Errorf("SmartCollection.Update returned error: %v", err)
	}

	smartCollectionTests(t, *returnedCollection)
}

func TestSmartCollectionDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/smart_collections/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.SmartCollection.Delete(1)
	if err != nil {
		t.Errorf("SmartCollection.Delete returned error: %v", err)
	}
}

func TestSmartCollectionListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.SmartCollection.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("SmartCollection.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("SmartCollection.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestSmartCollectionCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.SmartCollection.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("SmartCollection.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("SmartCollection.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.SmartCollection.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("SmartCollection.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("SmartCollection.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestSmartCollectionGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.SmartCollection.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("SmartCollection.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("SmartCollection.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestSmartCollectionCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.SmartCollection.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("SmartCollection.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestSmartCollectionUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.SmartCollection.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("SmartCollection.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestSmartCollectionDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.SmartCollection.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("SmartCollection.DeleteMetafield() returned error: %v", err)
	}
}
