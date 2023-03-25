package goshopify

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestCollectionGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d.json", client.pathPrefix, 1),
		httpmock.NewStringResponder(200,
			`{
				"collection": {
					"id": 25,
					"handle": "more-than-5",
					"title": "More than $5",
					"updated_at": "2020-07-23T15:12:12-04:00",
					"body_html": "<p>Items over $5</p>",
					"published_at": "2020-06-23T14:22:47-04:00",
					"sort_order": "best-selling",
					"template_suffix": "custom",
					"collection_type": "smart",
					"published_scope": "web",
					"image": {
						"created_at": "2020-02-27T15:01:45-05:00",
						"alt": null,
						"width": 1920,
						"height": 1279,
						"src": "https://example/image.jpg"
					}
				}
			}`))

	collection, err := client.Collection.Get(1, nil)
	if err != nil {
		t.Errorf("Collection.Get returned error: %v", err)
	}

	updatedAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:12:12-04:00")
	publishedAt, _ := time.Parse(time.RFC3339, "2020-06-23T14:22:47-04:00")

	imageCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-27T15:01:45-05:00")
	expected := &Collection{
		ID:             25,
		Handle:         "more-than-5",
		Title:          "More than $5",
		UpdatedAt:      &updatedAt,
		BodyHTML:       "<p>Items over $5</p>",
		SortOrder:      "best-selling",
		TemplateSuffix: "custom",
		PublishedAt:    &publishedAt,
		PublishedScope: "web",
		Image: Image{
			CreatedAt: &imageCreatedAt,
			Width:     1920,
			Height:    1279,
			Src:       "https://example/image.jpg",
		},
	}
	if !reflect.DeepEqual(collection, expected) {
		t.Errorf("Collection.Get returned %+v, expected %+v", collection, expected)
	}
}

