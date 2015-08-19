package cmds

import (
	bigv "bigv.io/client/lib"
	"bigv.io/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestListAccounts(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-group.test-account"})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccounts").Return([]*bigv.Account{&bigv.Account{ID: 1, Name: "Dr. Evil"}}).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ListAccounts([]string{})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListGroups(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"spooky-steve"})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("GetAccount", "spooky-steve").Return(&bigv.Account{
		Groups: []*bigv.Group{
			&bigv.Group{ID: 1, Name: "halloween-vms"},
			&bigv.Group{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ListGroups([]string{"spooky-steve"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestListVMs(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"halloween-vms.spooky-steve"})

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	groupname := bigv.GroupName{Group: "halloween-vms", Account: "spooky-steve"}
	c.When("ParseGroupName", "halloween-vms.spooky-steve").Return(groupname).Times(1)

	c.When("GetGroup", groupname).Return(&bigv.Group{
		VirtualMachines: []*bigv.VirtualMachine{
			&bigv.VirtualMachine{ID: 1, Name: "old-man-crumbles"},
			&bigv.VirtualMachine{ID: 23, Name: "jack-skellington"},
		},
	}).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ListVMs([]string{"halloween-vms.spooky-steve"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
