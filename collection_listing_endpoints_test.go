package goshopify

import (
	"context"
	"testing"
)

func TestCollectionListingAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "collection_listings.json",
			response: `{"collection_listings":[{"collection_id":7}]}`,
			want:     []CollectionListing{{CollectionId: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CollectionListing.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "collection_listings/7.json",
			response: `{"collection_listing":{"collection_id":7}}`,
			want:     &CollectionListing{CollectionId: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CollectionListing.Get(ctx, 7, nil)
			},
		},
		{
			name:   "Delete",
			method: "DELETE",
			path:   "collection_listings/7.json",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.CollectionListing.Delete(ctx, 7)
			},
		},
		{
			name:     "Publish",
			method:   "PUT",
			path:     "collection_listings/7.json",
			request:  `{"collection_listing":{"collection_id":7}}`,
			response: `{"collection_listing":{"collection_id":7}}`,
			want:     &CollectionListing{CollectionId: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CollectionListing.Publish(ctx, 7)
			},
		},
		{
			name:     "GetProductIds",
			method:   "GET",
			path:     "collection_listings/7/product_ids.json",
			response: `{"product_ids":[11,12]}`,
			want:     []uint64{11, 12},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CollectionListing.GetProductIds(ctx, 7, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestCollectionListingListWithPagination(t *testing.T) {
	testFeaturePagination(t, "collection_listings.json", `{"collection_listings":[{"collection_id":7}]}`, []CollectionListing{{CollectionId: 7}}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.CollectionListing.ListWithPagination(ctx, options)
	})
}

func TestCollectionListingGetProductIdsWithPagination(t *testing.T) {
	testFeaturePagination(t, "collection_listings/7/product_ids.json", `{"product_ids":[7]}`, []uint64{7}, func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
		return client.CollectionListing.GetProductIdsWithPagination(ctx, 7, options)
	})
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/collectionlisting
func TestCollectionListingDocumentedResponse(t *testing.T) {
	var value CollectionListingResource
	testFeatureResponse(t, "fixtures/collection_listing_documented.json", &value, true)
	if !(value.CollectionListing.CollectionId == 482865238 && value.CollectionListing.Image.Src != "") {
		t.Fatalf("incomplete response: %#v", value)
	}
}
