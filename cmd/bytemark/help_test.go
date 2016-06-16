package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
	"testing"
	"unicode"
)

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
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.BoolTFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.DurationFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.Float64Flag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.GenericFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.StringFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			case cli.StringSliceFlag:
				if f.Usage == "" {
					log.Logf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.Name)
					t.Fail()
				}
			}
		}
	})
}

func TestUsageStyleConformance(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		usage := []rune(c.Usage)
		if unicode.IsUpper(usage[0]) {
			log.Logf("Command %s's Usage begins with an uppercase letter. Please change it - Usages should be lowercase.\r\n", c.FullName())
			t.Fail()
		}

		if strings.Contains(c.Usage, ".") {
			log.Logf("Command %s's Usage has a full-stop. Get rid of it.\r\n", c.FullName())
			t.Fail()
		}
	})
}

// Tests for commands which have subcommands having the correct Description format
// the first line should start lowercase and end without a full stop, and the second
// should be blank
func TestSubcommandStyleConformance(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		if c.Subcommands == nil {
			return
		}
		if len(c.Subcommands) == 0 {
			return
		}
		lines := strings.Split(c.Description, "\n")
		desc := []rune(lines[0])
		if unicode.IsUpper(desc[0]) {
			log.Logf("Command %s's Description begins with an uppercase letter, but it has subcommands, so should be lowercase.\r\n", c.FullName())
			t.Fail()
		}
		if strings.Contains(lines[0], ".") {
			log.Logf("The first line of Command %s's Description contains a full stop. It shouldn't.\r\n", c.FullName())
			t.Fail()
		}
		if len(lines) > 1 {
			if len(strings.TrimSpace(lines[1])) > 0 {
				log.Logf("The second line of Command %s's Description should be blank.\r\n", c.FullName())
				t.Fail()
			}
		}

	})
}
