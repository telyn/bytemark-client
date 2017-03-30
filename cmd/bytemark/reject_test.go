package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestRejectVM(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	c.When("RejectVM", &vmName, "reason text").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "reject", "vm", "vm123.group.account", "reason text"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestRejectVMError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	rejectErr := fmt.Errorf("Error rejecting")
	vmName := lib.VirtualMachineName{VirtualMachine: "vm121", Group: "group", Account: "account"}
	c.When("RejectVM", &vmName, "reason text").Return(rejectErr).Times(1)

	err := global.App.Run([]string{"bytemark", "reject", "vm", "vm121.group.account", "reason text"})

	is.Equal(err, rejectErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
