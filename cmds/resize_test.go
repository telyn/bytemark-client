package cmds

import (
	bigv "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
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
	disc := bigv.Disc{
		Size:         25600,
		StorageGrade: "sata",
	}

	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	name := bigv.VirtualMachineName{VirtualMachine: "test-vm"}
	c.When("ParseVirtualMachineName", "test-vm", []bigv.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetDisc", name, "11").Return(&disc).Times(1)

	c.When("ResizeDisc", name, "11", 35*1024).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ResizeDisc(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
