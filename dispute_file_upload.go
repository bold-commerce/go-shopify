package goshopify

import (
	"context"
	"fmt"
	"time"
)

// DisputeFileUploadService provides access to DisputeFileUpload endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/dispute-file-upload
type DisputeFileUploadService interface {
	Create(ctx context.Context, disputeId uint64, upload DisputeFileUploadRequest) (*DisputeFileUpload, error)
	Delete(ctx context.Context, disputeId, fileId uint64) error
}

// DisputeFileUploadServiceOp implements DisputeFileUploadService.
type DisputeFileUploadServiceOp struct {
	client *Client
}

// DisputeFileUpload represents Shopify dispute file upload.
type DisputeFileUpload struct {
	Id                  uint64     `json:"id,omitempty"`
	ShopId              uint64     `json:"shop_id,omitempty"`
	FileSize            int64      `json:"file_size,omitempty"`
	FileType            string     `json:"file_type,omitempty"`
	OriginalFilename    string     `json:"original_filename,omitempty"`
	Filename            string     `json:"filename,omitempty"`
	DisputeEvidenceId   uint64     `json:"dispute_evidence_id,omitempty"`
	DisputeEvidenceType string     `json:"dispute_evidence_type,omitempty"`
	URL                 string     `json:"url,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

// DisputeFileUploadRequest represents Shopify dispute file upload request.
type DisputeFileUploadRequest struct {
	DocumentType string `json:"document_type"`
	Filename     string `json:"filename"`
	Mimetype     string `json:"mimetype"`
	// Data contains the raw file bytes; JSON encoding converts them to base64.
	Data []byte `json:"data"`
}

// DisputeFileUploadResource represents the dispute_file_upload response.
type DisputeFileUploadResource struct {
	DisputeFileUpload *DisputeFileUpload `json:"dispute_file_upload"`
}

// Create creates a dispute file upload record.
func (s *DisputeFileUploadServiceOp) Create(ctx context.Context, disputeId uint64, upload DisputeFileUploadRequest) (*DisputeFileUpload, error) {
	resource := new(DisputeFileUploadResource)
	err := s.client.Post(ctx, fmt.Sprintf("shopify_payments/disputes/%d/dispute_file_uploads.json", disputeId), struct {
		Upload DisputeFileUploadRequest `json:"dispute_file_upload"`
	}{upload}, resource)
	return resource.DisputeFileUpload, err
}

// Delete deletes a dispute file upload record.
func (s *DisputeFileUploadServiceOp) Delete(ctx context.Context, disputeId, fileId uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("shopify_payments/disputes/%d/dispute_file_uploads/%d.json", disputeId, fileId))
}
