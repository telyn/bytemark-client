package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestApproveVM(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, false).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "approve", "vm", "vm123.group.account"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestApproveVMAndPowerOn(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm122", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, true).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "approve", "vm", "vm122.group.account", "true"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestApproveVMError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	approveErr := fmt.Errorf("Error approving")
	vmName := lib.VirtualMachineName{VirtualMachine: "vm121", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, false).Return(approveErr).Times(1)

	err := global.App.Run([]string{"bytemark", "approve", "vm", "vm121.group.account"})

	is.Equal(err, approveErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
