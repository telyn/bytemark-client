package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	cf "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestUpdateConfig(t *testing.T) {
	tests := []struct {
		name        string
		args        string
		expectation func(*mocks.Config, *mocks.Client)
		shouldErr   bool
	}{
		{
			name: "SetUser",
			args: "--user fred",
			expectation: func(config *mocks.Config, _ *mocks.Client) {
				before := cf.Var{"user", "joan", ""}
				config.When("GetV", "user").Return(before).Times(1)
				config.When("SetPersistent", "user", "fred", "CMD set")
			},
		},
		{
			name: "UnsetUser",
			args: "--unset-user",
			expectation: func(config *mocks.Config, _ *mocks.Client) {
				config.When("Unset", "user")
			},
		},
		{
			name: "SetAccount",
			args: "--account smythe",
			expectation: func(config *mocks.Config, client *mocks.Client) {
				before := cf.Var{"account", "not-smythe", ""}
				smythe := lib.Account{}
				config.When("GetV", "account").Return(before).Times(1)
				config.When("GetIgnoreErr", "account").Return(before.Value).Times(1)
				client.When("GetAccount", "smythe").Return(smythe, nil).Times(1)
				config.When("SetPersistent", "account", "smythe", "CMD set").Times(1)
			},
		},
		{
			name: "SetAccountNoBilling",
			args: "--account smythe",
			expectation: func(config *mocks.Config, client *mocks.Client) {
				before := cf.Var{"account", "not-smythe", ""}
				smythe := lib.Account{}
				config.When("GetV", "account").Return(before).Times(1)
				config.When("GetIgnoreErr", "account").Return(before.Value).Times(1)
				client.When("GetAccount", "smythe").Return(smythe, lib.BillingAccountNotFound("smythe")).Times(1)
				config.When("SetPersistent", "account", "smythe", "CMD set").Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			test.expectation(config, client)
			args := strings.Split("bytemark update config "+test.args, " ")
			err := app.Run(args)
			if test.shouldErr && err == nil {
				t.Fatal("should error")
			} else if !test.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
