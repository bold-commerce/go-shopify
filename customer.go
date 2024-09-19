package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	customersBasePath     = "customers"
	customersResourceName = "customers"
)

// CustomerService is an interface for interfacing with the customers endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/customer
type CustomerService interface {
	List(context.Context, interface{}) ([]Customer, error)
	ListAll(context.Context, interface{}) ([]Customer, error)
	ListWithPagination(ctx context.Context, options interface{}) ([]Customer, *Pagination, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Customer, error)
	Search(context.Context, interface{}) ([]Customer, error)
	Create(context.Context, Customer) (*Customer, error)
	Update(context.Context, Customer) (*Customer, error)
	Delete(context.Context, uint64) error
	ListOrders(context.Context, uint64, interface{}) ([]Order, error)
	ListTags(context.Context, interface{}) ([]string, error)

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
	Id                        uint64                 `json:"id,omitempty"`
	Email                     string                 `json:"email,omitempty"`
	FirstName                 string                 `json:"first_name,omitempty"`
	LastName                  string                 `json:"last_name,omitempty"`
	State                     string                 `json:"state,omitempty"`
	Note                      string                 `json:"note,omitempty"`
	VerifiedEmail             bool                   `json:"verified_email,omitempty"`
	MultipassIdentifier       string                 `json:"multipass_identifier,omitempty"`
	OrdersCount               int                    `json:"orders_count,omitempty"`
	TaxExempt                 bool                   `json:"tax_exempt,omitempty"`
	TotalSpent                *decimal.Decimal       `json:"total_spent,omitempty"`
	Phone                     string                 `json:"phone,omitempty"`
	Currency                  string                 `json:"currency,omitempty"`
	AcceptsMarketing          bool                   `json:"accepts_marketing,omitempty"`
	AcceptsMarketingUpdatedAt *time.Time             `json:"accepts_marketing_updated_at,omitempty"`
	MarketingOptInLevel       bool                   `json:"marketing_opt_in_level,omitempty"`
	AdminGraphqlApiId         string                 `json:"admin_graphql_api_id,omitempty"`
	Tags                      string                 `json:"tags,omitempty"`
	LastOrderId               uint64                 `json:"last_order_id,omitempty"`
	LastOrderName             string                 `json:"last_order_name,omitempty"`
	EmailMarketingConsent     *EmailMarketingConsent `json:"email_marketing_consent"`
	SMSMarketingConsent       *SMSMarketingConsent   `json:"sms_marketing_consent"`
	DefaultAddress            *CustomerAddress       `json:"default_address,omitempty"`
	Addresses                 []*CustomerAddress     `json:"addresses,omitempty"`
	CreatedAt                 *time.Time             `json:"created_at,omitempty"`
	UpdatedAt                 *time.Time             `json:"updated_at,omitempty"`
	Metafields                []Metafield            `json:"metafields,omitempty"`
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

type EmailMarketingConsent struct {
	State            string     `json:"state"`
	OptInLevel       string     `json:"opt_in_level"`
	ConsentUpdatedAt *time.Time `json:"consent_updated_at"`
}

type SMSMarketingConsent struct {
	State                string     `json:"state"`
	OptInLevel           string     `json:"opt_in_level"`
	ConsentUpdatedAt     *time.Time `json:"consent_updated_at"`
	ConsentCollectedFrom string     `json:"consent_collected_from"`
}

// List customers
func (s *CustomerServiceOp) List(ctx context.Context, options interface{}) ([]Customer, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	resource := new(CustomersResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Customers, err
}

// ListAll Lists all customers, iterating over pages
func (s *CustomerServiceOp) ListAll(ctx context.Context, options interface{}) ([]Customer, error) {
	collector := []Customer{}

	for {
		entities, pagination, err := s.ListWithPagination(ctx, options)

		if err != nil {
			return collector, err
		}

		collector = append(collector, entities...)

		if pagination.NextPageOptions == nil {
			break
		}

		options = pagination.NextPageOptions
	}

	return collector, nil
}

// ListWithPagination lists customers and return pagination to retrieve next/previous results.
func (s *CustomerServiceOp) ListWithPagination(ctx context.Context, options interface{}) ([]Customer, *Pagination, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	resource := new(CustomersResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.Customers, pagination, nil
}

// Count customers
func (s *CustomerServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", customersBasePath)
	return s.client.Count(ctx, path, options)
}

// Get customer
func (s *CustomerServiceOp) Get(ctx context.Context, customerId uint64, options interface{}) (*Customer, error) {
	path := fmt.Sprintf("%s/%v.json", customersBasePath, customerId)
	resource := new(CustomerResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Customer, err
}

// Create a new customer
func (s *CustomerServiceOp) Create(ctx context.Context, customer Customer) (*Customer, error) {
	path := fmt.Sprintf("%s.json", customersBasePath)
	wrappedData := CustomerResource{Customer: &customer}
	resource := new(CustomerResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Customer, err
}

// Update an existing customer
func (s *CustomerServiceOp) Update(ctx context.Context, customer Customer) (*Customer, error) {
	path := fmt.Sprintf("%s/%d.json", customersBasePath, customer.Id)
	wrappedData := CustomerResource{Customer: &customer}
	resource := new(CustomerResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Customer, err
}

// Delete an existing customer
func (s *CustomerServiceOp) Delete(ctx context.Context, customerId uint64) error {
	path := fmt.Sprintf("%s/%d.json", customersBasePath, customerId)
	return s.client.Delete(ctx, path)
}

// Search customers
func (s *CustomerServiceOp) Search(ctx context.Context, options interface{}) ([]Customer, error) {
	path := fmt.Sprintf("%s/search.json", customersBasePath)
	resource := new(CustomersResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Customers, err
}

// ListOrders retrieves all orders from a customer
func (s *CustomerServiceOp) ListOrders(ctx context.Context, customerId uint64, options interface{}) ([]Order, error) {
	path := fmt.Sprintf("%s/%d/orders.json", customersBasePath, customerId)
	resource := new(OrdersResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Orders, err
}

// ListTags retrieves all unique tags across all customers
func (s *CustomerServiceOp) ListTags(ctx context.Context, options interface{}) ([]string, error) {
	path := fmt.Sprintf("%s/tags.json", customersBasePath)
	resource := new(CustomerTagsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Tags, err
}

// List metafields for a customer
func (s *CustomerServiceOp) ListMetafields(ctx context.Context, customerId uint64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.List(ctx, options)
}

// Count metafields for a customer
func (s *CustomerServiceOp) CountMetafields(ctx context.Context, customerId uint64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.Count(ctx, options)
}

// Get individual metafield for a customer
func (s *CustomerServiceOp) GetMetafield(ctx context.Context, customerId uint64, metafieldId uint64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.Get(ctx, metafieldId, options)
}

// Create a new metafield for a customer
func (s *CustomerServiceOp) CreateMetafield(ctx context.Context, customerId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.Create(ctx, metafield)
}

// Update an existing metafield for a customer
func (s *CustomerServiceOp) UpdateMetafield(ctx context.Context, customerId uint64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.Update(ctx, metafield)
}

// // Delete an existing metafield for a customer
func (s *CustomerServiceOp) DeleteMetafield(ctx context.Context, customerId uint64, metafieldId uint64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: customersResourceName, resourceId: customerId}
	return metafieldService.Delete(ctx, metafieldId)
}
