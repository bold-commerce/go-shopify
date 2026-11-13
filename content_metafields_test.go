package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestContentMetafieldArrays(t *testing.T) {
	setup()
	defer teardown()
	fields := []Metafield{{Key: "summary", Value: "hello", Type: MetafieldTypeSingleLineTextField, Namespace: "custom"}}
	for _, tc := range []struct {
		name, method, path string
		call               func() error
	}{
		{"article", "POST", "blogs/1/articles.json", func() error {
			_, err := client.Article.Create(context.Background(), 1, Article{Title: "Title", Metafields: fields})
			return err
		}},
		{"article", "PUT", "blogs/1/articles/2.json", func() error {
			_, err := client.Article.Update(context.Background(), 1, 2, Article{Id: 2, Metafields: fields})
			return err
		}},
		{"blog", "POST", "blogs.json", func() error {
			_, err := client.Blog.Create(context.Background(), Blog{Title: "Title", Metafields: fields})
			return err
		}},
		{"blog", "PUT", "blogs/2.json", func() error {
			_, err := client.Blog.Update(context.Background(), Blog{Id: 2, Metafields: fields})
			return err
		}},
	} {
		t.Run(tc.name+tc.method, func(t *testing.T) {
			httpmock.RegisterResponder(tc.method, fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, tc.path), func(req *http.Request) (*http.Response, error) {
				var got map[string]map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				want := []interface{}{map[string]interface{}{"key": "summary", "value": "hello", "type": "single_line_text_field", "namespace": "custom"}}
				if !reflect.DeepEqual(got[tc.name]["metafields"], want) {
					t.Errorf("metafields = %#v", got[tc.name]["metafields"])
				}
				if _, ok := got[tc.name]["metafield"]; ok {
					t.Error("unexpected singular metafield key")
				}
				return httpmock.NewStringResponse(200, `{}`), nil
			})
			if err := tc.call(); err != nil {
				t.Fatal(err)
			}
		})
	}
	for _, value := range []interface{}{Article{}, Blog{}} {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		if _, ok := fields["metafields"]; ok {
			t.Errorf("unset metafields must be omitted: %s", data)
		}
	}
}
