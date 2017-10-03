package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestDeleteServer(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("Force").Return(true)
	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	vm := getFixtureVM()

	c.When("GetVirtualMachine", name).Return(vm).Times(1)
	c.When("DeleteVirtualMachine", name, false).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark delete server --force test-server", " "))
	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
	c.Reset()

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(vm).Times(1)
	c.When("DeleteVirtualMachine", name, true).Return(nil).Times(1)

	err = app.Run(strings.Split("bytemark delete server --force --purge test-server", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("Force").Return(true)
	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account",
	}
	c.When("DeleteDisc", name, "666").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark delete disc --force test-server.test-group.test-account 666", " "))

	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteKey(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	usr := brain.User{
		Username: "test-user",
		Email:    "test-user@example.com",
		AuthorizedKeys: []string{
			"ssh-rsa AAAAFakeKey test-key-one",
			"ssh-rsa AAAAFakeKeyTwo test-key-two",
			"ssh-rsa AAAAFakeKeyThree test-key-two",
		},
	}

	config.When("Force").Return(true)
	c.When("GetUser", usr.Username).Return(usr)

	c.When("DeleteUserAuthorizedKey", "test-user", "ssh-rsa AAAAFakeKey test-key-one").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark delete key ssh-rsa AAAAFakeKey test-key-one", " "))

	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}

	config, c, app = testutil.BaseTestAuthSetup(t, false, commands)

	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "user").Return("test-user")

	c.When("AuthWithToken", "test-token").Return(nil)
	c.When("GetUser", usr.Username).Return(usr)
	kerr := new(lib.AmbiguousKeyError)
	c.When("DeleteUserAuthorizedKey", "test-user", "test-key-two").Return(kerr).Times(1)

	err = app.Run(strings.Split("bytemark delete key test-key-two", " "))

	is.Equal(kerr, err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}

func TestDeleteBackup(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("DeleteBackup", vmname, "test-disc", "test-backup").Return(nil).Times(1)

	err := app.Run([]string{
		"bytemark", "delete", "backup", "test-server", "test-disc", "test-backup",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteVLAN(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

	c.When("ReapVMs").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "reap", "servers"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteVLANError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

	c.When("ReapVMs").Return(fmt.Errorf("Could not delete VLAN")).Times(1)

	err := app.Run([]string{"bytemark", "reap", "servers"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
