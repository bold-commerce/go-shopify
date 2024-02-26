package goshopify

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func TransactionTests(t *testing.T, transaction Transaction) {
	// Check that the Id is assigned to the returned transaction
	expectedId := uint64(389404469)
	if transaction.Id != expectedId {
		t.Errorf("Transaction.Id returned %+v, expected %+v", transaction.Id, expectedId)
	}

	// Check that the OrderId value is assigned to the returned transaction
	expectedOrderId := uint64(450789469)
	if transaction.OrderId != expectedOrderId {
		t.Errorf("Transaction.OrderId returned %+v, expected %+v", transaction.OrderId, expectedOrderId)
	}

	// Check that the Amount value is assigned to the returned transaction
	expectedAmount, _ := decimal.NewFromString("409.94")
	if !transaction.Amount.Equals(expectedAmount) {
		t.Errorf("Transaction.Amount returned %+v, expected %+v", transaction.Amount, expectedAmount)
	}

	// Check that the Kind value is assigned to the returned transaction
	expectedKind := "authorization"
	if transaction.Kind != expectedKind {
		t.Errorf("Transaction.Kind returned %+v, expected %+v", transaction.Kind, expectedKind)
	}

	// Check that the Gateway value is assigned to the returned transaction
	expectedGateway := "bogus"
	if transaction.Gateway != expectedGateway {
		t.Errorf("Transaction.Gateway returned %+v, expected %+v", transaction.Gateway, expectedGateway)
	}

	// Check that the Status value is assigned to the returned transaction
	expectedStatus := "success"
	if transaction.Status != expectedStatus {
		t.Errorf("Transaction.Status returned %+v, expected %+v", transaction.Status, expectedStatus)
	}

	// Check that the Message value is assigned to the returned transaction
	expectedMessage := "Bogus Gateway: Forced success"
	if transaction.Message != expectedMessage {
		t.Errorf("Transaction.Message returned %+v, expected %+v", transaction.Message, expectedMessage)
	}

	// Check that the CreatedAt value is assigned to the returned transaction
	expectedCreatedAt := time.Date(2017, time.July, 24, 19, 9, 43, 0, time.UTC)
	if !expectedCreatedAt.Equal(*transaction.CreatedAt) {
		t.Errorf("Transaction.CreatedAt returned %+v, expected %+v", transaction.CreatedAt, expectedCreatedAt)
	}

	// Check that the Test value is assigned to the returned transaction
	expectedTest := true
	if transaction.Test != expectedTest {
		t.Errorf("Transaction.Test returned %+v, expected %+v", transaction.Test, expectedTest)
	}

	// Check that the Authorization value is assigned to the returned transaction
	expectedAuthorization := "authorization-key"
	if transaction.Authorization != expectedAuthorization {
		t.Errorf("Transaction.Authorization returned %+v, expected %+v", transaction.Authorization, expectedAuthorization)
	}

	// Check that the Currency value is assigned to the returned transaction
	expectedCurrency := "USD"
	if transaction.Currency != expectedCurrency {
		t.Errorf("Transaction.Currency returned %+v, expected %+v", transaction.Currency, expectedCurrency)
	}

	// Check that the LocationId value is assigned to the returned transaction
	var expectedLocationId *int64
	if transaction.LocationId != expectedLocationId {
		t.Errorf("Transaction.LocationId returned %+v, expected %+v", transaction.LocationId, expectedLocationId)
	}

	// Check that the UserId value is assigned to the returned transaction
	var expectedUserId *int64
	if transaction.UserId != expectedUserId {
		t.Errorf("Transaction.UserId returned %+v, expected %+v", transaction.UserId, expectedUserId)
	}

	// Check that the ParentId value is assigned to the returned transaction
	var expectedParentId *int64
	if transaction.ParentId != expectedParentId {
		t.Errorf("Transaction.ParentId returned %+v, expected %+v", transaction.ParentId, expectedParentId)
	}

	// Check that the DeviceId value is assigned to the returned transaction
	var expectedDeviceId *int64
	if transaction.DeviceId != expectedDeviceId {
		t.Errorf("Transacion.DeviceId returned %+v, expected %+v", transaction.DeviceId, expectedDeviceId)
	}

	// Check that the ErrorCode value is assigned to the returned transaction
	var expectedErrorCode string
	if transaction.ErrorCode != expectedErrorCode {
		t.Errorf("Transaction.ErrorCode returned %+v, expected %+v", transaction.ErrorCode, expectedErrorCode)
	}

	// Check that the SourceName value is assigned to the returned transaction
	expectedSourceName := "web"
	if transaction.SourceName != expectedSourceName {
		t.Errorf("Transaction.SourceName returned %+v, expected %+v", transaction.SourceName, expectedSourceName)
	}

	// Check that the PaymentDetails value is assigned to the returned transaction
	var nilString string
	expectedPaymentDetails := PaymentDetails{
		AVSResultCode:     nilString,
		CreditCardBin:     nilString,
		CVVResultCode:     nilString,
		CreditCardNumber:  "•••• •••• •••• 4242",
		CreditCardCompany: "Visa",
	}
	if transaction.PaymentDetails.AVSResultCode != expectedPaymentDetails.AVSResultCode {
		t.Errorf("Transaction.PaymentDetails.AVSResultCode returned %+v, expected %+v",
			transaction.PaymentDetails.AVSResultCode, expectedPaymentDetails.AVSResultCode)
	}
}

func TestTransactionList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transactions.json")))

	transactions, err := client.Transaction.List(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Transaction.List returned error: %v", err)
	}

	for _, transaction := range transactions {
		TransactionTests(t, transaction)
	}
}

func TestTransactionCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Transaction.Count(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Transaction.Count returned error: %v", err)
	}

	expected := 2
	if cnt != expected {
		t.Errorf("Transaction.Count returned %d, expected %d", cnt, expected)
	}
}

func TestTransactionGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transaction.json")))

	transaction, err := client.Transaction.Get(context.Background(), 1, 1, nil)
	if err != nil {
		t.Errorf("Transaction.Get returned error: %v", err)
	}

	TransactionTests(t, *transaction)
}

func TestTransactionCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transaction.json")))

	amount := decimal.NewFromFloat(409.94)

	transaction := Transaction{
		Amount: &amount,
	}
	result, err := client.Transaction.Create(context.Background(), 1, transaction)
	if err != nil {
		t.Errorf("Transaction.Create returned error: %+v", err)
	}
	TransactionTests(t, *result)
}
