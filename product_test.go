package goshopify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func productTests(t *testing.T, product Product) {
	// Check that Id is assigned to the returned product
	var expectedInt uint64 = 1071559748
	if product.Id != expectedInt {
		t.Errorf("Product.Id returned %+v, expected %+v", product.Id, expectedInt)
	}
}

func TestProductList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"products": [{"id":1},{"id":2}]}`))

	products, err := client.Product.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Product.List returned error: %v", err)
	}

	expected := []Product{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Product.List returned %+v, expected %+v", products, expected)
	}
}

func TestProductListFilterByIds(t *testing.T) {
	setup()
	defer teardown()

	params := map[string]string{"ids": "1,2,3"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"products": [{"id":1},{"id":2},{"id":3}]}`))

	listOptions := ListOptions{Ids: []uint64{1, 2, 3}}

	products, err := client.Product.List(context.Background(), listOptions)
	if err != nil {
		t.Errorf("Product.List returned error: %v", err)
	}

	expected := []Product{{Id: 1}, {Id: 2}, {Id: 3}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Product.List returned %+v, expected %+v", products, expected)
	}
}

func TestProductListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	products, err := client.Product.List(context.Background(), nil)
	if products != nil {
		t.Errorf("Product.List returned products, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Product.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestProductListAll(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix)

	cases := []struct {
		name                string
		expectedProducts    []Product
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
				`{"products": [{"id":1},{"id":2}]}`,
				`{"products": [{"id":3},{"id":4}]}`,
			},
			expectedProducts: []Product{{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4}},
			expectedErr:      nil,
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
				`{"products": [{"id":1}]}`,
			},
			expectedProducts: []Product{{Id: 1}},
			expectedErr:      nil,
		},
		{
			name: "Returns errors when required",
			expectedRequestURLs: []string{
				listURL,
			},
			expectedLinkHeaders: []string{
				`<http://valid.url?paage_info=pg2>; rel="previous"`,
			},
			expectedBodies: []string{
				`{"products": []}`,
			},
			expectedProducts: []Product{},
			expectedErr:      errors.New("page_info is missing"),
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

			products, err := client.Product.ListAll(context.Background(), nil)
			if !reflect.DeepEqual(products, c.expectedProducts) {
				t.Errorf("test %d Product.ListAll orders returned %+v, expected %+v", i, products, c.expectedProducts)
			}

			if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
				t.Errorf(
					"test %d Product.ListAll err returned %+v, expected %+v",
					i,
					err,
					c.expectedErr,
				)
			}
		})
	}
}

func TestProductListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedProducts   []Product
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"products": [{"id":1},{"id":2}]}`,
			"",
			[]Product{{Id: 1}, {Id: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Product(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Product(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Product(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Product(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Product(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"products": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Product{{Id: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"products": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Product{{Id: 2}},
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

		products, pagination, err := client.Product.ListWithPagination(context.Background(), nil)
		if !reflect.DeepEqual(products, c.expectedProducts) {
			t.Errorf("test %d Product.ListWithPagination products returned %+v, expected %+v", i, products, c.expectedProducts)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Product.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Product.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestProductCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Product.Count(context.Background(), nil)
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Product.Count(context.Background(), CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Product.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Product.Count returned %d, expected %d", cnt, expected)
	}
}

func TestProductGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"product": {"id":1}}`))

	product, err := client.Product.Get(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Product.Get returned error: %v", err)
	}

	expected := &Product{Id: 1}
	if !reflect.DeepEqual(product, expected) {
		t.Errorf("Product.Get returned %+v, expected %+v", product, expected)
	}
}

func TestProductCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("product.json")))

	product := Product{
		Title:       "Burton Custom Freestyle 151",
		BodyHTML:    "<strong>Good snowboard!<\\/strong>",
		Vendor:      "Burton",
		ProductType: "Snowboard",
	}

	returnedProduct, err := client.Product.Create(context.Background(), product)
	if err != nil {
		t.Errorf("Product.Create returned error: %v", err)
	}

	productTests(t, *returnedProduct)
}

func TestProductUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("product.json")))

	product := Product{
		Id:          1,
		ProductType: "Skateboard",
	}

	returnedProduct, err := client.Product.Update(context.Background(), product)
	if err != nil {
		t.Errorf("Product.Update returned error: %v", err)
	}

	productTests(t, *returnedProduct)
}

func TestProductDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Product.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Product.Delete returned error: %v", err)
	}
}

func TestProductListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Product.ListMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Product.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Product.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestProductCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Product.CountMetafields(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Product.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Product.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Product.CountMetafields(context.Background(), 1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Product.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Product.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestProductGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Product.GetMetafield(context.Background(), 1, 2, nil)
	if err != nil {
		t.Errorf("Product.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{Id: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Product.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestProductCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Product.CreateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Product.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestProductUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Id:        2,
		Key:       "app_key",
		Value:     "app_value",
		Type:      MetafieldTypeSingleLineTextField,
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Product.UpdateMetafield(context.Background(), 1, metafield)
	if err != nil {
		t.Errorf("Product.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestProductDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/products/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Product.DeleteMetafield(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("Product.DeleteMetafield() returned error: %v", err)
	}
}
