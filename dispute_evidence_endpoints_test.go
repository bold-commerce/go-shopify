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

func TestDisputeEvidenceSupplementalResponse(t *testing.T) {
	at := featureTime(t, "2026-11-01T09:00:00-07:00")
	day := featureTime(t, "2026-11-01T00:00:00Z")
	submitted := "2026-11-02"
	ids := []uint64{1, 2, 3, 4, 5, 9007199254740993}
	testFeatureEndpoint(t, featureEndpoint{
		name:     "PopulatedOptionalFields",
		method:   "GET",
		path:     "shopify_payments/disputes/11/dispute_evidences.json",
		response: string(loadFixture("dispute_evidence_expanded.json")),
		want: &DisputeEvidence{
			Id:                           7,
			PaymentsDisputeId:            11,
			AccessActivityLog:            "Logged in",
			CancellationPolicyDisclosure: "Policy shown",
			CancellationRebuttal:         "Not cancelled",
			CustomerEmailAddress:         "customer@example.com",
			CustomerFirstName:            "Zoë",
			CustomerLastName:             "Li",
			RefundPolicyDisclosure:       "Refund policy shown",
			RefundRefusalExplanation:     "No receipt",
			UncategorizedText:            "Additional details",
			ProductDescription:           "First line\nSecond line",
			ShippingAddress:              &Address{Address1: "1 Test St", CountryCode: "CA", Zip: "A1A 1A1"},
			CreatedAt:                    at,
			UpdatedAt:                    at,
			SubmittedByMerchantOn:        &submitted,
			Fulfillments: []DisputeEvidenceFulfillment{
				{ShippingCarrier: "UPS", ShippingTrackingNumber: "001ZABC", ShippingDate: &OnlyDate{Time: *day}},
			},
			DisputeEvidenceFiles: &DisputeEvidenceFiles{
				CancellationPolicyFileId:    &ids[0],
				CustomerCommunicationFileId: &ids[1],
				RefundPolicyFileId:          &ids[2],
				ServiceDocumentationFileId:  &ids[3],
				ShippingDocumentationFileId: &ids[4],
				UncategorizedFileId:         &ids[5],
			},
		},
		invalidResponses: []string{`{"dispute_evidence":{"fulfillments":[{"shipping_date":"invalid-date"}]}}`, `{"dispute_evidence":{"dispute_evidence_files":{"refund_policy_file_id":"not-an-id"}}}`},
		call:             func(ctx context.Context) (interface{}, error) { return client.DisputeEvidence.Get(ctx, 11) },
	})
}

func TestDisputeEvidencePopulatedUpdateFields(t *testing.T) {
	activity, policy, rebuttal := "Activity log", "Cancellation policy", "Not cancelled"
	email, first, last := "customer@example.com", "Zoë", "Li"
	refund, refusal, extra := "Refund policy", "No receipt", "More details"
	testFeatureEndpoint(t, featureEndpoint{
		method:   "PUT",
		path:     "shopify_payments/disputes/11/dispute_evidences.json",
		request:  `{"dispute_evidence":{"access_activity_log":"Activity log","cancellation_policy_disclosure":"Cancellation policy","cancellation_rebuttal":"Not cancelled","customer_email_address":"customer@example.com","customer_first_name":"Zoë","customer_last_name":"Li","refund_policy_disclosure":"Refund policy","refund_refusal_explanation":"No receipt","uncategorized_text":"More details"}}`,
		response: `{"dispute_evidence":{"id":7}}`,
		want:     &DisputeEvidence{Id: 7},
		call: func(ctx context.Context) (interface{}, error) {
			return client.DisputeEvidence.Update(ctx, 11, DisputeEvidenceUpdate{
				AccessActivityLog: &activity, CancellationPolicyDisclosure: &policy, CancellationRebuttal: &rebuttal,
				CustomerEmailAddress: &email, CustomerFirstName: &first, CustomerLastName: &last,
				RefundPolicyDisclosure: &refund, RefundRefusalExplanation: &refusal, UncategorizedText: &extra,
			})
		},
	})
}
