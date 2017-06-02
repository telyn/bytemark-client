package main

import (
	"fmt"
	"github.com/cheekybits/is"
	"testing"
)

func TestCancelDiscMigration(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CancelDiscMigration", 123).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "cancel", "migration", "disc", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelDiscMigrationError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CancelDiscMigration", 122).Return(fmt.Errorf("Error canceling migrations")).Times(1)

	err := global.App.Run([]string{"bytemark", "cancel", "migration", "disc", "122"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelVMMigration(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CancelVMMigration", 129).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "cancel", "migration", "server", "129"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelVMMigrationAlias(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CancelVMMigration", 128).Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "cancel", "migration", "vm", "128"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestCancelVMMigrationError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CancelVMMigration", 127).Return(fmt.Errorf("Error canceling migrations")).Times(1)

	err := global.App.Run([]string{"bytemark", "cancel", "migration", "server", "127"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
