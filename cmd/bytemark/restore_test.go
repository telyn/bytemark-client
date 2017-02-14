package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

//TODO(telyn): add test for restore server

func TestRestoreBackup(t *testing.T) {
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

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("RestoreBackup", vmname, "test-disc", "test-backup").Return(nil).Times(1)

	err := global.App.Run([]string{
		"bytemark", "restore", "backup", "test-server", "test-disc", "test-backup",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
