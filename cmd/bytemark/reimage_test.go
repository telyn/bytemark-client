package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"io/ioutil"
	"os"
	"testing"
)

func TestReimage(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	image := &brain.ImageInstall{
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

	err := global.App.Run([]string{"bytemark", "reimage", "--force", "--image", image.Distribution, "--root-password", image.RootPassword, "test-server.test-group.test-account"})

	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReimageFileFlags(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	image := &brain.ImageInstall{
		FirstbootScript: "i am the firstboot script! FEAR ME",
		PublicKeys:      "i am the authorized keys",
		Distribution:    "image",
		RootPassword:    "test-pass",
	}

	err := ioutil.WriteFile("firstboot", []byte("i am the firstboot script! FEAR ME"), 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("firstboot")
	err = ioutil.WriteFile("authorized-keys", []byte("i am the authorized keys"), 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("authorized-keys")

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)
	config.When("Force").Return(true)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("ReimageVirtualMachine", &vmname, image).Return(nil).Times(1)

	err = global.App.Run([]string{"bytemark", "reimage", "--force", "--image", "image", "--root-password", "test-pass", "--firstboot-script-file", "firstboot", "--authorized-keys-file", "authorized-keys", "test-server.test-group.test-account"})

	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
