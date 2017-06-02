package main

import (
	"fmt"
	"github.com/cheekybits/is"
	"testing"
)

func TestReapVMs(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "reap", "servers"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReapVMsAlias(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "reap", "vms"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReapVMsError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(fmt.Errorf("Error reaping VMs")).Times(1)

	err := global.App.Run([]string{"bytemark", "reap", "vms"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
