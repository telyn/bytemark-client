package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "hwprofiles",
		Usage:     "List hardware profiles available for cloud servers",
		UsageText: "bytemark hwprofiles", // TODO(telyn): it'd be cool if it had a JSON flag
		Description: `Hardware profiles are used by cloud servers and choosing between them can be thought of as 'which virtual motherboard should I use?'.
Generally bytemark provide two - virtio and compatibility. The virtio one has better performance but may not work with obscure operating systems, or without drivers (particularly important if you are installing windows from CD rather than our images`,
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for _, profile := range c.Definitions.HardwareProfiles {
				log.Log(profile)
			}
			return
		}),
	}, cli.Command{
		Name:        "images",
		Aliases:     []string{"distributions", "distros"},
		Usage:       "List images available for installation on all servers",
		UsageText:   "bytemark images", // TODO(telyn): it'd be cool if it had a JSON flag
		Description: "This command lists all the images that are available for installation on Bytemark servers.",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for distro, description := range c.Definitions.DistributionDescriptions {
				log.Logf("'%s': %s\r\n", distro, description)
			}
			return
		}),
	}, cli.Command{
		Name:        "storage",
		Usage:       "List available storage grades for cloud servers",
		UsageText:   "bytemark storage", // TODO(telyn): it'd be cool if it had a json flag.
		Description: "This outputs the available storage grades for cloud servers.",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for grade, description := range c.Definitions.StorageGradeDescriptions {
				log.Logf("'%s': %s\r\n", grade, description)
			}
			return
		}),
	}, cli.Command{
		Name:        "zones",
		Usage:       "List available zones for cloud servers",
		UsageText:   "bytemark zones",
		Description: "This outputs the zones available for cloud servers to be stored and started in. Note that it is not currently possible to migrate a server between zones.",
		Action: With(DefinitionsProvider, func(c *Context) (err error) {
			for _, zone := range c.Definitions.ZoneNames {
				log.Log(zone)
			}
			return
		}),
	})
}
