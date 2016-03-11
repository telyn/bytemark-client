package cmds

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestListAccounts(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-group.test-account"})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccounts").Return([]*lib.Account{&lib.Account{BrainID: 1, Name: "dr-evil"}}).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ListAccounts([]string{})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListDiscs(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"spooky-vm"})
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

	cmds := NewCommandSet(config, c)
	cmds.ListDiscs([]string{"spooky-vm"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListGroups(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "account").Return("spooky-steve-other-account")
	config.When("ImportFlags").Return([]string{"spooky-steve"})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Groups: []*lib.Group{
			&lib.Group{ID: 1, Name: "halloween-vms"},
			&lib.Group{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ListGroups([]string{"spooky-steve"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListServers(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"halloween-vms.spooky-steve"})
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

	cmds := NewCommandSet(config, c)
	cmds.ListServers([]string{"halloween-vms.spooky-steve"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
