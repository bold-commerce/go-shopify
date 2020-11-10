package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func productListingTests(t *testing.T, product ProductListing) {
	// Check that ID is assigned to the returned product
	var expectedInt int64 = 921728736
	if product.ID != expectedInt {
		t.Errorf("Product.ID returned %+v, expected %+v", product.ID, expectedInt)
	}
}

func TestProductListingList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"product_listings": [{"product_id":1},{"product_id":2}]}`))

	products, err := client.ProductListing.List(nil)
	if err != nil {
		t.Errorf("ProductListing.List returned error: %v", err)
	}

	expected := []ProductListing{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("ProductListing.List returned %+v, expected %+v", products, expected)
	}
}

func TestProductListingListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	products, err := client.ProductListing.List(nil)
	if products != nil {
		t.Errorf("ProductListing.List returned products, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("ProductListing.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestProductListingListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedProducts   []ProductListing
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"product_listings": [{"product_id":1},{"product_id":2}]}`,
			"",
			[]ProductListing{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]ProductListing(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]ProductListing(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]ProductListing(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]ProductListing(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]ProductListing(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"product_listings": [{"product_id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]ProductListing{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"product_listings": [{"product_id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]ProductListing{{ID: 2}},
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

		products, pagination, err := client.ProductListing.ListWithPagination(nil)
		if !reflect.DeepEqual(products, c.expectedProducts) {
			t.Errorf("test %d ProductListing.ListWithPagination products returned %+v, expected %+v", i, products, c.expectedProducts)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d ProductListing.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d ProductListing.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestProductListingsCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.ProductListing.Count(nil)
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.ProductListing.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}
}

func TestProductListingGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"product_listing": {"product_id":1}}`))

	product, err := client.ProductListing.Get(1, nil)
	if err != nil {
		t.Errorf("ProductListing.Get returned error: %v", err)
	}

	expected := &ProductListing{ID: 1}
	if !reflect.DeepEqual(product, expected) {
		t.Errorf("ProductListing.Get returned %+v, expected %+v", product, expected)
	}
}

func TestProductListingGetProductIDs(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/product_ids.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"product_ids": [1,2,3]}`))

	productIDs, err := client.ProductListing.GetProductIDs(nil)
	if err != nil {
		t.Errorf("ProductListing.Get returned error: %v", err)
	}

	expected := []int64{1, 2, 3}
	if !reflect.DeepEqual(productIDs, expected) {
		t.Errorf("ProductListing.Get returned %+v, expected %+v", productIDs, expected)
	}
}

func TestProductListingPublish(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/921728736.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("product_listing.json")))

	product := Product{
		ID:          921728736,
		ProductType: "Cult Products",
	}

	returnedProduct, err := client.ProductListing.Publish(product.ID)
	fmt.Printf("%+v", returnedProduct)
	if err != nil {
		t.Errorf("ProductListing.Publish returned error: %v", err)
	}

	productListingTests(t, *returnedProduct)
}

func TestProductListingDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/product_listings/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.ProductListing.Delete(1)
	if err != nil {
		t.Errorf("ProductListing.Delete returned error: %v", err)
	}
}
