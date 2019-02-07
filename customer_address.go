package goshopify

import "fmt"

const customerAddressResourceName = "customer-addresses"

// CustomerAddressService is an interface for interfacing with the customer address endpoints
// of the Shopify API.
// See: https://help.shopify.com/en/api/reference/customers/customer_address
type CustomerAddressService interface {
	List(int, interface{}) ([]CustomerAddress, error)
	Get(int, int, interface{}) (*CustomerAddress, error)
	Create(int, CustomerAddress) (*CustomerAddress, error)
	Update(int, CustomerAddress) (*CustomerAddress, error)
	Delete(int, int) error
}

// CustomerAddressServiceOp handles communication with the customer address related methods of
// the Shopify API.
type CustomerAddressServiceOp struct {
	client *Client
}

// CustomerAddress represents a Shopify customer address
type CustomerAddress struct {
	ID           int    `json:"id,omitempty"`
	CustomerID   int    `json:"customer_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Company      string `json:"company,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	Country      string `json:"country,omitempty"`
	Zip          string `json:"zip,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Name         string `json:"name,omitempty"`
	ProvinceCode string `json:"province_code,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
	CountryName  string `json:"country_name,omitempty"`
	Default      bool   `json:"default,omitempty"`
}

// CustomerAddressResoruce represents the result from the addresses/X.json endpoint
type CustomerAddressResource struct {
	Address *CustomerAddress `json:"customer_address"`
}

// CustomerAddressResoruce represents the result from the customers/X/addresses.json endpoint
type CustomerAddressesResource struct {
	Addresses []CustomerAddress `json:"addresses"`
}

// List addresses
func (s *CustomerAddressServiceOp) List(customerID int, options interface{}) ([]CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses.json", customersBasePath, customerID)
	resource := new(CustomerAddressesResource)
	err := s.client.Get(path, resource, options)
	return resource.Addresses, err
}

// Get address
func (s *CustomerAddressServiceOp) Get(customerID, addressID int, options interface{}) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerID, addressID)
	resource := new(CustomerAddressResource)
	err := s.client.Get(path, resource, options)
	return resource.Address, err
}

// Create a new address for given customer
func (s *CustomerAddressServiceOp) Create(customerID int, address CustomerAddress) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses.json", customersBasePath, customerID)
	wrappedData := CustomerAddressResource{Address: &address}
	resource := new(CustomerAddressResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Address, err
}

// Create a new address for given customer
func (s *CustomerAddressServiceOp) Update(customerID int, address CustomerAddress) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerID, address.ID)
	wrappedData := CustomerAddressResource{Address: &address}
	resource := new(CustomerAddressResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Address, err
}

// Delete an existing address
func (s *CustomerAddressServiceOp) Delete(customerID, addressID int) error {
	return s.client.Delete(fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerID, addressID))
}
