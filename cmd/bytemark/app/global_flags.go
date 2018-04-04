package app

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/urfave/cli"
)

// GlobalFlags returns a new set of global flags for the client.
// This is where they are defined.
func GlobalFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "account",
			Usage: "account name to use when no other accounts are specified",
		},
		cli.StringFlag{
			Name:  "api-endpoint",
			Usage: "URL where the domains service can be found. Set to blank in environments without a domains service.",
		},
		cli.StringFlag{
			Name:  "auth-endpoint",
			Usage: "URL where the auth service can be found",
		},
		cli.StringFlag{
			Name:  "billing-endpoint",
			Usage: "URL where bmbilling can be found. Set to blank in environments without bmbilling",
		},
		cli.BoolFlag{
			Name:  "admin",
			Usage: "allows admin commands in the client. see bytemark --admin --help",
		},
		cli.BoolFlag{
			Name:  "yubikey",
			Usage: "use a yubikey to authenticate",
		},
		cli.StringFlag{
			Name:  "impersonate",
			Usage: "a user to request impersonation of",
		},
		cli.IntFlag{
			Name:  "debug-level",
			Usage: "how much debug output to print to the terminal",
		},
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "URL of the brain",
		},
		cli.StringFlag{
			Name:  "config-dir",
			Usage: "directory in which bytemark-client's configuration resides. see bytemark help config, bytemark help profiles",
		},
		cli.StringFlag{
			Name:  "spp-endpoint",
			Usage: "URL of SPP. set to blank in environments without an SPP.",
		},
		cli.StringFlag{
			Name:  "output-format",
			Usage: "The output format to use. Currently defined output formats are human (default for most commands), json (machine readable format), table (human-readable table format)",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "user you wish to log in as",
		},
		cli.StringFlag{
			Name:  "yubikey-otp",
			Usage: "one-time password from your yubikey to use to login",
		},
		cli.IntFlag{
			Name:  "session-validity",
			Usage: "seconds until your session is automatically invalidated (max 3600)",
			Value: config.DefaultSessionValidity,
			// TODO(telyn): add more defaults to these flags
		},
	}
}
