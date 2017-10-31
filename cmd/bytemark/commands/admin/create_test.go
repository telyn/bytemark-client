package admin_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestCreateVLANGroup(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetGroup").Return(defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 0).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create vlan-group test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVLANGroupWithVLANNum(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetGroup").Return(defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 19).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create vlan-group test-group.test-account 19", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVLANGroupError(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetGroup").Return(defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 0).Return(fmt.Errorf("Group name already used")).Times(1)

	err := app.Run(strings.Split("bytemark create vlan-group test-group.test-account", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRange(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 14).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create ip range 192.168.3.0/28 14", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRangeError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 18).Return(fmt.Errorf("Error creating IP range")).Times(1)

	err := app.Run(strings.Split("bytemark create ip range 192.168.3.0/28 18", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateUser", "uname", "cluster_su").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUserError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateUser", "uname", "cluster_su").Return(fmt.Errorf("Error creating user")).Times(1)

	err := app.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
