package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func draftOrderTests(t *testing.T, draftOrder DraftOrder) {
	// Check that dates are parsed
	d := time.Date(2019, time.April, 9, 10, 02, 43, 0, time.UTC)
	if !d.Equal(*draftOrder.CreatedAt) {
		t.Errorf("Order.CreatedAt returned %+v, expected %+v", draftOrder.CreatedAt, d)
	}

	// Check null dates
	if draftOrder.UpdatedAt == nil {
		t.Errorf("DraftOrder.UpdatedAt returned %+v, expected %+v", draftOrder.UpdatedAt, nil)
	}

	// Check prices
	p := "206.25"
	if !(p == draftOrder.TotalPrice) {
		t.Errorf("draftOrder.TotalPrice returned %+v, expected %+v", draftOrder.TotalPrice, p)
	}

	// Check null prices, notice that prices are usually not empty.
	if draftOrder.TotalTax != "0.00" {
		t.Errorf("draftOrder.TotalTax returned %+v, expected %+v", draftOrder.TotalTax, nil)
	}

	//

	// Check customer
	if draftOrder.Customer == nil {
		t.Error("Expected Customer to not be nil")
	}
	if draftOrder.Customer.Email != "bob.norman@hostmail.com" {
		t.Errorf("Customer.Email, expected %v, actual %v", "bob.norman@hostmail.com", draftOrder.Customer.Email)
	}

	ptp := decimal.NewFromFloat(199)
	lineItem := draftOrder.LineItems[0]
	if !ptp.Equals(*lineItem.Price) {
		t.Errorf("DraftOrder.LineItems[0].Price, expected %v, actual %v", "199.00", lineItem.Price)
	}
}

func TestDraftOrderGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/994118539.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("draft_order.json")))

	draftOrder, err := client.DraftOrder.Get(994118539, nil)
	if err != nil {
		t.Errorf("DraftOrder.Get returned error: %v", err)
	}
	draftOrderTests(t, *draftOrder)
}

func TestDraftOrderCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders.json", globalApiPathPrefix),
		httpmock.NewStringResponder(201, `{"draft_order":{"id": 1}}`))

	draftOrder := DraftOrder{
		LineItems: []LineItem{
			LineItem{
				VariantID: 1,
				Quantity:  1,
			},
		},
	}

	d, err := client.DraftOrder.Create(draftOrder)
	if err != nil {
		t.Errorf("DraftOrder.Create returned error: %v", err)
	}

	expected := DraftOrder{ID: 1}
	if d.ID != expected.ID {
		t.Errorf("DraftOrder.Create returned id %d, expected %d", d.ID, expected.ID)
	}
}

func TestDraftOrderUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"draft_order":{"id": 1}}`))

	draftOrder := DraftOrder{
		ID:            1,
		Note:          "slow order",
		TaxesIncluded: true,
	}

	d, err := client.DraftOrder.Update(draftOrder)
	if err != nil {
		t.Errorf("DraftOrder.Create returned an error %v", err)
	}

	expected := DraftOrder{ID: 1}
	if d.ID != expected.ID {
		t.Errorf("DraftOrder.Update returned id %d, expected %d", d.ID, expected.ID)
	}

}

func TestDraftOrderCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 7}`))

	params := map[string]string{"status": "open"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.DraftOrder.Count(nil)
	if err != nil {
		t.Errorf("DraftOrder.Count returned an error: %v", err)
	}
	expected := 7
	if cnt != expected {
		t.Errorf("DraftOrder.Count returned %d, expected %d", cnt, expected)
	}

	status := "open"
	cnt, err = client.DraftOrder.Count(DraftOrderCountOptions{Status: status})
	if err != nil {
		t.Errorf("DraftOrder.Count returned an error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("DraftOrder.Count returned %d, expected %d", cnt, expected)
	}
}

func TestDraftOrderList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("draft_orders.json")))

	draftOrders, err := client.DraftOrder.List(nil)
	if err != nil {
		t.Errorf("DraftOrder.List returned error: %v", err)
	}

	if len(draftOrders) != 1 {
		t.Errorf("DraftOrder.List got %d orders, expected: 1", len(draftOrders))
	}
	draftOrder := draftOrders[0]
	draftOrderTests(t, draftOrder)
}

func TestDraftOrderListOptions(t *testing.T) {
	setup()
	defer teardown()
	params := map[string]string{
		"fields": "id,name",
		"limit":  "250",
		"status": "any",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders.json", globalApiPathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("draft_orders.json")))

	options := DraftOrderListOptions{
		Limit:  250,
		Status: "any",
		Fields: "id,name",
	}

	draftOrders, err := client.DraftOrder.List(options)
	if err != nil {
		t.Errorf("DraftOrder.List returned error: %v", err)
	}

	if len(draftOrders) != 1 {
		t.Errorf("DraftOrder.List got %d orders, expected: 1", len(draftOrders))
	}

	draftOrder := draftOrders[0]
	draftOrderTests(t, draftOrder)
}

func TestDraftOrderInvoice(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/send_invoice.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(201, loadFixture("invoice.json")))
	invoice := DraftOrderInvoice{
		To:   "first@example.com",
		From: "steve@apple.com",
		Bcc: []string{
			"steve@apple.com",
		},
		Subject:       "Apple Computer Invoice",
		CustomMessage: "Thank you for ordering!",
	}
	draftInvoice, err := client.DraftOrder.Invoice(1, invoice)
	if err != nil {
		t.Errorf("DraftOrder.Invoice returned an error: %v", err)
	}

	if !reflect.DeepEqual(*draftInvoice, invoice) {
		t.Errorf("DraftOrder.Invoice returned %+v, expected %+v,", draftInvoice, invoice)
	}
}

func TestDraftOrderDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, nil))

	err := client.DraftOrder.Delete(1)
	if err != nil {
		t.Errorf("DraftOrder.Delete returned an error %v", err)
	}
}
func TestDraftOrderComplete(t *testing.T) {
	setup()
	defer teardown()
	params := map[string]string{"payment_pending": "false"}
	httpmock.RegisterResponderWithQuery(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/complete.json", globalApiPathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("draft_order.json")))

	draftOrder, err := client.DraftOrder.Complete(1, false)
	if err != nil {
		t.Errorf("DraftOrder.Complete returned an error %v", err)
	}
	draftOrderTests(t, *draftOrder)
}

func TestDraftOrderCompletePending(t *testing.T) {
	setup()
	defer teardown()
	params := map[string]string{"payment_pending": "true"}
	httpmock.RegisterResponderWithQuery(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/complete.json", globalApiPathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("draft_order.json")))

	draftOrder, err := client.DraftOrder.Complete(1, true)
	if err != nil {
		t.Errorf("DraftOrder.Complete returned an error %v", err)
	}
	draftOrderTests(t, *draftOrder)
}

func TestDraftOrderListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.DraftOrder.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("DraftOrder.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Order.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestDraftOrderCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields/count.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.DraftOrder.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("Order.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Order.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.DraftOrder.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Order.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Order.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestDraftOrderGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.DraftOrder.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("Order.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Order.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestDraftOrderCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.DraftOrder.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Order.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestDraftOrderUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.DraftOrder.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Order.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestDraftOrderDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/draft_orders/1/metafields/2.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.DraftOrder.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Order.DeleteMetafield() returned error: %v", err)
	}
}
