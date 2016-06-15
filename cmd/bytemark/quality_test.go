package main

import (
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"testing"
)

var destructiveCommands = [...]string{
	"create server", // can increase cost
	"create discs",  // can increase cost
	"delete server", // can destroy data
	"delete group",  // can destroy data
	"delete disc",   // can destroy data
	"reimage",       // can destroy data
	"resize disc",   // can increase cost
	"set memory",    // can increase cost
	"set cores",     // can increase cost

}

type s struct {
	Seen     bool
	HasForce bool
}

// TODO(telyn): It would be lovely to write a test to ensure all destructiveCommands call PromptYesNo, but I don't see it happening any time soon. Callee data is only available from an AST, not via reflection. What we could do is try to mock all the functions that are called by all destructive commands so we can just loop over them.

// TestMain is used to work around weirdness with urfave/cli. See the large comment inside for details.
func TestMain(m *testing.M) {
	flag.Parse()
	for _, c := range commands {
		config, _ := baseTestSetup()
		config.When("Get", "token").Return("no-not-a-token")

		// the issue is that Command.FullName() is dependent on Command.commandNamePath.
		// Command.commandNamePath is filled in when the parent's Command.startApp is called
		// and startApp is only called when you actually try to run that command or one of
		// its subcommands. So we run "bytemark <command> help" on all commands that have
		// subcommands in order to get every subcommand to have a correct Command.commandPath

		if c.Subcommands != nil && len(c.Subcommands) > 0 {
			fmt.Fprintf(os.Stderr, c.Name)
			global.App.Run([]string{"bytemark.test", c.Name, "help"})
		}
	}
	os.Exit(m.Run())
}

// Ensure all destructive commands have a --force flag, to skip through prompting.
func TestDestructiveCommandsHaveForceFlags(t *testing.T) {
	// it would be nice to also check that they have prompting, but that can't be done via reflection, only by building an ast from the source or by running tests.
	cmds := make(map[string]*s)
	for _, cmd := range destructiveCommands {
		cmds[cmd] = &s{}
	}
	traverseAllCommands(commands, func(c cli.Command) {
		for _, cmd := range destructiveCommands {
			if c.FullName() == cmd {
				cmds[cmd].Seen = true
				for _, flag := range c.Flags {
					if flag.GetName() == "force" {
						cmds[cmd].HasForce = true
					}
				}

			}
		}
	})
	for cmd, points := range cmds {
		if !points.Seen {
			t.Errorf("%s not seen in commands", cmd)
		} else if !points.HasForce {
			t.Errorf("%s doesn't have force flag", cmd)
		}
	}

}
