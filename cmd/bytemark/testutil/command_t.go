package testutil

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

// CommandT is a representation of all the basic common stuff we tend to end
// up with in our single-command tests.
// This helps DRY up tests a bit and should help make writing them faster and
// less error prone
type CommandT struct {
	// Name is the name of the test - passed to t.Run
	Name string
	// Auth is whether the config, client and app should be prepped for
	// authentication
	Auth bool
	// Admin is what we want to stub config.GetBool("admin") to return - i.e.
	// whether the app is in admin-mode.
	Admin bool
	// ShouldErr is whether we're expecting an error to come out of the app.Run
	// call
	ShouldErr bool
	// Args are the command-line arguments to pass to app.Run (after "bytemark")
	Args string
	// Commands are the command set to use in the app. TODO: should I default to
	// main.Commands(Admin)?
	Commands []cli.Command
}

// Run runs the test. Before setup, it calls t.Run to make a subtest and sets up
// config, client and app. After setup, it calls app.Run with the args in the
// CommandT, then checks for errors and verifies the config and client got
// all the calls they were expecting.
func (test CommandT) Run(t *testing.T, setup func(*testing.T, *mocks.Config, *mocks.Client, *cli.App)) {
	t.Run(test.Name, func(t *testing.T) {
		config, client, app := BaseTestAuthSetup(t, test.Admin, test.Commands)
		if !test.Auth {
			config, client, app = BaseTestSetup(t, test.Admin, test.Commands)
		}

		setup(t, config, client, app)

		err := app.Run(strings.Split("bytemark "+test.Args, " "))
		if !test.ShouldErr && err != nil {
			t.Errorf("Unexpected error: %s", err)
		} else if test.ShouldErr && err == nil {
			t.Error("Expected error but did not get one")
		}

		if ok, err := config.Verify(); !ok {
			t.Error(err)
		}
		if ok, err := client.Verify(); !ok {
			t.Error(err)
		}
	})
}
