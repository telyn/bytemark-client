package cmds

import (
	util "bytemark.co.uk/client/cmds/util"

	bigv "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"github.com/cheekybits/is"
	"testing"
)

func TestDeleteVM(t *testing.T) {
	is := is.New(t)
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-vm"})
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	name := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	vm := getFixtureVM()

	c.When("ParseVirtualMachineName", "test-vm", []bigv.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", name, false).Return(nil).Times(1)
	cmds := NewCommandSet(config, c)

	is.Equal(util.E_SUCCESS, cmds.DeleteVM([]string{"test-vm"}))
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()

	c.When("ParseVirtualMachineName", "test-vm", []bigv.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", name, true).Return(nil).Times(1)

	is.Equal(util.E_SUCCESS, cmds.DeleteVM([]string{"--purge", "test-vm"}))
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account", "666"})
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	name := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account", []bigv.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("DeleteDisc", name, "666").Return(nil).Times(1)

	cmds := NewCommandSet(config, c)

	is.Equal(util.E_SUCCESS, cmds.DeleteDisc([]string{"--force", "test-vm.test-group.test-account", "666"}))
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
