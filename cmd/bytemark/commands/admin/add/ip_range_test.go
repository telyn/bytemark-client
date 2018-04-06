package add_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/cheekybits/is"
)

func TestCreateIPRange(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 14).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark add ip range 192.168.3.0/28 14", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRangeError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 18).Return(fmt.Errorf("Error creating IP range")).Times(1)

	err := app.Run(strings.Split("bytemark add ip range 192.168.3.0/28 18", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
