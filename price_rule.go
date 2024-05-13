package goshopify

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const priceRulesBasePath = "price_rules"

// PriceRuleService is an interface for interfacing with the price rule endpoints
// of the Shopify API.
// See: https://shopify.dev/docs/admin-api/rest/reference/discounts/pricerule
type PriceRuleService interface {
	Get(context.Context, uint64) (*PriceRule, error)
	Create(context.Context, PriceRule) (*PriceRule, error)
	Update(context.Context, PriceRule) (*PriceRule, error)
	List(context.Context) ([]PriceRule, error)
	Delete(context.Context, uint64) error
}

// PriceRuleServiceOp handles communication with the price rule related methods of the Shopify API.
type PriceRuleServiceOp struct {
	client *Client
}

// PriceRule represents a Shopify discount rule
type PriceRule struct {
	Id                                     uint64                                  `json:"id,omitempty"`
	Title                                  string                                  `json:"title,omitempty"`
	ValueType                              string                                  `json:"value_type,omitempty"`
	Value                                  *decimal.Decimal                        `json:"value,omitempty"`
	CustomerSelection                      string                                  `json:"customer_selection,omitempty"`
	TargetType                             string                                  `json:"target_type,omitempty"`
	TargetSelection                        string                                  `json:"target_selection,omitempty"`
	AllocationMethod                       string                                  `json:"allocation_method,omitempty"`
	AllocationLimit                        int                                     `json:"allocation_limit,omitempty"`
	OncePerCustomer                        bool                                    `json:"once_per_customer,omitempty"`
	UsageLimit                             int                                     `json:"usage_limit,omitempty"`
	StartsAt                               *time.Time                              `json:"starts_at,omitempty"`
	EndsAt                                 *time.Time                              `json:"ends_at,omitempty"`
	CreatedAt                              *time.Time                              `json:"created_at,omitempty"`
	UpdatedAt                              *time.Time                              `json:"updated_at,omitempty"`
	EntitledProductIds                     []uint64                                `json:"entitled_product_ids,omitempty"`
	EntitledVariantIds                     []uint64                                `json:"entitled_variant_ids,omitempty"`
	EntitledCollectionIds                  []uint64                                `json:"entitled_collection_ids,omitempty"`
	EntitledCountryIds                     []uint64                                `json:"entitled_country_ids,omitempty"`
	PrerequisiteProductIds                 []uint64                                `json:"prerequisite_product_ids,omitempty"`
	PrerequisiteVariantIds                 []uint64                                `json:"prerequisite_variant_ids,omitempty"`
	PrerequisiteCollectionIds              []uint64                                `json:"prerequisite_collection_ids,omitempty"`
	PrerequisiteSavedSearchIds             []uint64                                `json:"prerequisite_saved_search_ids,omitempty"`
	PrerequisiteCustomerIds                []uint64                                `json:"prerequisite_customer_ids,omitempty"`
	PrerequisiteSubtotalRange              *prerequisiteSubtotalRange              `json:"prerequisite_subtotal_range,omitempty"`
	PrerequisiteQuantityRange              *prerequisiteQuantityRange              `json:"prerequisite_quantity_range,omitempty"`
	PrerequisiteShippingPriceRange         *prerequisiteShippingPriceRange         `json:"prerequisite_shipping_price_range,omitempty"`
	PrerequisiteToEntitlementQuantityRatio *prerequisiteToEntitlementQuantityRatio `json:"prerequisite_to_entitlement_quantity_ratio,omitempty"`
}

type prerequisiteSubtotalRange struct {
	GreaterThanOrEqualTo string `json:"greater_than_or_equal_to,omitempty"`
}

type prerequisiteQuantityRange struct {
	GreaterThanOrEqualTo int `json:"greater_than_or_equal_to,omitempty"`
}

type prerequisiteShippingPriceRange struct {
	LessThanOrEqualTo string `json:"less_than_or_equal_to,omitempty"`
}

type prerequisiteToEntitlementQuantityRatio struct {
	PrerequisiteQuantity int `json:"prerequisite_quantity,omitempty"`
	EntitledQuantity     int `json:"entitled_quantity,omitempty"`
}

// PriceRuleResource represents the result from the price_rules/X.json endpoint
type PriceRuleResource struct {
	PriceRule *PriceRule `json:"price_rule"`
}

// PriceRulesResource represents the result from the price_rules.json endpoint
type PriceRulesResource struct {
	PriceRules []PriceRule `json:"price_rules"`
}

