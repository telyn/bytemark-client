package billing_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetDefinition(t *testing.T) {
	tests := []struct {
		billing.Definition
	}{}

	for i, test := range tests {
		rts := testutil.RequestTestSpec{
			Method:   "PUT",
			Endpoint: lib.BillingEndpoint,
			URL:      fmt.Sprintf("/api/v1/definitions/%s", test.Name),
			Response: test.Definition,
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			def, err := billingMethods.GetDefinition(client, test.Name)
			assert.Equal(t, testutil.Name(i), test.Definition, def)
			assert.Equal(t, testutil.Name(i), nil, err)
		})
	}
}
