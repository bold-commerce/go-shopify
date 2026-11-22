package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestVariantInventoryOnUpdate(t *testing.T) {
	setup()
	defer teardown()
	variants := []Variant{{Id: 7, InventoryQuantity: 5, OldInventoryQuantity: 9}}
	for _, resource := range []string{"variant", "product"} {
		url := fmt.Sprintf("https://fooshop.myshopify.com/%s/%ss/7.json", client.pathPrefix, resource)
		httpmock.RegisterResponder("PUT", url, func(r *http.Request) (*http.Response, error) {
			var wire struct {
				Variant map[string]interface{}
				Product struct{ Variants []map[string]interface{} }
			}
			if err := json.NewDecoder(r.Body).Decode(&wire); err != nil {
				return nil, err
			}
			fields := wire.Variant
			if resource == "product" {
				fields = wire.Product.Variants[0]
			}
			if fields["id"] != float64(7) || fields["inventory_quantity"] != nil || fields["old_inventory_quantity"] != nil {
				t.Errorf("sent fields %v", fields)
			}
			return httpmock.NewStringResponse(200, "{}"), nil
		})
		var err error
		if resource == "variant" {
			_, err = client.Variant.Update(context.Background(), variants[0])
		} else {
			_, err = client.Product.Update(context.Background(), Product{Id: 7, Variants: variants})
		}
		if err != nil || variants[0].InventoryQuantity != 5 {
			t.Fatal(variants, err)
		}
	}
}
