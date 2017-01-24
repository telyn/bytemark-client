package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"strings"
	"testing"
)

func TestGrantPrivilege(t *testing.T) {
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetVirtualMachine").Return(&lib.VirtualMachineName{})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account", []*lib.VirtualMachineName{{}}).Return(&lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}, nil)
	c.When("GetVirtualMachine", &lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}).Return(&brain.VirtualMachine{ID: 333}, nil).Times(1)

	c.When("GrantPrivilege", brain.Privilege{
		Username:         "test-user",
		Level:            brain.VMAdminPrivilege,
		VirtualMachineID: 333,
	}).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark grant vm_admin on test-vm.test-group.test-account to test-user", " "))
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
