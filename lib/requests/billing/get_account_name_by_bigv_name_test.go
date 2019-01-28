package billing_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetAccountByBigvName(t *testing.T) {
	tests := []struct {
		body      []billing.Account
		bigvName  string
		shouldErr bool
	}{
		{
			body:      []billing.Account{},
			bigvName:  "",
			shouldErr: true,
		},
		{
			body:      []billing.Account{},
			bigvName:  "notaperson",
			shouldErr: true,
		},
		{
			body: []billing.Account{{
				ID:   101,
				Name: "Geoff Jeff",
			}},
			bigvName: "bwagg",
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "GET",
			Endpoint:      lib.BillingEndpoint,
			URL:           "/api/v1/accounts",
			Response:      test.body,
			AssertRequest: assert.QueryValue("bigv_account_name", test.bigvName),
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			account, err := billingMethods.GetAccountByBigVName(client, test.bigvName)
			if len(test.body) > 0 {
				assert.Equal(t, testutil.Name(i), test.body[0], account)
			}
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
