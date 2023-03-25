package goshopify

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func collectTests(t *testing.T, collect Collect) {

	// Test a few fields
	cases := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", int64(18091352323), collect.ID},
		{"CollectionID", int64(241600835), collect.CollectionID},
		{"ProductID", int64(6654094787), collect.ProductID},
		{"Featured", false, collect.Featured},
		{"SortValue", "0000000002", collect.SortValue},
	}

	for _, c := range cases {
		if c.expected != c.actual {
			t.Errorf("Collect.%v returned %v, expected %v", c.field, c.actual, c.expected)
		}
	}
}

func TestCollectList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collects.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"collects": [{"id":1},{"id":2}]}`))

	collects, err := client.Collect.List(nil)
	if err != nil {
		t.Errorf("Collect.List returned error: %v", err)
	}

	expected := []Collect{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(collects, expected) {
		t.Errorf("Collect.List returned %+v, expected %+v", collects, expected)
	}
}

func TestCollectCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collects/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 5}`))

	params := map[string]string{"since_id": "123"}
	httpmock.RegisterResponderWithQuery("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/collects/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Collect.Count(nil)
	if err != nil {
		t.Errorf("Collect.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("Collect.Count returned %d, expected %d", cnt, expected)
	}

	cnt, err = client.Collect.Count(ListOptions{SinceID: 123})
	if err != nil {
		t.Errorf("Collect.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Collect.Count returned %d, expected %d", cnt, expected)
	}
}

func TestCollectGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collects/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"collect": {"id":1}}`))

	product, err := client.Collect.Get(1, nil)
	if err != nil {
		t.Errorf("Collect.Get returned error: %v", err)
	}

	expected := &Collect{ID: 1}
	if !reflect.DeepEqual(product, expected) {
		t.Errorf("Collect.Get returned %+v, expected %+v", product, expected)
	}
}

func TestCollectCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/collects.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("collect.json")))

	collect := Collect{
		CollectionID: 241600835,
		ProductID:    6654094787,
	}

	returnedCollect, err := client.Collect.Create(collect)
	if err != nil {
		t.Errorf("Collect.Create returned error: %v", err)
	}

	collectTests(t, *returnedCollect)
}

func TestCollectDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/collects/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Collect.Delete(1)
	if err != nil {
		t.Errorf("Collect.Delete returned error: %v", err)
	}
}
