package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const carrierBasePath = "carrier_services"

// CarrierServiceService is an interface for interfacing with the carrier service endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/admin-api/rest/reference/shipping-and-fulfillment/carrierservice
type CarrierServiceService interface {
	List() ([]CarrierService, error)
	Get(int64) (*CarrierService, error)
	Create(CarrierService) (*CarrierService, error)
	Update(CarrierService) (*CarrierService, error)
	Delete(int64) error
}

// CarrierServiceOp handles communication with the product related methods of
// the Shopify API.
type CarrierServiceOp struct {
	client *Client
}

// CarrierService represents a Shopify carrier service
type CarrierService struct {
	// Whether this carrier service is active.
	Active bool `json:"active,omitempty"`

	// The URL endpoint that Shopify needs to retrieve shipping rates. This must be a public URL.
	CallbackUrl string `json:"callback_url,omitempty"`

	// Distinguishes between API or legacy carrier services.
	CarrierServiceType string `json:"carrier_service_type,omitempty"`

	// The Id of the carrier service.
	Id int64 `json:"id,omitempty"`

	// The format of the data returned by the URL endpoint. Valid values: json and xml. Default value: json.
	Format string `json:"format,omitempty"`

	// The name of the shipping service as seen by merchants and their customers.
	Name string `json:"name,omitempty"`

	// Whether merchants are able to send dummy data to your service through the Shopify admin to see shipping rate examples.
	ServiceDiscovery bool `json:"service_discovery,omitempty"`

	AdminGraphqlAPIID string `json:"admin_graphql_api_id,omitempty"`
}

type SingleCarrierResource struct {
	CarrierService *CarrierService `json:"carrier_service"`
}

type ListCarrierResource struct {
	CarrierServices []CarrierService `json:"carrier_services"`
}

type ShippingRateRequest struct {
	Rate ShippingRateQuery `json:"rate"`
}

type ShippingRateQuery struct {
	Origin      ShippingRateAddress `json:"origin"`
	Destination ShippingRateAddress `json:"destination"`
	Items       []LineItem          `json:"items"`
	Currency    string              `json:"currency"`
	Locale      string              `json:"locale"`
}

// The address3, fax, address_type, and company_name fields are returned by specific ActiveShipping providers.
// For API-created carrier services, you should use only the following shipping address fields:
// * address1
// * address2
// * city
// * zip
// * province
// * country
// Other values remain as null and are not sent to the callback URL.
type ShippingRateAddress struct {
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Name        string `json:"name"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Address3    string `json:"address3"`
	Phone       string `json:"phone"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	AddressType string `json:"address_type"`
	CompanyName string `json:"company_name"`
}

// When Shopify requests shipping rates using your callback URL,
// the response object rates must be a JSON array of objects with the following fields.
// Required fields must be included in the response for the carrier service integration to work properly.
type ShippingRateResponse struct {
	Rates []ShippingRate `json:"rates"`
}

type ShippingRate struct {
	// The name of the rate, which customers see at checkout. For example: Expedited Mail.
	ServiceName string `json:"service_name"`

	// A description of the rate, which customers see at checkout. For example: Includes tracking and insurance.
	Description string `json:"description"`

	// A unique code associated with the rate. For example: expedited_mail.
	ServiceCode string `json:"service_code"`

	// The currency of the shipping rate.
	Currency string `json:"currency"`

	// The total price based on the shipping rate currency.
	// In cents unit. See https://github.com/Shopify/shipping-fulfillment-app/issues/15#issuecomment-725996936
	TotalPrice decimal.Decimal `json:"total_price"`

	// Whether the customer must provide a phone number at checkout.
	PhoneRequired bool `json:phone_required,omitempty"`

	// The earliest delivery date for the displayed rate.
	MinDeliveryDate *time.Time `json:"min_delivery_date"` // "2013-04-12 14:48:45 -0400"

	// The latest delivery date for the displayed rate to still be valid.
	MaxDeliveryDate *time.Time `json:"max_delivery_date"` // "2013-04-12 14:48:45 -0400"
}

// List carrier services
func (s *CarrierServiceOp) List() ([]CarrierService, error) {
	path := fmt.Sprintf("%s.json", carrierBasePath)
	resource := new(ListCarrierResource)
	err := s.client.Get(path, resource, nil)
	return resource.CarrierServices, err
}

// Get individual carrier resource by carrier resource ID
func (s *CarrierServiceOp) Get(id int64) (*CarrierService, error) {
	path := fmt.Sprintf("%s/%d.json", carrierBasePath, id)
	resource := new(SingleCarrierResource)
	err := s.client.Get(path, resource, nil)
	return resource.CarrierService, err
}

// Create a carrier service
func (s *CarrierServiceOp) Create(carrier CarrierService) (*CarrierService, error) {
	path := fmt.Sprintf("%s.json", carrierBasePath)
	body := SingleCarrierResource{
		CarrierService: &carrier,
	}
	resource := new(SingleCarrierResource)
	err := s.client.Post(path, body, resource)
	return resource.CarrierService, err
}

// Update a carrier service
func (s *CarrierServiceOp) Update(carrier CarrierService) (*CarrierService, error) {
	path := fmt.Sprintf("%s/%d.json", carrierBasePath, carrier.Id)
	body := SingleCarrierResource{
		CarrierService: &carrier,
	}
	resource := new(SingleCarrierResource)
	err := s.client.Put(path, body, resource)
	return resource.CarrierService, err
}

// Delete a carrier service
func (s *CarrierServiceOp) Delete(id int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", carrierBasePath, id))
}
