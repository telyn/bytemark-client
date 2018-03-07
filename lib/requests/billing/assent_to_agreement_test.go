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

func TestAssentToAgreement(t *testing.T) {
	tests := []struct {
		assent    billing.Assent
		expected  map[string]interface{}
		shouldErr bool
	}{
		{
			assent: billing.Assent{
				AgreementID: "jeff",
				AccountID:   123456,
				PersonID:    789101,
				Name:        "geoff",
				Email:       "geoff@bytemark.com",
			},
			expected: map[string]interface{}{
				"account_id": 123456.0,
				"person_id":  789101.0,
				"name":       "geoff",
				"email":      "geoff@bytemark.com",
			},
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "POST",
			Endpoint:      lib.BillingEndpoint,
			URL:           fmt.Sprintf("/api/v1/agreements/%s/assents", test.assent.AgreementID),
			AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := billingMethods.AssentToAgreement(client, test.assent)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
