package show

import (
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestShowApiKeys(t *testing.T) {
	tests := []struct {
		// CommandT provides all the basics to set up a config, client and app,
		// all the basic mocks on config, all the stuff to get fake
		// authentication going, and will run the app and assert a bunch of stuff too
		testutil.CommandT

		apiKeys brain.APIKeys
		user    brain.User
		err     error
	}{{
		CommandT: testutil.CommandT{
			Name: "api keys are listed",
			Args: "apikeys",
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("special-key.*group jeffgroup"),
			},
		},
		user: brain.User{
			ID:       100,
			Username: "jeff",
		},
		apiKeys: brain.APIKeys{{
			UserID: 100,
			Label:  "special-key",
			Privileges: brain.Privileges{brain.Privilege{
				Username:  "jeff",
				Level:     "group_admin",
				GroupName: "jeffgroup",
			}},
		}},
	}}
	for _, test := range tests {
		test.CommandT.Commands = Commands
		test.CommandT.Auth = true
		mockRequest := mocks.Request{
			StatusCode:     200,
			ResponseObject: test.apiKeys,
			Err:            nil,
		}
		test.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			client.When("BuildRequest", "GET", lib.BrainEndpoint, "/api_keys?view=overview", []string(nil)).Return(&mockRequest, nil)
			client.When("GetUser", "test-user").Return(test.user, nil)
			return
		})
	}
}
