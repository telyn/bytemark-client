package billing_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetAccountDeferredStatus(t *testing.T) {
	tests := []struct {
		body      []billing.DeferredStatus
		username  string
		shouldErr bool
	}{
		{
			body:      []billing.DeferredStatus{},
			username:  "",
			shouldErr: true,
		},
		{
			body: []billing.DeferredStatus{{
				ID:       139,
				Deferred: false,
			}},
			username:  "bwagg",
			shouldErr: false,
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "GET",
			Endpoint:      lib.BillingEndpoint,
			URL:           "/api/v1/accounts/status",
			Response:      test.body,
			AssertRequest: assert.QueryValue("username", test.username),
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			account, err := billingMethods.GetAccountDeferredStatus(client, test.username)
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
