package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestDeleteServer(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("Force").Return(true)
	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	vm := getFixtureVM()

	c.When("GetVirtualMachine", name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", name, false).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark delete server --force test-server", " "))
	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
	c.Reset()

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&vm).Times(1)
	c.When("DeleteVirtualMachine", name, true).Return(nil).Times(1)

	err = global.App.Run(strings.Split("bytemark delete server --force --purge test-server", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("Force").Return(true)
	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account",
	}
	c.When("DeleteDisc", name, "666").Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark delete disc --force test-server.test-group.test-account 666", " "))

	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteKey(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

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
	c.When("GetUser", usr.Username).Return(&usr)

	c.When("DeleteUserAuthorizedKey", "test-user", "ssh-rsa AAAAFakeKey test-key-one").Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark delete key ssh-rsa AAAAFakeKey test-key-one", " "))

	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
	c.Reset()
	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "2fa-otp").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")

	c.When("AuthWithToken", "test-token").Return(nil)
	c.When("GetUser", usr.Username).Return(&usr)
	kerr := new(lib.AmbiguousKeyError)
	c.When("DeleteUserAuthorizedKey", "test-user", "test-key-two").Return(kerr).Times(1)

	err = global.App.Run(strings.Split("bytemark delete key test-key-two", " "))

	is.Equal(kerr, err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}

func TestDeleteVLAN(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(nil).Times(1)

	err := global.App.Run([]string{"bytemark", "reap", "servers"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestDeleteVLANError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(fmt.Errorf("Could not delete VLAN")).Times(1)

	err := global.App.Run([]string{"bytemark", "reap", "servers"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
