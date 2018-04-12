package admin_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestSetDiscIOPSLimit(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("Force").Return(true)

	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	c.When("SetDiscIopsLimit", name, "disc-label", 100).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark set iops limit test-server disc-label 100", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
