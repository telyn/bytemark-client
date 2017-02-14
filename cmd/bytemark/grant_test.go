package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"strings"
	"testing"
)

// baseTestAuthSetup sets up a 'regular' test - with auth, no yubikey.
// user is test-user
func baseTestAuthSetup(t *testing.T, admin bool) (config *mocks.Config, c *mocks.Client) {
	config, c = baseTestSetup(t, admin)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return config, c
}

func TestGrantPrivilege(t *testing.T) {

	tests := []struct {
		Setup     func(config *mocks.Config, c *mocks.Client)
		ShouldErr bool
		Input     string
	}{
		{
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetGroup").Return(&defGroup)
				group := lib.GroupName{
					Group:   "test-group",
					Account: "test-account",
				}


				c.When("GetGroup", &group).Return(&brain.Group{
					ID: 303,
				}).Times(1)
				c.When("GrantPrivilege", brain.Privilege{
					Username: "test-user",
					Level:    brain.GroupAdminPrivilege,
					GroupID:  303,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant group_admin on test-group.test-account to test-user",
		},
		{
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetVirtualMachine").Return(&defVM)
				vm := lib.VirtualMachineName{
					VirtualMachine: "test-vm",
					Group:          "test-group",
					Account:        "test-account",
				}

				c.When("GetVirtualMachine", &vm).Return(&brain.VirtualMachine{ID: 333}, nil).Times(1)

				c.When("GrantPrivilege", brain.Privilege{
					Username:         "test-user",
					Level:            brain.VMAdminPrivilege,
					VirtualMachineID: 333,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant vm_admin on test-vm.test-group.test-account to test-user",
		},
		{
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console
				config.When("GetIgnoreErr", "account").Return("default-account")

				c.When("ParseAccountName", "test-account", []string{"default-account"}).Return("test-account")

				c.When("GetAccount", "test-account").Return(&lib.Account{
					BrainID: 32310,
				})

				c.When("GrantPrivilege", brain.Privilege{
					Username:  "test-user",
					Level:     brain.AccountAdminPrivilege,
					AccountID: 32310,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark grant account_admin on test-account to test-user",
		},
	}
	for i, test := range tests {
		config, c := baseTestAuthSetup(t, true)
		test.Setup(config, c)

		err := global.App.Run(strings.Split(test.Input, " "))
		if test.ShouldErr && err == nil {
			t.Errorf("TestGrantPrivilege %d should err and didn't", i)
		} else if !test.ShouldErr && err != nil {
			t.Errorf("TestGrantPrivilege %d shouldn't err, but: %s", i, err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