// SetPrerequisiteSubtotalRange sets or clears the subtotal range for which a cart must lie within to qualify for the price-rule
func (pr *PriceRule) SetPrerequisiteSubtotalRange(greaterThanOrEqualTo *string) error {
	if greaterThanOrEqualTo == nil {
		pr.PrerequisiteSubtotalRange = nil
	} else {
		if !validateMoney(*greaterThanOrEqualTo) {
			return fmt.Errorf("failed to parse value as Decimal, invalid prerequisite subtotal range")
		}

		pr.PrerequisiteSubtotalRange = &prerequisiteSubtotalRange{
			GreaterThanOrEqualTo: *greaterThanOrEqualTo,
		}
	}

	return nil
}

// SetPrerequisiteQuantityRange sets or clears the quantity range for which a cart must lie within to qualify for the price-rule
func (pr *PriceRule) SetPrerequisiteQuantityRange(greaterThanOrEqualTo *int) {
	if greaterThanOrEqualTo == nil {
		pr.PrerequisiteQuantityRange = nil
	} else {
		pr.PrerequisiteQuantityRange = &prerequisiteQuantityRange{
			GreaterThanOrEqualTo: *greaterThanOrEqualTo,
		}
	}
}

// SetPrerequisiteShippingPriceRange sets or clears the shipping price range for which a cart must lie within to qualify for the price-rule
func (pr *PriceRule) SetPrerequisiteShippingPriceRange(lessThanOrEqualTo *string) error {
	if lessThanOrEqualTo == nil {
		pr.PrerequisiteShippingPriceRange = nil
	} else {
		if !validateMoney(*lessThanOrEqualTo) {
			return fmt.Errorf("failed to parse value as Decimal, invalid prerequisite shipping price range")
		}

		pr.PrerequisiteShippingPriceRange = &prerequisiteShippingPriceRange{
			LessThanOrEqualTo: *lessThanOrEqualTo,
		}
	}

	return nil
}

// SetPrerequisiteToEntitlementQuantityRatio sets or clears the ratio between ordered items and entitled items (eg. buy X, get y free) for which a cart is eligible in the price-rule
func (pr *PriceRule) SetPrerequisiteToEntitlementQuantityRatio(prerequisiteQuantity *int, entitledQuantity *int) {
	if prerequisiteQuantity == nil && entitledQuantity == nil {
		pr.PrerequisiteToEntitlementQuantityRatio = nil
		return
	}

	var pQuant, eQuant int
	if prerequisiteQuantity != nil {
		pQuant = *prerequisiteQuantity
	}

	if entitledQuantity != nil {
		eQuant = *entitledQuantity
	}

	pr.PrerequisiteToEntitlementQuantityRatio = &prerequisiteToEntitlementQuantityRatio{
		PrerequisiteQuantity: pQuant,
		EntitledQuantity:     eQuant,
	}
}

// Get retrieves a single price rules
func (s *PriceRuleServiceOp) Get(ctx context.Context, priceRuleId uint64) (*PriceRule, error) {
	path := fmt.Sprintf("%s/%d.json", priceRulesBasePath, priceRuleId)
	resource := new(PriceRuleResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.PriceRule, err
}

// List retrieves a list of price rules
func (s *PriceRuleServiceOp) List(ctx context.Context) ([]PriceRule, error) {
	path := fmt.Sprintf("%s.json", priceRulesBasePath)
	resource := new(PriceRulesResource)
	err := s.client.Get(ctx, path, resource, nil)
	return resource.PriceRules, err
}

// Create creates a price rule
func (s *PriceRuleServiceOp) Create(ctx context.Context, pr PriceRule) (*PriceRule, error) {
	path := fmt.Sprintf("%s.json", priceRulesBasePath)
	resource := new(PriceRuleResource)
	wrappedData := PriceRuleResource{PriceRule: &pr}
	err := s.client.Post(ctx, path, wrappedData, resource)
	return resource.PriceRule, err
}

// Update updates an existing a price rule
func (s *PriceRuleServiceOp) Update(ctx context.Context, pr PriceRule) (*PriceRule, error) {
	path := fmt.Sprintf("%s/%d.json", priceRulesBasePath, pr.Id)
	resource := new(PriceRuleResource)
	wrappedData := PriceRuleResource{PriceRule: &pr}
	err := s.client.Put(ctx, path, wrappedData, resource)
	return resource.PriceRule, err
}

// Delete deletes a price rule
func (s *PriceRuleServiceOp) Delete(ctx context.Context, priceRuleId uint64) error {
	path := fmt.Sprintf("%s/%d.json", priceRulesBasePath, priceRuleId)
	err := s.client.Delete(ctx, path)
	return err
}

func validateMoney(v string) bool {
	_, err := decimal.NewFromString(v)
	return err == nil
}
