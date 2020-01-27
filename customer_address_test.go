package goshopify

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func verifyAddress(t *testing.T, address CustomerAddress) {
	expectedID := int64(1)
	if address.ID != expectedID {
		t.Errorf("CustomerAddress.ID returned %+v, expected %+v", address.ID, expectedID)
	}

	expectedCustomerID := int64(1)
	if address.CustomerID != expectedCustomerID {
		t.Errorf("CustomerAddress.CustomerID returned %+v, expected %+v", address.CustomerID, expectedCustomerID)
	}

	expectedFirstName := "Test"
	if address.FirstName != expectedFirstName {
		t.Errorf("CustomerAddress.FirstName returned %+v, expected %+v", address.FirstName, expectedFirstName)
	}

	expectedLastName := "Citizen"
	if address.LastName != expectedLastName {
		t.Errorf("CustomerAddress.LastName returned %+v, expected %+v", address.LastName, expectedLastName)
	}

	expectedCompany := "TestCo"
	if address.Company != expectedCompany {
		t.Errorf("CustomerAddress.Company returned %+v, expected %+v", address.Company, expectedCompany)
	}

	expectedAddress1 := "1 Smith St"
	if address.Address1 != expectedAddress1 {
		t.Errorf("CustomerAddress.Address1 returned %+v, expected %+v", address.Address1, expectedAddress1)
	}

	expectedAddress2 := ""
	if address.Address2 != expectedAddress2 {
		t.Errorf("CustomerAddress.Address2 returned %+v, expected %+v", address.Address2, expectedAddress2)
	}

	expectedCity := "BRISBANE"
	if address.City != expectedCity {
		t.Errorf("CustomerAddress.City returned %+v, expected %+v", address.City, expectedCity)
	}

	expectedProvince := "Queensland"
	if address.Province != expectedProvince {
		t.Errorf("CustomerAddress.Province returned %+v, expected %+v", address.Province, expectedProvince)
	}

	expectedCountry := "Australia"
	if address.Country != expectedCountry {
		t.Errorf("CustomerAddress.Country returned %+v, expected %+v", address.Country, expectedCountry)
	}

	expectedZip := "4000"
	if address.Zip != expectedZip {
		t.Errorf("CustomerAddress.Zip returned %+v, expected %+v", address.Zip, expectedZip)
	}

	expectedPhone := "1111 111 111"
	if address.Phone != expectedPhone {
		t.Errorf("CustomerAddress.Phone returned %+v, expected %+v", address.Phone, expectedPhone)
	}

	expectedName := "Test Citizen"
	if address.Name != expectedName {
		t.Errorf("CustomerAddress.Name returned %+v, expected %+v", address.Name, expectedName)
	}

	expectedProvinceCode := "QLD"
	if address.ProvinceCode != expectedProvinceCode {
		t.Errorf("CustomerAddress.ProvinceCode returned %+v, expected %+v", address.ProvinceCode, expectedProvinceCode)
	}

	expectedCountryCode := "AU"
	if address.CountryCode != expectedCountryCode {
		t.Errorf("CustomerAddress.CountryCode returned %+v, expected %+v", address.CountryCode, expectedCountryCode)
	}

	expectedCountryName := "Australia"
	if address.CountryName != expectedCountryName {
		t.Errorf("CustomerAddress.CountryName returned %+v, expected %+v", address.CountryName, expectedCountryName)
	}

	expectedDefault := true
	if address.Default != expectedDefault {
		t.Errorf("CustomerAddress.Default returned %+v, expected %+v", address.Default, expectedDefault)
	}
}

func TestList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/customers/1/addresses.json", globalApiPathPrefix), httpmock.NewBytesResponder(200, loadFixture("customer_addresses.json")))

	addresses, err := client.CustomerAddress.List(1, nil)
	if err != nil {
		t.Errorf("CustomerAddress.List returned error: %v", err)
	}

	if len(addresses) != 2 {
		t.Errorf("CustomerAddress.List got %v addresses, expected 2", len(addresses))
	}
	verifyAddress(t, addresses[0])
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/customers/1/addresses/1.json", globalApiPathPrefix), httpmock.NewBytesResponder(200, loadFixture("customer_address.json")))

	address, err := client.CustomerAddress.Get(1, 1, nil)
	if err != nil {
		t.Errorf("CustomerAddress.Get returned error: %v", err)
	}

	verifyAddress(t, *address)
}

func TestCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/customers/1/addresses.json", globalApiPathPrefix), httpmock.NewBytesResponder(200, loadFixture("customer_address.json")))

	address, err := client.CustomerAddress.Create(1, CustomerAddress{})
	if err != nil {
		t.Errorf("CustomerAddress.Create returned error: %v", err)
	}

	verifyAddress(t, *address)
}

func TestUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/customers/1/addresses/1.json", globalApiPathPrefix), httpmock.NewBytesResponder(200, loadFixture("customer_address.json")))

	address, err := client.CustomerAddress.Update(1, CustomerAddress{ID: 1})
	if err != nil {
		t.Errorf("CustomerAddress.Update returned error: %v", err)
	}

	verifyAddress(t, *address)
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/customers/1/addresses/1.json", globalApiPathPrefix), httpmock.NewStringResponder(200, "{}"))

	err := client.CustomerAddress.Delete(1, 1)
	if err != nil {
		t.Errorf("CustomerAddress.Update returned error: %v", err)
	}
}
