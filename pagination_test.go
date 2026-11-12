package goshopify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestPaginationFields(t *testing.T) {
	for _, fields := range []string{"id,title", "id%2Ctitle"} {
		for _, nextFirst := range []bool{false, true} {
			next := fmt.Sprintf(`<https://shop.myshopify.com/admin/products.json?page_info=a%%2Bb%%3D&limit=2&fields=%s>; rel="next"`, fields)
			previous := fmt.Sprintf(`<https://shop.myshopify.com/admin/products.json?fields=%s&limit=5&page_info=previous>; rel="previous"`, fields)
			header := previous + ", " + next
			if nextFirst {
				header = next + ", " + previous
			}
			got, err := extractPagination(header)
			if err != nil {
				t.Errorf("fields=%s nextFirst=%v: %v", fields, nextFirst, err)
				continue
			}
			want := &Pagination{
				NextPageOptions:     &ListOptions{PageInfo: "a+b=", Limit: 2, Fields: "id,title"},
				PreviousPageOptions: &ListOptions{PageInfo: "previous", Limit: 5, Fields: "id,title"},
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("fields=%s nextFirst=%v: got %#v, want %#v", fields, nextFirst, got, want)
			}
		}
	}
}

func TestPaginationMalformedLinks(t *testing.T) {
	valid := `<https://shop.myshopify.com/admin/products.json?page_info=next&fields=id,title>; rel="next"`
	for _, header := range []string{
		valid + ",", "," + valid, valid + ", invalid", valid + ",," + valid,
		`<https://shop.myshopify.com/admin/products.json?page_info=next&fields=id,title; rel="next"`,
	} {
		if got, err := extractPagination(header); err == nil || got != nil {
			t.Errorf("malformed header %q: got %#v, %v", header, got, err)
		}
	}
}

func TestProductListAllPreservesFields(t *testing.T) {
	setup()
	defer teardown()
	for _, fields := range []string{"id,title", "id%2Ctitle"} {
		requestURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/products.json", client.pathPrefix)
		calls := 0
		httpmock.RegisterResponder("GET", requestURL, func(req *http.Request) (*http.Response, error) {
			calls++
			want := url.Values{"limit": {"1"}, "fields": {"id,title"}}
			if calls == 1 {
				want.Set("vendor", "Acme")
			} else {
				want.Set("page_info", "a+b=")
			}
			if !reflect.DeepEqual(req.URL.Query(), want) {
				t.Errorf("page %d query = %v, want %v", calls, req.URL.Query(), want)
			}
			response := httpmock.NewStringResponse(200, fmt.Sprintf(`{"products":[{"id":%d,"title":"Product"}]}`, calls))
			if calls == 1 {
				response.Header.Set("Link", fmt.Sprintf(`<%s?page_info=a%%2Bb%%3D&limit=1&fields=%s>; rel="next"`, requestURL, fields))
			}
			return response, nil
		})
		products, err := client.Product.ListAll(context.Background(), ListOptions{Limit: 1, Fields: "id,title", Vendor: "Acme"})
		if err != nil {
			t.Errorf("fields=%s: %v", fields, err)
			continue
		}
		if calls != 2 || len(products) != 2 || products[0].Id != 1 || products[1].Id != 2 {
			t.Errorf("fields=%s: calls=%d products=%#v", fields, calls, products)
		}
	}
}
