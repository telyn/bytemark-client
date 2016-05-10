package main

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib"

	"github.com/cheekybits/is"
	"testing"
)

func TestDeleteServer(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account",
	}

	vm := getFixtureVM()

	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", &name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", &name, false).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark delete server test-server", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()

	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", &name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", &name, true).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark delete server --purge test-server", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account",
	}
	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("DeleteDisc", &name, "666").Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark delete disc test-server.test-group.test-account 666", " "))

	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteKey(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	usr := lib.User{
		Username: "test-user",
		Email:    "test-user@example.com",
		AuthorizedKeys: []string{
			"ssh-rsa AAAAFakeKey test-key-one",
			"ssh-rsa AAAAFakeKeyTwo test-key-two",
			"ssh-rsa AAAAFakeKeyThree test-key-two",
		},
	}

	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")
	c.When("AuthWithToken", "test-token").Return(nil)
	c.When("GetUser", usr.Username).Return(&usr)

	c.When("DeleteUserAuthorizedKey", "test-user", "ssh-rsa AAAAFakeKey test-key-one").Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark delete key ssh-rsa AAAAFakeKey test-key-one", " "))

	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")

	c.When("AuthWithToken", "test-token").Return(nil)
	c.When("GetUser", usr.Username).Return(&usr)
	kerr := new(lib.AmbiguousKeyError)
	c.When("DeleteUserAuthorizedKey", "test-user", "test-key-two").Return(kerr).Times(1)

	global.App.Run(strings.Split("bytemark delete key test-key-two", " "))

	is.Equal(kerr, global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}
