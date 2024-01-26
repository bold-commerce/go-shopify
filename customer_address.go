package goshopify

import (
	"context"
	"fmt"
)

// CustomerAddressService is an interface for interfacing with the customer address endpoints
// of the Shopify API.
// See: https://help.shopify.com/en/api/reference/customers/customer_address
type CustomerAddressService interface {
	List(context.Context, uint64, interface{}) ([]CustomerAddress, error)
	Get(context.Context, uint64, uint64, interface{}) (*CustomerAddress, error)
	Create(context.Context, uint64, CustomerAddress) (*CustomerAddress, error)
	Update(context.Context, uint64, CustomerAddress) (*CustomerAddress, error)
	Delete(context.Context, uint64, uint64) error
}

// CustomerAddressServiceOp handles communication with the customer address related methods of
// the Shopify API.
type CustomerAddressServiceOp struct {
	client *Client
}

// CustomerAddress represents a Shopify customer address
type CustomerAddress struct {
	Id           uint64 `json:"id,omitempty"`
	CustomerId   uint64 `json:"customer_id,omitempty"`
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
func (s *CustomerAddressServiceOp) List(ctx context.Context, customerId uint64, options interface{}) ([]CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses.json", customersBasePath, customerId)
	resource := new(CustomerAddressesResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Addresses, err
}

// Get address
func (s *CustomerAddressServiceOp) Get(ctx context.Context, customerId, addressId uint64, options interface{}) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerId, addressId)
	resource := new(CustomerAddressResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Address, err
}

// Create a new address for given customer
func (s *CustomerAddressServiceOp) Create(ctx context.Context, customerId uint64, address CustomerAddress) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses.json", customersBasePath, customerId)
	wrappedData := CustomerAddressResource{Address: &address}
	resource := new(CustomerAddressResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Address, err
}

// Create a new address for given customer
func (s *CustomerAddressServiceOp) Update(ctx context.Context, customerId uint64, address CustomerAddress) (*CustomerAddress, error) {
	path := fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerId, address.Id)
	wrappedData := CustomerAddressResource{Address: &address}
	resource := new(CustomerAddressResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Address, err
}

// Delete an existing address
func (s *CustomerAddressServiceOp) Delete(ctx context.Context, customerId, addressId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d/addresses/%d.json", customersBasePath, customerId, addressId))
}
