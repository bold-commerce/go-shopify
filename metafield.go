package goshopify

import (
	"context"
	"fmt"
	"time"
)

// MetafieldService is an interface for interfacing with the metafield endpoints
// of the Shopify API.
// https://help.shopify.com/api/reference/metafield
type MetafieldService interface {
	List(context.Context, interface{}) ([]Metafield, error)
	Count(context.Context, interface{}) (int, error)
	Get(context.Context, uint64, interface{}) (*Metafield, error)
	Create(context.Context, Metafield) (*Metafield, error)
	Update(context.Context, Metafield) (*Metafield, error)
	Delete(context.Context, uint64) error
}

// MetafieldsService is an interface for other Shopify resources
// to interface with the metafield endpoints of the Shopify API.
// https://help.shopify.com/api/reference/metafield
type MetafieldsService interface {
	ListMetafields(context.Context, uint64, interface{}) ([]Metafield, error)
	CountMetafields(context.Context, uint64, interface{}) (int, error)
	GetMetafield(context.Context, uint64, uint64, interface{}) (*Metafield, error)
	CreateMetafield(context.Context, uint64, Metafield) (*Metafield, error)
	UpdateMetafield(context.Context, uint64, Metafield) (*Metafield, error)
	DeleteMetafield(context.Context, uint64, uint64) error
}

// MetafieldServiceOp handles communication with the metafield
// related methods of the Shopify API.
type MetafieldServiceOp struct {
	client     *Client
	resource   string
	resourceId uint64
}

// MetafieldType
//
// Docs: https://shopify.dev/docs/api/admin-rest/2023-07/resources/metafield#resource-object
type MetafieldType string

// https://shopify.dev/docs/api/admin-rest/2023-07/resources/metafield#resource-object
const (
	// MetafieldTypeBoolean True or false.
	MetafieldTypeBoolean MetafieldType = "boolean"

	// MetafieldTypeColor A hexidecimal color code, #fff123.
	MetafieldTypeColor MetafieldType = "color" // #fff123

	// MetafieldTypeDate ISO601, YYYY-MM-DD.
	MetafieldTypeDate MetafieldType = "date"

	// MetafieldTypeDatetime ISO8601, YYYY-MM-DDTHH:MM:SS.
	MetafieldTypeDatetime MetafieldType = "date_time"

	// MetafieldTypeDimension JSON, {"value:" 25.0, "unit": "cm"}.
	MetafieldTypeDimension MetafieldType = "dimension"

	// MetafieldTypeJSON {"ingredient": "flour", "amount": 0.3}.
	MetafieldTypeJSON MetafieldType = "json"

	// MetafieldTypeMoney JSON, {"amount": 5.99, "currency_code": "CAD"}.
	MetafieldTypeMoney MetafieldType = "money"

	// MetafieldTypeMultiLineTextField lines of text separated with newline characters.
	MetafieldTypeMultiLineTextField MetafieldType = "multi_line_text_field"

	// MetafieldTypeNumberDecimal 10.4.
	MetafieldTypeNumberDecimal MetafieldType = "number_decimal"

	// MetafieldTypeNumberInteger 10.
	MetafieldTypeNumberInteger MetafieldType = "number_integer"

	// MetafieldTypeRating JSON, {"value": "3.5", "scale_min": "1.0", "scale_max": "5.0"}.
	MetafieldTypeRating MetafieldType = "rating"

	// MetafieldTypeRichTextField JSON, {"type": "root","children": [{"type": "paragraph","children": [{"type": "text","value": "Bold text.","bold": true}]}]}.
	MetafieldTypeRichTextField MetafieldType = "rich_text_field"

	// MetafieldTypeSingleLineTextField A single line of text. Do not use for numbers or dates, use correct type instead!
	MetafieldTypeSingleLineTextField MetafieldType = "single_line_text_field"

	// MetafieldTypeURL https, http, mailto, sms, or tel.
	MetafieldTypeURL MetafieldType = "url"

	// MetafieldTypeVolume JSON, {"value:" 20.0, "unit": "ml"}.
	MetafieldTypeVolume MetafieldType = "volume"

	// MetafieldTypeWeight JSON, {"value:" 2.5, "unit": "kg"}.
	MetafieldTypeWeight MetafieldType = "weight"
)

// Metafield represents a Shopify metafield.
type Metafield struct {
	CreatedAt         *time.Time    `json:"created_at,omitempty"`
	Description       string        `json:"description,omitempty"`    // Description of the metafield.
	Id                uint64        `json:"id,omitempty"`             // Assigned by Shopify, used for updating a metafield.
	Key               string        `json:"key,omitempty"`            // The unique identifier for a metafield within its namespace, 3-64 characters long.
	Namespace         string        `json:"namespace,omitempty"`      // The container for a group of metafields, 3-255 characters long.
	OwnerId           uint64        `json:"owner_id,omitempty"`       // The unique Id of the resource the metafield is for, i.e.: an Order Id.
	OwnerResource     string        `json:"owner_resource,omitempty"` // The type of reserouce the metafield is for, i.e.: and Order.
	UpdatedAt         *time.Time    `json:"updated_at,omitempty"`     //
	Value             interface{}   `json:"value,omitempty"`          // The data stored in the metafield. Always stored as a string, use Type field for actual data type.
	Type              MetafieldType `json:"type,omitempty"`           // One of Shopify's defined types, see MetafieldType.
	AdminGraphqlApiId string        `json:"admin_graphql_api_id,omitempty"`
}

// MetafieldResource represents the result from the metafields/X.json endpoint
type MetafieldResource struct {
	Metafield *Metafield `json:"metafield"`
}

// MetafieldsResource represents the result from the metafields.json endpoint
type MetafieldsResource struct {
	Metafields []Metafield `json:"metafields"`
}

// List metafields
func (s *MetafieldServiceOp) List(ctx context.Context, options interface{}) ([]Metafield, error) {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s.json", prefix)
	resource := new(MetafieldsResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Metafields, err
}

// Count metafields
func (s *MetafieldServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s/count.json", prefix)
	return s.client.Count(ctx, path, options)
}

// Get individual metafield
func (s *MetafieldServiceOp) Get(ctx context.Context, metafieldId uint64, options interface{}) (*Metafield, error) {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s/%d.json", prefix, metafieldId)
	resource := new(MetafieldResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Metafield, err
}

// Create a new metafield
func (s *MetafieldServiceOp) Create(ctx context.Context, metafield Metafield) (*Metafield, error) {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s.json", prefix)
	wrappedData := MetafieldResource{Metafield: &metafield}
	resource := new(MetafieldResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Metafield, err
}

// Update an existing metafield
func (s *MetafieldServiceOp) Update(ctx context.Context, metafield Metafield) (*Metafield, error) {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	path := fmt.Sprintf("%s/%d.json", prefix, metafield.Id)
	wrappedData := MetafieldResource{Metafield: &metafield}
	resource := new(MetafieldResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Metafield, err
}

// Delete an existing metafield
func (s *MetafieldServiceOp) Delete(ctx context.Context, metafieldId uint64) error {
	prefix := MetafieldPathPrefix(s.resource, s.resourceId)
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d.json", prefix, metafieldId))
}
