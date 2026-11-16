package goshopify

import (
	"context"
	"testing"
)

func TestCustomerAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "AccountActivationURL",
			method:   "POST",
			path:     "customers/11/account_activation_url.json",
			request:  `{}`,
			response: `{"account_activation_url":"https://example.com/activate"}`,
			want:     "https://example.com/activate",
			call: func(ctx context.Context) (interface{}, error) {
				return client.Customer.AccountActivationURL(ctx, 11)
			},
		},
		{
			name:     "SendInvite",
			method:   "POST",
			path:     "customers/11/send_invite.json",
			request:  `{"customer_invite":{"custom_message":"Welcome"}}`,
			response: `{"customer_invite":{"custom_message":"Welcome"}}`,
			want:     &CustomerInvite{CustomMessage: "Welcome"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Customer.SendInvite(ctx, 11, CustomerInvite{CustomMessage: "Welcome"})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}

func TestCustomerInviteOptions(t *testing.T) {
	for _, tc := range []featureEndpoint{
		{
			name:     "DefaultInvite",
			method:   "POST",
			path:     "customers/11/send_invite.json",
			request:  `{"customer_invite":{}}`,
			response: `{"customer_invite":{"to":"customer@example.com"}}`,
			want:     &CustomerInvite{To: "customer@example.com"},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Customer.SendInvite(ctx, 11, CustomerInvite{})
			},
		},
		{
			name:     "AllInviteFields",
			method:   "POST",
			path:     "customers/11/send_invite.json",
			request:  `{"customer_invite":{"to":"customer@example.com","from":"owner@example.com","bcc":["support@example.com","audit@example.com"],"subject":"Welcome & thanks","custom_message":"Hello, Zoë!\nPlease activate your account."}}`,
			response: `{"customer_invite":{"to":"customer@example.com","from":"owner@example.com","bcc":["support@example.com","audit@example.com"],"subject":"Welcome & thanks","custom_message":"Hello, Zoë!\nPlease activate your account."}}`,
			want: &CustomerInvite{
				To:            "customer@example.com",
				From:          "owner@example.com",
				Bcc:           []string{"support@example.com", "audit@example.com"},
				Subject:       "Welcome & thanks",
				CustomMessage: "Hello, Zoë!\nPlease activate your account.",
			},
			call: func(ctx context.Context) (interface{}, error) {
				return client.Customer.SendInvite(ctx, 11, CustomerInvite{To: "customer@example.com", From: "owner@example.com", Bcc: []string{"support@example.com", "audit@example.com"}, Subject: "Welcome & thanks", CustomMessage: "Hello, Zoë!\nPlease activate your account."})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
