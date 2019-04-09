package goshopify

import (
	"reflect"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestScriptTagList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags.json",
		httpmock.NewStringResponder(200, `{"script_tags": [{"id": 1},{"id": 2}]}`))

	scriptTags, err := client.ScriptTag.List(nil)
	if err != nil {
		t.Errorf("ScriptTag.List returned error: %v", err)
	}

	expected := []ScriptTag{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(scriptTags, expected) {
		t.Errorf("ScriptTag.List returned %+v, expected %+v", scriptTags, expected)
	}
}

func TestScriptTagCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags/count.json",
		httpmock.NewStringResponder(200, `{"count": 3}`))

	cnt, err := client.ScriptTag.Count(nil)
	if err != nil {
		t.Errorf("ScriptTag.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("ScriptTag.Count returned %d, expected %d", cnt, expected)
	}
}

func TestScriptTagGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags/1.json",
		httpmock.NewStringResponder(200, `{"script_tag": {"id": 1}}`))

	scriptTag, err := client.ScriptTag.Get(1, nil)
	if err != nil {
		t.Errorf("ScriptTag.Get returned error: %v", err)
	}

	expected := &ScriptTag{ID: 1}
	if !reflect.DeepEqual(scriptTag, expected) {
		t.Errorf("ScriptTag.Get returned %+v, expected %+v", scriptTag, expected)
	}
}

func scriptTagTests(t *testing.T, tag ScriptTag) {
	expected := int64(870402688)
	if tag.ID != expected {
		t.Errorf("tag.ID is %+v, expected %+v", tag.ID, expected)
	}
}

func TestScriptTagCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags.json",
		httpmock.NewBytesResponder(200, loadFixture("script_tags.json")))

	tag0 := ScriptTag{
		Src:          "https://djavaskripped.org/fancy.js",
		Event:        "onload",
		DisplayScope: "all",
	}

	returnedTag, err := client.ScriptTag.Create(tag0)
	if err != nil {
		t.Errorf("ScriptTag.Create returned error: %v", err)
	}
	scriptTagTests(t, *returnedTag)
}

func TestScriptTagUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags/1.json",
		httpmock.NewBytesResponder(200, loadFixture("script_tags.json")))

	tag := ScriptTag{
		ID:  1,
		Src: "https://djavaskripped.org/fancy.js",
	}

	returnedTag, err := client.ScriptTag.Update(tag)
	if err != nil {
		t.Errorf("ScriptTag.Update returned error: %v", err)
	}
	scriptTagTests(t, *returnedTag)
}

func TestScriptTagDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", "https://fooshop.myshopify.com/admin/api/2019-04/script_tags/1.json",
		httpmock.NewStringResponder(200, "{}"))

	if err := client.ScriptTag.Delete(1); err != nil {
		t.Errorf("ScriptTag.Delete returned error: %v", err)
	}
}
