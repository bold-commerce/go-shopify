package goshopify

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestApiPermissionsDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/current.json", client.pathPrefix, apiPermissionsBasePath),
		httpmock.NewStringResponder(200, ""))

	err := client.ApiPermissions.Delete(context.Background())
	if err != nil {
		t.Errorf("Theme.Delete returned error: %v", err)
	}
}
