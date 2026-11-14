package goshopify

import (
	"context"
	"fmt"
	"time"
)

// MobilePlatformApplicationService provides access to MobilePlatformApplication endpoints.
// See https://shopify.dev/docs/api/admin-rest/latest/resources/mobileplatformapplication
type MobilePlatformApplicationService interface {
	List(ctx context.Context, options interface{}) ([]MobilePlatformApplication, error)
	Get(ctx context.Context, id uint64, options interface{}) (*MobilePlatformApplication, error)
	Create(ctx context.Context, mobilePlatformApplication MobilePlatformApplication) (*MobilePlatformApplication, error)
	Update(ctx context.Context, mobilePlatformApplication MobilePlatformApplication) (*MobilePlatformApplication, error)
	Delete(ctx context.Context, id uint64) error
}

// MobilePlatformApplicationServiceOp implements MobilePlatformApplicationService.
type MobilePlatformApplicationServiceOp struct {
	client *Client
}

// MobilePlatformApplication represents Shopify mobile platform application.
type MobilePlatformApplication struct {
	Id                          uint64     `json:"id,omitempty"`
	ApplicationId               string     `json:"application_id,omitempty"`
	Platform                    string     `json:"platform,omitempty"`
	Sha256CertFingerprints      *[]string  `json:"sha256_cert_fingerprints,omitempty"`
	EnabledUniversalOrAppLinks  *bool      `json:"enabled_universal_or_app_links,omitempty"`
	EnabledSharedWebcredentials *bool      `json:"enabled_shared_webcredentials,omitempty"`
	EnabledAppClips             *bool      `json:"enabled_app_clips,omitempty"`
	AppClipApplicationId        string     `json:"app_clip_application_id,omitempty"`
	CreatedAt                   *time.Time `json:"created_at,omitempty"`
	UpdatedAt                   *time.Time `json:"updated_at,omitempty"`
}

// MobilePlatformApplicationsResource represents the mobile_platform_applications response.
type MobilePlatformApplicationsResource struct {
	MobilePlatformApplications []MobilePlatformApplication `json:"mobile_platform_applications"`
}

// MobilePlatformApplicationResource represents the mobile_platform_application response.
type MobilePlatformApplicationResource struct {
	MobilePlatformApplication *MobilePlatformApplication `json:"mobile_platform_application"`
}

// List retrieves mobile platform application records.
func (s *MobilePlatformApplicationServiceOp) List(ctx context.Context, options interface{}) ([]MobilePlatformApplication, error) {
	resource := new(MobilePlatformApplicationsResource)
	err := s.client.Get(ctx, "mobile_platform_applications.json", resource, options)
	return resource.MobilePlatformApplications, err
}

// Get retrieves a mobile platform application record.
func (s *MobilePlatformApplicationServiceOp) Get(ctx context.Context, id uint64, options interface{}) (*MobilePlatformApplication, error) {
	resource := new(MobilePlatformApplicationResource)
	err := s.client.Get(ctx, fmt.Sprintf("mobile_platform_applications/%d.json", id), resource, options)
	return resource.MobilePlatformApplication, err
}

// Create creates a mobile platform application record.
func (s *MobilePlatformApplicationServiceOp) Create(ctx context.Context, mobilePlatformApplication MobilePlatformApplication) (*MobilePlatformApplication, error) {
	resource := new(MobilePlatformApplicationResource)
	err := s.client.Post(ctx, "mobile_platform_applications.json", MobilePlatformApplicationResource{MobilePlatformApplication: &mobilePlatformApplication}, resource)
	return resource.MobilePlatformApplication, err
}

// Update updates a mobile platform application record.
func (s *MobilePlatformApplicationServiceOp) Update(ctx context.Context, mobilePlatformApplication MobilePlatformApplication) (*MobilePlatformApplication, error) {
	resource := new(MobilePlatformApplicationResource)
	err := s.client.Put(ctx, fmt.Sprintf("mobile_platform_applications/%d.json", mobilePlatformApplication.Id), MobilePlatformApplicationResource{MobilePlatformApplication: &mobilePlatformApplication}, resource)
	return resource.MobilePlatformApplication, err
}

// Delete deletes a mobile platform application record.
func (s *MobilePlatformApplicationServiceOp) Delete(ctx context.Context, id uint64) error {
	return s.client.Delete(ctx, fmt.Sprintf("mobile_platform_applications/%d.json", id))
}
