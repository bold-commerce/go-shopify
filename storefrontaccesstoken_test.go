package goshopify

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func storefrontAccessTokenTests(t *testing.T, StorefrontAccessToken StorefrontAccessToken) {
	expectedStr := "API Client Extension"
	if StorefrontAccessToken.Title != expectedStr {
		t.Errorf("StorefrontAccessToken.Title returned %+v, expected %+v", StorefrontAccessToken.Title, expectedStr)
	}

	expectedStr = "378d95641257a4ab3feff967ee234f4d"
	if StorefrontAccessToken.AccessToken != expectedStr {
		t.Errorf("StorefrontAccessToken.AccessToken returned %+v, expected %+v", StorefrontAccessToken.AccessToken, expectedStr)
	}

	expectedStr = "unauthenticated_read_product_listings"
	if StorefrontAccessToken.AccessScope != expectedStr {
		t.Errorf("StorefrontAccessToken.AccessScope returned %+v, expected %+v", StorefrontAccessToken.AccessScope, expectedStr)
	}

	expectedStr = "gid://shopify/StorefrontAccessToken/755357713"
	if StorefrontAccessToken.AdminGraphqlAPIID != expectedStr {
		t.Errorf("StorefrontAccessToken.AdminGraphqlAPIID returned %+v, expected %+v", StorefrontAccessToken.AdminGraphqlAPIID, expectedStr)
	}

	d := time.Date(2016, time.June, 1, 14, 10, 44, 0, time.UTC)
	if !d.Equal(*StorefrontAccessToken.CreatedAt) {
		t.Errorf("StorefrontAccessToken.CreatedAt returned %+v, expected %+v", StorefrontAccessToken.CreatedAt, d)
	}
}

func TestStorefrontAccessTokenList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/storefront_access_tokens.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("storefront_access_tokens.json")))

	storefrontAccessTokens, err := client.StorefrontAccessToken.List(nil)
	if err != nil {
		t.Errorf("StorefrontAccessToken.List returned error: %v", err)
	}

	if len(storefrontAccessTokens) != 1 {
		t.Errorf("StorefrontAccessToken.List got %v storefront access tokens, expected: 1", len(storefrontAccessTokens))
	}

	storefrontAccessTokenTests(t, storefrontAccessTokens[0])
}

func TestStorefrontAccessTokenCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/storefront_access_tokens.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("storefront_access_token.json")))

	storefrontAccessToken := StorefrontAccessToken{
		Title: "API Client Extension",
	}

	returnedStorefrontAccessToken, err := client.StorefrontAccessToken.Create(storefrontAccessToken)
	if err != nil {
		t.Errorf("StorefrontAccessToken.Create returned error: %v", err)
	}

	storefrontAccessTokenTests(t, *returnedStorefrontAccessToken)
}

func TestStorefrontAccessTokenDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/storefront_access_tokens/755357713.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.StorefrontAccessToken.Delete(755357713)
	if err != nil {
		t.Errorf("StorefrontAccessToken.Delete returned error: %v", err)
	}
}
