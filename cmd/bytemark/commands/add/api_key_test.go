package add_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/add"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
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
type privWithErr struct {
	brain.Privilege
	err bool
}

func TestAddApiKey(t *testing.T) {
	// here's all the test data for all the servers, groups and accounts we'll
	// use in this test
	serverNames := []lib.VirtualMachineName{
		{VirtualMachine: "myserver", Group: "default", Account: "default-account"},
		{VirtualMachine: "myserver2", Group: "test-group", Account: "test-account"},
	}
	servers := []brain.VirtualMachine{
		{ID: 1},
		{ID: 2},
	}
	groupNames := []pathers.GroupName{
		{Group: "", Account: "default-account"},
		{Group: "test-group", Account: "test-account"},
	}
	groups := []brain.Group{
		{ID: 11},
		{ID: 12},
	}
	accountNames := []string{
		"default-account",
		"test-account",
	}
	accounts := []lib.Account{
		{BrainID: 21},
		{BrainID: 22},
	}

	// ok let's define some tests. basing it on testutil.CommandT to TRY to DRY
	tests := []struct {
		// CommandT provides all the basics to set up a config, client and app,
		// all the basic mocks on config, all the stuff to get fake
		// authentication going, and will run the app and assert a bunch of stuff too
		testutil.CommandT

		// aaand here's all the stuff to assert on and return from Request
		// the map keys are used for the mock BuildRequest calls' URLs
		accounts   map[string]accountWithErr
		groups     map[pathers.GroupName]groupWithErr
		servers    map[lib.VirtualMachineName]serverWithErr
		privileges []privWithErr
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
				regexp.MustCompile("Successfully created an api key:\ntest-api-key\n  Expires: never\n  Key: apikey.ay1pee2eye3key\n"),
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
	}, {
		CommandT: testutil.CommandT{
			Name:      "one server spec",
			Args:      "--server myserver test-api-key",
			ShouldErr: false,
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile("creating 1 privilege"),
			},
		},
		user: brain.User{
			ID:       299,
			Username: "test-user",
		},
		apiKey: apiKeyWithErr{
			APIKey: brain.APIKey{
				ID:    9201,
				Label: "test-api-key",
			},
		},
		privileges: []privWithErr{{
			Privilege: brain.Privilege{
				ID:               9292,
				Username:         "test-user",
				Level:            "vm_admin",
				VirtualMachineID: 1,
			},
		}},
		servers: map[lib.VirtualMachineName]serverWithErr{
			serverNames[0]: serverWithErr{VirtualMachine: servers[0]},
		},
	}, {
		CommandT: testutil.CommandT{
			Name: "two of each type of privilege target",
			Args: "--server myserver --server myserver2.test-group.test-account --group . --group test-group.test-account --account-admin default-account --account-admin test-account test-api-key",
			OutputMustMatch: []*regexp.Regexp{
				regexp.MustCompile(`creating 6 privileges`),
			},
		},
		user: brain.User{
			ID:       299,
			Username: "test-user",
		},
		apiKey: apiKeyWithErr{
			APIKey: brain.APIKey{
				ID:    9202,
				Label: "test-api-key",
			},
		},
		servers: map[lib.VirtualMachineName]serverWithErr{
			serverNames[0]: serverWithErr{VirtualMachine: servers[0]},
			serverNames[1]: serverWithErr{VirtualMachine: servers[1]},
		},
		groups: map[pathers.GroupName]groupWithErr{
			groupNames[0]: groupWithErr{Group: groups[0]},
			groupNames[1]: groupWithErr{Group: groups[1]},
		},
		accounts: map[string]accountWithErr{
			accountNames[0]: accountWithErr{Account: accounts[0]},
			accountNames[1]: accountWithErr{Account: accounts[1]},
		},
		privileges: []privWithErr{{
			Privilege: brain.Privilege{
				ID:               14001,
				Username:         "test-user",
				Level:            "vm_admin",
				VirtualMachineID: 1,
			},
		}, {
			Privilege: brain.Privilege{
				ID:               14002,
				Username:         "test-user",
				Level:            "vm_admin",
				VirtualMachineID: 2,
			},
		}, {
			Privilege: brain.Privilege{
				ID:       14003,
				Username: "test-user",
				Level:    "group_admin",
				GroupID:  11,
			},
		}, {
			Privilege: brain.Privilege{
				ID:       14004,
				Username: "test-user",
				Level:    "group_admin",
				GroupID:  12,
			},
		}, {
			Privilege: brain.Privilege{
				ID:        14005,
				Username:  "test-user",
				Level:     "account_admin",
				AccountID: 21,
			},
		}, {
			Privilege: brain.Privilege{
				ID:        14006,
				Username:  "test-user",
				Level:     "account_admin",
				AccountID: 22,
			},
		}},
	}}

	for _, test := range tests {
		// these are common to all the tests so makes no sense to set them in
		// each test
		test.CommandT.Auth = true
		test.CommandT.Commands = add.Commands
		// hey this way I don't have to type out apikey 80bajillion times
		test.CommandT.Args = "apikey " + test.CommandT.Args

		// we define the mock requests here so that we can assert on them after
		// CommandT.Run
		responseAPIKey := test.apiKey.APIKey
		responseAPIKey.Privileges = brain.Privileges{}
		mockPostAPIKey := mocks.Request{
			StatusCode:     200,
			ResponseObject: responseAPIKey,
		}

		test.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			config.When("GetVirtualMachine").Return(lib.VirtualMachineName{Group: "default", Account: "default-account"})
			config.When("GetGroup").Return(pathers.GroupName{Group: "default", Account: "default-account"})
			config.When("GetIgnoreErr", "account").Return("default-account")
			client.When("GetUser", test.user.Username).Return(test.user, nil)

			if test.apiKey.err {
				mockPostAPIKey.Err = errors.New("fake error")
			}
			for _, privilege := range test.privileges {
				privRequest := privilege.Privilege
				privRequest.ID = 0
				privRequest.APIKeyID = test.apiKey.APIKey.ID
				privRequest.Username = test.user.Username
				err := error(nil)
				if privilege.err {
					err = errors.New("fake error")
				}
				client.When("GrantPrivilege", privRequest).Return(err).Times(1)

			}
			for serverName, serverWErr := range test.servers {
				err := error(nil)
				if serverWErr.err {
					err = errors.New("fake error")
				}
				client.When("GetVirtualMachine", serverName).Return(serverWErr.VirtualMachine, err).Times(1)
			}
			for groupName, groupWErr := range test.groups {
				err := error(nil)
				if groupWErr.err {
					err = errors.New("fake error")
				}
				client.When("GetGroup", groupName).Return(groupWErr.Group, err).Times(1)
			}
			for accountName, accountWErr := range test.accounts {
				err := error(nil)
				if accountWErr.err {
					err = errors.New("fake error")
				}
				client.When("GetAccount", accountName).Return(accountWErr.Account, err).Times(1)
			}
			// if apiKey label is defined (i.e. we expect to be sending a POST),
			// we should mock BuildRequest
			if test.apiKey.Label != "" {
				client.When("BuildRequest", "POST", lib.BrainEndpoint, "/api_keys", []string(nil)).Return(&mockPostAPIKey).Times(1)
			}
		})

		// now assert on all the mock requests
		if !test.ShouldErr {
			t.Run(test.CommandT.Name+" AssertRequestObjectEqual", func(t *testing.T) {
				mockPostAPIKey.T = t
				if test.apiKey.Label != "" {
					apiKeyRequest := test.apiKey.APIKey
					apiKeyRequest.APIKey = ""
					apiKeyRequest.ID = 0
					apiKeyRequest.UserID = test.user.ID
					mockPostAPIKey.AssertRequestObjectEqual(apiKeyRequest)
				}
			})
		}
	}
}
