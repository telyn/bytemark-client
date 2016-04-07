package main

import (
	auth3 "bytemark.co.uk/auth3/client"
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"flag"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

var client lib.Client
var commands = make([]cli.Command, 0)
var global = struct {
	Config util.ConfigManager
	Client lib.Client
	App    *cli.App
	Error  error
}{}

func baseAppSetup() {
	global.App = cli.NewApp()
	global.App.Commands = commands

}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			log.Error("\r\nCaught an interrupt - exiting.\r\n")
			os.Exit(int(util.E_TRAPPED_INTERRUPT))
		}

	}()

	baseAppSetup()

	// TODO(telyn): ok I haven't figured out a better way than this to integrate Config and stuff, but this way works for now.
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	configDir := flags.String("config-dir", "", "")
	flags.SetOutput(ioutil.Discard)

	flags.Parse(os.Args[1:])

	config, err := util.NewConfig(*configDir)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	global.Config = config
	global.App.Version = lib.GetVersion().String()

	global.Client, err = lib.New(global.Config.GetIgnoreErr("endpoint"), global.Config.GetIgnoreErr("billing-endpoint"))
	global.Client.SetDebugLevel(global.Config.GetDebugLevel())
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}
	//juggle the arguments in order to get the executable on the beginning
	flargs := flags.Args()
	newArgs := make([]string, len(flargs)+1)
	newArgs[0] = os.Args[0]
	copy(newArgs[1:], flargs)
	log.Logf("orig: %v\r\nflag: %v\r\n new: %v\r\n", os.Args, flag.Args(), newArgs)
	global.App.Run(newArgs)

	os.Exit(int(util.ProcessError(global.Error)))
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

			PromptForCredentials()
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
				global.Config.SetPersistent("token", global.Client.GetSessionToken(), "AUTH")
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
