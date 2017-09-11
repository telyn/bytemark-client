package main

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

func TestRegradeDisc(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("RegradeDisc", 111, "newg").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "regrade", "disc", "111", "newg"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestRegradeDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("RegradeDisc", 112, "newg").Return(fmt.Errorf("Could not regrade disc")).Times(1)

	err := app.Run([]string{"bytemark", "regrade", "disc", "112", "newg"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
