package goshopify

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestFlexDecimal_UnmarshalJSON_EmptyString(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`""`), &f)
	if err != nil {
		t.Fatalf("expected no error for empty string, got: %v", err)
	}
	if f.Decimal != nil {
		t.Fatalf("expected nil decimal for empty string, got: %v", f.Decimal)
	}
}

func TestFlexDecimal_UnmarshalJSON_Null(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`null`), &f)
	if err != nil {
		t.Fatalf("expected no error for null, got: %v", err)
	}
	if f.Decimal != nil {
		t.Fatalf("expected nil decimal for null, got: %v", f.Decimal)
	}
}

func TestFlexDecimal_UnmarshalJSON_ValidString(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`"12.50"`), &f)
	if err != nil {
		t.Fatalf("expected no error for valid string, got: %v", err)
	}
	if f.Decimal == nil {
		t.Fatal("expected non-nil decimal")
	}
	expected := decimal.NewFromFloat(12.50)
	if !f.Decimal.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, f.Decimal)
	}
}

func TestFlexDecimal_UnmarshalJSON_ValidNumber(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`0`), &f)
	if err != nil {
		t.Fatalf("expected no error for number, got: %v", err)
	}
	if f.Decimal == nil {
		t.Fatal("expected non-nil decimal")
	}
	if !f.Decimal.Equal(decimal.Zero) {
		t.Fatalf("expected 0, got %s", f.Decimal)
	}
}

func TestFlexDecimal_UnmarshalJSON_Invalid(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`"not-a-number"`), &f)
	if err == nil {
		t.Fatal("expected error for invalid decimal string, got nil")
	}
}

func TestFlexDecimal_MarshalJSON_Nil(t *testing.T) {
	f := FlexDecimal{Decimal: nil}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if string(data) != "null" {
		t.Fatalf("expected null, got %s", string(data))
	}
}

func TestFlexDecimal_MarshalJSON_Value(t *testing.T) {
	d := decimal.NewFromFloat(9.99)
	f := FlexDecimal{Decimal: &d}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if string(data) != `"9.99"` {
		t.Fatalf("expected \"9.99\", got %s", string(data))
	}
}

func TestFlexDecimal_UnmarshalJSON_SingleQuoteEmpty(t *testing.T) {
	var f FlexDecimal
	err := f.UnmarshalJSON([]byte(`''`))
	if err != nil {
		t.Fatalf("expected no error for single-quote empty, got: %v", err)
	}
	if f.Decimal != nil {
		t.Fatalf("expected nil decimal for single-quote empty, got: %v", f.Decimal)
	}
}

func TestFlexDecimal_UnmarshalJSON_InvalidValue(t *testing.T) {
	var f FlexDecimal
	err := json.Unmarshal([]byte(`"notanumber"`), &f)
	if err == nil {
		t.Fatal("expected error for invalid decimal value, got nil")
	}
}

func TestFlexDecimal_InStruct(t *testing.T) {
	// Simulates the actual Shopify API response with empty flat_modifier
	jsonData := `{
		"id": 815531950300,
		"carrier_service_id": 83880804572,
		"shipping_zone_id": 392866726108,
		"flat_modifier": "",
		"percent_modifier": 0,
		"service_filter": {"*": "+"}
	}`

	var provider CarrierShippingRateProvider
	err := json.Unmarshal([]byte(jsonData), &provider)
	if err != nil {
		t.Fatalf("expected no error unmarshalling carrier provider with empty flat_modifier, got: %v", err)
	}
	if provider.FlatModifier.Decimal != nil {
		t.Fatalf("expected nil flat_modifier, got: %v", provider.FlatModifier.Decimal)
	}
	if provider.PercentModifier.Decimal == nil {
		t.Fatal("expected non-nil percent_modifier")
	}
}
