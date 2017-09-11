package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestListAccounts(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, false)

	c.When("GetAccounts").Return([]lib.Account{lib.Account{BrainID: 1, Name: "dr-evil"}}).Times(1)

	err := app.Run(strings.Split("bytemark list accounts", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListDiscs(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup(t, false)

	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "spooky-vm",
		Group:          "default",
		Account:        "default-account",
	}

	vm := brain.VirtualMachine{
		ID:   4,
		Name: "spooky-vm",
		Discs: []brain.Disc{
			{StorageGrade: "sata", Size: 25600, Label: "vda"},
			{StorageGrade: "archive", Size: 666666, Label: "vdb"},
		},
	}
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)

	err := app.Run(strings.Split("bytemark list discs spooky-vm", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListGroups(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup(t, false)

	config.When("GetIgnoreErr", "account").Return("spooky-steve-other-account")

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Groups: []brain.Group{
			{ID: 1, Name: "halloween-vms"},
			{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	err := app.Run(strings.Split("bytemark list groups spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListServers(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup(t, false)

	config.When("GetIgnoreErr", "account").Return("spokny-stevn")
	config.When("GetGroup").Return(defGroup)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Name: "spooky-steve",
		Groups: []brain.Group{{
			Name: "default",
			VirtualMachines: []brain.VirtualMachine{
				{ID: 1, Name: "old-man-crumbles"},
				{ID: 23, Name: "jack-skellington"},
			},
		}},
	}).Times(1)

	err := app.Run(strings.Split("bytemark list servers spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListBackups(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("GetBackups", vmname, "test-disc").Return(nil).Times(1)

	err := app.Run([]string{
		"bytemark", "list", "backups", "test-server", "test-disc",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
