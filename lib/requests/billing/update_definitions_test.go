package billing_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestUpdateDefinitions(t *testing.T) {
	tests := []struct {
		Definitions billing.Definitions
		Expected    map[string]interface{}
		ShouldErr   bool
	}{{
		Definitions: billing.Definitions{},
		Expected:    map[string]interface{}{},
	}, {
		Definitions: billing.Definitions{
			TrialDays:  7,
			TrialPence: 5000,
		},
		Expected: map[string]interface{}{"trial_days": 7.0, "trial_pence": 5000.0},
	}}

	for i, test := range tests {
		rts := testutil.RequestTestSpec{
			Method:        "PUT",
			Endpoint:      lib.BillingEndpoint,
			URL:           "/api/v1/definitions",
			AssertRequest: assert.BodyUnmarshalEqual(test.Expected),
			Response:      nil,
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			billingMethods.UpdateDefinitions(client, test.Definitions)
		})
	}
}
