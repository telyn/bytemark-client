package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"net"
	"strings"
	"testing"
)

func TestAddKeyCommand(t *testing.T) {
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)

	err := global.App.Run(strings.Split("bytemark add key --user test-user ssh-rsa aaaaawhartevervAsde fake key", " "))
	if err != nil {
		t.Error(err)
	}

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}
func TestAddIPCommand(t *testing.T) {
	config, c := baseTestSetup()

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetVirtualMachine").Return(&defVM)

	vm := lib.VirtualMachineName{VirtualMachine: "test-server"}
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vm, nil).Times(1)

	ipcr := lib.IPCreateRequest{
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
