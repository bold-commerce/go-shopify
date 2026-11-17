package goshopify

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Optional settings must distinguish an omitted update from disabling a setting.
func TestOptionalBooleanSettings(t *testing.T) {
	for _, model := range []interface{}{Variant{}, Article{}, Customer{}, CustomerAddress{}, PriceRule{}, CarrierService{}, OrderRisk{}, FulfillmentServiceData{}, LineItem{}, Order{}, OrderCancelOptions{}, DraftOrder{}} {
		typ := reflect.TypeOf(model)
		t.Run(typ.Name(), func(t *testing.T) {
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if field.Type != reflect.TypeOf((*bool)(nil)) || !strings.Contains(field.Tag.Get("json"), ",omitempty") {
					continue
				}
				key := strings.Split(field.Tag.Get("json"), ",")[0]
				for _, value := range []*bool{nil, BoolPtr(false), BoolPtr(true)} {
					input := reflect.New(typ)
					input.Elem().Field(i).Set(reflect.ValueOf(value))
					b, err := json.Marshal(input.Interface())
					if err != nil {
						t.Fatal(err)
					}
					var wire map[string]json.RawMessage
					if err := json.Unmarshal(b, &wire); err != nil {
						t.Fatal(err)
					}
					raw, exists := wire[key]
					if value == nil {
						if exists {
							t.Fatalf("%s: omitted setting was sent", key)
						}
						continue
					}
					var got bool
					if err := json.Unmarshal(raw, &got); err != nil {
						t.Fatal(err)
					}
					if !exists || got != *value {
						t.Fatalf("%s=%s, want %v", key, raw, *value)
					}
					decoded := reflect.New(typ)
					if err := json.Unmarshal(b, decoded.Interface()); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(decoded.Elem().Field(i).Interface(), value) {
						t.Fatalf("%s lost response value", key)
					}
				}
			}
		})
	}
}

func TestDisableBooleanSettingsOverHTTP(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{name: "variant tax", method: "PUT", path: "variants/7.json", request: `{"variant":{"id":7,"taxable":false}}`, response: `{"variant":{"id":7,"taxable":false}}`, want: &Variant{Id: 7, Taxable: BoolPtr(false)}, call: func(ctx context.Context) (interface{}, error) {
			return client.Variant.Update(ctx, Variant{Id: 7, Taxable: BoolPtr(false)})
		}},
		{name: "customer tax", method: "PUT", path: "customers/7.json", request: `{"customer":{"id":7,"tax_exempt":false,"email_marketing_consent":null,"sms_marketing_consent":null}}`, response: `{"customer":{"id":7,"tax_exempt":false}}`, want: &Customer{Id: 7, TaxExempt: BoolPtr(false)}, call: func(ctx context.Context) (interface{}, error) {
			return client.Customer.Update(ctx, Customer{Id: 7, TaxExempt: BoolPtr(false)})
		}},
		{name: "article visibility", method: "PUT", path: "blogs/11/articles/7.json", request: `{"article":{"id":7,"published":false,"updated_at":"0001-01-01T00:00:00Z"}}`, response: `{"article":{"id":7,"published":false}}`, want: &Article{Id: 7, Published: BoolPtr(false)}, call: func(ctx context.Context) (interface{}, error) {
			return client.Article.Update(ctx, 11, 7, Article{Id: 7, Published: BoolPtr(false)})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
