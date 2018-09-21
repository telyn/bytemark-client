package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestDeleteServer(t *testing.T) {
	is := is.New(t)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	vm := getFixtureVM()

	t.Run("force delete", func(t *testing.T) {
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

		config.When("Force").Return(true)
		config.When("GetVirtualMachine").Return(defVM)

		c.When("GetVirtualMachine", name).Return(vm).Times(1)
		c.When("DeleteVirtualMachine", name, false).Return(nil).Times(1)

		err := app.Run(strings.Split("bytemark delete server --force test-server", " "))
		is.Nil(err)
		if ok, vErr := c.Verify(); !ok {
			t.Fatal(vErr)
		}
	})
	t.Run("force purge", func(t *testing.T) {
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		config.When("Force").Return(true)
		config.When("GetVirtualMachine").Return(defVM)

		c.When("GetVirtualMachine", name).Return(vm).Times(1)
		c.When("DeleteVirtualMachine", name, true).Return(nil).Times(1)

		err := app.Run(strings.Split("bytemark delete server --force --purge test-server", " "))
		is.Nil(err)
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	})
}

func TestDeleteDisc(t *testing.T) {
	t.Run("server and label", func(t *testing.T) {
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
	})
	t.Run("disc ID", func(t *testing.T) {
		is := is.New(t)
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

		config.When("Force").Return(true)

		c.When("BuildRequest", "DELETE", lib.BrainEndpoint, "/discs/%s?purge=true", []string{"666"}).Return(&mocks.Request{
			T:              t,
			StatusCode:     200,
			ResponseObject: nil,
		}).Times(1)

		err := app.Run(strings.Split("bytemark delete disc --force --id 666", " "))
		is.Nil(err)
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	})
}

func TestDeleteKey(t *testing.T) {
	usr := brain.User{
		Username: "test-user",
		Email:    "test-user@example.com",
		AuthorizedKeys: brain.Keys{
			brain.Key{Key: "ssh-rsa AAAAFakeKey test-key-one"},
			brain.Key{Key: "ssh-rsa AAAAFakeKeyTwo test-key-two"},
			brain.Key{Key: "ssh-rsa AAAAFakeKeyThree test-key-two"},
		},
	}
	t.Run("full key", func(t *testing.T) {
		is := is.New(t)
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

		config.When("Force").Return(true)
		c.When("GetUser", usr.Username).Return(usr)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}

		err := app.Run(strings.Split("bytemark delete key ssh-rsa AAAAFakeKey test-key-one", " "))

		is.Nil(err)
		if ok, vErr := c.Verify(); !ok {
			t.Fatal(vErr)
		}
	})

	t.Run("delete by ambiguous comment is err", func(t *testing.T) {
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

		config.When("Force").Return(true)
		config.When("GetIgnoreErr", "user").Return("test-user")

		c.When("AuthWithToken", "test-token").Return(nil)
		c.When("GetUser", usr.Username).Return(usr)

		err := app.Run(strings.Split("bytemark delete key test-key-two", " "))

		if err == nil {
			t.Error("expecting an error but didn't get one")
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
		c.Reset()
	})
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
