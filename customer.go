package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const customersBasePath = "customers"
const customersResourceName = "customers"

// CustomerService is an interface for interfacing with the customers endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/customer
type CustomerService interface {
	List(interface{}) ([]Customer, error)
	ListWithPagination(options interface{}) ([]Customer, *Pagination, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Customer, error)
	Search(interface{}) ([]Customer, error)
	Create(Customer) (*Customer, error)
	Update(Customer) (*Customer, error)
	Delete(int64) error
	ListOrders(int64, interface{}) ([]Order, error)
	ListTags(interface{}) ([]string, error)

	// MetafieldsService used for Customer resource to communicate with Metafields resource
	MetafieldsService
}

// CustomerServiceOp handles communication with the product related methods of
// the Shopify API.
type CustomerServiceOp struct {
	client *Client
}

// Customer represents a Shopify customer
type Customer struct {
	ID                  int64              `json:"id,omitempty"`
	Email               string             `json:"email,omitempty"`
	FirstName           string             `json:"first_name,omitempty"`
	LastName            string             `json:"last_name,omitempty"`
	State               string             `json:"state,omitempty"`
	Note                string             `json:"note,omitempty"`
	VerifiedEmail       bool               `json:"verified_email,omitempty"`
	MultipassIdentifier string             `json:"multipass_identifier,omitempty"`
	OrdersCount         int                `json:"orders_count,omitempty"`
	TaxExempt           bool               `json:"tax_exempt,omitempty"`
	TotalSpent          *decimal.Decimal   `json:"total_spent,omitempty"`
	Phone               string             `json:"phone,omitempty"`
	Tags                string             `json:"tags,omitempty"`
	LastOrderId         int64              `json:"last_order_id,omitempty"`
	LastOrderName       string             `json:"last_order_name,omitempty"`
	AcceptsMarketing    bool               `json:"accepts_marketing,omitempty"`
	DefaultAddress      *CustomerAddress   `json:"default_address,omitempty"`
	Addresses           []*CustomerAddress `json:"addresses,omitempty"`
	CreatedAt           *time.Time         `json:"created_at,omitempty"`
	UpdatedAt           *time.Time         `json:"updated_at,omitempty"`
	Metafields          []Metafield        `json:"metafields,omitempty"`
}

// Represents the result from the customers/X.json endpoint
type CustomerResource struct {
	Customer *Customer `json:"customer"`
}

// Represents the result from the customers.json endpoint
type CustomersResource struct {
	Customers []Customer `json:"customers"`
}

// Represents the result from the customers/tags.json endpoint
type CustomerTagsResource struct {
	Tags []string `json:"tags"`
}

// Represents the options available when searching for a customer
type CustomerSearchOptions struct {
	Page   int    `url:"page,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Fields string `url:"fields,omitempty"`
	Order  string `url:"order,omitempty"`
	Query  string `url:"query,omitempty"`
}

// List customers
func (s *CustomerServiceOp) List(options interface{}) ([]Customer, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	resource := new(CustomersResource)
	err := s.client.Get(path, resource, options)
	return resource.Customers, err
}

// ListWithPagination lists customers and return pagination to retrieve next/previous results.
func (s *CustomerServiceOp) ListWithPagination(options interface{}) ([]Customer, *Pagination, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	resource := new(CustomersResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Customers, pagination, nil
}

// Count customers
func (s *CustomerServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", customersBasePath)
	return s.client.Count(path, options)
}

// Get customer
func (s *CustomerServiceOp) Get(customerID int64, options interface{}) (*Customer, error) {
	path := fmt.Sprintf("%s/%v.json", customersBasePath, customerID)
	resource := new(CustomerResource)
	err := s.client.Get(path, resource, options)
	return resource.Customer, err
}

// Create a new customer
func (s *CustomerServiceOp) Create(customer Customer) (*Customer, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	wrappedData := CustomerResource{Customer: &customer}
	resource := new(CustomerResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Customer, err
}

// Update an existing customer
func (s *CustomerServiceOp) Update(customer Customer) (*Customer, error) {
	path := fmt.Sprintf("%s/%d.json", customersBasePath, customer.ID)
	wrappedData := CustomerResource{Customer: &customer}
	resource := new(CustomerResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Customer, err
}

// Delete an existing customer
func (s *CustomerServiceOp) Delete(customerID int64) error {
	path := fmt.Sprintf("%s/%d.json", customersBasePath, customerID)
	return s.client.Delete(path)
}

// Search customers
func (s *CustomerServiceOp) Search(options interface{}) ([]Customer, error) {
	path := fmt.Sprintf("%s/search.json", customersBasePath)
	resource := new(CustomersResource)
	err := s.client.Get(path, resource, options)
	return resource.Customers, err
}

// ListOrders retrieves all orders from a customer
func (s *CustomerServiceOp) ListOrders(customerID int64, options interface{}) ([]Order, error) {
	path := fmt.Sprintf("%s/%d/orders.json", customersBasePath, customerID)
	resource := new(OrdersResource)
	err := s.client.Get(path, resource, options)
	return resource.Orders, err
}

// ListTags retrieves all unique tags across all customers
func (s *CustomerServiceOp) ListTags(options interface{}) ([]string, error) {
	path := fmt.Sprintf("%s/tags.json", customersBasePath)
	resource := new(CustomerTagsResource)
	err := s.client.Get(path, resource, options)
	return resource.Tags, err
}

// List metafields for a customer
func (s *CustomerServiceOp) ListMetafields(customerID int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.List(options)
}

// Count metafields for a customer
func (s *CustomerServiceOp) CountMetafields(customerID int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.Count(options)
}

// Get individual metafield for a customer
func (s *CustomerServiceOp) GetMetafield(customerID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.Get(metafieldID, options)
}

// Create a new metafield for a customer
func (s *CustomerServiceOp) CreateMetafield(customerID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.Create(metafield)
}

// Update an existing metafield for a customer
func (s *CustomerServiceOp) UpdateMetafield(customerID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.Update(metafield)
}

// // Delete an existing metafield for a customer
func (s *CustomerServiceOp) DeleteMetafield(customerID int64, metafieldID int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceID: customerID}
	return metafieldService.Delete(metafieldID)
}
