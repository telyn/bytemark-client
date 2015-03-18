package cmd

import (
	client "bigv.io/client/lib"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// type Dispatcher is used to create API requests and direct output to views,
// except probably when those API requests don't require authorisation (e.g. /definitions, new user)
type Dispatcher struct {
	Config *Config
	Flags  *flag.FlagSet
	BigV   *client.Client
}

// This is trying to do a little too much
func NewDispatcher(config *Config) (d *Dispatcher) {
	d = new(Dispatcher)
	d.Config = config
	c, err := client.NewWithToken(config.Get("token"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to use token, trying credentials.\r\n\r\n")
		d.PromptForCredentials()
		credents := map[string]string{
			"username": config.Get("user"),
			"password": config.Get("pass"),
		}
		if config.Get("yubikey") != "" {
			credents["yubikey"] = config.Get("yubikey-otp")
		}

		c, err = client.NewWithCredentials(credents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to use credentials.\r\n")
			panic(err)
		}
	}
	d.BigV = c
	return d
}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs more for loop to ensure that they don't stay empty.
func (d *Dispatcher) PromptForCredentials() {
	buf := bufio.NewReader(os.Stdin)
	if d.Config.Get("user") == "" {
		fmt.Fprintf(os.Stderr, "User: ")
		user, _ := buf.ReadString('\n')
		d.Config.Set("user", strings.TrimSpace(user))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	if d.Config.Get("pass") == "" {
		fmt.Fprintf(os.Stderr, "Pass: ")
		pass, _ := buf.ReadString('\n')
		d.Config.Set("pass", strings.TrimSpace(pass))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	if d.Config.Get("yubikey") != "" && d.Config.Get("yubikey-otp") == "" {
		fmt.Fprintf(os.Stderr, "Press yubikey: ")
		yubikey, _ := buf.ReadString('\n')
		d.Config.Set("yubikey-otp", strings.TrimSpace(yubikey))
	}

}

// Do takes the command line arguments and figures out what to do
func (dispatch *Dispatcher) Do(args []string) {
	//	help := dispatch.Flags.Lookup("help")
	///	fmt.Printf("%+v", help)
	if dispatch.BigV.DebugLevel >= 1 {
		fmt.Fprintf(os.Stderr, "Args passed to Do: %#v\n", args)
	}

	if /*help == true || */ len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Printf("No command specified.\n\n")
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
	case "debug":
		dispatch.Debug(args[1:])
		return
	case "show-account":
		dispatch.ShowAccount(args[1:])
		return
	}
}
