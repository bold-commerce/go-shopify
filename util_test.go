package goshopify

import (
	"net/url"
	"testing"
	"time"
)

func TestShopFullName(t *testing.T) {
	cases := []struct {
		in, expected string
	}{
		{"myshop", "myshop.myshopify.com"},
		{"myshop.", "myshop.myshopify.com"},
		{" myshop", "myshop.myshopify.com"},
		{"myshop ", "myshop.myshopify.com"},
		{"myshop \n", "myshop.myshopify.com"},
		{"myshop.myshopify.com", "myshop.myshopify.com"},
	}

	for _, c := range cases {
		actual := ShopFullName(c.in)
		if actual != c.expected {
			t.Errorf("ShopFullName(%s): expected %s, actual %s", c.in, c.expected, actual)
		}
	}
}

func TestShopShortName(t *testing.T) {
	cases := []struct {
		in, expected string
	}{
		{"myshop", "myshop"},
		{"myshop.", "myshop"},
		{" myshop", "myshop"},
		{"myshop ", "myshop"},
		{"myshop \n", "myshop"},
		{"myshop.myshopify.com", "myshop"},
		{".myshop.myshopify.com.", "myshop"},
	}

	for _, c := range cases {
		actual := ShopShortName(c.in)
		if actual != c.expected {
			t.Errorf("ShopShortName(%s): expected %s, actual %s", c.in, c.expected, actual)
		}
	}
}

func TestShopBaseUrl(t *testing.T) {
	cases := []struct {
		in, expected string
	}{
		{"myshop", "https://myshop.myshopify.com"},
		{"myshop.", "https://myshop.myshopify.com"},
		{" myshop", "https://myshop.myshopify.com"},
		{"myshop ", "https://myshop.myshopify.com"},
		{"myshop \n", "https://myshop.myshopify.com"},
		{"myshop.myshopify.com", "https://myshop.myshopify.com"},
	}

	for _, c := range cases {
		actual := ShopBaseUrl(c.in)
		if actual != c.expected {
			t.Errorf("ShopBaseUrl(%s): expected %s, actual %s", c.in, c.expected, actual)
		}
	}
}

func TestMetafieldPathPrefix(t *testing.T) {
	cases := []struct {
		resource   string
		resourceID int64
		expected   string
	}{
		{"", 0, "metafields"},
		{"products", 123, "products/123/metafields"},
	}

	for _, c := range cases {
		actual := MetafieldPathPrefix(c.resource, c.resourceID)
		if actual != c.expected {
			t.Errorf("MetafieldPathPrefix(%s, %d): expected %s, actual %s", c.resource, c.resourceID, c.expected, actual)
		}
	}
}

func TestFulfillmentPathPrefix(t *testing.T) {
	cases := []struct {
		resource   string
		resourceID int64
		expected   string
	}{
		{"", 0, "fulfillments"},
		{"orders", 123, "orders/123/fulfillments"},
	}

	for _, c := range cases {
		actual := FulfillmentPathPrefix(c.resource, c.resourceID)
		if actual != c.expected {
			t.Errorf("FulfillmentPathPrefix(%s, %d): expected %s, actual %s", c.resource, c.resourceID, c.expected, actual)
		}
	}
}

func TestOnlyDateMarshal(t *testing.T) {
	cases := []struct {
		in       OnlyDate
		expected string
	}{
		{OnlyDate{time.Date(2023, 03, 31, 0, 0, 0, 0, time.Local)}, "\"2023-03-31\""},
		{OnlyDate{}, "\"0001-01-01\""},
	}

	for _, c := range cases {
		actual, _ := c.in.MarshalJSON()
		if string(actual) != c.expected {
			t.Errorf("MarshalJSON(%s): expected %s, actual %s", c.in.String(), c.expected, string(actual))
		}
	}
}

func TestOnlyDateUnmarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected OnlyDate
	}{
		{"\"2023-03-31\"", OnlyDate{time.Date(2023, 03, 31, 0, 0, 0, 0, time.Local)}},
		{"\"0001-01-01\"", OnlyDate{}},
		{"\"\"", OnlyDate{}},
	}

	for _, c := range cases {
		newDate := OnlyDate{}
		_ = newDate.UnmarshalJSON([]byte(c.in))
		if newDate.String() != c.expected.String() {
			t.Errorf("UnmarshalJSON(%s): expected %s, actual %s", c.in, newDate.String(), c.expected.String())
		}
	}
}

func TestOnlyDateEncode(t *testing.T) {
	cases := []struct {
		in       OnlyDate
		expected string
	}{
		{OnlyDate{time.Date(2023, 03, 31, 0, 0, 0, 0, time.Local)}, "\"2023-03-31\""},
		{OnlyDate{}, "\"0001-01-01\""},
	}

	for _, c := range cases {
		urlVal := url.Values{}
		_ = c.in.EncodeValues("date", &urlVal)
		if urlVal.Get("date") != c.expected {
			t.Errorf("EncodeValues(%s): expected %s, actual %s", c.in.String(), c.expected, urlVal.Get("date"))
		}
	}
}
