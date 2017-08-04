package main

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

func TestEmptyStoragePool(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("EmptyStoragePool", "pool1").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "empty", "storage_pool", "pool1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestEmptyStoragePoolError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("EmptyStoragePool", "pool1").Return(fmt.Errorf("Could not empty storage pool")).Times(1)

	err := global.App.Run([]string{"bytemark", "empty", "storage_pool", "pool1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestEmptyHead(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("EmptyHead", "pool1").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "empty", "head", "pool1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestEmptyHeadError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("EmptyHead", "pool1").Return(fmt.Errorf("Could not empty storage pool")).Times(1)

	err := global.App.Run([]string{"bytemark", "empty", "head", "pool1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
