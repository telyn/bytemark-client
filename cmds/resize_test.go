package cmds

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestResizeDisk(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")

	args := []string{"test-server", "11", "35"}
	disc := lib.Disc{
		Size:         25600,
		StorageGrade: "sata",
	}

	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	name := lib.VirtualMachineName{VirtualMachine: "test-server"}
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetDisc", name, "11").Return(&disc).Times(1)

	c.When("ResizeDisc", name, "11", 35*1024).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ResizeDisc(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
