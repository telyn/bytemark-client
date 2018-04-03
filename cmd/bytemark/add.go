package main

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	publicKeyFile := util.FileFlag{}
	commands = append(commands, cli.Command{
		Name:        "add",
		Usage:       "add SSH keys to a user / IPs to a server",
		UsageText:   "add key|ip",
		Description: "add SSH keys to a user or IPs to a server",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "key",
			Usage:       "add public SSH keys to a Bytemark user",
			UsageText:   "add key [--user <user>] [--public-key-file <filename>] <key>",
			Description: `Add the given public key to the given user (or the default user). This will allow them to use that key to access management IPs they have access to using that key. To remove a key, use the remove key command. --public-key-file will be ignored if a public key is specified in the arguments`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "Which user to add the key to. Defaults to the username you log in as.",
				},
				cli.StringFlag{
					Name:  "public-key",
					Usage: "the text of a public key to add. If set, is used in preference to --public-key-file",
				},
				cli.GenericFlag{
					Name:  "public-key-file",
					Usage: "The public key file to add to the account",
					Value: &publicKeyFile,
				},
			},
			Action: app.Action(args.Join("public-key"), with.Auth, func(ctx *app.Context) (err error) {
				user := ctx.String("user")
				if user == "" {
					user = ctx.Config().GetIgnoreErr("user")
				}

				key := strings.TrimSpace(ctx.String("public-key"))
				if key == "" {
					if publicKeyFile.Value == "" {
						return ctx.Help("Please specify a key")
					}
					key = publicKeyFile.Value
				} else {
					// if public-key is not blank, try to use it as a filename
					// FileFlag does some nice ~-substitution which is why we use it rather than the infinitely more normal-looking ioutil.ReadFile
					publicKeyFile = util.FileFlag{FileName: key}
					if err := publicKeyFile.Set(key); err == nil {
						key = publicKeyFile.Value
					}
				}

				if strings.Contains(key, "PRIVATE KEY") {
					return ctx.Help("The key needs to be a public key, not a private key")
				}

				err = ctx.Client().AddUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key added successfully")
				}
				return
			}),
		}, {
			Name:        "ips",
			Aliases:     []string{"ip"},
			Usage:       "add extra IP addresses to a server",
			UsageText:   "add ips [--ipv4 | --ipv6] [--ips <number>] <server name>",
			Description: `Add an extra IP to the given server. The IP will be chosen by the brain and output to standard out.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ipv4",
					Usage: "If set, requests IPv4 addresses. This is the default",
				},
				cli.BoolFlag{
					Name:  "ipv6",
					Usage: "If set, requests IPv6 addresses.",
				},
				cli.IntFlag{
					Name:  "ips",
					Usage: "How many IPs to add (1 to 4) - defaults to one.",
				},
				cli.StringFlag{
					Name:  "reason",
					Usage: "Reason for adding the IP. If not set, will prompt.",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "The server to add IPs to",
					Value: new(app.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
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
						reason = c.Prompter().Prompt("Enter the purpose for this extra IP: ")
					} else {
						reason = c.Prompter().Prompt("Enter the purpose for these extra IPs: ")
					}
				}
				ipcr := brain.IPCreateRequest{
					Addresses:  addrs,
					Family:     family,
					Reason:     reason,
					Contiguous: c.Bool("contiguous"),
				}
				vmName := c.VirtualMachineName("server")
				ips, err := c.Client().AddIP(vmName, ipcr)
				if err != nil {
					return err
				}
				log.Log("IPs added:")
				log.Output(ips.String(), "\r\n")
				return nil
			}),
		}},
	})
}
