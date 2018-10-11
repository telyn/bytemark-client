package billing_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetPerson(t *testing.T) {
	tests := []struct {
		body      []billing.Person
		username  string
		shouldErr bool
	}{
		{
			body:      []billing.Person{},
			username:  "bwagg",
			shouldErr: true,
		},
		{
			body: []billing.Person{{
				ID:        101,
				FirstName: "Geoff",
				LastName:  "Jeff",
				Email:     "geoff@bytemark.com",
			}},
			username: "bwagg",
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "GET",
			Endpoint:      lib.BillingEndpoint,
			URL:           "/api/v1/people",
			Response:      test.body,
			AssertRequest: assert.QueryValue("username", test.username),
		}
		rts.Run(t, testutil.Name(i), true, func(client lib.Client) {
			person, err := billingMethods.GetPerson(client, test.username)
			if len(test.body) > 0 {
				assert.Equal(t, testutil.Name(i), test.body[0], person)
			}
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
