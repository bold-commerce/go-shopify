package goshopify

import (
	"fmt"
	"strings"
)

// Return the full shop name, including .myshopify.com
func ShopFullName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if strings.Contains(name, "myshopify.com") {
		return name
	}
	return name + ".myshopify.com"
}

// Return the short shop name, excluding .myshopify.com
func ShopShortName(name string) string {
	// Convert to fullname and remove the myshopify part. Perhaps not the most
	// performant solution, but then we don't have to repeat all the trims here
	// :-)
	return strings.Replace(ShopFullName(name), ".myshopify.com", "", -1)
}

// Return the Shop's base url.
func ShopBaseUrl(name string) string {
	name = ShopFullName(name)
	return fmt.Sprintf("https://%s", name)
}

// Return the prefix for an order path
func OrderPathPrefix(prefix string, resource string, resourceID int64) string {
	if resource != "" {
		return fmt.Sprintf("%s/%d/%s", resource, resourceID, prefix)
	}

	return prefix
}

// Return the prefix for a metafield path
func MetafieldPathPrefix(resource string, resourceID int64) string {
	return OrderPathPrefix("metafields", resource, resourceID)
}

// Return the prefix for a fulfillment path
func FulfillmentPathPrefix(resource string, resourceID int64) string {
	return OrderPathPrefix("fulfillments", resource, resourceID)
}

// Return the prefix for a fulfillment order path
func FulfillmentOrderPathPrefix(resource string, resourceID int64) string {
	return OrderPathPrefix("fulfillment_orders", resource, resourceID)
}
