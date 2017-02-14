package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestListAccounts(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccounts").Return([]*lib.Account{&lib.Account{BrainID: 1, Name: "dr-evil"}}).Times(1)

	err := global.App.Run(strings.Split("bytemark list accounts", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListDiscs(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "spooky-vm",
		Group:          "default",
		Account:        "default-account",
	}
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := brain.VirtualMachine{
		ID:   4,
		Name: "spooky-vm",
		Discs: []*brain.Disc{
			&brain.Disc{StorageGrade: "sata", Size: 25600, Label: "vda"},
			&brain.Disc{StorageGrade: "archive", Size: 666666, Label: "vdb"},
		},
	}
	c.When("GetVirtualMachine", &name).Return(&vm).Times(1)

	err := global.App.Run(strings.Split("bytemark list discs spooky-vm", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListGroups(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "account").Return("spooky-steve-other-account")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Groups: []*brain.Group{
			&brain.Group{ID: 1, Name: "halloween-vms"},
			&brain.Group{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	err := global.App.Run(strings.Split("bytemark list groups spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListServers(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "account").Return("spokny-stevn")
	config.When("GetGroup").Return(&defGroup)

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Name: "spooky-steve",
		Groups: []*brain.Group{{
			Name: "default",
			VirtualMachines: []*brain.VirtualMachine{
				&brain.VirtualMachine{ID: 1, Name: "old-man-crumbles"},
				&brain.VirtualMachine{ID: 23, Name: "jack-skellington"},
			},
		}},
	}).Times(1)

	err := global.App.Run(strings.Split("bytemark list servers spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListBackups(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "",
		Account:        "",
	}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&lib.VirtualMachineName{"", "", ""})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetBackups", vmname, "test-disc").Return(nil).Times(1)

	err := global.App.Run([]string{
		"bytemark", "list", "backups", "test-server", "test-disc",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
