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

func TestUpdateDefinitions(t *testing.T) {
	tests := []struct {
		billing.Definition
		Expected  map[string]interface{}
		ShouldErr bool
	}{{
		Definition: billing.Definition{
			Name:           "test-def",
			Value:          "test-val",
			UpdateGroupReq: "staff",
		},
		Expected: map[string]interface{}{
			"name":             "test-def",
			"value":            "test-val",
			"update_group_req": "staff",
		},
	}, {
		ShouldErr: true,
	}}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "PUT",
			Endpoint:      lib.BillingEndpoint,
			URL:           fmt.Sprintf("/api/v1/definitions/%s", test.Name),
			AssertRequest: assert.BodyUnmarshalEqual(test.Expected),
			Response:      nil,
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := billingMethods.UpdateDefinition(client, test.Definition)
			if test.ShouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
