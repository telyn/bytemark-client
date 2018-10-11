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

func TestRevokePrivilege(t *testing.T) {

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
			Input:     "bytemark revoke privilege",
		}, {
			Name: "BadPrivilege",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// reset to get rid of the default AuthWithToken expectation
				c.Reset()
			},
			ShouldErr: true,
			Input:     "bytemark revoke privilege smedly",
		}, {
			Name: "MissingAuth",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// reset to get rid of the default AuthWithToken expectation
				c.Reset()
			},
			ShouldErr: true,
			Input:     "bytemark revoke privilege cluster_admin on bucholic to no-one",
		}, {
			Name: "RevokeVMAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetVirtualMachine").Return(defVM)
				vm := lib.VirtualMachineName{
					VirtualMachine: "test-vm",
					Group:          "test-group",
					Account:        "test-account",
				}

				c.When("GetVirtualMachine", vm).Return(brain.VirtualMachine{ID: 333}, nil).Times(1)
				c.When("GetPrivilegesForVirtualMachine", vm).Return(brain.Privileges{
					{
						ID:               2342,
						Username:         "burt",
						Level:            brain.VMAdminPrivilege,
						VirtualMachineID: 333,
					}, {
						ID:               9823,
						Username:         "test-user",
						Level:            brain.VMAdminPrivilege,
						VirtualMachineID: 333,
					},
				})

				c.When("RevokePrivilege", brain.Privilege{
					ID:               9823,
					Username:         "test-user",
					Level:            brain.VMAdminPrivilege,
					VirtualMachineID: 333,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark revoke privilege vm_admin on test-vm.test-group.test-account from test-user",
		}, {
			Name: "RevokeGroupAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetGroup").Return(defGroup)
				group := lib.GroupName{
					Group:   "test-group",
					Account: "test-account",
				}

				c.When("GetGroup", group).Return(brain.Group{ID: 953}, nil).Times(1)
				c.When("GetPrivilegesForGroup", group).Return(brain.Privileges{
					{
						ID:       4354,
						Username: "burt",
						Level:    brain.GroupAdminPrivilege,
						GroupID:  953,
					}, {
						ID:       32647,
						Username: "test-user",
						Level:    brain.GroupAdminPrivilege,
						GroupID:  953,
					},
				})

				c.When("RevokePrivilege", brain.Privilege{
					ID:       32647,
					Username: "test-user",
					Level:    brain.GroupAdminPrivilege,
					GroupID:  953,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark revoke privilege group_admin on test-group.test-account from test-user",
		}, {
			Name: "RevokeAccountAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetIgnoreErr", "account").Return("default-account")
				acc := "default-account"

				c.When("GetAccount", acc).Return(lib.Account{BrainID: 223435}, nil).Times(1)
				c.When("GetPrivilegesForAccount", acc).Return(brain.Privileges{
					{
						ID:        12412,
						Username:  "burt",
						Level:     brain.AccountAdminPrivilege,
						AccountID: 223435,
					}, {
						ID:        129865,
						Username:  "test-user",
						Level:     brain.AccountAdminPrivilege,
						AccountID: 223435,
					},
				})

				c.When("RevokePrivilege", brain.Privilege{
					ID:        129865,
					Username:  "test-user",
					Level:     brain.AccountAdminPrivilege,
					AccountID: 223435,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			// secret feature: '.' as alias for default account
			Input: "bytemark revoke privilege account_admin on . to test-user",
		}, {
			Name: "RevokeClusterAdmin",
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				c.When("GetPrivileges", "test-user").Return(brain.Privileges{
					{
						ID:        43634,
						Username:  "test-user",
						Level:     brain.AccountAdminPrivilege,
						AccountID: 223435,
					}, {
						ID:              3245,
						Username:        "test-user",
						YubikeyRequired: true,
						Level:           brain.ClusterAdminPrivilege,
					},
				})

				c.When("RevokePrivilege", brain.Privilege{
					ID:              3245,
					Username:        "test-user",
					Level:           brain.ClusterAdminPrivilege,
					YubikeyRequired: true,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark revoke privilege --yubikey-required cluster_admin to test-user",
		},
	}
	for i, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			test.Setup(config, c)

			err := app.Run(strings.Split(test.Input, " "))
			if test.ShouldErr && err == nil {
				t.Errorf("TestRevokePrivilege %d should err and didn't", i)
			} else if !test.ShouldErr && err != nil {
				t.Errorf("TestRevokePrivilege %d shouldn't err, but: %s", i, err)
			}
			if ok, err := c.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