func TestCollectionListProducts(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d/products.json", client.pathPrefix, 1),
		httpmock.NewStringResponder(200,
			`{
				"products": [
					{
						"id": 632910392,
						"title": "The Best Product",
						"body_html": "<p>The best product available</p>",
						"vendor": "local-vendor",
						"product_type": "Best Products",
						"created_at": "2020-07-23T15:12:10-04:00",
						"handle": "the-best-product",
						"updated_at": "2020-07-23T15:13:26-04:00",
						"published_at": "2020-07-23T15:12:11-04:00",
						"template_suffix": "special",
						"published_scope": "web",
						"tags": "Best",
						"admin_graphql_api_id": "gid://shopify/Location/4688969785",
						"options": [
							{
								"id": 6519940513924,
								"product_id": 632910392,
								"name": "Title",
								"position": 1
							}
						],
						"images": [
							{
								"id": 14601766043780,
								"product_id": 632910392,
								"position": 1,
								"created_at": "2020-02-27T13:21:52-05:00",
								"updated_at": "2020-02-27T13:21:52-05:00",
								"alt": null,
								"width": 480,
								"height": 720,
								"src": "https://example/image.jpg",
								"variant_ids": [
									32434329944196,
									32434531893380
								]
							}
						],
						"image": null
					}
				]
			}`))

	products, err := client.Collection.ListProducts(1, nil)
	if err != nil {
		t.Errorf("Collection.ListProducts returned error: %v", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:12:10-04:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:13:26-04:00")
	publishedAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:12:11-04:00")
	imageCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-27T13:21:52-05:00")
	imageUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-27T13:21:52-05:00")

	expected := []Product{
		{
			ID:             632910392,
			Title:          "The Best Product",
			BodyHTML:       "<p>The best product available</p>",
			Vendor:         "local-vendor",
			ProductType:    "Best Products",
			Handle:         "the-best-product",
			CreatedAt:      &createdAt,
			UpdatedAt:      &updatedAt,
			PublishedAt:    &publishedAt,
			PublishedScope: "web",
			Tags:           "Best",
			Options: []ProductOption{
				{
					ID:        6519940513924,
					ProductID: 632910392,
					Name:      "Title",
					Position:  1,
					Values:    nil,
				},
			},
			Variants: nil,
			Images: []Image{
				{
					ID:         14601766043780,
					ProductID:  632910392,
					Position:   1,
					CreatedAt:  &imageCreatedAt,
					UpdatedAt:  &imageUpdatedAt,
					Width:      480,
					Height:     720,
					Src:        "https://example/image.jpg",
					VariantIds: []int64{32434329944196, 32434531893380},
				},
			},
			TemplateSuffix:    "special",
			AdminGraphqlAPIID: "gid://shopify/Location/4688969785",
		},
	}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Collection.ListProducts returned %+v, expected %+v", products, expected)
	}
}
func TestCollectionListProductsError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d/products.json", client.pathPrefix, 1),
		httpmock.NewStringResponder(200,
			`{
				"products": [
					{
						some invalid json
			}`))

	products, err := client.Collection.ListProducts(1, nil)

	if len(products) > 0 {
		t.Errorf("Collection.ListProducts returned products %v, expected no products to be returned", products)
	}

	expectedError := fmt.Errorf("invalid character 's' looking for beginning of object key string")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Collection.ListProducts err returned %v, expected %v", err, expectedError)
	}
}
func TestListProductsWithPagination(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d/products.json", client.pathPrefix, 1),
		httpmock.ResponderFromResponse(&http.Response{
			StatusCode: 200,
			Body: httpmock.NewRespBodyFromString(`{
				"products": [
					{
						"id": 632910392,
						"title": "The Best Product",
						"body_html": "<p>The best product available</p>",
						"vendor": "local-vendor",
						"product_type": "Best Products",
						"created_at": "2020-07-23T15:12:10-04:00",
						"handle": "the-best-product",
						"updated_at": "2020-07-23T15:13:26-04:00",
						"published_at": "2020-07-23T15:12:11-04:00",
						"template_suffix": "special",
						"published_scope": "web",
						"tags": "Best",
						"admin_graphql_api_id": "gid://shopify/Location/4688969785",
						"options": [
							{
								"id": 6519940513924,
								"product_id": 632910392,
								"name": "Title",
								"position": 1
							}
						],
						"images": [
							{
								"id": 14601766043780,
								"product_id": 632910392,
								"position": 1,
								"created_at": "2020-02-27T13:21:52-05:00",
								"updated_at": "2020-02-27T13:21:52-05:00",
								"alt": null,
								"width": 480,
								"height": 720,
								"src": "https://example/image.jpg",
								"variant_ids": [
									32434329944196,
									32434531893380
								]
							}
						],
						"image": null
					}
				]
			}`),
			Header: http.Header{
				"Link": {`<http://valid.url?limit=1&page_info=pageInfoCode>; rel="next"`},
			},
		}))

	products, page, err := client.Collection.ListProductsWithPagination(1, nil)
	if err != nil {
		t.Errorf("Collection.ListProductsWithPagination returned error: %v", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:12:10-04:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:13:26-04:00")
	publishedAt, _ := time.Parse(time.RFC3339, "2020-07-23T15:12:11-04:00")
	imageCreatedAt, _ := time.Parse(time.RFC3339, "2020-02-27T13:21:52-05:00")
	imageUpdatedAt, _ := time.Parse(time.RFC3339, "2020-02-27T13:21:52-05:00")

	expectedProducts := []Product{
		{
			ID:             632910392,
			Title:          "The Best Product",
			BodyHTML:       "<p>The best product available</p>",
			Vendor:         "local-vendor",
			ProductType:    "Best Products",
			Handle:         "the-best-product",
			CreatedAt:      &createdAt,
			UpdatedAt:      &updatedAt,
			PublishedAt:    &publishedAt,
			PublishedScope: "web",
			Tags:           "Best",
			Options: []ProductOption{
				{
					ID:        6519940513924,
					ProductID: 632910392,
					Name:      "Title",
					Position:  1,
					Values:    nil,
				},
			},
			Variants: nil,
			Images: []Image{
				{
					ID:         14601766043780,
					ProductID:  632910392,
					Position:   1,
					CreatedAt:  &imageCreatedAt,
					UpdatedAt:  &imageUpdatedAt,
					Width:      480,
					Height:     720,
					Src:        "https://example/image.jpg",
					VariantIds: []int64{32434329944196, 32434531893380},
				},
			},
			TemplateSuffix:    "special",
			AdminGraphqlAPIID: "gid://shopify/Location/4688969785",
		},
	}
	if !reflect.DeepEqual(products, expectedProducts) {
		t.Errorf("Collection.ListProductsWithPagination returned %+v, expected %+v", products, expectedProducts)
	}

	expectedPage := &Pagination{
		NextPageOptions: &ListOptions{
			PageInfo:     "pageInfoCode",
			Page:         0,
			Limit:        1,
			SinceID:      0,
			CreatedAtMin: time.Time{},
			CreatedAtMax: time.Time{},
			UpdatedAtMin: time.Time{},
			UpdatedAtMax: time.Time{},
			Order:        "",
			Fields:       "",
			Vendor:       "",
			IDs:          nil,
		},
		PreviousPageOptions: nil,
	}
	if !reflect.DeepEqual(page, expectedPage) {
		t.Errorf("Collection.ListProductsWithPagination returned %+v, expected %+v", page, expectedPage)
	}
}

func TestCollectionListProductsWithPaginationRequestError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d/products.json", client.pathPrefix, 1),
		httpmock.NewStringResponder(200,
			`{
				"products": [
					{
						some invalid json
			}`))

	products, pagination, err := client.Collection.ListProductsWithPagination(1, nil)

	if len(products) > 0 {
		t.Errorf("Collection.ListProductsWithPagination returned products %v, expected no products to be returned", products)
	}

	if pagination != nil {
		t.Errorf("Collection.ListProductsWithPagination returned pagination %v, expected nil", products)
	}

	expectedError := fmt.Errorf("invalid character 's' looking for beginning of object key string")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Collection.ListProductsWithPagination err returned %v, expected %v", err, expectedError)
	}
}

func TestCollectionListProductsWithPaginationExtractionError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/collections/%d/products.json", client.pathPrefix, 1),
		httpmock.ResponderFromResponse(&http.Response{
			StatusCode: 200,
			Body: httpmock.NewRespBodyFromString(`{
				"products": []
			}`),
			Header: http.Header{
				"Link": {`invalid link`},
			},
		}))

	products, pagination, err := client.Collection.ListProductsWithPagination(1, nil)
	if len(products) > 0 {
		t.Errorf("Collection.ListProductsWithPagination returned products %v, expected no products to be returned", products)
	}

	if pagination != nil {
		t.Errorf("Collection.ListProductsWithPagination returned pagination %v, expected nil", products)
	}

	expectedError := fmt.Errorf("could not extract pagination link header")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Collection.ListProductsWithPagination err returned %v, expected %v", err, expectedError)
	}
}
