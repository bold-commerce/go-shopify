package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func MetafieldTests(t *testing.T, metafield Metafield) {
	// Check that ID is assigned to the returned metafield
	expectedInt := int64(721389482)
	if metafield.ID != expectedInt {
		t.Errorf("Metafield.ID returned %+v, expected %+v", metafield.ID, expectedInt)
	}
}

func TestMetafieldList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Metafield.List(nil)
	if err != nil {
		t.Errorf("Metafield.List returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Metafield.List returned %+v, expected %+v", metafields, expected)
	}
}

func TestMetafieldCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Metafield.Count(nil)
	if err != nil {
		t.Errorf("Metafield.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Metafield.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Metafield.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Metafield.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Metafield.Count returned %d, expected %d", cnt, expected)
	}
}

func TestMetafieldGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":1}}`))

	metafield, err := client.Metafield.Get(1, nil)
	if err != nil {
		t.Errorf("Metafield.Get returned error: %v", err)
	}

	expected := &Metafield{ID: 1}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Metafield.Get returned %+v, expected %+v", metafield, expected)
	}
}

func TestMetafieldCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Namespace: "inventory",
		Key:       "warehouse",
		Value:     "25",
		ValueType: "integer",
	}

	returnedMetafield, err := client.Metafield.Create(metafield)
	if err != nil {
		t.Errorf("Metafield.Create returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestMetafieldUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        1,
		Value:     "something new",
		ValueType: "string",
	}

	returnedMetafield, err := client.Metafield.Update(metafield)
	if err != nil {
		t.Errorf("Metafield.Update returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestMetafieldDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Metafield.Delete(1)
	if err != nil {
		t.Errorf("Metafield.Delete returned error: %v", err)
	}
}
