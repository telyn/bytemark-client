package main

import (
	"fmt"
	"github.com/cheekybits/is"
	mock "github.com/maraino/go-mock"
	"testing"
)

func TestMigrateDiscWithNewStoragePool(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("MigrateDisc", 123, "t6-sata1").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscWithoutNewStoragePool(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("MigrateDisc", 123, "").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "disc", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	migrateErr := fmt.Errorf("Error migrating")
	c.When("MigrateDisc", 123, "t6-sata1").Return(migrateErr).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Equal(err, migrateErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineWithNewHead(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("MigrateVirtualMachine", mock.Any, "stg-h1").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "vm", "123", "stg-h1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineWithoutNewHead(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("MigrateVirtualMachine", mock.Any, "").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "vm", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateVirtualMachineError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, true)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	migrateErr := fmt.Errorf("Error migrating")
	c.When("MigrateVirtualMachine", mock.Any, "stg-h2").Return(migrateErr).Times(1)

	err := global.App.Run([]string{"bytemark", "migrate", "vm", "123", "stg-h2"})

	is.Equal(err, migrateErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
