package show

import (
	"errors"
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
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
			Name:      "err is returned",
			ShouldErr: true,
		},
		err: errors.New("whatever"),
	}, {
		CommandT: testutil.CommandT{
			Name:            "api keys are listed",
			OutputMustMatch: []*regexp.Regexp{},
		},
		user: brain.User{
			ID:       100,
			Username: "jeff",
		},
	}}
	for _, test := range tests {
		test.Run(t, func(config *mocks.Config, client *mocks.Client, app *cli.App) error {
			return nil
		})
	}
}
