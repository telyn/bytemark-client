package main

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

func TestReifyDisc(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("ReifyDisc", 111).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "reify", "disc", "111"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReifyDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("ReifyDisc", 112).Return(fmt.Errorf("Could not reify disc")).Times(1)

	err := app.Run([]string{"bytemark", "reify", "disc", "112"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
