package goshopify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestScriptTagList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/script_tags.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"script_tags": [{"id": 1},{"id": 2}]}`))

	scriptTags, err := client.ScriptTag.List(context.Background(), nil)
	if err != nil {
		t.Errorf("ScriptTag.List returned error: %v", err)
	}

	expected := []ScriptTag{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(scriptTags, expected) {
		t.Errorf("ScriptTag.List returned %+v, expected %+v", scriptTags, expected)
	}
}

func TestScriptTagCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/script_tags/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	cnt, err := client.ScriptTag.Count(context.Background(), nil)
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/script_tags/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"script_tag": {"id": 1}}`))

	scriptTag, err := client.ScriptTag.Get(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("ScriptTag.Get returned error: %v", err)
	}

	expected := &ScriptTag{Id: 1}
	if !reflect.DeepEqual(scriptTag, expected) {
		t.Errorf("ScriptTag.Get returned %+v, expected %+v", scriptTag, expected)
	}
}

func TestScriptTagWritesUnsupported(t *testing.T) {
	setup()
	defer teardown()

	for _, write := range []func(context.Context, ScriptTag) (*ScriptTag, error){client.ScriptTag.Create, client.ScriptTag.Update} {
		tag, err := write(context.Background(), ScriptTag{Id: 1})
		if tag != nil || err != ErrScriptTagWriteUnsupported {
			t.Fatalf("got %v, %v", tag, err)
		}
	}
	if calls := httpmock.GetTotalCallCount(); calls != 0 {
		t.Fatalf("sent %d requests", calls)
	}
}

func TestScriptTagDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/script_tags/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.ScriptTag.Delete(context.Background(), 1); err != nil {
		t.Errorf("ScriptTag.Delete returned error: %v", err)
	}
}
