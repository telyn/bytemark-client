package main

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestUpdateVMMigrationWithSpeedAndDowntime(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup (t, true)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	speed := int64(8500000000000)
	downtime := 15
	c.When("UpdateVMMigration", vmName, &speed, &downtime).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "vm", "migration", "vm123.group.account", "8500000000000", "15"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateVMMigrationError(t *testing.T) {
	is := is.New(t)
	config, c, app := baseTestAuthSetup (t, true)

	config.When("GetVirtualMachine").Return(defVM)

	err := app.Run([]string{"bytemark", "update", "vm", "migration", "vm124.group.account"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateHead(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateHead{}

	c.When("UpdateHead", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "head", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateHeadError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateHead{}

	c.When("UpdateHead", "1", options).Return(fmt.Errorf("Could not update head")).Times(1)

	err := app.Run([]string{"bytemark", "update", "head", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateTail(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateTail{}

	c.When("UpdateTail", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "tail", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateTailError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateTail{}

	c.When("UpdateTail", "1", options).Return(fmt.Errorf("Could not update tail")).Times(1)

	err := app.Run([]string{"bytemark", "update", "tail", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateStoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage_pool", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateStoragePoolError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup (t, true)

	options := lib.UpdateStoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(fmt.Errorf("Could not update storage pool")).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage_pool", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
