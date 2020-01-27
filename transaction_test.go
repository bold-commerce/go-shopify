package goshopify

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func TransactionTests(t *testing.T, transaction Transaction) {
	// Check that the ID is assigned to the returned transaction
	expectedID := int64(389404469)
	if transaction.ID != expectedID {
		t.Errorf("Transaction.ID returned %+v, expected %+v", transaction.ID, expectedID)
	}

	// Check that the OrderID value is assigned to the returned transaction
	expectedOrderID := int64(450789469)
	if transaction.OrderID != expectedOrderID {
		t.Errorf("Transaction.OrderID returned %+v, expected %+v", transaction.OrderID, expectedOrderID)
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

	// Check that the LocationID value is assigned to the returned transaction
	var expectedLocationID *int64
	if transaction.LocationID != expectedLocationID {
		t.Errorf("Transaction.LocationID returned %+v, expected %+v", transaction.LocationID, expectedLocationID)
	}

	// Check that the UserID value is assigned to the returned transaction
	var expectedUserID *int64
	if transaction.UserID != expectedUserID {
		t.Errorf("Transaction.UserID returned %+v, expected %+v", transaction.UserID, expectedUserID)
	}

	// Check that the ParentID value is assigned to the returned transaction
	var expectedParentID *int64
	if transaction.ParentID != expectedParentID {
		t.Errorf("Transaction.ParentID returned %+v, expected %+v", transaction.ParentID, expectedParentID)
	}

	// Check that the DeviceID value is assigned to the returned transaction
	var expectedDeviceID *int64
	if transaction.DeviceID != expectedDeviceID {
		t.Errorf("Transacion.DeviceID returned %+v, expected %+v", transaction.DeviceID, expectedDeviceID)
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transactions.json")))

	transactions, err := client.Transaction.List(1, nil)
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Transaction.Count(1, nil)
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transaction.json")))

	transaction, err := client.Transaction.Get(1, 1, nil)
	if err != nil {
		t.Errorf("Transaction.Get returned error: %v", err)
	}

	TransactionTests(t, *transaction)
}

func TestTransactionCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/transactions.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("transaction.json")))

	amount := decimal.NewFromFloat(409.94)

	transaction := Transaction{
		Amount: &amount,
	}
	result, err := client.Transaction.Create(1, transaction)
	if err != nil {
		t.Errorf("Transaction.Create returned error: %+v", err)
	}
	TransactionTests(t, *result)
}
