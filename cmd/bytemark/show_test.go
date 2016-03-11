package main

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestShowGroupCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-group.test-account"})
	config.When("GetGroup").Return(lib.GroupName{})

	c.When("ParseGroupName", "test-group.test-account", []lib.GroupName{{}}).Return(lib.GroupName{Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	group := getFixtureGroup()
	c.When("GetGroup", lib.GroupName{Group: "test-group", Account: "test-account"}).Return(&group, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowGroup([]string{"test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	vm := getFixtureVM()
	c.When("GetVirtualMachine", lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}).Return(&vm, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowServer([]string{"test-server.test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
