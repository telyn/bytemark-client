package admin_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestMigrateDiscWithNewStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("MigrateDisc", 123, "t6-sata1").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscWithoutNewStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("MigrateDisc", 123, "").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	migrateErr := fmt.Errorf("Error migrating")
	c.When("MigrateDisc", 123, "t6-sata1").Return(migrateErr).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Equal(err, migrateErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineWithNewHead(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	c.When("MigrateVirtualMachine", vmName, "stg-h1").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "vm", "vm123.group.account", "stg-h1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineWithoutNewHead(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm122", Group: "group", Account: "account"}
	c.When("MigrateVirtualMachine", vmName, "").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "vm", "vm122.group.account"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineError(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	migrateErr := fmt.Errorf("Error migrating")
	vmName := lib.VirtualMachineName{VirtualMachine: "vm121", Group: "group", Account: "account"}
	c.When("MigrateVirtualMachine", vmName, "stg-h2").Return(migrateErr).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "vm", "vm121.group.account", "stg-h2"})

	is.Equal(err, migrateErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
