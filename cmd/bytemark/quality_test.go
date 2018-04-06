// do not add the quality build tag to this file - it requires access to the
// unexported commands variable, so has to be in package main rather than
// main_test. For some reason this seems to mess up build tags.
// This test doesn't add any prerequisites for running 'go test' anyway, unlike
// the other ones, and doesn't take long to run. So it's fine to keep it in the
// main test job.
package main

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/urfave/cli"
)

var destructiveCommands = [...]string{
	// "add server",    // can increase cost
	// "add discs",     // can increase cost
	"delete server", // can destroy data
	"delete group",  // can destroy data
	"delete disc",   // can destroy data
	"reimage",       // can destroy data
	"resize disc",   // can increase cost
	// FullName() on commands from admin/ and commands/ fails
	//"update server", // can increase cost
}

type s struct {
	Seen     bool
	HasForce bool
}

// TODO(telyn): It would be lovely to write a test to ensure all destructiveCommands call PromptYesNo, but I don't see it happening any time soon. Callee data is only available from an AST, not via reflection. What we could do is try to mock all the functions that are called by all destructive commands so we can just loop over them.
// Actually, just making PromptYesNo a variable in util will do. Set it to the actual impl of PromptYesNo by default, mock it in the test.

// Ensure all destructive commands have a --force flag, to skip through prompting.
func TestDestructiveCommandsHaveForceFlags(t *testing.T) {
	// it would be nice to also check that they have prompting, but that can't be done via reflection, only by building an ast from the source or by running tests.
	cmds := make(map[string]*s)
	for _, cmd := range destructiveCommands {
		cmds[cmd] = &s{}
	}
	testutil.TraverseAllCommands(Commands(true), func(c cli.Command) {
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
