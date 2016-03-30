package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestShowGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetGroup").Return(lib.GroupName{})

	c.When("ParseGroupName", "test-group.test-account", []lib.GroupName{{}}).Return(lib.GroupName{Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	group := getFixtureGroup()
	c.When("GetGroup", lib.GroupName{Group: "test-group", Account: "test-account"}).Return(&group, nil).Times(1)

	global.App.Run(strings.Split("bytemark show group test-group.test-account", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	vm := getFixtureVM()
	c.When("GetVirtualMachine", lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}).Return(&vm, nil).Times(1)

	global.App.Run(strings.Split("bytemark show server test-server.test-group.test-account", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
