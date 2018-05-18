package main

import (
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "profiles",
		Usage:     "information on using multiple bytemark client configurations",
		UsageText: "profiles",
		Action:    cli.ShowSubcommandHelp,
		Description: `Having multiple configurations with Bytemark client is useful if you regularly log in as two different users,
or to different instances of the Bytemark API. You can set up and use different configurations with the --config-dir global flag.

For example, to set up your default configuration to log in by default as 'alice', and one configuration where you log in as 'bob'
with a yubikey, run the following commands (windows users, note that in cmd.exe you'll need to use --config-dir="%HOME%\.bob"):

    bytemark config set user alice
    bytemark --config-dir="$HOME/.bob" update config --user bob --yubikey
	
Almost every one of the global flags can be set into the config in this way - see 'bytemark update config --help' for the full list of accepted flags.

At this point you can set up an alias to use your 'bob' configuration. Say you use bash/zsh, add the following to your bashrc/zshrc:
    alias bytemark-bob='bytemark --config-dir="$HOME/.bob"'

Alternatively you could create a shim in your PATH like the following:

  #!/bin/sh
  bytemark --config-dir="$HOME/.bob" "$@"

Either way, now you can run 'bytemark-bob show servers' to list all the servers in bob's default account and 'bytemark show servers' to do the same for alice.`,
	}, cli.Command{
		Name:      "scripting",
		Usage:     "information on scripting with the client",
		UsageText: "scripting",
		Action:    cli.ShowSubcommandHelp,
		Description: `The Bytemark client has been programmed from the beginning to attempt to make it easy for users to script with it.

Some particularly relevant notes: 
    * Set the global flag --output-format=list to receive output on stdout, one item per line. This is particularly useful on plural show commands (e.g. show servers). Fields can be filtered with the --table-fields per-command flag.
    * The reimage and create server commands only print the root password to stdout - all other output is sent to stderr.
    * If you're a fan of jq or you want to otherwise script against bytemark-client, you can get json output from show commands (and a few others) with --json, or by setting the global flag --output-format=json.

Here are just a couple of tricks I've been able to come up with.

To output the uptime for all your machines in the "critical" group:
  for i in $(bytemark show servers critical | tail -n +1); do echo "${i%%.*}:"; ssh $i uptime; done

To add 10GB of space to each archive grade disk in your "storage" server:
  for disc in $(bytemark show discs storage | grep "archive grade"); do bytemark update disc --new-size +10G --server storage --disc $(awk '{print $2}'); done

To list all my servers that have a disc bigger than the default 25GiB:
    bytemark show account --json telyn | jq '[.groups[].virtual_machines[] | select(.discs[] | .size > 25600) | .hostname ]' | uniq
`,
	})
}
