package main

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// This test ensures that all commands have an Action, Description, Usage and UsageText
// and that all their subcommands do too.
func TestCommandsComplete(t *testing.T) {

	// TODO: Add descriptions to admin commands. it's necessary now
	t.Skip("Need to add descriptions for admin commands.")
	traverseAllCommands(Commands(true), func(c cli.Command) {
		emptyThings := make([]string, 0, 4)
		if c.Name == "" {
			log.Log("There is a command with an empty Name.")
			t.Fail()
		}
		// if a command is only usable via its sub commands, and its usage is built from the
		// subcommands usage, its not necessary to check it.
		// incredibly hacky because this asks for the name of the method, and if it matches, just ignore it
		f := runtime.FuncForPC(reflect.ValueOf(c.Action).Pointer()).Name()
		if f == "github.com/BytemarkHosting/bytemark-client/vendor/github.com/urfave/cli.ShowSubcommandHelp" {
			return
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
			t.Errorf("Command %s has empty %s.\r\n", c.FullName(), strings.Join(emptyThings, ", "))
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
	if s == "" {
		return false
	}

	runes := []rune(s)
	return unicode.IsUpper(runes[0])
}

func hasFullStop(s string) bool {
	return strings.Contains(s, ".")
}

func TestFlagsHaveUsage(t *testing.T) {
	traverseAllCommands(Commands(true), func(c cli.Command) {
		for _, f := range c.Flags {
			if checkFlagUsage(f, isEmpty) {
				t.Errorf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.GetName())
			}
		}
	})
	for _, flag := range app.GlobalFlags() {
		if checkFlagUsage(flag, isEmpty) {
			t.Errorf("Global flag %s doesn't have usage.", flag.GetName())
		}
	}
}

func TestUsageStyleConformance(t *testing.T) {
	traverseAllCommandsWithContext(Commands(true), "", func(name string, c cli.Command) {
		t.Run(name, func(t *testing.T) {
			if firstIsUpper(c.Usage) {
				t.Error("Usage should be lowercase but begins with an uppercase letter")
			}

			if hasFullStop(c.Usage) {
				t.Errorf("Usage should not have full stop")
			}

			if c.UsageText == "" {
				t.Error("UsageText is blank")
			}
			if strings.HasPrefix(c.UsageText, "bytemark ") {
				t.Error("UsageText starts with 'bytemark' - shouldn't anymore")
			}
			if hasPositionalArgBeforeFlag(c.UsageText) {
				t.Errorf("UsageText has a positional argument before a flag. Positional arguments must be after ALL flags.")
			}
			if c.Description == "" {
				t.Errorf("Description is blank")
			}
		})
	})
}

func hasPositionalArgBeforeFlag(usageText string) bool {
	var inSquare, inAngle int
	var seenPositional bool
	var last = '\x00'
	for _, c := range usageText {
		switch c {
		case '[':
			inSquare++
		case ']':
			inSquare--
		case '<':
			inAngle++
		case '>':
			inAngle--
		}
		if inAngle > 0 && inSquare == 0 {
			seenPositional = true
		}
		// check to see if we're starting a flag
		switch last {
		case '[', ' ', '\x00':
			// if we've already seen a positional arg, and we're starting a flag then we have our answer
			if c == '-' && seenPositional {
				return true
			}
		}
		last = c
	}
	return false
}
