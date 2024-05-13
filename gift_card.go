package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const giftCardsBasePath = "gift_cards"

// giftCardService is an interface for interfacing with the gift card endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/api/admin-rest/2023-04/resources/gift-card
type GiftCardService interface {
	Get(context.Context, uint64) (*GiftCard, error)
	Create(context.Context, GiftCard) (*GiftCard, error)
	Update(context.Context, GiftCard) (*GiftCard, error)
	List(context.Context) ([]GiftCard, error)
	Disable(context.Context, uint64) (*GiftCard, error)
	Count(context.Context, interface{}) (int, error)
}

// giftCardServiceOp handles communication with the gift card related methods of the Shopify API.
type GiftCardServiceOp struct {
	client *Client
}

// giftCard represents a Shopify discount rule
type GiftCard struct {
	Id             uint64           `json:"id,omitempty"`
	ApiClientId    uint64           `json:"api_client_id,omitempty"`
	Balance        *decimal.Decimal `json:"balance,omitempty"`
	InitalValue    *decimal.Decimal `json:"initial_value,omitempty"`
	Code           string           `json:"code,omitempty"`
	Currency       string           `json:"currency,omitempty"`
	CustomerId     *CustomerId      `json:"customer_id,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	DisabledAt     *time.Time       `json:"disabled_at,omitempty"`
	ExpiresOn      string           `json:"expires_on,omitempty"`
	LastCharacters string           `json:"last_characters,omitempty"`
	LineItemId     uint64           `json:"line_item_id,omitempty"`
	Note           string           `json:"note,omitempty"`
	OrderId        uint64           `json:"order_id,omitempty"`
	TemplateSuffix string           `json:"template_suffix,omitempty"`
	UserId         uint64           `json:"user_id,omitempty"`
	UpdatedAt      *time.Time       `json:"updated_at,omitempty"`
}

type CustomerId struct {
	CustomerId uint64 `json:"customer_id,omitempty"`
}

// giftCardResource represents the result from the gift_cards/X.json endpoint
type GiftCardResource struct {
	GiftCard *GiftCard `json:"gift_card"`
}

// giftCardsResource represents the result from the gift_cards.json endpoint
type GiftCardsResource struct {
	GiftCards []GiftCard `json:"gift_cards"`
}

// Get retrieves a single gift cards
func (s *GiftCardServiceOp) Get(ctx context.Context, giftCardId uint64) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%d.json", giftCardsBasePath, giftCardId)
	resource := new(GiftCardResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.GiftCard, err
}

// List retrieves a list of gift cards
func (s *GiftCardServiceOp) List(ctx context.Context) ([]GiftCard, error) {
	path := fmt.Sprintf("%s.json", giftCardsBasePath)
	resource := new(GiftCardsResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.GiftCards, err
}

// Create creates a gift card
func (s *GiftCardServiceOp) Create(ctx context.Context, pr GiftCard) (*GiftCard, error) {
	path := fmt.Sprintf("%s.json", giftCardsBasePath)
	resource := new(GiftCardResource)
	wrappedData := GiftCardResource{GiftCard: &pr}
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.GiftCard, err
}

// Update updates an existing a gift card
func (s *GiftCardServiceOp) Update(ctx context.Context, pr GiftCard) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%d.json", giftCardsBasePath, pr.Id)
	resource := new(GiftCardResource)
	wrappedData := GiftCardResource{GiftCard: &pr}
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.GiftCard, err
}

// Disable disables an existing a gift card
func (s *GiftCardServiceOp) Disable(ctx context.Context, giftCardId uint64) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%d/disable.json", giftCardsBasePath, giftCardId)
	resource := new(GiftCardResource)
	wrappedData := GiftCardResource{GiftCard: &GiftCard{Id: giftCardId}}
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.GiftCard, err
}

// Count retrieves the number of gift cards
func (s *GiftCardServiceOp) Count(ctx context.Context, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", giftCardsBasePath)
	return s.client.Count(ctx, path, options)
}
