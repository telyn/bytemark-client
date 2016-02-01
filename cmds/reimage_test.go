package cmds

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"github.com/cheekybits/is"
	"testing"
)

func TestReimage(t *testing.T) {
	is := is.New(t)
	c := &mocks.Client{}
	config := &mocks.Config{}

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}

	image := &lib.ImageInstall{
		Distribution:    "symbiosis",
		FirstbootScript: "",
		RootPassword:    "gNFgYYIgayyDOjkV",
		PublicKeys:      "",
	}
	args := []string{"--image", image.Distribution, "--root-password", image.RootPassword, "test-vm.test-group.test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", args[0], []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("ReimageVirtualMachine", vmname, image).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	is.Equal(0, cmds.Reimage(args))

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
