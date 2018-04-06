package show_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestAdminShowRecentServersCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetRecentVMs").Return(&vms, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show recent servers", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
