package goshopify

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestDiscountCodeList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		"https://fooshop.myshopify.com/admin/price_rules/507328175/discount_codes.json",
		httpmock.NewStringResponder(
			200,
			`{"discount_codes":[{"id":507328175,"price_rule_id":507328175,"code":"SUMMERSALE10OFF","usage_count":0,"created_at":"2018-07-05T12:41:00-04:00","updated_at":"2018-07-05T12:41:00-04:00"}]}`,
		),
	)

	codes, err := client.DiscountCode.List(507328175)
	if err != nil {
		t.Errorf("DiscountCode.List returned error: %v", err)
	}

	expected := []PriceRuleDiscountCode{{ID: 507328175}}
	if expected[0].ID != codes[0].ID {
		t.Errorf("DiscountCode.List returned %+v, expected %+v", codes, expected)
	}

}

func TestDiscountCodeGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		"https://fooshop.myshopify.com/admin/price_rules/507328175/discount_codes/507328175.json",
		httpmock.NewStringResponder(
			200,
			`{"discount_code":{"id":507328175,"price_rule_id":507328175,"code":"SUMMERSALE10OFF","usage_count":0,"created_at":"2018-07-05T12:41:00-04:00","updated_at":"2018-07-05T12:41:00-04:00"}}`,
		),
	)

	dc, err := client.DiscountCode.Get(507328175, 507328175)
	if err != nil {
		t.Errorf("DiscountCode.Get returned error: %v", err)
	}

	expected := &PriceRuleDiscountCode{ID: 507328175}

	if dc.ID != expected.ID {
		t.Errorf("DiscountCode.Get returned %+v, expected %+v", dc, expected)
	}

}

func TestDiscountCodeCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		"https://fooshop.myshopify.com/admin/price_rules/507328175/discount_codes.json",
		httpmock.NewBytesResponder(
			201,
			loadFixture("discount_code.json"),
		),
	)

	dc := PriceRuleDiscountCode{
		Code: "SUMMERSALE10OFF",
	}

	returnedDC, err := client.DiscountCode.Create(507328175, dc)
	if err != nil {
		t.Errorf("DiscountCode.Create returned error: %v", err)
	}

	expectedInt := int64(1054381139)
	if returnedDC.ID != expectedInt {
		t.Errorf("DiscountCode.ID returned %+v, expected %+v", returnedDC.ID, expectedInt)
	}

}

func TestDiscountCodeUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		"https://fooshop.myshopify.com/admin/price_rules/507328175/discount_codes/1054381139.json",
		httpmock.NewBytesResponder(
			200,
			loadFixture("discount_code.json"),
		),
	)

	dc := PriceRuleDiscountCode{
		ID:   int64(1054381139),
		Code: "SUMMERSALE10OFF",
	}

	returnedDC, err := client.DiscountCode.Update(507328175, dc)
	if err != nil {
		t.Errorf("DiscountCode.Update returned error: %v", err)
	}

	expectedInt := int64(1054381139)
	if returnedDC.ID != expectedInt {
		t.Errorf("DiscountCode.ID returned %+v, expected %+v", returnedDC.ID, expectedInt)
	}
}

func TestDiscountCodeDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", "https://fooshop.myshopify.com/admin/price_rules/507328175/discount_codes/507328175.json",
		httpmock.NewStringResponder(204, "{}"))

	err := client.DiscountCode.Delete(507328175, 507328175)
	if err != nil {
		t.Errorf("DiscountCode.Delete returned error: %v", err)
	}
}
