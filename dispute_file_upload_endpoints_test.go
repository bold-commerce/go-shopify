package goshopify

import (
	"context"
	"testing"
)

func TestDisputeFileUploadAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "Create",
			method:   "POST",
			path:     "shopify_payments/disputes/11/dispute_file_uploads.json",
			request:  `{"dispute_file_upload":{"document_type":"uncategorized_file","filename":"test.pdf","mimetype":"application/pdf","data":"dGVzdA=="}}`,
			response: `{"dispute_file_upload":{"id":7,"file_size":4,"file_type":"application/pdf"}}`,
			want:     &DisputeFileUpload{Id: 7, FileSize: 4, FileType: "application/pdf"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.DisputeFileUpload.Create(ctx, 11, DisputeFileUploadRequest{DocumentType: "uncategorized_file", Filename: "test.pdf", Mimetype: "application/pdf", Data: []byte("test")})
			},
		},
		{
			name:   "Delete",
			method: "DELETE",
			path:   "shopify_payments/disputes/11/dispute_file_uploads/7.json",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.DisputeFileUpload.Delete(ctx, 11, 7)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/dispute-file-upload
func TestDisputeFileUploadDocumentedResponse(t *testing.T) {
	var value DisputeFileUploadResource
	testFeatureResponse(t, "fixtures/dispute_file_upload_documented.json", &value, true)
	if !(value.DisputeFileUpload.FileSize == 34509 && value.DisputeFileUpload.DisputeEvidenceId == 819974671) {
		t.Fatalf("incomplete response: %#v", value)
	}
}
