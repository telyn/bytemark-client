package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
	publicKeyFile := util.FileFlag{}
	commands = append(commands, cli.Command{
		Name:        "add",
		Usage:       "add SSH keys to a user / IPs to a server",
		UsageText:   "bytemark add key|ip",
		Description: "add SSH keys to a user or IPs to a server",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "key",
			Usage:       "add public SSH keys to a Bytemark user",
			UsageText:   "bytemark add key [--user <user>] [--public-key-file <filename>] <key>",
			Description: `Add the given public key to the given user (or the default user). This will allow them to use that key to access management IPs they have access to using that key. To remove a key, use the remove key command. --public-key-file will be ignored if a public key is specified in the arguments`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "Which user to add the key to. Defaults to the username you log in as.",
				},
				cli.GenericFlag{
					Name:  "public-key-file",
					Usage: "The public key file to add to the account",
					Value: &publicKeyFile,
				},
			},
			Action: With(AuthProvider, func(ctx *Context) (err error) {
				user := ctx.String("user")
				if user == "" {
					user = global.Config.GetIgnoreErr("user")
				}

				key := strings.TrimSpace(strings.Join(ctx.Args(), " "))
				if key == "" {
					if publicKeyFile.Value == "" {
						return ctx.Help("Please specify a key")
					}
					key = publicKeyFile.Value
				} else {
					publicKeyFile = util.FileFlag{FileName: key}
					if err := publicKeyFile.Set(key); err == nil {
						key = publicKeyFile.Value
					}
				}

				if strings.Contains(key, "PRIVATE KEY") {
					return ctx.Help("The key needs to be a public key, not a private key")
				}

				err = global.Client.AddUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key added successfully")
				}
				return
			}),
		}, {
			Name:        "ips",
			Aliases:     []string{"ip"},
			Usage:       "add extra IP addresses to a server",
			UsageText:   "bytemark add ips [--ipv4 | --ipv6] [--ips <number>] <server name>",
			Description: `Add an extra IP to the given server. The IP will be chosen by the brain and output to standard out.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ipv4",
					Usage: "Add an IPv4 address. This is the default",
				},
				cli.BoolFlag{
					Name:  "ipv6",
					Usage: "Add an IPv6 address.",
				},
				cli.IntFlag{
					Name:  "ips",
					Usage: "How many IPs to add (1 to 4) - defaults to one.",
				},
				cli.StringFlag{
					Name:  "reason",
					Usage: "Reason for adding the IP. If not set, will prompt.",
				},
			},
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				addrs := c.Int("ips")
				if addrs < 1 {
					addrs = 1
				}
				family := "ipv4"
				if c.Bool("ipv6") {
					if c.Bool("ipv4") {
						return c.Help("--ipv4 cannot be specified at the same time as --ipv6")
					}
					family = "ipv6"
				}
				reason := c.String("reason")
				if reason == "" {
					if addrs == 1 {
						reason = util.Prompt("Enter the purpose for this extra IP: ")
					} else {
						reason = util.Prompt("Enter the purpose for these extra IPs: ")
					}
				}
				ipcr := brain.IPCreateRequest{
					Addresses:  addrs,
					Family:     family,
					Reason:     reason,
					Contiguous: c.Bool("contiguous"),
				}
				ips, err := global.Client.AddIP(c.VirtualMachineName, &ipcr)
				if err != nil {
					return err
				}
				log.Log("IPs addded:")
				log.Output(ips.String(), "\r\n")
				return nil
			}),
		}},
	})
}
