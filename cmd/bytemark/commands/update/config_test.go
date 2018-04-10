package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	cf "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestUpdateConfig(t *testing.T) {
	tests := []struct {
		name        string
		args        string
		expectation func(*mocks.Config)
		shouldErr   bool
	}{
		{
			name: "SetUser",
			args: "--user fred",
			expectation: func(config *mocks.Config) {
			    before := cf.Var{"user", "joan", ""}
			    config.When("GetV", "user").Return(before).Times(1)
			    config.When("SetPersistent", "user", "fred", "CMD set")
			},
		},
		{
			name: "UnsetUser",
			args: "--unset-user",
			expectation: func(config *mocks.Config) {
			    config.When("Unset", "user")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, _, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			test.expectation(config)
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
