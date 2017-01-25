package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"strings"
	"testing"
)

func TestRevokePrivilege(t *testing.T) {
	tests := []struct {
		Setup     func(config *mocks.Config, c *mocks.Client)
		ShouldErr bool
		Input     string
	}{
		{
			Setup: func(config *mocks.Config, c *mocks.Client) {
				// specific to vm_admin/vm_console

				config.When("GetVirtualMachine").Return(&defVM)
				vm := lib.VirtualMachineName{
					VirtualMachine: "test-vm",
					Group:          "test-group",
					Account:        "test-account",
				}

				c.When("ParseVirtualMachineName", "test-vm.test-group.test-account", []*lib.VirtualMachineName{{}}).Return(&vm, nil)
				c.When("GetVirtualMachine", &vm).Return(&brain.VirtualMachine{ID: 333}, nil).Times(1)

				c.When("RevokePrivilege", brain.Privilege{
					Username:         "test-user",
					Level:            brain.VMAdminPrivilege,
					VirtualMachineID: 333,
				}).Return(nil).Times(1)
			},
			ShouldErr: false,
			Input:     "bytemark revoke vm_admin on test-vm.test-group.test-account to test-user",
		},
	}
	for i, test := range tests {
		config, c := baseTestAuthSetup(t)
		test.Setup(config, c)

		err := global.App.Run(strings.Split(test.Input, " "))
		if test.ShouldErr && err == nil {
			t.Errorf("TestRevokePrivilege %d should err and didn't", i)
		} else if !test.ShouldErr && err != nil {
			t.Errorf("TestRevokePrivilege %d shouldn't err, but: %s", i, err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
