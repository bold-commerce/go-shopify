package goshopify

import (
	"context"
	"fmt"
)

// TransactionService is an interface for interfacing with the transactions endpoints of
// the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/transaction
type TransactionService interface {
	List(context.Context, uint64, interface{}) ([]Transaction, error)
	Count(context.Context, uint64, interface{}) (int, error)
	Get(context.Context, uint64, uint64, interface{}) (*Transaction, error)
	Create(context.Context, uint64, Transaction) (*Transaction, error)
}

// TransactionServiceOp handles communication with the transaction related methods of the
// Shopify API.
type TransactionServiceOp struct {
	client *Client
}

// TransactionResource represents the result from the orders/X/transactions/Y.json endpoint
type TransactionResource struct {
	Transaction *Transaction `json:"transaction"`
}

// TransactionsResource represents the result from the orders/X/transactions.json endpoint
type TransactionsResource struct {
	Transactions []Transaction `json:"transactions"`
}

// List transactions
func (s *TransactionServiceOp) List(ctx context.Context, orderId uint64, options interface{}) ([]Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions.json", ordersBasePath, orderId)
	resource := new(TransactionsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Transactions, err
}

// Count transactions
func (s *TransactionServiceOp) Count(ctx context.Context, orderId uint64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/transactions/count.json", ordersBasePath, orderId)
	return s.client.Count(ctx, path, options)
}

// Get individual transaction
func (s *TransactionServiceOp) Get(ctx context.Context, orderId uint64, transactionId uint64, options interface{}) (*Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions/%d.json", ordersBasePath, orderId, transactionId)
	resource := new(TransactionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Transaction, err
}

// Create a new transaction
func (s *TransactionServiceOp) Create(ctx context.Context, orderId uint64, transaction Transaction) (*Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions.json", ordersBasePath, orderId)
	wrappedData := TransactionResource{Transaction: &transaction}
	resource := new(TransactionResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Transaction, err
}
