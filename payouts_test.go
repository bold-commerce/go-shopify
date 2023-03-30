package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestPayoutList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shopify_payments/payouts.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"payouts": [{"id":1},{"id":2}]}`))

	payouts, err := client.Payout.List(nil)
	if err != nil {
		t.Errorf("Payouts.List returned error: %v", err)
	}

	expected := []Payout{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(payouts, expected) {
		t.Errorf("Payout.List returned %+v, expected %+v", payouts, expected)
	}
}

func TestPayoutListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shopify_payments/payouts.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	payouts, err := client.Payout.List(nil)
	if payouts != nil {
		t.Errorf("Payout.List returned payouts, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Payout.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestPayoutListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/shopify_payments/payouts.json", client.pathPrefix)

	cases := []struct {
		body               string
		linkHeader         string
		expectedPayouts    []Payout
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"payouts": [{"id":1},{"id":2}]}`,
			"",
			[]Payout{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Payout(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Payout(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Payout(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Payout(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Payout(nil),
			nil,
			errors.New(`strconv.Atoi: parsing "invalid": invalid syntax`),
		},
		// Valid link header responses
		{
			`{"payouts": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Payout{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"payouts": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Payout{{ID: 2}},
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

		payouts, pagination, err := client.Payout.ListWithPagination(nil)
		if !reflect.DeepEqual(payouts, c.expectedPayouts) {
			t.Errorf("test %d Payout.ListWithPagination payouts returned %+v, expected %+v", i, payouts, c.expectedPayouts)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Payout.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Payout.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestPayoutGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shopify_payments/payouts/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"payout": {"id":1}}`))

	payout, err := client.Payout.Get(1, nil)
	if err != nil {
		t.Errorf("Payout.Get returned error: %v", err)
	}

	expected := &Payout{ID: 1}
	if !reflect.DeepEqual(payout, expected) {
		t.Errorf("Payout.Get returned %+v, expected %+v", payout, expected)
	}
}
