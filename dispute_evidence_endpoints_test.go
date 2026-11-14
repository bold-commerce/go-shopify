package goshopify

import (
	"context"
	"testing"
	"time"
)

func TestDisputeEvidenceAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Get",
			method:   "GET",
			path:     "shopify_payments/disputes/11/dispute_evidences.json",
			response: `{"dispute_evidence":{"id":7,"payments_dispute_id":11,"product_description":"Draft"}}`,
			want:     &DisputeEvidence{Id: 7, PaymentsDisputeId: 11, ProductDescription: "Draft"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DisputeEvidence.Get(ctx, 11)
			},
		},
		{
			name:     "Update",
			method:   "PUT",
			path:     "shopify_payments/disputes/11/dispute_evidences.json",
			request:  `{"dispute_evidence":{"shipping_address_attributes":{"address1":"1 Infinite Loop","country_code":"US"}}}`,
			response: `{"dispute_evidence":{"id":7,"payments_dispute_id":11,"product_description":"Draft"}}`,
			want:     &DisputeEvidence{Id: 7, PaymentsDisputeId: 11, ProductDescription: "Draft"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DisputeEvidence.Update(ctx, 11, DisputeEvidenceUpdate{ShippingAddressAttributes: &Address{Address1: "1 Infinite Loop", CountryCode: "US"}})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/dispute-evidence
func TestDisputeEvidenceDocumentedResponse(t *testing.T) {
	var value DisputeEvidenceResource
	testFeatureResponse(t, "fixtures/dispute_evidence_documented.json", &value, true)
	if !(value.DisputeEvidence.ProductDescription != "" && len(value.DisputeEvidence.Fulfillments) == 2 && value.DisputeEvidence.Fulfillments[1].ShippingDate.Format("2006-01-02") == "2017-01-02") {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestDisputeEvidenceUpdateFields(t *testing.T) {
	at := time.Date(2017, 1, 2, 13, 0, 0, 0, time.UTC)
	empty := ""
	testFeatureEndpoint(t, featureEndpoint{
		name:     "Update",
		method:   "PUT",
		path:     "shopify_payments/disputes/11/dispute_evidences.json",
		request:  `{"dispute_evidence":{"uncategorized_text":"","fulfillments_attributes":[{"shipping_carrier":"FedEx","shipping_tracking_number":"4321","shipping_date":"2017-01-02T13:00:00Z"}]}}`,
		response: `{"dispute_evidence":{"id":7}}`,
		want:     &DisputeEvidence{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.DisputeEvidence.Update(ctx, 11, DisputeEvidenceUpdate{UncategorizedText: &empty, FulfillmentsAttributes: []DisputeEvidenceFulfillmentInput{{ShippingCarrier: "FedEx", ShippingTrackingNumber: "4321", ShippingDate: &at}}})
		},
	})
}
