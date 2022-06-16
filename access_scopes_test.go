package goshopify

import (
	"testing"
	"fmt"
	"reflect"

	"github.com/jarcoal/httpmock"
)

func TestAccessScopesServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/oauth/access_scopes.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("access_scopes.json")),
	)

	scopeResponse, err := client.AccessScopes.List(nil)
	if err != nil {
		t.Errorf("AccessScopes.List returned an error: %v", err)
	}

	expected := []AccessScope{
		{
			Handle: "scope_1",
		},
		{
			Handle: "scope_2",
		},
	}
	if !reflect.DeepEqual(scopeResponse, expected) {
		t.Errorf("AccessScopes.List returned %+v, expected %+v", expected, expected)
	}
}
