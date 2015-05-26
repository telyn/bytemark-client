package main

import (
	bigv "bigv.io/client/lib"
	"testing"
	//"github.com/cheekybits/is"
)

func getFixtureVM() bigv.VirtualMachine {
	return bigv.VirtualMachine{
		Name:    "test-vm",
		GroupID: 1,
	}
}

func TestCommandConfig(t *testing.T) {
	config := &mockConfig{}

	config.When("GetV", "user").Return(ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")
	config.When("Get", "silent").Return("true")

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	cmds := NewCommandSet(config, nil)
	cmds.Config([]string{"set", "user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}

// everything else is going to involve making a mock client
func TestShowVMCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	config.When("Get", "token").Return("test-token")
	config.When("Get", "silent").Return("true")

	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	vm := getFixtureVM()
	c.When("GetVirtualMachine", bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}).Return(&vm, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowVM([]string{"test-vm.test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
