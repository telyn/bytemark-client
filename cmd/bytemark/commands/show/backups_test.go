package show_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestShowBackups(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	config.When("GetVirtualMachine").Return(testutil.DefVM)

	c.When("GetBackups", vmname, "test-disc").Return(nil).Times(1)

	err := app.Run([]string{
		"bytemark", "show", "backups", "test-server", "test-disc",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
