package cmds

import (
	bigv "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestShowGroupCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-group.test-account"})

	c.When("ParseGroupName", "test-group.test-account").Return(bigv.GroupName{Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	group := getFixtureGroup()
	c.When("GetGroup", bigv.GroupName{Group: "test-group", Account: "test-account"}).Return(&group, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowGroup([]string{"test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowVMCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})

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
