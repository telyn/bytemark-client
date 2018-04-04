package add_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestCreateGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

	config.When("GetGroup").Return(defGroup)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "default-account",
	}
	c.When("CreateGroup", group).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark add group test-group", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
