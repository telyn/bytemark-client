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
		bigv_name string
		shouldErr bool
	}{
		{
			body:      []billing.Account{},
			bigv_name: "",
			shouldErr: true,
		},
		{
			body:      []billing.Account{},
			bigv_name: "notaperson",
			shouldErr: true,
		},
		{
			body: []billing.Account{{
				ID:   101,
				Name: "Geoff Jeff",
			}},
			bigv_name: "bwagg",
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "GET",
			Endpoint:      lib.BillingEndpoint,
			URL:           "/api/v1/accounts",
			Response:      test.body,
			AssertRequest: assert.QueryValue("bigv_account_name", test.bigv_name),
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			account, err := billingMethods.GetAccountByBigVName(client, test.bigv_name)
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
