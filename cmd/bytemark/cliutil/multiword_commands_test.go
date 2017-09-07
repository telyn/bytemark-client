package cliutil

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"
	"github.com/urfave/cli"
)

func mkTestCommand() cli.Command {
	return cli.Command{
		Name:        "deep test command",
		Description: "this is a test",
		UsageText:   "--test",
		Usage:       "this is a real good test",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "real-good-flag",
			},
		},
		Action: func(c *cli.Context) error {
			return fmt.Errorf("Test error!")
		},
	}
}

func TestCreateMultiwordCommand(t *testing.T) {
	is := is.New(t)
	cmd := mkTestCommand()

	newCmd := CreateMultiwordCommand(cmd)
	is.Equal("deep", newCmd.Name)

	subcmds := newCmd.Subcommands

	if len(subcmds) != 1 {
		t.Fatalf("Wrong number of subcommands %d", len(subcmds))
	}
	is.Equal("test", subcmds[0].Name)
	subcmds = subcmds[0].Subcommands

	if len(subcmds) != 1 {
		t.Fatalf("Wrong number of subsubcommands %d", len(subcmds))
	}
	is.Equal("command", subcmds[0].Name)

	if len(subcmds[0].Flags) != 1 {
		t.Fatalf("Wrong number of flags %d", len(subcmds[0].Flags))
	}
	is.Equal("real-good-flag", subcmds[0].Flags[0].GetName())
	is.Equal("this is a test", subcmds[0].Description)
	is.Equal("--test", subcmds[0].UsageText)
	is.Equal("this is a real good test", subcmds[0].Usage)

}

func TestCreateMultiwordCommands(t *testing.T) {
	is := is.New(t)

	cmd := mkTestCommand()

	before := []cli.Command{{
		Name:        "outer",
		Subcommands: []cli.Command{cmd},
	}}
	expected := cli.Command{
		Name: "outer",
		Subcommands: []cli.Command{
			cmd,
			CreateMultiwordCommand(cmd),
		},
	}

	after := CreateMultiwordCommands(before)
	is.Equal(1, len(after))
	outer := after[0]
	is.Equal(expected.Name, outer.Name)
	is.Equal(len(expected.Subcommands), len(outer.Subcommands))

	foundDeep := false
	foundDeepTest := false
	foundDeepTestCommand := false

	// traverse into after looking for our new multiword commands.
	for _, afterDeep := range outer.Subcommands {
		if afterDeep.Name == "deep" {
			foundDeep = true
			is.Equal(1, len(afterDeep.Subcommands))

			for _, afterDeepTest := range afterDeep.Subcommands {
				is.Equal(1, len(afterDeepTest.Subcommands))
				if afterDeepTest.Name == "test" {
					foundDeepTest = true

					for _, afterDeepTestCommand := range afterDeepTest.Subcommands {
						if afterDeepTestCommand.Name == "command" {
							foundDeepTestCommand = true
						}
					}
				}
			}
		} else {
			is.Equal(cmd.Name, afterDeep.Name)
			is.Equal(cmd.Usage, afterDeep.Usage)
			is.Equal(cmd.UsageText, afterDeep.UsageText)
			is.Equal(cmd.Description, afterDeep.Description)
		}
	}

	is.True(foundDeep)
	is.True(foundDeepTest)
	is.True(foundDeepTestCommand)
}
