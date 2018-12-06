package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestAddKeyCommand(t *testing.T) {

	err := ioutil.WriteFile("testkey.pub", []byte("ssh-rsa aaaaawhartevervAsde fake key"), 0600)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("testkey", []byte("-----BEGIN RSA PRIVATE KEY-----\nfake key\n-----END RSA PRIVATE KEY-----"), 0600)
	if err != nil {
		t.Error(err)
	}

	t.Run("Key in command line", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		_, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		c.When("GetUser", "test-user").Return(brain.User{Username: "test-user"}).Times(1)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}

		err = app.Run(strings.Split("bytemark add key --user test-user ssh-rsa aaaaawhartevervAsde fake key", " "))
		if err != nil {
			t.Error(err)
		}
		c.MockRequest.AssertRequestObjectEqual(brain.User{
			Username:       "test-user",
			AuthorizedKeys: brain.Keys{brain.Key{Key: "ssh-rsa aaaaawhartevervAsde fake key"}},
		})
		if ok, err = c.Verify(); !ok {
			t.Fatal(err)
		}

	})
	t.Run("Key in file", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		_, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		c.When("GetUser", "test-user").Return(brain.User{Username: "test-user"}).Times(1)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}

		err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "testkey.pub"})
		if err != nil {
			t.Error(err)
		}
		c.MockRequest.AssertRequestObjectEqual(brain.User{
			Username:       "test-user",
			AuthorizedKeys: brain.Keys{brain.Key{Key: "ssh-rsa aaaaawhartevervAsde fake key"}},
		})
		if ok, err = c.Verify(); !ok {
			t.Fatal(err)
		}

	})

	t.Run("Key in file using flag", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		_, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		c.When("GetUser", "test-user").Return(brain.User{Username: "test-user"}).Times(1)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}
		err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "--public-key-file", "testkey.pub"})
		if err != nil {
			t.Error(err)
		}
		c.MockRequest.AssertRequestObjectEqual(brain.User{
			Username:       "test-user",
			AuthorizedKeys: brain.Keys{brain.Key{Key: "ssh-rsa aaaaawhartevervAsde fake key"}},
		})
		if ok, vErr := c.Verify(); !ok {
			t.Fatal(vErr)
		}

	})

	t.Run("dont allow private key", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		_, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		err = app.Run([]string{"bytemark", "add", "key", "--user", "test-user", "--public-key-file", "testkey"})
		if err == nil {
			t.Error("Expected an error")
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	})

	_ = os.Remove("testkey.pub")
	_ = os.Remove("testkey")
}

func TestAddIPCommand(t *testing.T) {
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

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
