package main

import (
	"errors"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
	"testing"
)

func TestMergeCommand(t *testing.T) {
	is := is.New(t)
	fail := func(c *cli.Context) error {
		return errors.New("fail")
	}
	succeed := func(c *cli.Context) error {
		return nil
	}
	base := cli.Command{
		Name:        "test",
		Usage:       "base-usage",
		UsageText:   "base-usage-text",
		Description: "base-description",
		Flags:       []cli.Flag{cli.BoolFlag{}},
		Action:      fail,
	}

	tests := []struct {
		cmd                cli.Command
		exUsage            string
		exUsageText        string
		exDescription      string
		exFlags            int
		exActionReturnsNil bool
	}{
		{
			cmd:                cli.Command{},
			exUsage:            "base-usage",
			exUsageText:        "base-usage-text",
			exDescription:      "base-description",
			exFlags:            1,
			exActionReturnsNil: false,
		}, {
			cmd: cli.Command{
				Usage: "new-usage",
				Flags: []cli.Flag{
					cli.BoolFlag{},
				},
				Action: succeed,
			},
			exUsage:            "new-usage",
			exUsageText:        "base-usage-text",
			exDescription:      "base-description",
			exFlags:            2,
			exActionReturnsNil: true,
		},
	}

	for _, test := range tests {
		result := base
		mergeCommand(&result, test.cmd)
		is.Equal(base.Name, result.Name)
		is.Equal(test.exUsage, result.Usage)
		is.Equal(test.exUsageText, result.UsageText)
		is.Equal(test.exDescription, result.Description)
		is.Equal(test.exFlags, len(result.Flags))
		act, ok := result.Action.(func(c *cli.Context) error)
		if !ok {
			t.Fatal("result.Action was not a func(c *cli.Context) error")
		}
		if test.exActionReturnsNil {
			is.Nil(act(nil))
		} else {
			is.NotNil(act(nil))
		}
	}
}

func TestMergeCommands(t *testing.T) {

}
