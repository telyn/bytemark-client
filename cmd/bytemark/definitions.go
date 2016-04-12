package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "hwprofiles",
		Usage:     "List hardware profiles available for cloud servers",
		UsageText: "bytemark hwprofiles [--json]",
		Description: `Hardware profiles are used by cloud servers and choosing between them can be thought of as 'which virtual motherboard should I use?'.
Generally bytemark provide two - virtio and compatibility. The virtio one has better performance but may not work with obscure operating systems, or without drivers (particularly important if you are installing windows from CD rather than our images`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output the list as a JSON array",
			},
		},
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.IfNotMarshalJSON(c.Definitions.HardwareProfiles, func() error {
				for _, profile := range c.Definitions.HardwareProfiles {
					log.Log(profile)
				}
				return nil
			})
		}),
	}, cli.Command{
		Name:        "images",
		Aliases:     []string{"distributions", "distros"},
		Usage:       "List images available for installation on all servers",
		UsageText:   "bytemark images",
		Description: "This command lists all the images that are available for installation on Bytemark servers.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output the list as a JSON array",
			},
		},
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.IfNotMarshalJSON(c.Definitions.DistributionDescriptions, func() error {
				for distro, description := range c.Definitions.DistributionDescriptions {
					log.Logf("'%s': %s\r\n", distro, description)
				}
				return nil
			})
		}),
	}, cli.Command{
		Name:        "storage",
		Usage:       "List available storage grades for cloud servers",
		UsageText:   "bytemark storage",
		Description: "This outputs the available storage grades for cloud servers.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output the list as a JSON array",
			},
		},
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.IfNotMarshalJSON(c.Definitions.StorageGradeDescriptions, func() error {
				for grade, description := range c.Definitions.StorageGradeDescriptions {
					log.Logf("'%s': %s\r\n", grade, description)
				}
				return nil
			})
		}),
	}, cli.Command{
		Name:        "zones",
		Usage:       "List available zones for cloud servers",
		UsageText:   "bytemark zones",
		Description: "This outputs the zones available for cloud servers to be stored and started in. Note that it is not currently possible to migrate a server between zones.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output the list as a JSON array",
			},
		},
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.IfNotMarshalJSON(c.Definitions.ZoneNames, func() error {
				for _, zone := range c.Definitions.ZoneNames {
					log.Log(zone)
				}
				return nil
			})
		}),
	})
}
