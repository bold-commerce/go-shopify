package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func TestOrderListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	orders, err := client.Order.List(nil)
	if orders != nil {
		t.Errorf("Order.List returned orders, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Order.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestOrderListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedOrders     []Order
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"orders": [{"id":1},{"id":2}]}`,
			"",
			[]Order{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Order(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Order(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Order(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Order(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Order(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"orders": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Order{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"orders": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Order{{ID: 2}},
			&Pagination{
				NextPageOptions:     &ListOptions{PageInfo: "foo"},
				PreviousPageOptions: &ListOptions{PageInfo: "bar"},
			},
			nil,
		},
	}
	for i, c := range cases {
		response := &http.Response{
			StatusCode: 200,
			Body:       httpmock.NewRespBodyFromString(c.body),
			Header: http.Header{
				"Link": {c.linkHeader},
			},
		}

		httpmock.RegisterResponder("GET", listURL, httpmock.ResponderFromResponse(response))

		orders, pagination, err := client.Order.ListWithPagination(nil)
		if !reflect.DeepEqual(orders, c.expectedOrders) {
			t.Errorf("test %d Order.ListWithPagination orders returned %+v, expected %+v", i, orders, c.expectedOrders)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Order.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Order.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", client.pathPrefix),
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
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", client.pathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("orders.json")))

	options := OrderListOptions{
		ListOptions: ListOptions{
			Page:   10,
			Limit:  250,
			Fields: "id,name",
		},

		Status: "any",
	}

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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 7}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/count.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders.json", client.pathPrefix),
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1.json", client.pathPrefix),
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

func TestOrderCancel(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456/cancel.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("order_with_transaction.json")))

	order, err := client.Order.Cancel(123456, nil)
	if err != nil {
		t.Errorf("Order.Update returned error: %v", err)
	}
	// Check that dates are parsed
	timezone, _ := time.LoadLocation("America/New_York")

	d := time.Date(2016, time.May, 17, 4, 14, 36, 0, timezone)
	if !d.Equal(*order.CancelledAt) {
		t.Errorf("Order.CancelledAt returned %+v, expected %+v", order.CancelledAt, d)
	}

	orderTests(t, *order)
}

func TestOrderClose(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456/close.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"order":{"closed_at":"2016-05-17T04:14:36-04:00"}}`))

	order, err := client.Order.Close(123456)
	if err != nil {
		t.Errorf("Order.Update returned error: %v", err)
	}
	// Check that dates are parsed
	timezone, _ := time.LoadLocation("America/New_York")

	d := time.Date(2016, time.May, 17, 4, 14, 36, 0, timezone)
	if !d.Equal(*order.ClosedAt) {
		t.Errorf("Order.ClosedAt returned %+v, expected %+v", order.ClosedAt, d)
	}
}

