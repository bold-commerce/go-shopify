package goshopify

import (
	"context"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAccessScopesServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		"https://fooshop.myshopify.com/admin/oauth/access_scopes.json",
		httpmock.NewBytesResponder(200, loadFixture("access_scopes.json")),
	)

	scopeResponse, err := client.AccessScopes.List(context.Background(), nil)
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
