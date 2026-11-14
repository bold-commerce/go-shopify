package goshopify

import (
	"context"
	"fmt"
)

// CustomerAddressService is an interface for interfacing with the customer address endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/customer_address
type CustomerAddressService interface {
	List(context.Context, uint64, interface{}) ([]CustomerAddress, error)
	ListWithPagination(ctx context.Context, customerId uint64, options interface{}) ([]CustomerAddress, *Pagination, error)
	ListAll(ctx context.Context, customerId uint64, options interface{}) ([]CustomerAddress, error)
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

// ListWithPagination returns one page and its next/previous page options.
func (s *CustomerAddressServiceOp) ListWithPagination(ctx context.Context, customerId uint64, options interface{}) ([]CustomerAddress, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/addresses.json", customersBasePath, customerId)
	resource := new(CustomerAddressesResource)
	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}
	return resource.Addresses, pagination, nil
}

// ListAll retrieves all pages. On error, it returns the items from completed pages.
func (s *CustomerAddressServiceOp) ListAll(ctx context.Context, customerId uint64, options interface{}) ([]CustomerAddress, error) {
	items := []CustomerAddress{}
	for {
		page, pagination, err := s.ListWithPagination(ctx, customerId, options)
		if err != nil {
			return items, err
		}
		items = append(items, page...)
		if pagination.NextPageOptions == nil {
			return items, nil
		}
		options = pagination.NextPageOptions
	}
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
