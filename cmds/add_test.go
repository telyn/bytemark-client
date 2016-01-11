package cmds

import (
	"bytemark.co.uk/client/mocks"
	"testing"
)

func TestAddKeyCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	args := []string{"test-user", "ssh-rsa", "aaaaawhartevervAsde", "fake key"}
	config.When("ImportFlags").Return(args).Times(1)
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("AddUserAuthorizedKey", "test-user", "ssh-rsa aaaaawhartevervAsde fake key").Times(1)

	cmds := NewCommandSet(config, c)
	cmds.AddKey(args)

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
	c.Reset()
}
