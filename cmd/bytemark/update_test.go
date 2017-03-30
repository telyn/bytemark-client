package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"testing"
)

func TestUpdateVMMigrationWithSpeedAndDowntime(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	speed := int64(8500000000000)
	downtime := 15
	c.When("UpdateVMMigration", &vmName, &speed, &downtime).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "update", "vm", "migration", "vm123.group.account", "8500000000000", "15"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateVMMigrationError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetVirtualMachine").Return(&defVM)

	err := global.App.Run([]string{"bytemark", "update", "vm", "migration", "vm124.group.account"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
