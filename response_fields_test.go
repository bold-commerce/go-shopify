package goshopify

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

// The fixture contains previously missing fields from Shopify REST 2026-10
// response examples: https://shopify.dev/docs/api/admin-rest/2026-10/resources
// Keep the literal question marks in the two draft-order response keys.
func TestDocumentedResponseFields(t *testing.T) {
	models := map[string]reflect.Type{
		"AbandonedCheckout":              reflect.TypeOf(AbandonedCheckout{}),
		"AbandonedCheckoutShippingLine":  reflect.TypeOf(AbandonedCheckoutShippingLine{}),
		"ApplicationCharge":              reflect.TypeOf(ApplicationCharge{}),
		"Article":                        reflect.TypeOf(Article{}),
		"Asset":                          reflect.TypeOf(Asset{}),
		"AssignedFulfillmentOrder":       reflect.TypeOf(AssignedFulfillmentOrder{}),
		"Collection":                     reflect.TypeOf(Collection{}),
		"CustomCollection":               reflect.TypeOf(CustomCollection{}),
		"DiscountAllocations":            reflect.TypeOf(DiscountAllocations{}),
		"DraftOrder":                     reflect.TypeOf(DraftOrder{}),
		"Fulfillment":                    reflect.TypeOf(Fulfillment{}),
		"FulfillmentEvent":               reflect.TypeOf(FulfillmentEvent{}),
		"FulfillmentOrder":               reflect.TypeOf(FulfillmentOrder{}),
		"FulfillmentOrderDeliveryMethod": reflect.TypeOf(FulfillmentOrderDeliveryMethod{}),
		"GiftCard":                       reflect.TypeOf(GiftCard{}),
		"LineItem":                       reflect.TypeOf(LineItem{}),
		"Location":                       reflect.TypeOf(Location{}),
		"Order":                          reflect.TypeOf(Order{}),
		"Page":                           reflect.TypeOf(Page{}),
		"PaymentDetails":                 reflect.TypeOf(PaymentDetails{}),
		"PaymentSchedule":                reflect.TypeOf(PaymentSchedule{}),
		"PaymentTerms":                   reflect.TypeOf(PaymentTerms{}),
		"PriceRule":                      reflect.TypeOf(PriceRule{}),
		"PriceRuleDiscountCode":          reflect.TypeOf(PriceRuleDiscountCode{}),
		"ProductListing":                 reflect.TypeOf(ProductListing{}),
		"Refund":                         reflect.TypeOf(Refund{}),
		"ScriptTag":                      reflect.TypeOf(ScriptTag{}),
		"ShippingLines":                  reflect.TypeOf(ShippingLines{}),
		"Shop":                           reflect.TypeOf(Shop{}),
		"SmartCollection":                reflect.TypeOf(SmartCollection{}),
		"TaxLine":                        reflect.TypeOf(TaxLine{}),
		"Variant":                        reflect.TypeOf(Variant{}),
		"Webhook":                        reflect.TypeOf(Webhook{}),
	}
	data, err := os.ReadFile("fixtures/response_fields.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for name, input := range fixtures {
		t.Run(name, func(t *testing.T) {
			typ, ok := models[name]
			if !ok {
				t.Fatalf("no model registered for %s", name)
			}
			value := reflect.New(typ).Interface()
			dec := json.NewDecoder(bytes.NewReader(input))
			dec.DisallowUnknownFields()
			if err := dec.Decode(value); err != nil {
				t.Fatal(err)
			}
			var expected map[string]interface{}
			if err := json.Unmarshal(input, &expected); err != nil {
				t.Fatal(err)
			}
			// Null and empty example values must also have a named model field.
			for key := range expected {
				found := false
				for i := 0; i < typ.NumField(); i++ {
					if strings.Split(typ.Field(i).Tag.Get("json"), ",")[0] == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing field %q", key)
				}
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]interface{}
			if err := json.Unmarshal(encoded, &actual); err != nil {
				t.Fatal(err)
			}
			assertDocumentedFields(t, name, expected, actual)
		})
	}
}

