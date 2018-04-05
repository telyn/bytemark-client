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

func TestShowGroups(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	config.When("GetIgnoreErr", "account").Return("spooky-steve-other-account")

	c.When("GetAccount", "spooky-steve").Return(&lib.Account{
		Groups: []brain.Group{
			{ID: 1, Name: "halloween-vms"},
			{ID: 200, Name: "gravediggers-biscuits"},
		},
	}).Times(1)

	err := app.Run(strings.Split("bytemark show groups spooky-steve", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
