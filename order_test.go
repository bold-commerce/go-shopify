package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func orderTests(t *testing.T, order Order) {
	// Check that dates are parsed
	d := time.Date(2016, time.May, 17, 4, 14, 36, 0, time.UTC)
	if !d.Equal(*order.CreatedAt) {
		t.Errorf("Order.CreatedAt returned %+v, expected %+v", order.CreatedAt, d)
	}

	// Check null dates
	if order.ProcessedAt != nil {
		t.Errorf("Order.ProcessedAt returned %+v, expected %+v", order.ProcessedAt, nil)
	}

	// Check prices
	p := decimal.NewFromFloat(10)
	if !p.Equals(*order.TotalPrice) {
		t.Errorf("Order.TotalPrice returned %+v, expected %+v", order.TotalPrice, p)
	}

	// Check null prices, notice that prices are usually not empty.
	if order.TotalTax != nil {
		t.Errorf("Order.TotalTax returned %+v, expected %+v", order.TotalTax, nil)
	}

	// Check customer
	if order.Customer == nil {
		t.Error("Expected Customer to not be nil")
	}
	if order.Customer.Email != "john@test.com" {
		t.Errorf("Customer.Email, expected %v, actual %v", "john@test.com", order.Customer.Email)
	}

	ptp := decimal.NewFromFloat(9)
	lineItem := order.LineItems[0]
	if !ptp.Equals(*lineItem.PreTaxPrice) {
		t.Errorf("Order.LineItems[0].PreTaxPrice, expected %v, actual %v", "9.00", lineItem.PreTaxPrice)
	}
}

func transactionTest(t *testing.T, transaction Transaction) {
	// Check that dates are parsed
	d := time.Date(2017, time.October, 9, 19, 26, 23, 0, time.UTC)
	if !d.Equal(*transaction.CreatedAt) {
		t.Errorf("Transaction.CreatedAt returned %+v, expected %+v", transaction.CreatedAt, d)
	}

	// Check null value
	if transaction.LocationID != nil {
		t.Error("Expected Transaction.LocationID to be nil")
	}

	if transaction.PaymentDetails == nil {
		t.Error("Expected Transaction.PaymentDetails to not be nil")
	}
}

func TestOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("orders.json")))

	orders, err := client.Order.List(nil)
	if err != nil {
		t.Errorf("Order.List returned error: %v", err)
	}

	// Check that orders were parsed
	if len(orders) != 1 {
		t.Errorf("Order.List got %v orders, expected: 1", len(orders))
	}

	order := orders[0]
	orderTests(t, order)
}

func TestOrderListOptions(t *testing.T) {
	setup()
	defer teardown()
	params := map[string]string{
		"fields": "id,name",
		"limit":  "250",
		"page":   "10",
		"status": "any",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", globalApiPathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("orders.json")))

	options := OrderListOptions{
		Page:   10,
		Limit:  250,
		Status: "any",
		Fields: "id,name"}

	orders, err := client.Order.List(options)
	if err != nil {
		t.Errorf("Order.List returned error: %v", err)
	}

	// Check that orders were parsed
	if len(orders) != 1 {
		t.Errorf("Order.List got %v orders, expected: 1", len(orders))
	}

	order := orders[0]
	orderTests(t, order)
}

func TestOrderGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("order.json")))

	order, err := client.Order.Get(123456, nil)
	if err != nil {
		t.Errorf("Order.List returned error: %v", err)
	}

	// Check that dates are parsed
	timezone, _ := time.LoadLocation("America/New_York")

	d := time.Date(2016, time.May, 17, 4, 14, 36, 0, timezone)
	if !d.Equal(*order.CancelledAt) {
		t.Errorf("Order.CancelledAt returned %+v, expected %+v", order.CancelledAt, d)
	}

	orderTests(t, *order)
}

func TestOrderGetWithTransactions(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("order_with_transaction.json")))

	options := struct {
		ApiFeatures string `url:"_apiFeatures"`
	}{"include-transactions"}

	order, err := client.Order.Get(123456, options)
	if err != nil {
		t.Errorf("Order.List returned error: %v", err)
	}

	orderTests(t, *order)

	// Check transactions is not nil
	if order.Transactions == nil {
		t.Error("Expected Transactions to not be nil")
	}
	if len(order.Transactions) != 1 {
		t.Errorf("Expected Transactions to have 1 transaction but received %v", len(order.Transactions))
	}

	transactionTest(t, order.Transactions[0])
}

func TestOrderCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 7}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Order.Count(nil)
	if err != nil {
		t.Errorf("Order.Count returned error: %v", err)
	}

	expected := 7
	if cnt != expected {
		t.Errorf("Order.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Order.Count(OrderCountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Order.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Order.Count returned %d, expected %d", cnt, expected)
	}
}

func TestOrderCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", globalApiPathPrefix),
		httpmock.NewStringResponder(201, `{"order":{"id": 1}}`))

	order := Order{
		LineItems: []LineItem{
			LineItem{
				VariantID: 1,
				Quantity:  1,
			},
		},
	}

	o, err := client.Order.Create(order)
	if err != nil {
		t.Errorf("Order.Create returned error: %v", err)
	}

	expected := Order{ID: 1}
	if o.ID != expected.ID {
		t.Errorf("Order.Create returned id %d, expected %d", o.ID, expected.ID)
	}
}

func TestOrderUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(201, `{"order":{"id": 1}}`))

	order := Order{
		ID:                1,
		FinancialStatus:   "paid",
		FulfillmentStatus: "fulfilled",
	}

	o, err := client.Order.Update(order)
	if err != nil {
		t.Errorf("Order.Update returned error: %v", err)
	}

	expected := Order{ID: 1}
	if o.ID != expected.ID {
		t.Errorf("Order.Update returned id %d, expected %d", o.ID, expected.ID)
	}
}

func TestOrderListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Order.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("Order.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Order.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestOrderCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Order.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("Order.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Order.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Order.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Order.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Order.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestOrderGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Order.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("Order.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Order.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestOrderCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Order.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Order.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestOrderUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Order.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Order.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestOrderDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Order.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Order.DeleteMetafield() returned error: %v", err)
	}
}

func TestOrderListFulfillments(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillments": [{"id":1},{"id":2}]}`))

	fulfillments, err := client.Order.ListFulfillments(1, nil)
	if err != nil {
		t.Errorf("Order.ListFulfillments() returned error: %v", err)
	}

	expected := []Fulfillment{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(fulfillments, expected) {
		t.Errorf("Order.ListFulfillments() returned %+v, expected %+v", fulfillments, expected)
	}
}

func TestOrderCountFulfillments(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Order.CountFulfillments(1, nil)
	if err != nil {
		t.Errorf("Order.CountFulfillments() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Order.CountFulfillments() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Order.CountFulfillments(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Order.CountFulfillments() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Order.CountFulfillments() returned %d, expected %d", cnt, expected)
	}
}

func TestOrderGetFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"fulfillment": {"id":2}}`))

	fulfillment, err := client.Order.GetFulfillment(1, 2, nil)
	if err != nil {
		t.Errorf("Order.GetFulfillment() returned error: %v", err)
	}

	expected := &Fulfillment{ID: 2}
	if !reflect.DeepEqual(fulfillment, expected) {
		t.Errorf("Order.GetFulfillment() returned %+v, expected %+v", fulfillment, expected)
	}
}

func TestOrderCreateFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillment := Fulfillment{
		LocationID:     905684977,
		TrackingNumber: "123456789",
		TrackingUrls: []string{
			"https://shipping.xyz/track.php?num=123456789",
			"https://anothershipper.corp/track.php?code=abc",
		},
		NotifyCustomer: true,
	}

	returnedFulfillment, err := client.Order.CreateFulfillment(1, fulfillment)
	if err != nil {
		t.Errorf("Order.CreateFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestOrderUpdateFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/1022782888.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	fulfillment := Fulfillment{
		ID:             1022782888,
		TrackingNumber: "987654321",
	}
	returnedFulfillment, err := client.Order.UpdateFulfillment(1, fulfillment)
	if err != nil {
		t.Errorf("Order.UpdateFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestOrderCompleteFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/complete.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	returnedFulfillment, err := client.Order.CompleteFulfillment(1, 2)
	if err != nil {
		t.Errorf("Order.CompleteFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestOrderTransitionFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/open.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	returnedFulfillment, err := client.Order.TransitionFulfillment(1, 2)
	if err != nil {
		t.Errorf("Order.TransitionFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

func TestOrderCancelFulfillment(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/cancel.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	returnedFulfillment, err := client.Order.CancelFulfillment(1, 2)
	if err != nil {
		t.Errorf("Order.CancelFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}
