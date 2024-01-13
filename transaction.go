package goshopify

import (
	"context"
	"fmt"
)

// TransactionService is an interface for interfacing with the transactions endpoints of
// the Shopify API.
// See: https://help.shopify.com/api/reference/transaction
type TransactionService interface {
	List(context.Context, int64, interface{}) ([]Transaction, error)
	Count(context.Context, int64, interface{}) (int, error)
	Get(context.Context, int64, int64, interface{}) (*Transaction, error)
	Create(context.Context, int64, Transaction) (*Transaction, error)
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
func (s *TransactionServiceOp) List(ctx context.Context, orderID int64, options interface{}) ([]Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions.json", ordersBasePath, orderID)
	resource := new(TransactionsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Transactions, err
}

// Count transactions
func (s *TransactionServiceOp) Count(ctx context.Context, orderID int64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/transactions/count.json", ordersBasePath, orderID)
	return s.client.Count(ctx, path, options)
}

// Get individual transaction
func (s *TransactionServiceOp) Get(ctx context.Context, orderID int64, transactionID int64, options interface{}) (*Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions/%d.json", ordersBasePath, orderID, transactionID)
	resource := new(TransactionResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Transaction, err
}

// Create a new transaction
func (s *TransactionServiceOp) Create(ctx context.Context, orderID int64, transaction Transaction) (*Transaction, error) {
	path := fmt.Sprintf("%s/%d/transactions.json", ordersBasePath, orderID)
	wrappedData := TransactionResource{Transaction: &transaction}
	resource := new(TransactionResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Transaction, err
}
