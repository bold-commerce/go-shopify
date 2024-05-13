package goshopify

import (
	"context"
	"fmt"
	"time"
)

// ImageService is an interface for interacting with the image endpoints
// of the Shopify API.
// See https://help.shopify.com/api/reference/product_image
type ImageService interface {
	List(context.Context, uint64, interface{}) ([]Image, error)
	Count(context.Context, uint64, interface{}) (int, error)
	Get(context.Context, uint64, uint64, interface{}) (*Image, error)
	Create(context.Context, uint64, Image) (*Image, error)
	Update(context.Context, uint64, Image) (*Image, error)
	Delete(context.Context, uint64, uint64) error
}

// ImageServiceOp handles communication with the image related methods of
// the Shopify API.
type ImageServiceOp struct {
	client *Client
}

// Image represents a Shopify product's image.
type Image struct {
	Id                uint64     `json:"id,omitempty"`
	ProductId         uint64     `json:"product_id,omitempty"`
	Position          int        `json:"position,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Width             int        `json:"width,omitempty"`
	Height            int        `json:"height,omitempty"`
	Src               string     `json:"src,omitempty"`
	Attachment        string     `json:"attachment,omitempty"`
	Filename          string     `json:"filename,omitempty"`
	Alt               string     `json:"alt,omitempty"`
	VariantIds        []uint64   `json:"variant_ids,omitempty"`
	AdminGraphqlApiId string     `json:"admin_graphql_api_id,omitempty"`
}

// ImageResource represents the result form the products/X/images/Y.json endpoint
type ImageResource struct {
	Image *Image `json:"image"`
}

// ImagesResource represents the result from the products/X/images.json endpoint
type ImagesResource struct {
	Images []Image `json:"images"`
}

// List images
func (s *ImageServiceOp) List(ctx context.Context, productId uint64, options interface{}) ([]Image, error) {
	path := fmt.Sprintf("%s/%d/images.json", productsBasePath, productId)
	resource := new(ImagesResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Images, err
}

// Count images
func (s *ImageServiceOp) Count(ctx context.Context, productId uint64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/images/count.json", productsBasePath, productId)
	return s.client.Count(ctx, path, options)
}

// Get individual image
func (s *ImageServiceOp) Get(ctx context.Context, productId uint64, imageId uint64, options interface{}) (*Image, error) {
	path := fmt.Sprintf("%s/%d/images/%d.json", productsBasePath, productId, imageId)
	resource := new(ImageResource)
	err := s.client.Get(ctx, path, resource, options)
	return resource.Image, err
}

// Create a new image
//
// There are 2 methods of creating an image in Shopify:
// 1. Src
// 2. Filename and Attachment
//
// If both Image.Filename and Image.Attachment are supplied,
// then Image.Src is not needed.  And vice versa.
//
// If both Image.Attachment and Image.Src are provided,
// Shopify will take the attachment.
//
// Shopify will accept Image.Attachment without Image.Filename.
func (s *ImageServiceOp) Create(ctx context.Context, productId uint64, image Image) (*Image, error) {
	path := fmt.Sprintf("%s/%d/images.json", productsBasePath, productId)
	wrappedData := ImageResource{Image: &image}
	resource := new(ImageResource)
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.Image, err
}

// Update an existing image
func (s *ImageServiceOp) Update(ctx context.Context, productId uint64, image Image) (*Image, error) {
	path := fmt.Sprintf("%s/%d/images/%d.json", productsBasePath, productId, image.Id)
	wrappedData := ImageResource{Image: &image}
	resource := new(ImageResource)
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.Image, err
}

// Delete an existing image
func (s *ImageServiceOp) Delete(ctx context.Context, productId uint64, imageId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("%s/%d/images/%d.json", productsBasePath, productId, imageId))
}
