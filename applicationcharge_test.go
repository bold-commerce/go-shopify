package goshopify

import (
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/jarcoal/httpmock.v1"
)

// applicationChargeTests tests if the fields are properly parsed.
func applicationChargeTests(t *testing.T, charge ApplicationCharge) {
	var nilTest *bool
	cases := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", int64(1017262355), charge.ID},
		{"Name", "Super Duper Expensive action", charge.Name},
		{"APIClientID", 755357713, charge.APIClientID},
		{"Price", decimal.NewFromFloat(100.00).String(), charge.Price.String()},
		{"Status", "pending", charge.Status},
		{"ReturnURL", "http://super-duper.shopifyapps.com/", charge.ReturnURL},
		{"Test", nilTest, charge.Test},
		{"CreatedAt", "2018-07-05T13:11:28-04:00", charge.CreatedAt.Format(time.RFC3339)},
		{"UpdatedAt", "2018-07-05T13:11:28-04:00", charge.UpdatedAt.Format(time.RFC3339)},
		{
			"DecoratedReturnURL",
			"http://super-duper.shopifyapps.com/?charge_id=1017262355",
			charge.DecoratedReturnURL,
		},
		{
			"ConfirmationURL",
			"https://apple.myshopify.com/admin/charges/1017262355/confirm_application_charge?sign" +
				"ature=BAhpBBMxojw%3D--1139a82a3433b1a6771786e03f02300440e11883",
			charge.ConfirmationURL,
		},
	}

	for _, c := range cases {
		if c.expected != c.actual {
			t.Errorf("ApplicationCharge.%s returned %v, expected %v", c.field, c.actual, c.expected)
		}
	}
}

func TestApplicationChargeServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		"https://fooshop.myshopify.com/admin/application_charges.json",
		httpmock.NewBytesResponder(200, loadFixture("applicationcharge.json")),
	)

	p := decimal.NewFromFloat(100.00)
	charge := ApplicationCharge{
		Name:      "Super Duper Expensive action",
		Price:     &p,
		ReturnURL: "http://super-duper.shopifyapps.com",
	}

	returnedCharge, err := client.ApplicationCharge.Create(charge)
	if err != nil {
		t.Errorf("ApplicationCharge.Create returned an error: %v", err)
	}

	applicationChargeTests(t, *returnedCharge)
}

func TestApplicationChargeServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		"https://fooshop.myshopify.com/admin/application_charges/1.json",
		httpmock.NewStringResponder(200, `{"application_charge": {"id":1}}`),
	)

	charge, err := client.ApplicationCharge.Get(1, nil)
	if err != nil {
		t.Errorf("ApplicationCharge.Get returned an error: %v", err)
	}

	expected := &ApplicationCharge{ID: 1}
	if !reflect.DeepEqual(charge, expected) {
		t.Errorf("ApplicationCharge.Get returned %+v, expected %+v", charge, expected)
	}
}

func TestApplicationChargeServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		"https://fooshop.myshopify.com/admin/application_charges.json",
		httpmock.NewStringResponder(200, `{"application_charges": [{"id":1},{"id":2}]}`),
	)

	charges, err := client.ApplicationCharge.List(nil)
	if err != nil {
		t.Errorf("ApplicationCharge.List returned an error: %v", err)
	}

	expected := []ApplicationCharge{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(charges, expected) {
		t.Errorf("ApplicationCharge.List returned %+v, expected %+v", charges, expected)
	}
}

func TestApplicationChargeServiceOp_Activate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		"https://fooshop.myshopify.com/admin/application_charges/455696195/activate.json",
		httpmock.NewStringResponder(
			200,
			`{"application_charge":{"id":455696195,"status":"active"}}`,
		),
	)

	charge := ApplicationCharge{
		ID:     455696195,
		Status: "accepted",
	}

	returnedCharge, err := client.ApplicationCharge.Activate(charge)
	if err != nil {
		t.Errorf("ApplicationCharge.Activate returned an error: %v", err)
	}

	expected := &ApplicationCharge{ID: 455696195, Status: "active"}
	if !reflect.DeepEqual(returnedCharge, expected) {
		t.Errorf("ApplicationCharge.Activate returned %+v, expected %+v", charge, expected)
	}
}
