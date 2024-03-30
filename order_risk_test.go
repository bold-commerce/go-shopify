package goshopify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestOrderRiskListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	orders, err := client.OrderRisk.List(context.Background(), 450789469, nil)
	if orders != nil {
		t.Errorf("OrderRisk.List returned orders, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("OrderRisk.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestOrderRiskListAll(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix)

	cases := []struct {
		name                string
		expectedOrderRisks  []OrderRisk
		expectedRequestURLs []string
		expectedLinkHeaders []string
		expectedBodies      []string
		expectedErr         error
	}{
		{
			name: "Pulls the next page",
			expectedRequestURLs: []string{
				listURL,
				fmt.Sprintf("%s?page_info=pg2", listURL),
			},
			expectedLinkHeaders: []string{
				`<http://valid.url?page_info=pg2>; rel="next"`,
				`<http://valid.url?page_info=pg1>; rel="previous"`,
			},
			expectedBodies: []string{
				`{"risks": [{"id":1},{"id":2}]}`,
				`{"risks": [{"id":3},{"id":4}]}`,
			},
			expectedOrderRisks: []OrderRisk{{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4}},
			expectedErr:        nil,
		},
		{
			name: "Stops when there is not a next page",
			expectedRequestURLs: []string{
				listURL,
			},
			expectedLinkHeaders: []string{
				`<http://valid.url?page_info=pg2>; rel="previous"`,
			},
			expectedBodies: []string{
				`{"risks": [{"id":1}]}`,
			},
			expectedOrderRisks: []OrderRisk{{Id: 1}},
			expectedErr:        nil,
		},
	}

	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.expectedRequestURLs) != len(c.expectedLinkHeaders) {
				t.Errorf(
					"test case must have the same number of expected request urls (%d) as expected link headers (%d)",
					len(c.expectedRequestURLs),
					len(c.expectedLinkHeaders),
				)

				return
			}

			if len(c.expectedRequestURLs) != len(c.expectedBodies) {
				t.Errorf(
					"test case must have the same number of expected request urls (%d) as expected bodies (%d)",
					len(c.expectedRequestURLs),
					len(c.expectedBodies),
				)

				return
			}

			for i := range c.expectedRequestURLs {
				response := &http.Response{
					StatusCode: 200,
					Body:       httpmock.NewRespBodyFromString(c.expectedBodies[i]),
					Header: http.Header{
						"Link": {c.expectedLinkHeaders[i]},
					},
				}

				httpmock.RegisterResponder("GET", c.expectedRequestURLs[i], httpmock.ResponderFromResponse(response))
			}

			risks, err := client.OrderRisk.ListAll(context.Background(), 450789469, nil)
			if !reflect.DeepEqual(risks, c.expectedOrderRisks) {
				t.Errorf("test %d OrderRisk.ListAll orders returned %+v, expected %+v", i, risks, c.expectedOrderRisks)
			}

			if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
				t.Errorf(
					"test %d OrderRisk.ListAll err returned %+v, expected %+v",
					i,
					err,
					c.expectedErr,
				)
			}
		})
	}
}

func TestOrderRiskListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedOrders     []OrderRisk
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"risks": [{"id":1},{"id":2}]}`,
			"",
			[]OrderRisk{{Id: 1}, {Id: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]OrderRisk(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]OrderRisk(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]OrderRisk(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]OrderRisk(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]OrderRisk(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"risks": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]OrderRisk{{Id: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"risks": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]OrderRisk{{Id: 2}},
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

		orderRisks, pagination, err := client.OrderRisk.ListWithPagination(context.Background(), 450789469, nil)
		if !reflect.DeepEqual(orderRisks, c.expectedOrders) {
			t.Errorf("test %d OrderRisk.ListWithPagination OrderRisk returned %+v, expected %+v", i, orderRisks, c.expectedOrders)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d OrderRisk.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d OrderRisk.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestOrderRiskList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("order_risks.json")))

	orderRisks, err := client.OrderRisk.List(context.Background(), 450789469, nil)
	if err != nil {
		t.Errorf("OrderRisk.List returned error: %v", err)
	}

	expected := []OrderRisk{
		{
			Id:              284138680,
			CheckoutId:      0,
			OrderId:         450789469,
			CauseCancel:     true,
			Display:         true,
			MerchantMessage: "This order was placed from a proxy IP",
			Message:         "This order was placed from a proxy IP",
			Score:           "1.0",
			Source:          "External",
			Recommendation:  OrderRecommendationCancel,
		},
		{
			Id:              1029151489,
			CheckoutId:      901414060,
			OrderId:         450789469,
			CauseCancel:     true,
			Display:         true,
			MerchantMessage: "This order came from an anonymous proxy",
			Message:         "This order came from an anonymous proxy",
			Score:           "1.0",
			Source:          "External",
			Recommendation:  OrderRecommendationCancel,
		},
	}

	if !reflect.DeepEqual(orderRisks, expected) {
		t.Errorf("OrderRisks.List returned %+v, expected %+v", orderRisks, expected)
	}
}

func TestOrderRiskListOptions(t *testing.T) {
	setup()
	defer teardown()
	params := map[string]string{
		"fields": "id",
		"limit":  "250",
		"page":   "10",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("order_risks.json")))

	options := OrderRiskListOptions{
		ListOptions: ListOptions{
			Page:   10,
			Limit:  250,
			Fields: "id",
		},
	}

	orderRisks, err := client.OrderRisk.List(context.Background(), 450789469, options)
	if err != nil {
		t.Errorf("OrderRisk.List returned error: %v", err)
	}
	expected := []OrderRisk{
		{
			Id:              284138680,
			CheckoutId:      0,
			OrderId:         450789469,
			CauseCancel:     true,
			Display:         true,
			MerchantMessage: "This order was placed from a proxy IP",
			Message:         "This order was placed from a proxy IP",
			Score:           "1.0",
			Source:          "External",
			Recommendation:  OrderRecommendationCancel,
		},
		{
			Id:              1029151489,
			CheckoutId:      901414060,
			OrderId:         450789469,
			CauseCancel:     true,
			Display:         true,
			MerchantMessage: "This order came from an anonymous proxy",
			Message:         "This order came from an anonymous proxy",
			Score:           "1.0",
			Source:          "External",
			Recommendation:  OrderRecommendationCancel,
		},
	}

	if !reflect.DeepEqual(orderRisks, expected) {
		t.Errorf("OrderRisks.List returned %+v, expected %+v", orderRisks, expected)
	}
}

func TestOrderRiskGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks/284138680.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("order_risk.json")))

	orderRisk, err := client.OrderRisk.Get(context.Background(), 450789469, 284138680, nil)
	if err != nil {
		t.Errorf("OrderRisk.List returned error: %v", err)
	}
	expected := &OrderRisk{
		Id:              284138680,
		CheckoutId:      0,
		OrderId:         450789469,
		CauseCancel:     true,
		Display:         true,
		MerchantMessage: "This order was placed from a proxy IP",
		Message:         "This order was placed from a proxy IP",
		Score:           "1.0",
		Source:          "External",
		Recommendation:  OrderRecommendationCancel,
	}
	if !reflect.DeepEqual(orderRisk, expected) {
		t.Errorf("OrderRisks.Get returned %+v, expected %+v", orderRisk, expected)
	}
}

func TestOrderRiskCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/450789469/risks.json", client.pathPrefix),
		httpmock.NewStringResponder(201, `{"risk":{"id": 1}}`))

	orderRisk := OrderRisk{
		Id: 1,
	}

	o, err := client.OrderRisk.Create(context.Background(), 450789469, orderRisk)
	if err != nil {
		t.Errorf("OrderRisk.Create returned error: %v", err)
	}

	expected := OrderRisk{Id: 1}
	if o.Id != expected.Id {
		t.Errorf("OrderRisk.Create returned id %d, expected %d", o.Id, expected.Id)
	}
}

func TestOrderRiskUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/risks/2.json", client.pathPrefix),
		httpmock.NewStringResponder(201, `{"risk":{"id": 1,"order_id": 2}}`))

	orderRisk := OrderRisk{
		Id:             1,
		OrderId:        2,
		Recommendation: OrderRecommendationAccept,
	}

	o, err := client.OrderRisk.Update(context.Background(), 1, 2, orderRisk)
	if err != nil {
		t.Errorf("Order.Update returned error: %v", err)
	}

	expected := OrderRisk{Id: 1, OrderId: 2, Recommendation: OrderRecommendationAccept}
	if o.Id != expected.Id && o.OrderId != expected.OrderId && o.Recommendation == expected.Recommendation {
		t.Errorf("Order.Update returned id %d, expected %d, expected %d", o.Id, expected.Id, expected.OrderId)
	}
}

func TestOrderRiskDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/orders/1/risks/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.OrderRisk.Delete(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("Order.Delete returned error: %v", err)
	}
}
