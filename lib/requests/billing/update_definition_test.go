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
		Name       string
		Definition billing.Definition
		Expected   map[string]interface{}
		ShouldErr  bool
	}{{
		Name: "works",
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
		Name:      "Errors when blank definition provided",
		ShouldErr: true,
	}}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				Method:        "PUT",
				Endpoint:      lib.BillingEndpoint,
				URL:           fmt.Sprintf("/api/v1/definitions/%s", test.Definition.Name),
				AssertRequest: assert.BodyUnmarshalEqual(test.Expected),
				Response:      nil,
				NoVerify:      test.ShouldErr,
			}
			rts.Run(t, "", true, func(client lib.Client) {
				err := billingMethods.UpdateDefinition(client, test.Definition)
				if test.ShouldErr {
					assert.NotEqual(t, "", nil, err)
				} else {
					assert.Equal(t, "", nil, err)
				}
			})
		})
	}
}
