package show_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestShowDiscs(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	config.When("GetVirtualMachine").Return(testutil.DefVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "spooky-vm",
		Group:          "default",
		Account:        "default-account",
	}

	vm := brain.VirtualMachine{
		ID:   4,
		Name: "spooky-vm",
		Discs: []brain.Disc{
			{StorageGrade: "sata", Size: 25600, Label: "vda"},
			{StorageGrade: "archive", Size: 666666, Label: "vdb"},
		},
	}
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)

	err := app.Run(strings.Split("bytemark show discs spooky-vm", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
