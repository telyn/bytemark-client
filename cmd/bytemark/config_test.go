package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestCommandConfigSet(t *testing.T) {
	is := is.New(t)
	config, _ := baseTestSetup(t, false)

	config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
	config.When("GetIgnoreErr", "user").Return("old-test-user")

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	err := global.App.Run(strings.Split("bytemark config set user test-user", " "))
	is.Nil(err)

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}

	err = global.App.Run(strings.Split("bytemark config set flimflam test-user", " "))
	is.NotNil(err)
}
