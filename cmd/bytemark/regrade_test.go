package main

import (
	"fmt"
	"github.com/cheekybits/is"
	"testing"
)

func TestRegradeDisc(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("RegradeDisc", 111, "newg").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "regrade", "disc", "111", "newg"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestRegradeDiscError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("RegradeDisc", 112, "newg").Return(fmt.Errorf("Could not regrade disc")).Times(1)

	err := global.App.Run([]string{"bytemark", "regrade", "disc", "112", "newg"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
