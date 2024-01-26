package goshopify

import (
	"context"
	"fmt"
)

const (
	fulfillmentServiceBasePath = "fulfillment_services"
)

// FulfillmentServiceService is an interface for interfacing with the fulfillment service of the Shopify API.
// https://help.shopify.com/api/reference/fulfillmentservice
type FulfillmentServiceService interface {
	List(context.Context, interface{}) ([]FulfillmentServiceData, error)
	Get(context.Context, uint64, interface{}) (*FulfillmentServiceData, error)
	Create(context.Context, FulfillmentServiceData) (*FulfillmentServiceData, error)
	Update(context.Context, FulfillmentServiceData) (*FulfillmentServiceData, error)
	Delete(context.Context, uint64) error
}

type FulfillmentServiceData struct {
	Id                     uint64 `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	Email                  string `json:"email,omitempty"`
	ServiceName            string `json:"service_name,omitempty"`
	Handle                 string `json:"handle,omitempty"`
	FulfillmentOrdersOptIn bool   `json:"fulfillment_orders_opt_in,omitempty"`
	IncludePendingStock    bool   `json:"include_pending_stock,omitempty"`
	ProviderId             uint64 `json:"provider_id,omitempty"`
	LocationId             uint64 `json:"location_id,omitempty"`
	CallbackURL            string `json:"callback_url,omitempty"`
	TrackingSupport        bool   `json:"tracking_support,omitempty"`
	InventoryManagement    bool   `json:"inventory_management,omitempty"`
	AdminGraphqlApiId      string `json:"admin_graphql_api_id,omitempty"`
	PermitsSkuSharing      bool   `json:"permits_sku_sharing,omitempty"`
	RequiresShippingMethod bool   `json:"requires_shipping_method,omitempty"`
	Format                 string `json:"format,omitempty"`
}

type FulfillmentServiceResource struct {
	FulfillmentService *FulfillmentServiceData `json:"fulfillment_service,omitempty"`
}

type FulfillmentServicesResource struct {
	FulfillmentServices []FulfillmentServiceData `json:"fulfillment_services,omitempty"`
}

type FulfillmentServiceOptions struct {
	Scope string `url:"scope,omitempty"`
}

// FulfillmentServiceServiceOp handles communication with the FulfillmentServices
// related methods of the Shopify API
type FulfillmentServiceServiceOp struct {
	client *Client
}

// List Receive a list of all FulfillmentServiceData
func (s *FulfillmentServiceServiceOp) List(ctx context.Context, options interface{}) ([]FulfillmentServiceData, error) {
	path := fmt.Sprintf("%s.json", fulfillmentServiceBasePath)
	resource := new(FulfillmentServicesResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.FulfillmentServices, err
}

// Get Receive a single FulfillmentServiceData
func (s *FulfillmentServiceServiceOp) Get(ctx context.Context, fulfillmentServiceId uint64, options interface{}) (*FulfillmentServiceData, error) {
	path := fmt.Sprintf("%s/%d.json", fulfillmentServiceBasePath, fulfillmentServiceId)
	resource := new(FulfillmentServiceResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.FulfillmentService, err
}

// Create Create a new FulfillmentServiceData
func (s *FulfillmentServiceServiceOp) Create(ctx context.Context, fulfillmentService FulfillmentServiceData) (*FulfillmentServiceData, error) {
	path := fmt.Sprintf("%s.json", fulfillmentServiceBasePath)
	wrappedData := FulfillmentServiceResource{FulfillmentService: &fulfillmentService}
	resource := new(FulfillmentServiceResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.FulfillmentService, err
}

// Update Modify an existing FulfillmentServiceData
func (s *FulfillmentServiceServiceOp) Update(ctx context.Context, fulfillmentService FulfillmentServiceData) (*FulfillmentServiceData, error) {
	path := fmt.Sprintf("%s/%d.json", fulfillmentServiceBasePath, fulfillmentService.Id)
	wrappedData := FulfillmentServiceResource{FulfillmentService: &fulfillmentService}
	resource := new(FulfillmentServiceResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.FulfillmentService, err
}

// Delete Remove an existing FulfillmentServiceData
func (s *FulfillmentServiceServiceOp) Delete(ctx context.Context, fulfillmentServiceId uint64) error {
	path := fmt.Sprintf("%s/%d.json", fulfillmentServiceBasePath, fulfillmentServiceId)
	return s.client.Delete(ctx, path)
}
