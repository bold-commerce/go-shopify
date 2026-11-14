package goshopify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

// Exercise public service interfaces against explicit endpoint and response contracts.
func TestServicePagination(t *testing.T) {
	setup()
	defer teardown()
	cases := []struct {
		name, path, key string
		page            func(context.Context, interface{}) (interface{}, *Pagination, error)
		all             func(context.Context, interface{}) (interface{}, error)
	}{
		{name: "AbandonedCheckout.List", path: "checkouts.json", key: "checkouts",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.AbandonedCheckout.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.AbandonedCheckout.ListAll(ctx, options)
			},
		},
		{name: "Article.List", path: "blogs/42/articles.json", key: "articles",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Article.ListWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Article.ListAll(ctx, 42, options)
			},
		},
		{name: "Blog.List", path: "blogs.json", key: "blogs",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Blog.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Blog.ListAll(ctx, options)
			},
		},
		{name: "Collect.List", path: "collects.json", key: "collects",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Collect.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Collect.ListAll(ctx, options)
			},
		},
		{name: "CustomCollection.List", path: "custom_collections.json", key: "custom_collections",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.CustomCollection.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.CustomCollection.ListAll(ctx, options)
			},
		},
		{name: "Customer.Search", path: "customers/search.json", key: "customers",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Customer.SearchWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Customer.SearchAll(ctx, options)
			},
		},
		{name: "CustomerAddress.List", path: "customers/42/addresses.json", key: "addresses",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.CustomerAddress.ListWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.CustomerAddress.ListAll(ctx, 42, options)
			},
		},
		{name: "DiscountCode.List", path: "price_rules/42/discount_codes.json", key: "discount_codes",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.DiscountCode.ListWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.DiscountCode.ListAll(ctx, 42, options)
			},
		},
		{name: "DraftOrder.List", path: "draft_orders.json", key: "draft_orders",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.DraftOrder.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.DraftOrder.ListAll(ctx, options)
			},
		},
		{name: "GiftCard.List", path: "gift_cards.json", key: "gift_cards",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.GiftCard.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.GiftCard.ListAll(ctx, options)
			},
		},
		{name: "InventoryItem.List", path: "inventory_items.json", key: "inventory_items",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.InventoryItem.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.InventoryItem.ListAll(ctx, options)
			},
		},
		{name: "InventoryLevel.List", path: "inventory_levels.json", key: "inventory_levels",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.InventoryLevel.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.InventoryLevel.ListAll(ctx, options)
			},
		},
		{name: "Location.List", path: "locations.json", key: "locations",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Location.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Location.ListAll(ctx, options)
			},
		},
		{name: "Metafield.List", path: "metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Metafield.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Metafield.ListAll(ctx, options)
			},
		},
		{name: "Page.List", path: "pages.json", key: "pages",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Page.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Page.ListAll(ctx, options)
			},
		},
		{name: "PriceRule.List", path: "price_rules.json", key: "price_rules",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.PriceRule.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.PriceRule.ListAll(ctx, options)
			},
		},
		{name: "Redirect.List", path: "redirects.json", key: "redirects",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Redirect.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Redirect.ListAll(ctx, options)
			},
		},
		{name: "ScriptTag.List", path: "script_tags.json", key: "script_tags",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.ScriptTag.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.ScriptTag.ListAll(ctx, options)
			},
		},
		{name: "SmartCollection.List", path: "smart_collections.json", key: "smart_collections",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.SmartCollection.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.SmartCollection.ListAll(ctx, options)
			},
		},
		{name: "Variant.List", path: "products/42/variants.json", key: "variants",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Variant.ListWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Variant.ListAll(ctx, 42, options)
			},
		},
		{name: "Webhook.List", path: "webhooks.json", key: "webhooks",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Webhook.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Webhook.ListAll(ctx, options)
			},
		},
		{name: "Article.ListTags", path: "articles/tags.json", key: "tags",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Article.ListTagsWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Article.ListTagsAll(ctx, options)
			},
		},
		{name: "Article.ListBlogTags", path: "blogs/42/articles/tags.json", key: "tags",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Article.ListBlogTagsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Article.ListBlogTagsAll(ctx, 42, options)
			},
		},
		{name: "ProductListing.GetProductIds", path: "product_listings/product_ids.json", key: "product_ids",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.ProductListing.GetProductIdsWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.ProductListing.GetProductIdsAll(ctx, options)
			},
		},
		{name: "Customer.ListOrders", path: "customers/42/orders.json", key: "orders",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Customer.ListOrdersWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Customer.ListOrdersAll(ctx, 42, options)
			},
		},
		{name: "Collection.ListProducts", path: "collections/42/products.json", key: "products",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Collection.ListProductsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Collection.ListProductsAll(ctx, 42, options)
			},
		},
		{name: "Order.ListFulfillments", path: "orders/42/fulfillments.json", key: "fulfillments",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Order.ListFulfillmentsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Order.ListFulfillmentsAll(ctx, 42, options)
			},
		},
		{name: "Product.List", path: "products.json", key: "products",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Product.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Product.ListAll(ctx, options)
			},
		},
		{name: "Customer.List", path: "customers.json", key: "customers",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Customer.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Customer.ListAll(ctx, options)
			},
		},
		{name: "Order.List", path: "orders.json", key: "orders",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Order.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Order.ListAll(ctx, options)
			},
		},
		{name: "ProductListing.List", path: "product_listings.json", key: "product_listings",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.ProductListing.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.ProductListing.ListAll(ctx, options)
			},
		},
		{name: "Payouts.List", path: "shopify_payments/payouts.json", key: "payouts",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Payouts.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Payouts.ListAll(ctx, options)
			},
		},
		{name: "PaymentsTransactions.List", path: "shopify_payments/balance/transactions.json", key: "transactions",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.PaymentsTransactions.ListWithPagination(ctx, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.PaymentsTransactions.ListAll(ctx, options)
			},
		},
		{name: "OrderRisk.List", path: "orders/42/risks.json", key: "risks",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.OrderRisk.ListWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.OrderRisk.ListAll(ctx, 42, options)
			},
		},
		{name: "DraftOrder.ListMetafields", path: "draft_orders/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.DraftOrder.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.DraftOrder.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Page.ListMetafields", path: "pages/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Page.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Page.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Order.ListMetafields", path: "orders/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Order.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Order.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "SmartCollection.ListMetafields", path: "collections/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.SmartCollection.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.SmartCollection.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Product.ListMetafields", path: "products/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Product.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Product.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "CustomCollection.ListMetafields", path: "collections/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.CustomCollection.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.CustomCollection.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Location.ListMetafields", path: "locations/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Location.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Location.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Variant.ListMetafields", path: "variants/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Variant.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Variant.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Shop.ListMetafields", path: "metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Shop.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Shop.ListMetafieldsAll(ctx, 42, options)
			},
		},
		{name: "Customer.ListMetafields", path: "customers/42/metafields.json", key: "metafields",
			page: func(ctx context.Context, options interface{}) (interface{}, *Pagination, error) {
				return client.Customer.ListMetafieldsWithPagination(ctx, 42, options)
			},
			all: func(ctx context.Context, options interface{}) (interface{}, error) {
				return client.Customer.ListMetafieldsAll(ctx, 42, options)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("https://fooshop.myshopify.com/%s/%s", client.pathPrefix, tc.path)
			first := &ListOptions{Limit: 1, Fields: "id", SinceId: uint64Pointer(7)}
			for _, mode := range []string{"pages", "empty", "http_error", "json_error", "link_error", "cancel"} {
				t.Run(mode, func(t *testing.T) {
					calls := 0
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					httpmock.RegisterResponder("GET", endpoint, func(req *http.Request) (*http.Response, error) {
						calls++
						want := url.Values{"limit": {"1"}, "fields": {"id"}}
						if calls == 1 {
							want.Set("since_id", "7")
						} else {
							want.Set("page_info", "next+cursor,=")
						}
						if !reflect.DeepEqual(req.URL.Query(), want) {
							t.Errorf("query = %v, want %v", req.URL.Query(), want)
						}
						if calls > 2 {
							return nil, fmt.Errorf("pagination did not stop after the last page")
						}
						if calls == 2 {
							switch mode {
							case "http_error":
								return httpmock.NewStringResponse(400, `{"errors":"bad request"}`), nil
							case "cancel":
								return nil, context.Canceled
							case "json_error":
								return httpmock.NewStringResponse(200, `{`), nil
							}
						}
						body := paginationResponse(tc.key, calls)
						if mode == "empty" {
							body = fmt.Sprintf(`{"%s":[]}`, tc.key)
						}
						resp := httpmock.NewStringResponse(200, body)
						if calls == 1 && mode != "empty" {
							resp.Header.Set("Link", `<`+endpoint+`?page_info=next%2Bcursor%2C%3D&limit=1&fields=id>; rel="next"`)
						} else if calls == 2 {
							resp.Header.Set("Link", `<`+endpoint+`?page_info=previous&limit=1&fields=id>; rel="previous"`)
						}
						if calls == 2 && mode == "link_error" {
							resp.Header.Set("Link", "invalid")
						}
						return resp, nil
					})
					items, err := tc.all(ctx, first)
					want := []string{"1", "2"}
					if mode == "empty" {
						want = []string{}
					} else if mode != "pages" {
						want = []string{"1"}
					}
					if (mode == "pages" || mode == "empty") && err != nil {
						t.Fatal(err)
					}
					if mode != "pages" && mode != "empty" && err == nil {
						t.Fatal("expected error")
					}
					if mode == "cancel" && !errors.Is(err, context.Canceled) {
						t.Fatalf("error = %v, want cancellation", err)
					}
					if got := paginationIDs(t, items); !reflect.DeepEqual(got, want) {
						t.Fatalf("items = %v, want %v", got, want)
					}
					if first.PageInfo != "" || first.Limit != 1 || first.Fields != "id" || *first.SinceId != 7 {
						t.Fatal("input options were changed")
					}
				})
			}
			// A caller can use both returned cursors on this exact service method.
			for _, direction := range []string{"next", "previous"} {
				calls := 0
				httpmock.RegisterResponder("GET", endpoint, func(req *http.Request) (*http.Response, error) {
					calls++
					resp := httpmock.NewStringResponse(200, paginationResponse(tc.key, calls))
					if calls == 1 {
						resp.Header.Add("Link", `<`+endpoint+`?page_info=back&limit=1&fields=id>; rel="previous"`)
						resp.Header.Add("Link", `<`+endpoint+`?page_info=forward&limit=1&fields=id>; rel="next"`)
					} else {
						cursor := "forward"
						if direction == "previous" {
							cursor = "back"
						}
						want := url.Values{"page_info": {cursor}, "limit": {"1"}, "fields": {"id"}}
						if !reflect.DeepEqual(req.URL.Query(), want) {
							t.Errorf("query=%v, want %v", req.URL.Query(), want)
						}
					}
					return resp, nil
				})
				_, p, err := tc.page(context.Background(), nil)
				if err != nil {
					t.Fatal(err)
				}
				if p == nil || p.NextPageOptions == nil || p.PreviousPageOptions == nil {
					t.Fatalf("missing cursors: %#v", p)
				}
				next := p.NextPageOptions
				if direction == "previous" {
					next = p.PreviousPageOptions
				}
				items, p, err := tc.page(context.Background(), next)
				if err != nil {
					t.Fatal(err)
				}
				if got := paginationIDs(t, items); !reflect.DeepEqual(got, []string{"2"}) {
					t.Fatalf("page=%v", got)
				}
				if p.NextPageOptions != nil || p.PreviousPageOptions != nil {
					t.Fatal("unexpected continuation")
				}
			}
		})
	}
}

