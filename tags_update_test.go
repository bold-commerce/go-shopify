package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestUpdateTags(t *testing.T) {
	setup()
	defer teardown()
	for _, resource := range []string{"product", "order", "customer"} {
		for _, mode := range []string{"omit", "clear", "existing"} {
			t.Run(resource+"/"+mode, func(t *testing.T) {
				endpoint := fmt.Sprintf("https://fooshop.myshopify.com/%s/%ss/7.json", client.pathPrefix, resource)
				httpmock.RegisterResponder("PUT", endpoint, func(r *http.Request) (*http.Response, error) {
					var body map[string]map[string]json.RawMessage
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						return nil, err
					}
					tags := string(body[resource]["tags"])
					expected := map[string]string{"omit": "", "clear": `""`, "existing": `"vip"`}[mode]
					if tags != expected {
						t.Errorf("tags = %q, want %q", tags, expected)
					}
					return httpmock.NewStringResponse(200, "{}"), nil
				})
				var tags string
				if mode == "existing" {
					tags = "vip"
				}
				p := Product{Id: 7, Tags: tags}
				o := Order{Id: 7, Tags: tags}
				c := Customer{Id: 7, Tags: tags}
				if mode == "clear" {
					p.SetTags("")
					o.SetTags("")
					c.SetTags("")
				}
				var err error
				switch resource {
				case "product":
					_, err = client.Product.Update(context.Background(), p)
				case "order":
					_, err = client.Order.Update(context.Background(), o)
				case "customer":
					_, err = client.Customer.Update(context.Background(), c)
				}
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
