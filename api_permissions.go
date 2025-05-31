package goshopify

import (
	"context"
	"fmt"
)

const apiPermissionsBasePath = "api_permissions"

// ApiPermissionsService is an interface for interfacing with the API
// permissions endpoints of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/latest/resources/theme
type ApiPermissionsService interface {
	Delete(context.Context) error
}

// ApiPermissionsServiceOp handles communication with the theme related methods of
// the Shopify API.
type ApiPermissionsServiceOp struct {
	client *Client
}

// Uninstall an app.
func (s *ApiPermissionsServiceOp) Delete(ctx context.Context) error {
	path := fmt.Sprintf("%s/current.json", apiPermissionsBasePath)
	return s.client.Delete(ctx, path)
}
