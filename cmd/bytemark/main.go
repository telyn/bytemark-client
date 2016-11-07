package main

import (
	"flag"
	"fmt"
	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/bgentry/speakeasy"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

var forceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}

//commands is assembled during init()
var commands = make([]cli.Command, 0)
var global = struct {
	Config util.ConfigManager
	Client lib.Client
	App    *cli.App
}{}

func baseAppSetup() (app *cli.App, err error) {
	app = cli.NewApp()
	app.Version = lib.Version

	app.Commands = commands
	return

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
			os.Exit(int(util.E_TRAPPED_INTERRUPT))
		}

	}()

	overrideHelp()
	args := prepConfig()
	app, err := baseAppSetup()
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	global.App = app

	// build ourselves an api client
	auth, err := auth3.New(global.Config.GetIgnoreErr("auth-endpoint"))
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	cli := lib.NewWithAuth(global.Config.GetIgnoreErr("endpoint"), global.Config.GetIgnoreErr("billing-endpoint"), global.Config.GetIgnoreErr("spp-endpoint"), auth)
	global.Client = cli
	global.Client.SetDebugLevel(global.Config.GetDebugLevel())

	err = global.App.Run(args)

	os.Exit(int(util.ProcessError(err)))
}

// EnsureAuth authenticates with the Bytemark authentication server, prompting for credentials if necessary.
func EnsureAuth() error {
	token, err := global.Config.Get("token")

	err = global.Client.AuthWithToken(token)
	if err != nil {
		if aErr, ok := err.(*auth3.Error); ok {
			if _, ok := aErr.Err.(*url.Error); ok {
				return aErr
			}
		}
		log.Error("Please log in to Bytemark\r\n")
		attempts := 3

		for err != nil {
			attempts--

			err = PromptForCredentials()
			if err != nil {
				return err
			}
			credents := map[string]string{
				"username": global.Config.GetIgnoreErr("user"),
				"password": global.Config.GetIgnoreErr("pass"),
			}
			if useKey, _ := global.Config.GetBool("yubikey"); useKey {
				credents["yubikey"] = global.Config.GetIgnoreErr("yubikey-otp")
			}

			err = global.Client.AuthWithCredentials(credents)
			if err == nil {
				// sucess!
				// it doesn't _really_ matter if we can't write the token to the token place, right?
				_ = global.Config.SetPersistent("token", global.Client.GetSessionToken(), "AUTH")
				break
			} else {
				if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") {
					if attempts > 0 {
						log.Errorf("Invalid credentials, please try again\r\n")
						global.Config.Set("user", global.Config.GetIgnoreErr("user"), "PRIOR INTERACTION")
						global.Config.Set("pass", "", "INVALID")
						global.Config.Set("yubikey-otp", "", "INVALID")
					} else {
						return err
					}
				} else {
					return err
				}

			}
		}
	}
	if global.Config.GetIgnoreErr("yubikey") != "" {
		factors := global.Client.GetSessionFactors()
		for _, f := range factors {
			if f == "yubikey" {
				return nil
			}
		}
		// if still executing, we didn't have yubikey factor
		global.Config.Set("token", "", "FLAG yubikey")
		return EnsureAuth()
	}
	return nil
}

// mergeCommand merges src into dst, only copying non-nil fields of src,
// and calling mergeCommands upon the .Subcommands
// and appending all .Flags
func mergeCommand(dst *cli.Command, src cli.Command) {
	if src.Usage != "" {
		dst.Usage = src.Usage
	}
	if src.UsageText != "" {
		dst.UsageText = src.UsageText
	}
	if src.Description != "" {
		dst.Description = src.Description
	}
	if src.Action != nil {
		dst.Action = src.Action
	}
	if src.Flags != nil {
		for _, f := range src.Flags {
			dst.Flags = append(dst.Flags, f)
		}
	}
	if src.Subcommands != nil {
		dst.Subcommands = mergeCommands(dst.Subcommands, src.Subcommands)
	}
}

// mergeCommands copies over all the commands from base to result,
// then puts all the commands from extras in too, overwriting any provided fields.
func mergeCommands(base []cli.Command, extras []cli.Command) (result []cli.Command) {
	result = make([]cli.Command, len(base))
	copy(result, base)

	for _, cmd := range extras {
		found := false
		for idx := range result {
			if result[idx].Name == cmd.Name {
				mergeCommand(&result[idx], cmd)
				found = true
			}
		}
		if !found {
			result = append(result, cmd)
		}
	}
	return
}

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

func prepConfig() (args []string) {
	// set up our global flags because we need some config before we can set up our App
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	configDir := flags.String("config-dir", "", "")
	help := flags.Bool("help", false, "")
	h := flags.Bool("h", false, "")
	flags.Bool("yubikey", false, "")
	flags.Int("debug-level", 0, "")
	flags.String("user", "", "")
	flags.String("account", "", "")
	flags.String("endpoint", "", "")
	flags.String("billing-endpoint", "", "")
	flags.String("spp-endpoint", "", "")
	flags.String("auth-endpoint", "", "")
	flags.String("yubikey-otp", "", "")

	flags.SetOutput(ioutil.Discard)

	err := flags.Parse(os.Args[1:])
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	config, err := util.NewConfig(*configDir)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	// import the flags into config
	flargs := config.ImportFlags(flags)
	if config.GetIgnoreErr("endpoint") == "https://int.bigv.io" {
		config.Set("billing-endpoint", "", "CODE nullify billing-endpoint when using bigv-int")
		config.Set("spp-endpoint", "", "CODE nullify spp-endpoint when using bigv-int")
	}
	global.Config = config

	//juggle the arguments in order to get the executable on the beginning
	args = make([]string, len(flargs)+1)
	if len(flargs) > 0 && flargs[0] == "help" {
		copy(args[1:], flargs[1:])
		args[len(args)-1] = "--help"
	} else {
		copy(args[1:], flargs)
	}
	args[0] = os.Args[0]
	log.Debugf(log.LvlFlags, "orig: %v\r\nflag: %v\r\n new: %v\r\n", os.Args, flargs, args)

	if *help || *h {
		helpArgs := make([]string, len(args)+1)
		helpArgs[len(args)] = "--help"
		copy(helpArgs, args)
		args = helpArgs
	}
	return
}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
// returns nil on success or an error on failure
func PromptForCredentials() error {
	userVar, _ := global.Config.GetV("user")
	for userVar.Value == "" || userVar.Source != "INTERACTION" {
		if userVar.Value != "" {
			user := util.Prompt(fmt.Sprintf("User [%s]: ", userVar.Value))
			if strings.TrimSpace(user) == "" {
				global.Config.Set("user", userVar.Value, "INTERACTION")
			} else {
				global.Config.Set("user", strings.TrimSpace(user), "INTERACTION")
			}
		} else {
			user := util.Prompt("User: ")
			global.Config.Set("user", strings.TrimSpace(user), "INTERACTION")
		}
		userVar, _ = global.Config.GetV("user")
	}

	for global.Config.GetIgnoreErr("pass") == "" {
		pass, err := speakeasy.FAsk(os.Stderr, "Pass: ")

		if err != nil {
			return err
		}
		global.Config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
	}

	if global.Config.GetIgnoreErr("yubikey") != "" {
		for global.Config.GetIgnoreErr("yubikey-otp") == "" {
			yubikey := util.Prompt("Press yubikey: ")
			global.Config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}
	log.Log("")
	return nil
}
