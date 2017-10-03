package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	bmapp "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/bgentry/speakeasy"
	"github.com/urfave/cli"
)

// forceFlag is common to a bunch of commands and can have a generic Usage.
var forceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}

//commands is assembled during init()
var commands = make([]cli.Command, 0)

//adminCommands is assembled during init() and has the commands that're only available when --admin is specified.
// it gets merged in to commands
var adminCommands = make([]cli.Command, 0)

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

	overrideHelp()
	flags, args, config := prepConfig()

	// add admin commands if --admin is set
	wantAdminCmds, err := config.GetBool("admin")
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	myCommands := commands
	if wantAdminCmds {
		myCommands = cliutil.MergeCommands(commands, adminCommands)
	}

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

	app.Metadata = map[string]interface{}{
		"client": client,
		"config": config,
	}

	outputDebugInfo(config)

	err = app.Run(args)

	os.Exit(int(util.ProcessError(err)))
}

func outputDebugInfo(config util.ConfigManager) {
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
func overrideHelp() {
	cli.SubcommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

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
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
}

func prepConfig() (flags []cli.Flag, args []string, config util.ConfigManager) {
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
	config, err = util.NewConfig(configDir)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	flargs := config.ImportFlags(flagset)

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

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
// returns nil on success or an error on failure
func PromptForCredentials(config util.ConfigManager) error {
	userVar, _ := config.GetV("user")
	for userVar.Value == "" || userVar.Source != "INTERACTION" {
		if userVar.Value != "" {
			user := util.Prompt(fmt.Sprintf("User [%s]: ", userVar.Value))
			if strings.TrimSpace(user) == "" {
				config.Set("user", userVar.Value, "INTERACTION")
			} else {
				config.Set("user", strings.TrimSpace(user), "INTERACTION")
			}
		} else {
			user := util.Prompt("User: ")
			config.Set("user", strings.TrimSpace(user), "INTERACTION")
		}
		userVar, _ = config.GetV("user")
	}

	for config.GetIgnoreErr("pass") == "" {
		pass, err := speakeasy.FAsk(os.Stderr, "Pass: ")

		if err != nil {
			return err
		}
		config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
	}

	if config.GetIgnoreErr("yubikey") != "" {
		for config.GetIgnoreErr("yubikey-otp") == "" {
			yubikey := util.Prompt("Press yubikey: ")
			config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}
	log.Log("")
	return nil
}
