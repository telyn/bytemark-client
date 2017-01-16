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

type stringPredicate func(string) bool

func checkFlagUsage(f cli.Flag, predicate stringPredicate) bool {
	switch f := f.(type) {
	case cli.BoolFlag:
		return predicate(f.Usage)
	case cli.BoolTFlag:
		return predicate(f.Usage)
	case cli.DurationFlag:
		return predicate(f.Usage)
	case cli.Float64Flag:
		return predicate(f.Usage)
	case cli.GenericFlag:
		return predicate(f.Usage)
	case cli.StringFlag:
		return predicate(f.Usage)
	case cli.StringSliceFlag:
		return predicate(f.Usage)
	}
	return false
}

func isEmpty(s string) bool {
	return s == ""
}
func notEmpty(s string) bool {
	return s != ""
}

func firstIsUpper(s string) bool {
	runes := []rune(s)
	return unicode.IsUpper(runes[0])
}

func hasFullStop(s string) bool {
	return strings.Contains(s, ".")
}

func TestFlagsHaveUsage(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		for _, f := range c.Flags {
			if checkFlagUsage(f, isEmpty) {
				t.Errorf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.GetName())
			}
		}
	})
	for _, flag := range globalFlags() {
		if checkFlagUsage(flag, isEmpty) {
			t.Errorf("Global flag %s doesn't have usage.", flag.GetName())
		}
	}
}

func TestUsageStyleConformance(t *testing.T) {
	traverseAllCommands(commands, func(c cli.Command) {
		if firstIsUpper(c.Usage) {
			t.Errorf("Command %s's Usage begins with an uppercase letter. Please change it - Usages should be lowercase.\r\n", c.FullName())
		}

		if hasFullStop(c.Usage) {
			t.Errorf("Command %s's Usage has a full-stop. Get rid of it.\r\n", c.FullName())
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
