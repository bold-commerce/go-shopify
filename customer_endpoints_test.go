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
