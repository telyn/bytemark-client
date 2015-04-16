package cmd

import (
	client "bigv.io/client/lib"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type Dispatcher is used to create API requests and direct output to views,
// except probably when those API requests don't require authorisation (e.g. /definitions, new user)
type Dispatcher struct {
	Config *Config
	Flags  *flag.FlagSet
	BigV   client.Client
}

// NewDispatcher creates a new Dispatcher given a config.
func NewDispatcher(config *Config) (d *Dispatcher) {
	d = new(Dispatcher)
	d.Config = config
	return d
}

// EnsureAuth makes sure a valid token is stored in config.
// This should be called by anything that needs auth.
func (d *Dispatcher) EnsureAuth() {
	c, err := client.NewWithToken(d.Config.Get("token"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to use token, trying credentials.\r\n\r\n")
		d.PromptForCredentials()
		credents := map[string]string{
			"username": d.Config.Get("user"),
			"password": d.Config.Get("pass"),
		}
		if d.Config.Get("yubikey") != "" {
			credents["yubikey"] = d.Config.Get("yubikey-otp")
		}

		c, err = client.NewWithCredentials(credents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to use credentials.\r\n")
			panic(err)
		}
	}
	debugLevel, err := strconv.ParseInt(d.Config.Get("debug-level"), 10, 0)
	if err == nil {
		c.SetDebugLevel(int(debugLevel))
	}
	d.BigV = c

	d.Config.SetPersistent("token", d.BigV.GetSessionToken())
}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
func (d *Dispatcher) PromptForCredentials() {
	buf := bufio.NewReader(os.Stdin)
	for d.Config.Get("user") == "" {
		fmt.Fprintf(os.Stderr, "User: ")
		user, _ := buf.ReadString('\n')
		d.Config.Set("user", strings.TrimSpace(user))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	for d.Config.Get("pass") == "" {
		fmt.Fprintf(os.Stderr, "Pass: ")
		pass, _ := buf.ReadString('\n')
		d.Config.Set("pass", strings.TrimSpace(pass))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	if d.Config.Get("yubikey") != "" {
		for d.Config.Get("yubikey-otp") == "" {
			fmt.Fprintf(os.Stderr, "Press yubikey: ")
			yubikey, _ := buf.ReadString('\n')
			d.Config.Set("yubikey-otp", strings.TrimSpace(yubikey))
		}
	}

}

// TODO(telyn): Write a test for Do. Somehow.

// Do takes the command line arguments and figures out what to do
func (dispatch *Dispatcher) Do(args []string) {
	//	help := dispatch.Flags.Lookup("help")
	///	fmt.Printf("%+v", help)
	debugLevel, err := strconv.ParseInt(dispatch.Config.Get("debug-level"), 10, 0)
	if err != nil {
		debugLevel = 0
	}

	if debugLevel >= 1 {
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
	case "set":
		dispatch.Set(args[1:])
		return
	case "show-account":
		dispatch.ShowAccount(args[1:])
		return
	case "show-vm":
		dispatch.ShowVM(args[1:])
		return
	case "unset":
		dispatch.Unset(args[1:])
		return
	}
}
