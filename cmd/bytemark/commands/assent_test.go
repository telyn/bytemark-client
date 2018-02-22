package commands_test

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"strings"
	"testing"
	// billingRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	// // "github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestAssent(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		account   string
		shouldErr bool
	}{
		{
			name:      "MissingArguments",
			input:     "",
			account:   "accountFromConfig",
			shouldErr: true,
		},
		{
			name:    "SuccessfullyAssentsWithAccountID",
			input:   "--agreement 1 --person bwagg --accountid 101 --name BryanWagg --email geoff@jeff.com",
			account: "bwagg",
		},
		{
			name:    "SuccessfullyAssentsWithAccount",
			input:   "--agreement 1 --person bwagg --account bwagg --name BryanWagg --email geoff@jeff.com",
			account: "bwagg",
		},
		{
			name:      "AmbiguousAccount",
			input:     "--agreement 1 --person bwagg --account bwagg --accountid 1234",
			account:   "bwagg",
			shouldErr: true,
		},
		{
			name:    "SuccessfullyAssentsWithAccountFromConfig",
			account: "accountFromConfig",
			input:   "--agreement 1 --person bwagg --name BryanWagg --email geoff@jeff.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetIgnoreErr", "account").Return("accountFromConfig")
			client.When("BuildRequest", "GET", lib.BillingEndpoint, "/api/v1/accounts?bigv_account_name=%s", []string{test.account}).Return(&mocks.Request{
				T:          t,
				StatusCode: 200,
				ResponseObject: []billing.Account{{
					Name: "bwagg",
					ID:   101,
				},
				},
			})
			client.When("BuildRequest", "GET", lib.BillingEndpoint, "/api/v1/people?username=%s", []string{"bwagg"}).Return(&mocks.Request{
				T:          t,
				StatusCode: 200,
				ResponseObject: []billing.Person{{
					FirstName: "Bryan",
					LastName:  "Wagg",
					ID:        201,
					Email:     "geoff@jeff.com",
				},
				},
			})
			client.When("BuildRequest", "POST", lib.BillingEndpoint, "/api/v1/agreements/%s/assents", []string{"1"}).Return(&mocks.Request{
				T:          t,
				StatusCode: 200,
				RequestObject: billing.Assent{
					AgreementID: "1",
					AccountID:   101,
					PersonID:    201,
					Name:        "BryanWagg",
					Email:       "geoff@jeff.com",
				},
			}).Times(1)

			args := fmt.Sprintf("bytemark assent %s", test.input)
			err := app.Run(strings.Split(args, " "))
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Errorf("should err, but didn't")
			}
			if ok, err := client.Verify(); !ok && !test.shouldErr {
				t.Fatal(err)
			}
		})
	}

}
