package goshopify

import (
	"fmt"
)

const ordersRiskBasePath = "orders"
const ordersRiskResourceName = "risks"

// OrderRiskService is an interface for interfacing with the orders Risk endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-10/resources/order-risk
type OrderRiskService interface {
	List(int64, interface{}) ([]OrderRisk, error)
	ListWithPagination(int64, interface{}) ([]OrderRisk, *Pagination, error)
	Get(int64, int64, interface{}) (*OrderRisk, error)
	Create(int64, OrderRisk) (*OrderRisk, error)
	Update(int64, int64, OrderRisk) (*OrderRisk, error)
	Delete(int64, int64) error
}

// OrderRiskServiceOp handles communication with the order related methods of the
// Shopify API.
type OrderRiskServiceOp struct {
	client *Client
}

// Represents the result from the orders-risk/X.json endpoint
type OrderRiskResource struct {
	OrderRisk *OrderRisk `json:"risk"`
}

// Represents the result from the orders-risk.json endpoint
type OrdersRisksResource struct {
	OrderRisk []OrderRisk `json:"risks"`
}
type orderRiskRecommendation string

const (
	//order is fraudulent.
	OrderRecommendationCancel orderRiskRecommendation = "cancel"

	//medium level of risk that this order is fraudulent.
	OrderRecommendationInvestigate orderRiskRecommendation = "investigate"

	//level of risk that this order is fraudulent.
	OrderRecommendationAccept orderRiskRecommendation = "accept"
)

// A struct for all available order Risk list options.
// See: https://shopify.dev/docs/api/admin-rest/2023-10/resources/order-risk#index
type OrderRiskListOptions struct {
	ListOptions
}

// OrderRisk represents a Shopify order risk
type OrderRisk struct {
	ID              int64                   `json:"id,omitempty"`
	CheckoutID      int64                   `json:"checkout_id,omitempty"`
	OrderID         int64                   `json:"order_id,omitempty"`
	CauseCancel     bool                    `json:"cause_cancel,omitempty"`
	Display         bool                    `json:"display,omitempty"`
	MerchantMessage string                  `json:"merchant_message,omitempty"`
	Message         string                  `json:"message,omitempty"`
	Score           string                  `json:"score,omitempty"`
	Source          string                  `json:"source,omitempty"`
	Recommendation  orderRiskRecommendation `json:"recommendation,omitempty"`
}

// List OrderRisk
func (s *OrderRiskServiceOp) List(orderId int64, options interface{}) ([]OrderRisk, error) {
	orders, _, err := s.ListWithPagination(orderId, options)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderRiskServiceOp) ListWithPagination(orderId int64, options interface{}) ([]OrderRisk, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/%s.json", ordersRiskBasePath, orderId, ordersRiskResourceName)
	resource := new(OrdersRisksResource)

	pagination, err := s.client.ListWithPagination(path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.OrderRisk, pagination, nil
}

// Get individual order
func (s *OrderRiskServiceOp) Get(orderID int64, riskID int64, options interface{}) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderID, ordersRiskResourceName, riskID)
	resource := new(OrderRiskResource)
	err := s.client.Get(path, resource, options)
	return resource.OrderRisk, err
}

// Create order
func (s *OrderRiskServiceOp) Create(orderID int64, orderRisk OrderRisk) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s.json", ordersRiskBasePath, orderID, ordersRiskResourceName)
	wrappedData := OrderRiskResource{OrderRisk: &orderRisk}
	resource := new(OrderRiskResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.OrderRisk, err
}

// Update order
func (s *OrderRiskServiceOp) Update(orderID int64, riskID int64, orderRisk OrderRisk) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderID, ordersRiskResourceName, riskID)
	wrappedData := OrderRiskResource{OrderRisk: &orderRisk}
	resource := new(OrderRiskResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.OrderRisk, err
}

// Delete order
func (s *OrderRiskServiceOp) Delete(orderID int64, riskID int64) error {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderID, ordersRiskResourceName, riskID)
	err := s.client.Delete(path)
	return err
}
