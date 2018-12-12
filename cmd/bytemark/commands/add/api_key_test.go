package add_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/add"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

type apiKeyWithErr struct {
	brain.APIKey
	err bool
}
type accountWithErr struct {
	lib.Account
	err bool
}
type serverWithErr struct {
	brain.VirtualMachine
	err bool
}
type groupWithErr struct {
	brain.Group
	err bool
}

func TestAddApiKey(t *testing.T) {
	tests := []struct {
		testutil.CommandT

		// aaand here's all the stuff to assert on and return from Request
		// the map keys are used for the mock BuildRequest calls' URLs
		accounts map[string]accountWithErr
		groups   map[lib.GroupName]groupWithErr
		servers  map[lib.VirtualMachineName]serverWithErr
		// apiKey is both the apiKey to return and some of its fields will be
		// used to define the mock
		apiKey apiKeyWithErr
		user   brain.User
	}{{
		CommandT: testutil.CommandT{
			Name:      "all defaults no label",
			Args:      "",
			ShouldErr: true,
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("Couldn't make a specification for the API key"),
			},
		},
		user: brain.User{
			ID: 299,
			// username not specified in the input - test-user is the default
			// returned by config.GetUser() (see testutil.BaseTestAuthSetup)
			Username: "test-user",
		},
		apiKey: apiKeyWithErr{},
	}, {
		CommandT: testutil.CommandT{
			Name: "all defaults",
			Args: "test-api-key",
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("Successfully created an api key:\ntest-api-key\n  Expires: never\n  Key: apikey.apikey.ay1pee2eye3key\n"),
			},
		},
		user: brain.User{
			ID: 299,
			// username not specified in the input - test-user is the default
			// returned by config.GetUser() (see testutil.BaseTestAuthSetup)
			Username: "test-user",
		},
		apiKey: apiKeyWithErr{APIKey: brain.APIKey{
			ID:        1,
			UserID:    299,
			APIKey:    "apikey.ay1pee2eye3key",
			Label:     "test-api-key",
			ExpiresAt: "",
		}},
	}, {
		CommandT: testutil.CommandT{
			Name: "no privilege specs & POST fails",
			Args: "test-api-key",
			// output is expected to be blank at this stage since it's
			// ProcessError which outputs the error message - which is called
			// from main.main, not from the command's action.
			ShouldErr: true,
		},
		user: brain.User{
			ID:       299,
			Username: "test-user",
		},
		apiKey: apiKeyWithErr{
			APIKey: brain.APIKey{
				Label: "bad-key",
			},
			err: true,
		},
	}}

	for _, test := range tests {
		test.CommandT.Auth = true
		test.CommandT.Commands = add.Commands
		test.CommandT.Args = "apikey " + test.CommandT.Args

		mockPostAPIKey := mocks.Request{
			StatusCode:     200,
			ResponseObject: test.apiKey.APIKey,
		}
		apiKeyRequest := test.apiKey.APIKey
		apiKeyRequest.ID = 0
		apiKeyRequest.APIKey = ""
		test.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			defer func() {
				if err := recover(); err != nil {
					t.Fatal(err)
				}
			}()

			client.When("GetUser", test.user.Username).Return(test.user, nil)

			if test.apiKey.err {
				mockPostAPIKey.Err = errors.New("fake error")
			}
			// if apiKey label is defined (i.e. we expect to be sending a POST),
			// we should mock BuildRequest
			if test.apiKey.Label != "" {
				client.When("BuildRequest", "POST", lib.BrainEndpoint, "/api_keys", []string(nil)).Return(&mockPostAPIKey)
			}
		})
		if !test.ShouldErr {
			t.Run(test.CommandT.Name+" AssertRequestObjectEqual", func(t *testing.T) {
				mockPostAPIKey.T = t
				if test.apiKey.Label != "" {
					mockPostAPIKey.AssertRequestObjectEqual(apiKeyRequest)
				}
			})
		}
	}
}
