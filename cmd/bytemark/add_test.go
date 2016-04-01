package main

import (
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

	global.App.Run(strings.Split("bytemark add key --user test-user ssh-rsa aaaaawhartevervAsde fake key", " "))

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}
