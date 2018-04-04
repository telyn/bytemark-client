package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sort"
	"strings"

	bmapp "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	commandsPkg "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// forceFlag is common to a bunch of commands and can have a generic Usage.
var forceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}

//commands is assembled during init()
var commands = make([]cli.Command, 0)

// Commands returns the full list of commands that are available to a user
// (including admin commands if requested), sorts them into alphabetical order,
// and then generates the help section for each command
func Commands(wantAdminCmds bool) []cli.Command {
	myCommands := cliutil.AssembleCommands(commands, commandsPkg.Commands)
	if wantAdminCmds {
		myCommands = cliutil.MergeCommands(myCommands, admin.Commands)
	}
	generateHelp(myCommands)
	sorted := cli.CommandsByName(myCommands)
	sort.Sort(sorted)
	return sorted
}

func main() {
	// watch for interrupts (Ctrl-C) and exit "gracefully" if they are encountered.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			log.Error("\r\nCaught an interrupt - exiting.\r\n")
			// os.Exit is not actually graceful but WHATEVER I don't
			// actually have a better way since bytemark-client has no
			// main-loop or anything - it's just a one-shot
			os.Exit(int(util.ExitCodeTrappedInterrupt))
		}

	}()

	name := os.Args[0]
	overrideHelp(name)
	flags, args, config := prepConfig()

	// add admin commands if --admin is set
	wantAdminCmds, err := config.GetBool("admin")
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	myCommands := Commands(wantAdminCmds)

	app, err := bmapp.BaseAppSetup(flags, myCommands)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	client, err := lib.NewWithURLs(lib.EndpointURLs{
		Brain:   config.GetIgnoreErr("endpoint"),
		API:     config.GetIgnoreErr("api-endpoint"),
		Billing: config.GetIgnoreErr("billing-endpoint"),
		SPP:     config.GetIgnoreErr("spp-endpoint"),
		Auth:    config.GetIgnoreErr("auth-endpoint"),
	})
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	client.SetDebugLevel(config.GetDebugLevel())

	bmapp.SetClientAndConfig(app, client, config)

	outputDebugInfo(config)

	err = app.Run(args)

	os.Exit(int(util.ProcessError(err)))
}

func outputDebugInfo(config config.Manager) {
	log.Debugf(log.LvlOutline, "bytemark-client %s\r\n\r\n", lib.Version)
	// assemble a string of config vars (excluding token)
	vars, err := config.GetAll()
	if err != nil {
		log.Debugf(log.LvlFlags, "(not a real problem maybe): had trouble getting all config vars: %s\r\n", err.Error())
	}

	log.Debugf(log.LvlFlags, "reading config from %s\r\n\r\n", config.ConfigDir())
	log.Debug(log.LvlFlags, "config vars:")
	for _, v := range vars {
		if v.Name == "token" {
			log.Debugf(log.LvlFlags, "  %s (%s): not printed for security\r\n", v.Name, v.Source)
			continue
		}
		log.Debugf(log.LvlFlags, "  %s (%s): '%s'\r\n", v.Name, v.Source, v.Value)
	}
	log.Debug(log.LvlFlags, "")

	log.Debugf(log.LvlFlags, "invocation: %s\r\n\r\n", strings.Join(os.Args, " "))
}

// overrideHelp writes our own help templates into urfave/cli
func overrideHelp(name string) {
	cli.SubcommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   ` + name + ` {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   ` + name + ` {{.UsageText}}{{if .Description}}

   {{.Description}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
}

func prepConfig() (flags []cli.Flag, args []string, conf config.Manager) {
	// set up our global flags because we need some config before we can set up our App
	flagset := flag.NewFlagSet("flags", flag.ContinueOnError)
	help := flagset.Bool("help", false, "")
	h := flagset.Bool("h", false, "")
	version := flagset.Bool("version", false, "")
	v := flagset.Bool("v", false, "")

	flags = bmapp.GlobalFlags()
	for _, f := range flags {
		f.Apply(flagset)
	}

	flagset.SetOutput(ioutil.Discard)

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	configDir := flagset.Lookup("config-dir").Value.String()
	conf, err = config.New(configDir)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	flargs := conf.ImportFlags(flagset)

	//juggle the arguments in order to get the executable on the beginning
	args = make([]string, len(flargs)+1)
	if len(flargs) > 0 && flargs[0] == "help" {
		copy(args[1:], flargs[1:])
		args[len(args)-1] = "--help"
	} else if len(flargs) > 0 && flargs[0] == "version" {
		copy(args[1:], flargs[1:])
		args[len(args)-1] = "--version"
	} else {
		copy(args[1:], flargs)
	}
	args[0] = os.Args[0]

	if *help || *h {
		helpArgs := make([]string, len(args)+1)
		helpArgs[len(args)] = "--help"
		copy(helpArgs, args)
		args = helpArgs
	} else if *version || *v {
		versionArgs := make([]string, len(args)+1)
		versionArgs[len(args)] = "--version"
		copy(versionArgs, args)
		args = versionArgs
	}
	return
}
