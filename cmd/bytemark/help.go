package main

import (
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "profiles",
		Usage:     "Information on using multiple configurations",
		UsageText: "bytemark profiles",
		Action:    cli.ShowSubcommandHelp,
		Description: `Having multiple configurations with Bytemark client is useful if you regularly log in as two different users,
or to different instances of the Bytemark API. One can set up and use different configurations with the --config-dir global flag.

For example, to set up your default configuration to log in by default as 'alice', and one configuration where you log in as 'bob'
with a yubikey, run the following commands (windows users, note that you'll need to use --config-dir="%HOME%\.bob"):

    bytemark config set user alice
    bytemark --config-dir="$HOME/.bob" set user bob
    bytemark --config-dir="$HOME/.bob" set yubikey

At this point you can set up an alias to use your 'bob' configuration. Say you use bash/zsh, add the following to your bashrc/zshrc:
    alias bytemark-bob='bytemark --config-dir="$HOME/.bob"'

Now you can run 'bytemark-bob list servers' to list all the servers in bob's default account and 'bytemark list servers' to do the same for alice.

Sorted.`,
	}, cli.Command{
		Name:      "scripting",
		Usage:     "Information on scripting with the client",
		UsageText: "bytemark scripting",
		Action:    cli.ShowSubcommandHelp,
		Description: `The Bytemark client has been programmed from the beginning to attempt to make it easy for users to script with it.

Some particularly relavent notes: 
    * The 'list' command is entirely designed for scripting - it outputs to stdout, one item per line.
    * The reimage and create server commands only print the root password to stdout - all other output is sent to stderr.
    * All exit codes are documented - see the help topic exit codes.
    * If you're a fan of jq or you want to use bytemark-client within a OO scripting language, you can get json output from several commands like show using --json.

Here are just a couple of tricks I've been able to come up with.

To output the uptime for all your machines in the "critical" group:
  for i in $(bytemark list servers critical); do echo "${i%%.*}:"; ssh $i uptime; done

To add 10GB of space to each archive grade disk in your "storage" server:
  for disc in $(bytemark list discs storage | grep "archive grade"); do bytemark resize disc --size +10G $machine $(awk '{print $2}'); done

To list all my servers that have a disc bigger than the default 25GiB:
    bytemark show account --json telyn | jq '[.groups[].virtual_machines[] | select(.discs[] | .size > 25600) | .hostname ]' | uniq
`,
	})
}
