package goshopify

import (
	"context"
	"testing"
)

func TestCustomerAddressAddedEndpoints(t *testing.T) {
	tests := []featureEndpoint{
		{
			name:     "SetDefault",
			method:   "PUT",
			path:     "customers/11/addresses/7/default.json",
			response: `{"customer_address":{"id":7,"default":true}}`,
			want:     &CustomerAddress{Id: 7, Default: true},
			call: func(ctx context.Context) (interface{}, error) {
				return client.CustomerAddress.SetDefault(ctx, 11, 7)
			},
		},
		{
			name:   "BulkDestroy",
			method: "PUT",
			path:   "customers/11/addresses/set.json",
			query:  "operation=destroy&address_ids[]=7&address_ids[]=8",
			want:   nil,
			call: func(ctx context.Context) (interface{}, error) {
				return nil, client.CustomerAddress.BulkDestroy(ctx, 11, []uint64{7, 8})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFeatureEndpoint(t, tc) })
	}
}
