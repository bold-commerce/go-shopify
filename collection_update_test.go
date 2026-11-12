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

func TestCollectionPartialUpdates(t *testing.T) {
	setup()
	defer teardown()

	unpublished := false
	empty := ""
	for _, tc := range []struct {
		name   string
		body   string
		update func() error
	}{
		{"custom_collection", `{"id":1,"title":"New title"}`, func() error {
			_, err := client.CustomCollection.Update(context.Background(), CustomCollection{Id: 1, Title: "New title"})
			return err
		}},
		{"smart_collection", `{"id":1,"title":"New title"}`, func() error {
			_, err := client.SmartCollection.Update(context.Background(), SmartCollection{Id: 1, Title: "New title"})
			return err
		}},
		{"custom_collection", `{"id":1,"body_html":"","template_suffix":""}`, func() error {
			_, err := client.CustomCollection.Update(context.Background(), CustomCollection{Id: 1, BodyHTML: &empty, TemplateSuffix: &empty})
			return err
		}},
		{"smart_collection", `{"id":1,"body_html":"","template_suffix":""}`, func() error {
			_, err := client.SmartCollection.Update(context.Background(), SmartCollection{Id: 1, BodyHTML: &empty, TemplateSuffix: &empty})
			return err
		}},
		{"custom_collection", `{"id":1,"published":false}`, func() error {
			_, err := client.CustomCollection.Update(context.Background(), CustomCollection{Id: 1, Published: &unpublished})
			return err
		}},
		{"smart_collection", `{"id":1,"published":false,"disjunctive":false}`, func() error {
			_, err := client.SmartCollection.Update(context.Background(), SmartCollection{Id: 1, Published: &unpublished, Disjunctive: &unpublished})
			return err
		}},
	} {
		t.Run(tc.name+tc.body, func(t *testing.T) {
			url := fmt.Sprintf("https://fooshop.myshopify.com/%s/%ss/1.json", client.pathPrefix, tc.name)
			httpmock.RegisterResponder("PUT", url, func(req *http.Request) (*http.Response, error) {
				var got, want map[string]interface{}
				if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(`{"`+tc.name+`":`+tc.body+`}`), &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("request = %#v, want %#v", got, want)
				}
				return httpmock.NewStringResponse(200, `{}`), nil
			})
			if err := tc.update(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
