package cmds

import (
	bigv "bigv.io/client/lib"
	"bigv.io/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestResizeDisk(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}
	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)

	config.When("ImportFlags").Return([]string{"test-vm", "archive:35"})
	name := bigv.VirtualMachineName{VirtualMachine: "test-vm"}
	c.When("ParseVirtualMachineName", "test-vm").Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&bigv.VirtualMachine{Hostname: "test-vm.default.test-user.endpoint"})

	c.When("ResizeDisc", name, 11, 22).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateDiscs([]string{"test-vm", "archive:35"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
