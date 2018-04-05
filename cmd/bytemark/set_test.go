package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestSetCDROM(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("SetVirtualMachineCDROM", vmname, "test-cdrom").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark set cdrom test-server.test-group.test-account test-cdrom", " "))
	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
}
