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

	args := []string{"test-vm", "11", "35"}

	config.When("ImportFlags").Return(args)
	name := bigv.VirtualMachineName{VirtualMachine: "test-vm"}
	c.When("ParseVirtualMachineName", "test-vm").Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&bigv.VirtualMachine{Hostname: "test-vm.default.test-user.endpoint"})

	c.When("ResizeDisc", name, 11, 35*1024).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ResizeDisc(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
