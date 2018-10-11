package migrate_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/cheekybits/is"
)

func TestMigrateDiscWithNewStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("MigrateDisc", 123, "t6-sata1").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscWithoutNewStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("MigrateDisc", 123, "").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestMigrateDiscError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	migrateErr := fmt.Errorf("Error migrating")
	c.When("MigrateDisc", 123, "t6-sata1").Return(migrateErr).Times(1)

	err := app.Run([]string{"bytemark", "migrate", "disc", "123", "t6-sata1"})

	is.Equal(err, migrateErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