func uint64Pointer(v uint64) *uint64 { return &v }

func paginationResponse(key string, id int) string {
	switch key {
	case "tags":
		return fmt.Sprintf(`{"tags":["%d"]}`, id)
	case "product_ids":
		return fmt.Sprintf(`{"product_ids":[%d]}`, id)
	case "product_listings":
		return fmt.Sprintf(`{"product_listings":[{"product_id":%d}]}`, id)
	case "inventory_levels":
		return fmt.Sprintf(`{"inventory_levels":[{"inventory_item_id":%d}]}`, id)
	default:
		return fmt.Sprintf(`{"%s":[{"id":%d}]}`, key, id)
	}
}

func paginationIDs(t *testing.T, items interface{}) []string {
	t.Helper()
	values := reflect.ValueOf(items)
	ids := []string{}
	for i := 0; i < values.Len(); i++ {
		v := values.Index(i)
		if v.Kind() == reflect.Struct {
			id := v.FieldByName("Id")
			if !id.IsValid() {
				id = v.FieldByName("InventoryItemId")
			}
			if !id.IsValid() {
				t.Fatalf("unrecognized item type %s", v.Type())
			}
			v = id
		}
		ids = append(ids, fmt.Sprint(v.Interface()))
	}
	return ids
}

// Cancel a real in-flight second page and retain only the completed first page.
func TestPaginationContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stopped := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page_info") == "next" {
			cancel()
			<-r.Context().Done()
			close(stopped)
			return
		}
		w.Header().Set("Link", `<https://shop.myshopify.com/admin/gift_cards.json?page_info=next>; rel="next"`)
		fmt.Fprint(w, `{"gift_cards":[{"id":1}]}`)
	}))
	defer server.Close()
	c, err := NewClient(App{}, "shop", "token", WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	c.baseURL, err = url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	items, err := c.GiftCard.ListAll(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want cancellation", err)
	}
	if len(items) != 1 || items[0].Id != 1 {
		t.Fatalf("completed pages = %#v", items)
	}
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("HTTP request was not cancelled")
	}
}
