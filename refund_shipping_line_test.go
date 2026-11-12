package goshopify

import (
	"encoding/json"
	"testing"
)

func TestRefundShippingLine(t *testing.T) {
	data := []byte(`{"refund_shipping_lines":[{"id":1,"shipping_line_id":2,"shipping_line":{"id":2,"title":"Standard","price":"5.00"}}]}`)
	var refund Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		t.Fatal(err)
	}
	if len(refund.RefundShippingLines) != 1 {
		t.Fatalf("refund = %#v", refund)
	}
	line := refund.RefundShippingLines[0]
	if line.ShippingLine == nil || line.ShippingLine.Id != 2 || line.ShippingLine.Title != "Standard" || line.ShippingLine.Price.String() != "5" {
		t.Fatalf("shipping line = %#v", line.ShippingLine)
	}
	encoded, err := json.Marshal(line)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["shipping_line"]; !ok {
		t.Fatalf("missing shipping_line: %s", encoded)
	}
	if _, ok := fields["shipping_lines"]; ok {
		t.Fatalf("unexpected shipping_lines: %s", encoded)
	}
}