// Compare populated response fields; omitempty may remove null/empty values,
// and decimal.Decimal may normalize the scale of a money string.
func assertDocumentedFields(t *testing.T, path string, want, got interface{}) {
	t.Helper()
	if want == nil {
		return
	}
	switch expected := want.(type) {
	case map[string]interface{}:
		actual, _ := got.(map[string]interface{})
		for key, value := range expected {
			assertDocumentedFields(t, path+"."+key, value, actual[key])
		}
	case []interface{}:
		if len(expected) == 0 {
			return
		}
		actual, ok := got.([]interface{})
		if !ok || len(actual) != len(expected) {
			t.Errorf("%s: got %v, want %v", path, got, want)
			return
		}
		for i, value := range expected {
			assertDocumentedFields(t, path+"[]", value, actual[i])
		}
	default:
		if reflect.DeepEqual(want, got) {
			return
		}
		if a, ok := want.(string); ok {
			if b, ok := got.(string); ok {
				x, e1 := decimal.NewFromString(a)
				y, e2 := decimal.NewFromString(b)
				if e1 == nil && e2 == nil && x.Equal(y) {
					return
				}
			}
		}
		t.Errorf("%s: got %#v, want %#v", path, got, want)
	}
}

func TestPopulatedResponseFields(t *testing.T) {
	var item LineItem
	if err := json.Unmarshal([]byte(`{"attributed_staffs":[{"id":"gid://shopify/StaffMember/2108187","quantity":1}],"compare_at_price":"123456789012345678.90","duties":[{"id":9007199254740993}]}`), &item); err != nil {
		t.Fatal(err)
	}
	if len(item.AttributedStaffs) != 1 || item.AttributedStaffs[0].Id != "gid://shopify/StaffMember/2108187" || item.AttributedStaffs[0].Quantity != 1 {
		t.Fatalf("staff attribution lost: %+v", item.AttributedStaffs)
	}
	if item.CompareAtPrice == nil || item.CompareAtPrice.String() != "123456789012345678.9" {
		t.Fatalf("price lost precision: %v", item.CompareAtPrice)
	}
	// Opaque fields preserve JSON without rounding IDs through float64.
	if len(item.Duties) != 1 || string(item.Duties[0]) != `{"id":9007199254740993}` {
		t.Fatalf("duty JSON changed: %s", item.Duties)
	}
	var details PaymentDetails
	if err := json.Unmarshal([]byte(`{"credit_card_expiration_month":10,"credit_card_expiration_year":2028,"credit_card_wallet":"shopify_pay","credit_card_name":"John Smith","buyer_action_info":{"multibanco":{"entity":"12345","reference":"123456789"}}}`), &details); err != nil {
		t.Fatal(err)
	}
	if details.CreditCardExpirationMonth == nil || *details.CreditCardExpirationMonth != 10 || details.CreditCardExpirationYear == nil || *details.CreditCardExpirationYear != 2028 || details.CreditCardWallet != "shopify_pay" || details.CreditCardName != "John Smith" || len(details.BuyerActionInfo) == 0 {
		t.Fatalf("payment details lost: %+v", details)
	}
	var rule PriceRule
	if err := json.Unmarshal([]byte(`{"prerequisite_to_entitlement_purchase":{"prerequisite_amount":"80.00"},"customer_segment_prerequisite_ids":[9007199254740993]}`), &rule); err != nil {
		t.Fatal(err)
	}
	if rule.PrerequisiteToEntitlementPurchase == nil || rule.PrerequisiteToEntitlementPurchase.PrerequisiteAmount == nil || rule.PrerequisiteToEntitlementPurchase.PrerequisiteAmount.String() != "80" || len(rule.CustomerSegmentPrerequisiteIds) != 1 || rule.CustomerSegmentPrerequisiteIds[0] != 9007199254740993 {
		t.Fatalf("price rule fields lost: %+v", rule)
	}
}
