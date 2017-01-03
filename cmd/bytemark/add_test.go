package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
)

func TestAddKeyCommand(t *testing.T) {
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")

	err := ioutil.WriteFile("testkey.pub", []byte("ssh-rsa aaaaawhartevervAsde fake key"), 0600)
	if err != nil {
		t.Error(err)
	}

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = global.App.Run(strings.Split("bytemark add key --user test-user ssh-rsa aaaaawhartevervAsde fake key", " "))
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	c.Reset()
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = global.App.Run([]string{"bytemark", "add", "key", "--user", "test-user", "testkey.pub"})
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	c.Reset()
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = global.App.Run([]string{"bytemark", "add", "key", "--user", "test-user", "--public-key-file", "testkey.pub"})
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	_ = os.Remove("testkey.pub")
}

func TestAddIPCommand(t *testing.T) {
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetVirtualMachine").Return(&defVM)

	vm := lib.VirtualMachineName{VirtualMachine: "test-server"}
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vm, nil).Times(1)

	ipcr := brain.IPCreateRequest{
		Addresses:  1,
		Family:     "ipv4",
		Reason:     "testing",
		Contiguous: false,
	}

	ip := net.ParseIP("10.10.10.10")

	ipcres := ipcr
	ipcres.IPs = []*net.IP{&ip}

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddIP", &vm, &ipcr).Return(&ipcres, nil)

	err := global.App.Run(strings.Split("bytemark add ip --reason testing test-server", " "))
	if err != nil {
		t.Error(err)
	}

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
