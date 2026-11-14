package goshopify

import (
	"context"
	"testing"
)

func TestMobilePlatformApplicationAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "List",
			method:   "GET",
			path:     "mobile_platform_applications.json",
			response: `{"mobile_platform_applications":[{"id":7}]}`,
			want:     []MobilePlatformApplication{{Id: 7}},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MobilePlatformApplication.List(ctx, nil)
			},
		},
		{
			name:     "Get",
			method:   "GET",
			path:     "mobile_platform_applications/7.json",
			response: `{"mobile_platform_application":{"id":7}}`,
			want:     &MobilePlatformApplication{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MobilePlatformApplication.Get(ctx, 7, nil)
			},
		},
		{
			name:   "Delete",
			method: "DELETE",
			path:   "mobile_platform_applications/7.json",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.MobilePlatformApplication.Delete(ctx, 7)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

// Example response from https://shopify.dev/docs/api/admin-rest/2026-10/resources/mobileplatformapplication
func TestMobilePlatformApplicationDocumentedResponse(t *testing.T) {
	var value MobilePlatformApplicationsResource
	testFeatureResponse(t, "fixtures/mobile_platform_application_documented.json", &value, true)
	if !(len(value.MobilePlatformApplications) == 2 && value.MobilePlatformApplications[0].EnabledAppClips != nil && !*value.MobilePlatformApplications[0].EnabledAppClips) {
		t.Fatalf("incomplete response: %#v", value)
	}
}

func TestMobilePlatformApplicationWriteFields(t *testing.T) {
	disabled := false
	fingerprints := []string{}
	for _, tc := range []featureEndpoint{
		{
			name:     "Create",
			method:   "POST",
			path:     "mobile_platform_applications.json",
			request:  `{"mobile_platform_application":{"application_id":"X1Y2.ca.domain.app","platform":"ios"}}`,
			response: `{"mobile_platform_application":{"id":7}}`,
			want:     &MobilePlatformApplication{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MobilePlatformApplication.Create(ctx, MobilePlatformApplication{ApplicationId: "X1Y2.ca.domain.app", Platform: "ios"})
			},
		},
		{
			name:     "UpdateFalseAndEmpty",
			method:   "PUT",
			path:     "mobile_platform_applications/7.json",
			request:  `{"mobile_platform_application":{"id":7,"enabled_universal_or_app_links":false,"enabled_shared_webcredentials":false,"enabled_app_clips":false,"sha256_cert_fingerprints":[]}}`,
			response: `{"mobile_platform_application":{"id":7}}`,
			want:     &MobilePlatformApplication{Id: 7},
			call: func(ctx context.Context) (interface{}, error) {
				return client.MobilePlatformApplication.Update(ctx, MobilePlatformApplication{Id: 7, EnabledUniversalOrAppLinks: &disabled, EnabledSharedWebcredentials: &disabled, EnabledAppClips: &disabled, Sha256CertFingerprints: &fingerprints})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
