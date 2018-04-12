package show_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	cf "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
)

func TestShowConfig(t *testing.T) {
	config, c, app := testutil.BaseTestSetup(t, false, commands.Commands)
	config.When("GetAll").Return(cf.Vars{
		{
			Name:   "endpoint",
			Value:  "https://uk0.bigv.io/",
			Source: "CODE",
		},
	})

	args := strings.Split("bytemark show config", " ")
	err := app.Run(args)
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
