package mocks_test

import (
	"os"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

var t = fakeTestingT{}

func ExampleRequest_aSingleEndpoint() {
	_, client, cliApp := testutil.BaseTestSetup(t, false, []cli.Command{{
		Name:        "example",
		Description: "example command, gets the test bigv account",
		Action: app.Action(func(c *app.Context) error {
			account, err := billingRequests.GetAccountByBigVName(c.Client(), "test")
			if err != nil {
				return err
			}
			c.Log("Account: %d", account.ID)
			return nil
		}),
	}})

	r := mocks.Request{
		// This is an example test - t is a *testing.T
		T: t,
		ResponseObject: []billing.Account{{
			ID: 2,
		}},
	}

	// This chunk undoes the test-related writer redirection in BaseTestSetup - do not use in an actual test!
	cliApp.Metadata["debugWriter"] = nil
	cliApp.Writer = os.Stdout
	cliApp.ErrWriter = os.Stdout
	// end of the chunk not to use in a real test

	// Note that when using MockRequest there is no way to make an assertion
	// based on the URL. If you wish to assert anything about the URL, use the
	// client.When("BuildRequest", ...).Returns(Request) method, shown in the
	// MultipleEndpoints example below.
	client.MockRequest = &r

	cliApp.Run([]string{"bytemark", "example"})

	// Output: Account: 2
}

func ExampleRequest_multipleEndpoints() {
	// TODO(telyn): write this example
	_, client, cliApp := testutil.BaseTestSetup(t, false, []cli.Command{{
		Name:        "example",
		Description: "example command, gets the test bigv account",
		Action: app.Action(func(c *app.Context) error {
			account, err := billingRequests.AssentToAgreement(c.Client(), billing.Assent{})
			if err != nil {
				return err
			}
			c.Log("Success!")
			return nil
		}),
	}})

	r := mocks.Request{
		// This is an example test - t is a *testing.T
		T: t,
		ResponseObject: []billing.Account{{
			ID: 2,
		}},
	}

	// This chunk undoes the test-related writer redirection in BaseTestSetup - do not use in an actual test!
	cliApp.Metadata["debugWriter"] = nil
	cliApp.Writer = os.Stdout
	cliApp.ErrWriter = os.Stdout
	// end of the chunk not to use in a real test

	client.MockRequest = &r

	cliApp.Run([]string{"bytemark", "example"})

	// Output: Account: 2

}

func ExampleRequest_withError() {
	// TODO(telyn): write this example
	// This is an example test - t is a pretend *testing.T
}
