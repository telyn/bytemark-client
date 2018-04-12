package admin_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/cheekybits/is"
)

func TestReifyDisc(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("ReifyDisc", 111).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "reify", "disc", "111"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReifyDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("ReifyDisc", 112).Return(fmt.Errorf("Could not reify disc")).Times(1)

	err := app.Run([]string{"bytemark", "reify", "disc", "112"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
