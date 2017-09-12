package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestAddKeyCommand(t *testing.T) {
	_, c, app := baseTestAuthSetup (t, false)

	err := ioutil.WriteFile("testkey.pub", []byte("ssh-rsa aaaaawhartevervAsde fake key"), 0600)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("testkey", []byte("-----BEGIN RSA PRIVATE KEY-----\nfake key\n-----END RSA PRIVATE KEY-----"), 0600)
	if err != nil {
		t.Error(err)
	}

	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = app.Run(strings.Split("bytemark add key --user test-user ssh-rsa aaaaawhartevervAsde fake key", " "))
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	c.Reset()
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "testkey.pub"})
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	c.Reset()
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)
	err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "--public-key-file", "testkey.pub"})
	if err != nil {
		t.Error(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	c.Reset()
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "--public-key-file", "testkey"})
	if err == nil {
		t.Error("Expected an error")
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	_ = os.Remove("testkey.pub")
	_ = os.Remove("testkey")
}

func TestAddIPCommand(t *testing.T) {
	config, c, app := baseTestAuthSetup (t, false)

	config.When("GetVirtualMachine").Return(defVM)

	vm := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "default", Account: "default-account"}

	ipcr := brain.IPCreateRequest{
		Addresses:  1,
		Family:     "ipv4",
		Reason:     "testing",
		Contiguous: false,
	}

	ip := net.ParseIP("10.10.10.10")

	ipcres := ipcr
	ipcres.IPs = []net.IP{ip}

	c.When("AddIP", vm, ipcr).Return(&ipcres, nil)

	err := app.Run(strings.Split("bytemark add ip --reason testing test-server", " "))
	if err != nil {
		t.Error(err)
	}

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