func TestOrderOpen(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/123456/open.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"order":{"closed_at":null}}`))

	order, err := client.Order.Open(123456)
	if err != nil {
		t.Errorf("Order.Update returned error: %v", err)
	}

	if order.ClosedAt != nil {
		t.Errorf("Order.ClosedAt returned %+v, expected nil", order.ClosedAt)
	}
}

func TestOrderListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/count.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Type:      "single_line_text_field",
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

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Order.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Order.DeleteMetafield() returned error: %v", err)
	}
}

func TestOrderListFulfillments(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/count.json", client.pathPrefix),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments.json", client.pathPrefix),
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/1022782888.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/complete.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/open.json", client.pathPrefix),
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

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/fulfillments/2/cancel.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("fulfillment.json")))

	returnedFulfillment, err := client.Order.CancelFulfillment(1, 2)
	if err != nil {
		t.Errorf("Order.CancelFulfillment() returned error: %v", err)
	}

	FulfillmentTests(t, *returnedFulfillment)
}

// TestLineItemUnmarshalJSON tests unmarsalling a LineItem from json
func TestLineItemUnmarshalJSON(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/valid.json"))
	if err != nil {
		t.Errorf("LineItem.UnmarshalJSON returned error: %v", err)
	}

	expected := validLineItem()

	testLineItem(t, expected, actual)
}

// TestLineItemUnmarshalJSONInvalid0 tests unmarsalling a LineItem from invalid json
func TestLineItemUnmarshalJSONInvalid0(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_invalid0.json"))
	if err == nil || !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Errorf("LineItem.UnmarshalJSON expected unexpected end of JSON input error got %v", err)
	}
}

// TestLineItemUnmarshalJSONInvalid1 tests unmarsalling a LineItem with properties that are a struct with invalid
// values
func TestLineItemUnmarshalJSONInvalid1(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_invalid1.json"))
	if err == nil || !strings.Contains(err.Error(), "cannot unmarshal number") {
		t.Errorf("LineItem.UnmarshalJSON expected cannot unmarshal number into string error got %v", err)
	}
}

// TestLineItemUnmarshalJSONInvalid2 tests unmarsalling a LineItem with properties that are an array with invalid
// values
func TestLineItemUnmarshalJSONInvalid2(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_invalid2.json"))
	if err == nil || !strings.Contains(err.Error(), "cannot unmarshal number") {
		t.Errorf("LineItem.UnmarshalJSON expected cannot unmarshal number into string error got %v", err)
	}
}

// TestLineItemUnmarshalJSONPropertiesOmitted tests unmarshalling a LineItem where the "properties" property is omitted
func TestLineItemUnmarshalJSONPropertiesOmitted(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_ommited.json"))
	if err != nil {
		t.Errorf("LineItem.UnmarshalJSON returned error: %v", err)
	}

	expected := propertiesEmptyStructLientItem()

	testLineItem(t, expected, actual)
}

// TestLineItemUnmarshalJSONPropertiesEmptyObject tests unmarsalling a LineItem from json which has properties as an empty json object
func TestLineItemUnmarshalJSONPropertiesEmptyObject(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_empty_object.json"))
	if err != nil {
		t.Errorf("LineItem.UnmarshalJSON returned error: %v", err)
	}

	expected := propertiesEmptyStructLientItem()

	testLineItem(t, expected, actual)
}

// TestLineItemUnmarshalJSONPropertiesObject tests unmarsalling a LineItem from json which has properties as an json object
func TestLineItemUnmarshalJSONPropertiesObject(t *testing.T) {
	setup()
	defer teardown()

	actual := LineItem{}

	err := actual.UnmarshalJSON(loadFixture("orderlineitems/properties_object.json"))
	if err != nil {
		t.Errorf("LineItem.UnmarshalJSON returned error: %v", err)
	}

	expected := propertiesStructLientItem()

	testLineItem(t, expected, actual)
}

// TestShippingLines tests unmarshalling ShippingLines.RequestFulfillmentServiceID from a JSON string
func TestShippingLines_UnmarshallJSON(t *testing.T) {
	setup()
	defer teardown()

	actual := ShippingLines{}

	err := actual.UnmarshalJSON(loadFixture("shippinglines/valid.json"))
	if err != nil {
		t.Errorf("ShippingLines.UnmarshalJSON returned error: %v", err)
	}

	expected := validShippingLines()

	testShippingLines(t, expected, actual)
}

// TestShippingLines tests unmarshalling ShippingLines.RequestFulfillmentServiceID from a JSON Number
func TestShippingLines_UnmarshallJSON_RequestFulfillmentServiceIDNumber(t *testing.T) {
	setup()
	defer teardown()

	actual := ShippingLines{}

	err := actual.UnmarshalJSON(loadFixture("shippinglines/requested_fulfillment_service_id_number.json"))
	if err != nil {
		t.Errorf("ShippingLines.UnmarshalJSON returned error: %v", err)
	}

	expected := validShippingLines()
	expected.RequestedFulfillmentServiceID = "123456"

	testShippingLines(t, expected, actual)
}

// TestShippingLines tests unmarshalling ShippingLines.RequestFulfillmentServiceID from a JSON null
func TestShippingLines_UnmarshallJSON_RequestFulfillmentServiceIDNull(t *testing.T) {
	setup()
	defer teardown()

	actual := ShippingLines{}

	err := actual.UnmarshalJSON(loadFixture("shippinglines/requested_fulfillment_service_id_null.json"))
	if err != nil {
		t.Errorf("ShippingLines.UnmarshalJSON returned error: %v", err)
	}

	expected := validShippingLines()
	expected.RequestedFulfillmentServiceID = ""

	testShippingLines(t, expected, actual)
}

// TestShippingLines tests unmarshalling ShippingLines.RequestFulfillmentServiceID from a malformed JSON
func TestShippingLines_UnmarshallJSON_RequestFulfillmentServiceIDInvalid(t *testing.T) {
	setup()
	defer teardown()

	actual := ShippingLines{}

	err := actual.UnmarshalJSON(loadFixture("shippinglines/requested_fulfillment_service_id_invalid.json"))
	if err == nil || !strings.Contains(err.Error(), "invalid character '}'") {
		t.Errorf("ShippingLines.UnmarshalJSON expected invalid character '}' %v", err)
	}
}

func testLineItem(t *testing.T, expected, actual LineItem) {
	if actual.ID != expected.ID {
		t.Errorf("LineItem.ID should be (%v), was (%v)", expected.ID, actual.ID)
	}

	if actual.ProductID != expected.ProductID {
		t.Errorf("LineItem.ProductID should be (%v), was (%v)", expected.ProductID, actual.ProductID)
	}

	if actual.VariantID != expected.VariantID {
		t.Errorf("LineItem.VariantID should be (%v), was (%v)", expected.VariantID, actual.VariantID)
	}

	if actual.Quantity != expected.Quantity {
		t.Errorf("LineItem.Quantity should be (%v), was (%v)", expected.Quantity, actual.Quantity)
	}

	if actual.Price == nil {
		if actual.Price != expected.Price {
			t.Errorf("LineItem.Price should be (%s), was (%s)", expected.Price, actual.Price)
		}
	} else {
		if !actual.Price.Equals(*expected.Price) {
			t.Errorf("LineItem.Price should be (%s), was (%s)", expected.Price, actual.Price)
		}
	}

	if actual.TotalDiscount == nil {
		if actual.TotalDiscount != expected.TotalDiscount {
			t.Errorf("LineItem.TotalDiscount should be (%s), was (%s)", expected.TotalDiscount, actual.TotalDiscount)
		}
	} else {
		if !actual.TotalDiscount.Equals(*expected.TotalDiscount) {
			t.Errorf("LineItem.TotalDiscount should be (%s), was (%s)", expected.TotalDiscount, actual.TotalDiscount)
		}
	}

	if actual.Title != expected.Title {
		t.Errorf("LineItem.Title should be (%v), was (%v)", expected.Title, actual.Title)
	}

	if actual.VariantTitle != expected.VariantTitle {
		t.Errorf("LineItem.VariantTitle should be (%v), was (%v)", expected.VariantTitle, actual.VariantTitle)
	}

	if actual.Name != expected.Name {
		t.Errorf("LineItem.Name should be (%v), was (%v)", expected.Name, actual.Name)
	}

	if actual.SKU != expected.SKU {
		t.Errorf("LineItem.SKU should be (%v), was (%v)", expected.SKU, actual.SKU)
	}

	if actual.Vendor != expected.Vendor {
		t.Errorf("LineItem.Vendor should be (%v), was (%v)", expected.Vendor, actual.Vendor)
	}

	if actual.GiftCard != expected.GiftCard {
		t.Errorf("LineItem.GiftCard should be (%v), was (%v)", expected.GiftCard, actual.GiftCard)
	}

	if actual.Taxable != expected.Taxable {
		t.Errorf("LineItem.Taxable should be (%v), was (%v)", expected.Taxable, actual.Taxable)
	}

	if actual.FulfillmentService != expected.FulfillmentService {
		t.Errorf("LineItem.FulfillmentService should be (%v), was (%v)", expected.FulfillmentService, actual.FulfillmentService)
	}

	if actual.RequiresShipping != expected.RequiresShipping {
		t.Errorf("LineItem.RequiresShipping should be (%v), was (%v)", expected.RequiresShipping, actual.RequiresShipping)
	}

	if actual.VariantInventoryManagement != expected.VariantInventoryManagement {
		t.Errorf("LineItem.VariantInventoryManagement should be (%v), was (%v)", expected.VariantInventoryManagement, actual.VariantInventoryManagement)
	}

	if actual.PreTaxPrice == nil {
		if actual.PreTaxPrice != expected.PreTaxPrice {
			t.Errorf("LineItem.PreTaxPrice should be (%v), was (%v)", expected.PreTaxPrice, actual.PreTaxPrice)
		}
	} else {
		if !actual.PreTaxPrice.Equals(*expected.PreTaxPrice) {
			t.Errorf("LineItem.PreTaxPrice should be (%v), was (%v)", expected.PreTaxPrice, actual.PreTaxPrice)
		}
	}

	testProperties(t, expected.Properties, actual.Properties)

	if actual.ProductExists != expected.ProductExists {
		t.Errorf("LineItem.ProductExists should be (%v), was (%v)", expected.ProductExists, actual.ProductExists)
	}

	if actual.FulfillableQuantity != expected.FulfillableQuantity {
		t.Errorf("LineItem.FulfillableQuantity should be (%v), was (%v)", expected.FulfillableQuantity, actual.FulfillableQuantity)
	}

	if actual.Grams != expected.Grams {
		t.Errorf("LineItem.Grams should be (%v), was (%v)", expected.Grams, actual.Grams)
	}

	if actual.FulfillmentStatus != expected.FulfillmentStatus {
		t.Errorf("LineItem.FulfillmentStatus should be (%v), was (%v)", expected.FulfillmentStatus, actual.FulfillmentStatus)
	}

	testTaxLines(t, expected.TaxLines, actual.TaxLines)

	if actual.OriginLocation == nil {
		if actual.OriginLocation != expected.OriginLocation {
			t.Errorf("LineItem.OriginLocation should be (%v), was (%v)", expected.OriginLocation, actual.OriginLocation)
		}
	} else {
		if *actual.OriginLocation != *expected.OriginLocation {
			t.Errorf("LineItem.OriginLocation should be (%v), was (%v)", expected.OriginLocation, actual.OriginLocation)
		}
	}

	if actual.DestinationLocation == nil {
		if actual.DestinationLocation != expected.DestinationLocation {
			t.Errorf("LineItem.DestinationLocation should be (%v), was (%v)", expected.DestinationLocation, actual.DestinationLocation)
		}
	} else {
		if *actual.DestinationLocation != *expected.DestinationLocation {
			t.Errorf("LineItem.DestinationLocation should be (%v), was (%v)", expected.DestinationLocation, actual.DestinationLocation)
		}
	}

	if actual.AppliedDiscount == nil {
		if actual.AppliedDiscount != expected.AppliedDiscount {
			t.Errorf("LineItem.AppliedDiscount should be (%v), was (%v)", expected.AppliedDiscount, actual.AppliedDiscount)
		}
	} else {
		if *actual.AppliedDiscount != *expected.AppliedDiscount {
			t.Errorf("LineItem.AppliedDiscount should be (%v), was (%v)", expected.AppliedDiscount, actual.AppliedDiscount)
		}
	}
}

func testProperties(t *testing.T, expected, actual []NoteAttribute) {
	if len(expected) != len(actual) {
		t.Errorf("LineItem.Properties expected len (%d) actual (%d)", len(expected), len(actual))
	} else {
		for i := 0; i < len(actual); i++ {
			a := actual[i]
			if a.Name != expected[i].Name {
				t.Errorf("LineItem.Properties[%d].Name should be (%s), was (%s)", i, expected[i].Name, a.Name)
			}
			if a.Value != expected[i].Value {
				t.Errorf("LineItem.Properties[%d].Value should be (%s), was (%s)", i, expected[i].Value, a.Value)
			}
		}
	}
}

func testTaxLines(t *testing.T, expected, actual []TaxLine) {
	if len(expected) != len(actual) {
		t.Errorf("LineItem.TaxLines expected len (%d) actual (%d)", len(expected), len(actual))
	} else {
		for i := 0; i < len(actual); i++ {
			a := actual[i]
			e := expected[i]
			if !a.Price.Equals(*e.Price) {
				t.Errorf("LineItem.TaxLine[%d].Price should be (%s), was (%s)", i, e.Price, a.Price)
			}
			if !a.Rate.Equals(*e.Rate) {
				t.Errorf("LineItem.TaxLine[%d].Rate should be (%s), was (%s)", i, e.Rate, a.Rate)
			}
			if a.Title != e.Title {
				t.Errorf("LineItem.TaxLine[%d].Title should be (%s), was (%s)", i, e.Title, a.Title)
			}
		}
	}
}

func testShippingLines(t *testing.T, expected, actual ShippingLines) {
	if actual.ID != expected.ID {
		t.Errorf("ShippingLines.ID should be (%v), was (%v)", expected.ID, actual.ID)
	}

	if actual.Title != expected.Title {
		t.Errorf("ShippingLines.Title should be (%v), was (%v)", expected.Title, actual.Title)
	}

	if !actual.Price.Equals(*expected.Price) {
		t.Errorf("ShippingLines.Price should be (%v), was (%v)", expected.Price, actual.Price)
	}

	if actual.Code != expected.Code {
		t.Errorf("ShippingLines.Code should be (%v), was (%v)", expected.Code, actual.Code)
	}

	if actual.Source != expected.Source {
		t.Errorf("ShippingLines.Source should be (%v), was (%v)", expected.Source, actual.Source)
	}

	if actual.Phone != expected.Phone {
		t.Errorf("ShippingLines.Phone should be (%v), was (%v)", expected.Phone, actual.Phone)
	}

	if actual.RequestedFulfillmentServiceID != expected.RequestedFulfillmentServiceID {
		t.Errorf("ShippingLines.RequestedFulfillmentServiceID should be (%v), was (%v)", expected.RequestedFulfillmentServiceID, actual.RequestedFulfillmentServiceID)
	}

	if actual.DeliveryCategory != expected.DeliveryCategory {
		t.Errorf("ShippingLines.DeliveryCategory should be (%v), was (%v)", expected.DeliveryCategory, actual.DeliveryCategory)
	}

	if actual.CarrierIdentifier != expected.CarrierIdentifier {
		t.Errorf("ShippingLines.CarrierIdentifier should be (%v), was (%v)", expected.CarrierIdentifier, actual.CarrierIdentifier)
	}

	testTaxLines(t, expected.TaxLines, actual.TaxLines)
}

func propertiesEmptyStructLientItem() LineItem {
	return LineItem{
		Properties: []NoteAttribute{},
	}
}

func propertiesStructLientItem() LineItem {
	return LineItem{
		Properties: []NoteAttribute{
			NoteAttribute{
				Name:  "property 1",
				Value: float64(3),
			},
		},
	}
}

func validLineItem() LineItem {
	price := decimal.New(1234, -2)
	totalDiscount := decimal.New(123, -2)
	preTaxPrice := decimal.New(900, -2)
	tl1Price := decimal.New(1350, -2)
	tl1Rate := decimal.New(6, -2)
	tl2Price := decimal.New(1250, -2)
	tl2Rate := decimal.New(5, -2)
	discountAllocationAmount := decimal.New(55, -1)
	return LineItem{
		ID:                         int64(254721536),
		ProductID:                  int64(111475476),
		VariantID:                  int64(1234),
		Quantity:                   1,
		Price:                      &price,
		TotalDiscount:              &totalDiscount,
		Title:                      "Soda Title",
		VariantTitle:               "Test Variant",
		Name:                       "Soda",
		SKU:                        "sku-123",
		Vendor:                     "Test Vendor",
		GiftCard:                   true,
		Taxable:                    true,
		FulfillmentService:         "manual",
		RequiresShipping:           true,
		VariantInventoryManagement: "shopify",
		PreTaxPrice:                &preTaxPrice,
		Properties: []NoteAttribute{
			NoteAttribute{
				Name:  "note 1",
				Value: "one",
			},
			NoteAttribute{
				Name:  "note 2",
				Value: float64(2),
			},
		},
		ProductExists:       true,
		FulfillableQuantity: 1,
		Grams:               100,
		FulfillmentStatus:   "partial",
		TaxLines: []TaxLine{
			TaxLine{
				Title: "State tax",
				Price: &tl1Price,
				Rate:  &tl1Rate,
			},
			TaxLine{
				Title: "Federal tax",
				Price: &tl2Price,
				Rate:  &tl2Rate,
			},
		},
		OriginLocation: &Address{
			ID:           123,
			Address1:     "100 some street",
			Address2:     "",
			City:         "Winnipeg",
			Company:      "Acme Corporation",
			Country:      "Canada",
			CountryCode:  "CA",
			FirstName:    "Bob",
			LastName:     "Smith",
			Latitude:     49.811550,
			Longitude:    -97.189480,
			Name:         "test address",
			Phone:        "8675309",
			Province:     "Manitoba",
			ProvinceCode: "MB",
			Zip:          "R3Y 0L6",
		},
		DestinationLocation: &Address{
			ID:           124,
			Address1:     "200 some street",
			Address2:     "",
			City:         "Winnipeg",
			Company:      "Acme Corporation",
			Country:      "Canada",
			CountryCode:  "CA",
			FirstName:    "Bob",
			LastName:     "Smith",
			Latitude:     49.811550,
			Longitude:    -97.189480,
			Name:         "test address",
			Phone:        "8675309",
			Province:     "Manitoba",
			ProvinceCode: "MB",
			Zip:          "R3Y 0L6",
		},
		AppliedDiscount: &AppliedDiscount{
			Title:       "test discount",
			Description: "my test discount",
			Value:       "0.05",
			ValueType:   "percent",
			Amount:      "25.00",
		},
		DiscountAllocations: []DiscountAllocations{
			{
				Amount: &discountAllocationAmount,
				AmountSet: AmountSet{
					ShopMoney: AmountSetEntry{
						Amount:       &discountAllocationAmount,
						CurrencyCode: "EUR",
					},
					PresentmentMoney: AmountSetEntry{
						Amount:       &discountAllocationAmount,
						CurrencyCode: "EUR",
					},
				},
			},
		},
	}
}

func validShippingLines() ShippingLines {
	price := decimal.New(400, -2)
	tl1Price := decimal.New(1350, -2)
	tl1Rate := decimal.New(6, -2)
	tl2Price := decimal.New(1250, -2)
	tl2Rate := decimal.New(5, -2)

	return ShippingLines{
		ID:                            int64(254721542),
		Title:                         "Small Packet International Air",
		Price:                         &price,
		Code:                          "INT.TP",
		Source:                        "canada_post",
		Phone:                         "",
		RequestedFulfillmentServiceID: "third_party_fulfillment_service_id",
		DeliveryCategory:              "",
		CarrierIdentifier:             "third_party_carrier_identifier",
		TaxLines: []TaxLine{
			{
				Title: "State tax",
				Price: &tl1Price,
				Rate:  &tl1Rate,
			},
			{
				Title: "Federal tax",
				Price: &tl2Price,
				Rate:  &tl2Rate,
			},
		},
	}
}
