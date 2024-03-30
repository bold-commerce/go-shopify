package goshopify

import (
	"context"
	"fmt"
)

const (
	ordersRiskBasePath     = "orders"
	ordersRiskResourceName = "risks"
)

// OrderRiskService is an interface for interfacing with the orders Risk endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-10/resources/order-risk
type OrderRiskService interface {
	List(context.Context, uint64, interface{}) ([]OrderRisk, error)
	ListAll(context.Context, uint64, interface{}) ([]OrderRisk, error)
	ListWithPagination(context.Context, uint64, interface{}) ([]OrderRisk, *Pagination, error)
	Get(context.Context, uint64, uint64, interface{}) (*OrderRisk, error)
	Create(context.Context, uint64, OrderRisk) (*OrderRisk, error)
	Update(context.Context, uint64, uint64, OrderRisk) (*OrderRisk, error)
	Delete(context.Context, uint64, uint64) error
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
	// order is fraudulent.
	OrderRecommendationCancel orderRiskRecommendation = "cancel"

	// medium level of risk that this order is fraudulent.
	OrderRecommendationInvestigate orderRiskRecommendation = "investigate"

	// level of risk that this order is fraudulent.
	OrderRecommendationAccept orderRiskRecommendation = "accept"
)

// A struct for all available order Risk list options.
// See: https://shopify.dev/docs/api/admin-rest/2023-10/resources/order-risk#index
type OrderRiskListOptions struct {
	ListOptions
}

// OrderRisk represents a Shopify order risk
type OrderRisk struct {
	Id              uint64                  `json:"id,omitempty"`
	CheckoutId      uint64                  `json:"checkout_id,omitempty"`
	OrderId         uint64                  `json:"order_id,omitempty"`
	CauseCancel     bool                    `json:"cause_cancel,omitempty"`
	Display         bool                    `json:"display,omitempty"`
	MerchantMessage string                  `json:"merchant_message,omitempty"`
	Message         string                  `json:"message,omitempty"`
	Score           string                  `json:"score,omitempty"`
	Source          string                  `json:"source,omitempty"`
	Recommendation  orderRiskRecommendation `json:"recommendation,omitempty"`
}

// List OrderRisk
func (s *OrderRiskServiceOp) List(ctx context.Context, orderId uint64, options interface{}) ([]OrderRisk, error) {
	orders, _, err := s.ListWithPagination(ctx, orderId, options)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// ListAll Lists all OrderRisk, iterating over pages
func (s *OrderRiskServiceOp) ListAll(ctx context.Context, orderId uint64, options interface{}) ([]OrderRisk, error) {
	collector := []OrderRisk{}

	for {
		entities, pagination, err := s.ListWithPagination(ctx, orderId, options)

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

func (s *OrderRiskServiceOp) ListWithPagination(ctx context.Context, orderId uint64, options interface{}) ([]OrderRisk, *Pagination, error) {
	path := fmt.Sprintf("%s/%d/%s.json", ordersRiskBasePath, orderId, ordersRiskResourceName)
	resource := new(OrdersRisksResource)

	pagination, err := s.client.ListWithPagination(ctx, path, resource, options)
	if err != nil {
		return nil, nil, err
	}

	return resource.OrderRisk, pagination, nil
}

// Get individual order
func (s *OrderRiskServiceOp) Get(ctx context.Context, orderId uint64, riskId uint64, options interface{}) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderId, ordersRiskResourceName, riskId)
	resource := new(OrderRiskResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.OrderRisk, err
}

// Create order
func (s *OrderRiskServiceOp) Create(ctx context.Context, orderId uint64, orderRisk OrderRisk) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s.json", ordersRiskBasePath, orderId, ordersRiskResourceName)
	wrappedData := OrderRiskResource{OrderRisk: &orderRisk}
	resource := new(OrderRiskResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.OrderRisk, err
}

// Update order
func (s *OrderRiskServiceOp) Update(ctx context.Context, orderId uint64, riskId uint64, orderRisk OrderRisk) (*OrderRisk, error) {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderId, ordersRiskResourceName, riskId)
	wrappedData := OrderRiskResource{OrderRisk: &orderRisk}
	resource := new(OrderRiskResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.OrderRisk, err
}

// Delete order
func (s *OrderRiskServiceOp) Delete(ctx context.Context, orderId uint64, riskId uint64) error {
	path := fmt.Sprintf("%s/%d/%s/%d.json", ordersRiskBasePath, orderId, ordersRiskResourceName, riskId)
	err := s.client.Delete(ctx, path)
	return err
}
