package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestListAccounts(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccounts").Return([]*lib.Account{&lib.Account{BrainID: 1, Name: "dr-evil"}}).Times(1)

	global.App.Run(strings.Split("bytemark list accounts", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListDiscs(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	name := lib.VirtualMachineName{
		VirtualMachine: "spooky-vm",
		Group:          "halloween-vms",
		Account:        "",
	}
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("ParseVirtualMachineName", "spooky-vm", []lib.VirtualMachineName{{}}).Return(name).Times(1)

	vm := lib.VirtualMachine{
		ID:   4,
		Name: "spooky-vm",
		Discs: []*lib.Disc{
			&lib.Disc{StorageGrade: "sata", Size: 25600},
			&lib.Disc{StorageGrade: "archive", Size: 666666},
		},
	}
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)

	global.App.Run(strings.Split("bytemark list discs spooky-vm", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListGroups(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "account").Return("spooky-steve-other-account")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Groups: []*lib.Group{
			&lib.Group{ID: 1, Name: "halloween-vms"},
			&lib.Group{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	global.App.Run(strings.Split("bytemark list groups spooky-steve", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListServers(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetGroup").Return(lib.GroupName{})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	groupname := lib.GroupName{Group: "halloween-vms", Account: "spooky-steve"}
	c.When("ParseGroupName", "halloween-vms.spooky-steve", []lib.GroupName{{}}).Return(groupname).Times(1)

	c.When("GetGroup", groupname).Return(&lib.Group{
		VirtualMachines: []*lib.VirtualMachine{
			&lib.VirtualMachine{ID: 1, Name: "old-man-crumbles"},
			&lib.VirtualMachine{ID: 23, Name: "jack-skellington"},
		},
	}).Times(1)

	global.App.Run(strings.Split("bytemark list servers halloween-vms.spooky-steve", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
