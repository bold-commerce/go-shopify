package goshopify

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// These payloads exercise fields from the REST response examples. Checking the
// JSON after decoding also checks nested fields that encoding/json would silently drop.
func TestResponseModelFields(t *testing.T) {
	money := `{"shop_money":{"amount":"70.25","currency_code":"CAD"},"presentment_money":{"amount":"51.23","currency_code":"USD"}}`
	terms := `{"id":1,"amount":"70.25","currency":"CAD","payment_terms_name":"NET_30","payment_terms_type":"NET","due_in_days":30,"payment_schedules":[{"id":2,"amount":"70.25","currency":"CAD","issued_at":"2026-01-01T13:00:00-05:00","due_at":"2026-01-31T13:00:00-05:00","completed_at":"2026-01-30T13:00:00-05:00","expected_payment_method":"shopify_payments"}]}`
	order := fmt.Sprintf(`{"company":{"id":182772,"location_id":753095},"confirmation_number":"TLQDA3NN9","merchant_business_entity_id":"MTE","merchant_of_record_app_id":42,"po_number":"PO-123","payment_terms":%s,"estimated_taxes":false,"duties_included":false}`, terms)
	var orderFields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(order), &orderFields); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"current_total_additional_fees_set", "current_total_duties_set", "current_total_price_set", "current_subtotal_price_set", "original_total_additional_fees_set", "original_total_duties_set", "subtotal_price_set", "total_cash_rounding_payment_adjustment_set", "total_cash_rounding_refund_adjustment_set"} {
		orderFields[field] = json.RawMessage(money)
	}
	orderJSON, err := json.Marshal(orderFields)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, body string
		target     interface{}
	}{
		{"payout", `{"summary":{"adjustments_fee_amount":"0.12","adjustments_gross_amount":"2.13","charges_fee_amount":"1.32","charges_gross_amount":"45.52","refunds_fee_amount":"-0.23","refunds_gross_amount":"-3.54","reserved_funds_fee_amount":"0","reserved_funds_gross_amount":"0","retried_payouts_fee_amount":"0","retried_payouts_gross_amount":"0"}}`, &Payout{}},
		{"inventory item", `{"requires_shipping":false}`, &InventoryItem{}},
		{"customer", `{"tax_exemptions":["CA_STATUS_CARD_EXEMPTION","CA_BC_RESELLER_EXEMPTION"]}`, &Customer{}},
		{"customer empty exemptions", `{"tax_exemptions":[]}`, &Customer{}},
		{"order", string(orderJSON), &Order{}},
		{"draft order", `{"source_name":"web","admin_graphql_api_id":"gid://shopify/DraftOrder/1","payment_terms":` + terms + `}`, &DraftOrder{}},
		{"transaction", `{
   "processed_at":"2026-01-01T13:00:00-05:00",
   "authorization_expires_at":"2026-01-08T13:00:00-05:00",
   "extended_authorization_attributes":{"standard_authorization_expires_at":"2026-01-08T13:00:00-05:00","extended_authorization_expires_at":"2026-01-30T13:00:00-05:00"},
   "payments_refund_attributes":{"status":"success","acquirer_reference_number":"123456789012345678901234"},
   "amount_rounding":"0.02","manual_payment_gateway":false,"admin_graphql_api_id":"gid://shopify/OrderTransaction/1",
   "total_unsettled_set":{"presentment_money":{"amount":"171.8","currency":"USD"},"shop_money":{"amount":"171.8","currency":"USD"}},
   "currency_exchange_adjustment":{"id":1,"adjustment":"-0.01","original_amount":"-53.62","final_amount":"-53.63","currency":"CAD"}
  }`, &Transaction{}},
		{"refund", `{"processed_at":"2026-01-01T13:00:00-05:00","admin_graphql_api_id":"gid://shopify/Refund/1","duties":[{"duty_id":1,"amount_set":` + money + `}],"refund_line_items":[{"restock_type":"return","location_id":905684977}]}`, &Refund{}},
		{"balance transaction", `{"adjustment_reason":"channel_credit","adjustment_order_transactions":[{"id":12,"amount":"102.53","fee":"2.53","net":"100.00","order":{"id":1,"name":"#1001"}}]}`, &PaymentsTransactions{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tc.body), tc.target); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(tc.target)
			if err != nil {
				t.Fatal(err)
			}
			var got, want interface{}
			if err = json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal([]byte(tc.body), &want); err != nil {
				t.Fatal(err)
			}
			assertResponseFields(t, got, want, tc.name)
		})
	}
}

// Assert that every supplied response field survives decoding, without requiring
// unrelated fields on the resource to be omitted when it is encoded again.
func assertResponseFields(t *testing.T, got, want interface{}, path string) {
	t.Helper()
	switch want := want.(type) {
	case map[string]interface{}:
		actual, ok := got.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: expected object, got %#v", path, got)
		}
		for k, v := range want {
			value, ok := actual[k]
			if !ok {
				t.Errorf("%s.%s: response field was dropped", path, k)
				continue
			}
			assertResponseFields(t, value, v, path+"."+k)
		}
	case []interface{}:
		actual, ok := got.([]interface{})
		if !ok || len(actual) != len(want) {
			t.Fatalf("%s: got %#v, want %#v", path, got, want)
		}
		for i, v := range want {
			assertResponseFields(t, actual[i], v, fmt.Sprintf("%s[%d]", path, i))
		}
	default:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %#v, want %#v", path, got, want)
		}
	}
}

func TestResponseModelNullFields(t *testing.T) {
	for _, tc := range []struct {
		body   string
		target interface{}
	}{
		{`{"summary":null}`, &Payout{}},
		{`{"requires_shipping":null}`, &InventoryItem{}},
		{`{"tax_exemptions":null}`, &Customer{}},
		{`{"company":null,"payment_terms":null,"current_total_price_set":null,"current_total_duties_set":null}`, &Order{}},
		{`{"payment_terms":null}`, &DraftOrder{}},
		{`{"processed_at":null,"authorization_expires_at":null,"extended_authorization_attributes":null,"payments_refund_attributes":null,"total_unsettled_set":null,"currency_exchange_adjustment":null}`, &Transaction{}},
		{`{"processed_at":null,"duties":null}`, &Refund{}},
		{`{"adjustment_order_transactions":null}`, &PaymentsTransactions{}},
		{`{"completed_at":null,"due_at":null,"issued_at":null}`, &PaymentSchedule{}},
	} {
		if err := json.Unmarshal([]byte(tc.body), tc.target); err != nil {
			t.Fatal(err)
		}
		if !reflect.ValueOf(tc.target).Elem().IsZero() {
			t.Errorf("null response populated %T: %+v", tc.target, tc.target)
		}
	}
}

func TestPaymentTermsNumericAmounts(t *testing.T) {
	var terms PaymentTerms
	if err := json.Unmarshal([]byte(`{"amount":70.25,"due_in_days":0,"payment_schedules":[{"amount":70.25,"completed_at":null}]}`), &terms); err != nil {
		t.Fatal(err)
	}
	if terms.Amount == nil || terms.Amount.String() != "70.25" || len(terms.PaymentSchedules) != 1 || terms.PaymentSchedules[0].Amount == nil || terms.PaymentSchedules[0].Amount.String() != "70.25" || terms.DueInDays == nil || *terms.DueInDays != 0 {
		t.Fatalf("unexpected payment terms: %+v", terms)
	}
}
