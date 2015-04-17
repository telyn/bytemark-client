package cmd

import (
	"testing"
	//"github.com/cheekybits/is"
)

func TestCommandSet(t *testing.T) {
	config := &mockConfig{}
	config.When("SetPersistent", "user", "test-user").Times(1)
	config.When("Get", "user").Return("old-test-user")
	cmds := NewCommandSet(config, nil)
	cmds.Set([]string{"user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCommandUnset(t *testing.T) {
	config := &mockConfig{}
	config.When("Unset", "user").Times(1)
	config.When("Get", "user").Return("old-test-user")
	cmds := NewCommandSet(config, nil)
	cmds.Unset([]string{"user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}

// everything else is going to involve making a mock client
// TODO(telyn): make a mock client
