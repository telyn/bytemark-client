package admin_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/cheekybits/is"
)

func TestRegradeDisc(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("RegradeDisc", 111, "newg").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "regrade", "disc", "111", "newg"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestRegradeDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("RegradeDisc", 112, "newg").Return(fmt.Errorf("Could not regrade disc")).Times(1)

	err := app.Run([]string{"bytemark", "regrade", "disc", "112", "newg"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
