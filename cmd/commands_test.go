package main

import (
	"testing"
	//"github.com/cheekybits/is"
)

func TestCommandConfig(t *testing.T) {
	config := &mockConfig{}

	config.When("GetV", "user").Return(ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	cmds := NewCommandSet(config, nil)
	cmds.Config([]string{"set", "user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}

// everything else is going to involve making a mock client
// TODO(telyn): make a mock client
