package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestReimage(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	image := &lib.ImageInstall{
		Distribution:    "symbiosis",
		FirstbootScript: "",
		RootPassword:    "gNFgYYIgayyDOjkV",
		PublicKeys:      "",
	}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)
	config.When("Force").Return(true)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("ReimageVirtualMachine", &vmname, image).Return(nil).Times(1)

	global.App.Run([]string{"bytemark", "reimage", "--image", image.Distribution, "--root-password", image.RootPassword, "test-server.test-group.test-account"})

	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
