package add_test

import (
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
)

func TestCreateDiskCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "default", Account: "default-account"}
	c.When("GetVirtualMachine", name).Return(&brain.VirtualMachine{Hostname: "test-server.default.default-account.endpoint"})

	disc := brain.Disc{Size: 35 * 1024, StorageGrade: "archive"}

	c.When("CreateDisc", name, disc).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create disc --force --disc archive:35 test-server", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
