package goshopify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
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

// Return the prefix for a metafield path
func MetafieldPathPrefix(resource string, resourceId uint64) string {
	prefix := "metafields"
	if resource != "" {
		prefix = fmt.Sprintf("%s/%d/metafields", resource, resourceId)
	}
	return prefix
}

// Return the prefix for a fulfillment path
func FulfillmentPathPrefix(resource string, resourceId uint64) string {
	prefix := "fulfillments"
	if resource != "" {
		prefix = fmt.Sprintf("%s/%d/fulfillments", resource, resourceId)
	}
	return prefix
}

type OnlyDate struct {
	time.Time
}

func (c *OnlyDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		*c = OnlyDate{time.Time{}}
		return nil
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*c = OnlyDate{t}
	return nil
}

func (c OnlyDate) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}

// EncodeValues encodes the date as YYYY-MM-DD without JSON quotes.
func (c *OnlyDate) EncodeValues(key string, v *url.Values) error {
	v.Add(key, c.Format("2006-01-02"))
	return nil
}

func (c *OnlyDate) String() string {
	return `"` + c.Format("2006-01-02") + `"`
}

func TimePtr(v time.Time) *time.Time {
	return &v
}

// BoolPtr returns a pointer for an optional boolean request field.
func BoolPtr(v bool) *bool { return &v }

// stringOrNumber preserves the text of fields Shopify can send as a string or a JSON number.
type stringOrNumber string

func (s *stringOrNumber) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		return json.Unmarshal(data, (*string)(s))
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*s = stringOrNumber(n)
	return nil
}

// IntPtr returns a pointer to an integer, including zero.
func IntPtr(v int) *int { return &v }
