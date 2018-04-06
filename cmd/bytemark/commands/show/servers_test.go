package show_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestShowServers(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	config.When("GetIgnoreErr", "account").Return("spokny-stevn")
	config.When("GetGroup").Return(testutil.DefGroup)

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Name: "spooky-steve",
		Groups: []brain.Group{{
			Name: "default",
			VirtualMachines: []brain.VirtualMachine{
				{ID: 1, Name: "old-man-crumbles"},
				{ID: 23, Name: "jack-skellington"},
			},
		}},
	}).Times(1)

	err := app.Run(strings.Split("bytemark show servers spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
