package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"testing"
	//"github.com/cheekybits/is"
)

func TestCommandConfigSet(t *testing.T) {
	config, _ := baseTestSetup()

	config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	global.App.Run([]string{"set", "user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}
