package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:  "hwprofiles",
		Usage: "List hardware profiles available for cloud servers",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for _, profile := range c.Definitions.HardwareProfiles {
				log.Log(profile)
			}
			return
		}),
	}, cli.Command{
		Name:    "images",
		Aliases: []string{"distributions", "distros"},
		Usage:   "List images available for installation on all servers",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for distro, description := range c.Definitions.DistributionDescriptions {
				log.Logf("'%s': %s\r\n", distro, description)
			}
			return
		}),
	}, cli.Command{
		Name:  "storage",
		Usage: "List available storage grades for cloud servers",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for grade, description := range c.Definitions.StorageGradeDescriptions {
				log.Logf("'%s': %s\r\n", grade, description)
			}
			return
		}),
	}, cli.Command{
		Name:  "zones",
		Usage: "List available zones for cloud servers",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for _, zone := range c.Definitions.ZoneNames {
				log.Log(zone)
			}
			return
		}),
	})
}
