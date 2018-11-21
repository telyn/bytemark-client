package commands_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestGrantPrivilege(t *testing.T) {

	defVM := lib.VirtualMachineName{Group: "default", Account: "default-account"}
	defGroup := lib.GroupName{Group: "default", Account: "default-account"}

	tests := []struct {
		Name      string
		Setup     func(config *mocks.Config, c *mocks.Client)
		ShouldErr bool
		Input     string
	}{
		{
			Name: "MissingArguments",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// reset to get rid of the default AuthWithToken expectation
				c.Reset()
			},
			ShouldErr: true,
			Input:     "bytemark grant privilege",
		}, {
			Name: "BadPrivilege",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// reset to get rid of the default AuthWithToken expectation
				c.Reset()
			},
			ShouldErr: true,
			Input:     "bytemark grant privilege smedly",
		}, {
			Name: "MissingAuth",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// reset to get rid of the default AuthWithToken expectation
				c.Reset()
			},
			ShouldErr: true,
			Input:     "bytemark grant privilege cluster_admin on bucholic to no-one",
		}, {
			Name: "GrantGroupAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetGroup").Return(defGroup)
				group := lib.GroupName{
					Group:   "test-group",
					Account: "test-account",
				}

				c.When("GetGroup", group).Return(brain.Group{
					ID: 303,
				}).Times(1)
				c.When("GrantPrivilege", brain.Privilege{
					Username: "test-user",
					Level:    brain.GroupAdminPrivilege,
					GroupID:  303,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant privilege group_admin test-group.test-account test-user",
		},
		{
			Name: "GrantGroupAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetGroup").Return(defGroup)
				group := lib.GroupName{
					Group:   "test-group",
					Account: "test-account",
				}

				c.When("GetGroup", group).Return(brain.Group{
					ID: 303,
				}).Times(1)
				c.When("GrantPrivilege", brain.Privilege{
					Username: "test-user",
					Level:    brain.GroupAdminPrivilege,
					GroupID:  303,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant privilege group_admin on test-group.test-account to test-user",
		},
		{
			Name: "GrantVMAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetVirtualMachine").Return(defVM)
				vm := lib.VirtualMachineName{
					VirtualMachine: "test-vm",
					Group:          "test-group",
					Account:        "test-account",
				}

				c.When("GetVirtualMachine", vm).Return(brain.VirtualMachine{ID: 333}, nil).Times(1)

				c.When("GrantPrivilege", brain.Privilege{
					Username:         "test-user",
					Level:            brain.VMAdminPrivilege,
					VirtualMachineID: 333,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant privilege vm_admin on test-vm.test-group.test-account to test-user",
		},
		{
			Name: "AccountVMAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console
				config.When("GetIgnoreErr", "account").Return("default-account")

				c.When("GetAccount", "test-account").Return(lib.Account{
					BrainID: 32310,
				})

				c.When("GrantPrivilege", brain.Privilege{
					Username:  "test-user",
					Level:     brain.AccountAdminPrivilege,
					AccountID: 32310,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant privilege account_admin on test-account to test-user",
		}, {
			Name: "ApiKey",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				config.When("GetIgnoreErr", "account").Return("default-account")
				c.When("GetAccount", "account").Return(lib.Account{BrainID: 32310})
				c.When("GrantPrivilege", brain.Privilege{
					Username:  "user",
					AccountID: 32310,
					Level:     brain.AccountAdminPrivilege,
					APIKeyID:  4,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant privilege --api-key-id 4 account_admin on account to user",
		},
	}
	for i, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			test.Setup(config, c)

			err := app.Run(strings.Split(test.Input, " "))
			if test.ShouldErr && err == nil {
				t.Errorf("TestGrantPrivilege %d should err and didn't", i)
			} else if !test.ShouldErr && err != nil {
				t.Errorf("TestGrantPrivilege %d shouldn't err, but: %s", i, err)
			}
			if ok, err := c.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
