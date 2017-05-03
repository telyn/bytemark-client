package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestMove(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	oldName := lib.VirtualMachineName{
		VirtualMachine: "old-name",
		Group:          "old-group",
		Account:        "old-account",
	}
	newName := lib.VirtualMachineName{
		VirtualMachine: "new-name",
		Group:          "new-group",
		Account:        "new-account"}

	config.When("GetVirtualMachine").Return(defVM)
	config.When("Force").Return(true)

	c.When("MoveVirtualMachine", oldName, newName).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "move", "old-name.old-group.old-account", "new-name.new-group.new-account"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
