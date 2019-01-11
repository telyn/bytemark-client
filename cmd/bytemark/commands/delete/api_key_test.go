package delete_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/delete"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestDeleteAPIKey(t *testing.T) {
	tests := []struct {
		testutil.CommandT

		apiKeys brain.APIKeys
		id      int
		label   string
	}{{
		CommandT: testutil.CommandT{
			Name: "by id",
			Args: "2",
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("Successfully deleted 2"),
			},
		},
		id: 2,
	}, {
		CommandT: testutil.CommandT{
			Name: "by label",
			Args: "jeff",
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("Successfully deleted jeff"),
			},
		},
		apiKeys: brain.APIKeys{{
			ID:    1,
			Label: "jeffo",
		}, {
			ID:    9,
			Label: "jeff",
		}, {
			ID:    3,
			Label: "jeffery",
		}},
		id:    9,
		label: "jeff",
	}}

	for _, test := range tests {
		test.Args = "apikey " + test.Args
		test.Commands = delete.Commands
		test.Auth = true
		test.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			if test.label != "" {
				apiKeysReq := mocks.Request{
					StatusCode:     200,
					ResponseObject: test.apiKeys,
				}
				client.When("BuildRequest", "GET", lib.BrainEndpoint, "/api_keys?view=overview", []string(nil)).Return(&apiKeysReq, nil)
			}
			deleteReq := mocks.Request{
				StatusCode: 200,
			}
			client.When("BuildRequest", "DELETE", lib.BrainEndpoint, "/api_keys/%s", []string{strconv.Itoa(test.id)}).Return(&deleteReq, nil)
		})
	}
}
