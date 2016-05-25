package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
	"testing"
)

func traverseAllCommands(cmds []cli.Command, fn func(cli.Command)) {
	if cmds == nil {
		return
	}
	for _, c := range cmds {
		fn(c)
		traverseAllCommands(c.Subcommands, fn)
	}
}

// This test ensures that all commands have an Action, Description, Usage and UsageText
// and that all their subcommands do too.
func TestCommandsComplete(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		emptyThings := make([]string, 0, 4)
		if c.Name == "" {
			log.Log("There is a command with an empty Name.")
			t.Fail()
		}
		if c.Usage == "" {
			emptyThings = append(emptyThings, "Usage")
		}
		if c.UsageText == "" {
			emptyThings = append(emptyThings, "UsageText")
		}
		if c.Description == "" {
			emptyThings = append(emptyThings, "Description")
		}
		if c.Action == nil {
			emptyThings = append(emptyThings, "Action")
		}
		if len(emptyThings) > 0 {
			t.Fail()
			log.Logf("Command %s has empty %s.\r\n", c.FullName(), strings.Join(emptyThings, ", "))
		}
	})

}

func TestFlagsHaveUsage(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		for _, f := range c.Flags {
			switch f := f.(type) {
			case cli.BoolFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.BoolTFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.DurationFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.Float64Flag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.GenericFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.StringFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.StringSliceFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage", c.FullName(), f.Name)
					t.Fail()
				}
			}
		}
	})
}
