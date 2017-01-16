package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestRestoreSnapshot(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "",
		Account:        "",
	}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&lib.VirtualMachineName{"", "", ""})

	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("RestoreSnapshot", vmname, "test-disc", "test-snapshot").Return(nil).Times(1)

	err := global.App.Run([]string{
		"bytemark", "restore", "snapshot", "test-server", "test-disc", "test-snapshot",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
