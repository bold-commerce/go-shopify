package goshopify

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAssetScopesList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/oauth/access_scopes.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"access_scopes":[{"handle":"read_products"},{"handle":"write_orders"}]}`))

	scopes, err := client.AccessScope.List(nil)
	if err != nil {
		t.Errorf("AccessScope.List returned error: %v", err)
	}

	expected := []AccessScope{{Handle: "read_products"}, {Handle: "write_orders"}}
	if !reflect.DeepEqual(scopes, expected) {
		t.Errorf("AccessScope.List returned %+v, expected %+v", scopes, expected)
	}
}
