package admin_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/cheekybits/is"
)

func TestCancelDiscMigration(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CancelDiscMigration", 123).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "cancel", "migration", "--disc", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelDiscMigrationError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CancelDiscMigration", 122).Return(fmt.Errorf("Error canceling migrations")).Times(1)

	err := app.Run([]string{"bytemark", "cancel", "migration", "--disc", "122"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelVMMigration(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CancelVMMigration", 129).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "cancel", "migration", "--server", "129"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelVMMigrationError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CancelVMMigration", 127).Return(fmt.Errorf("Error canceling migrations")).Times(1)

	err := app.Run([]string{"bytemark", "cancel", "migration", "--server", "127"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCancelBothError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	err := app.Run([]string{"bytemark", "cancel", "migration", "--disc", "1", "--server", "127"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
