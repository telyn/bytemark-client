package add_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/cheekybits/is"
)

func TestCreateDiskCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	config.When("GetVirtualMachine").Return(testutil.DefVM)

	name := pathers.VirtualMachineName{VirtualMachine: "test-server", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}}
	c.When("GetVirtualMachine", name).Return(&brain.VirtualMachine{Hostname: "test-server.default.default-account.endpoint"})

	disc := brain.Disc{Size: 35 * 1024, StorageGrade: "archive"}

	c.When("CreateDisc", name, disc).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark add disc --force --disc archive:35 test-server", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
