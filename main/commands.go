package main

import (
	bigv "bigv.io/client/lib"
	"bufio"
	"fmt"
	"github.com/bgentry/speakeasy"
	"os"
	"strings"
)

// Commands represents the available commands in the BigV client. Each command should have its own function defined here, with a corresponding HelpFor* function too.
type Commands interface {
	EnsureAuth()

	Config([]string)
	CreateGroup([]string)
	DeleteVM([]string)
	Debug([]string)
	Help([]string)
	ShowAccount([]string)
	ShowVM([]string)
	UndeleteVM([]string)

	HelpForConfig()
	HelpForDebug()
	HelpForDelete()
	HelpForHelp()
	HelpForShow()
}

// CommandSet is the main implementation of the Commands interface
type CommandSet struct {
	bigv   bigv.Client
	config ConfigManager
}

// NewCommandSet creates a CommandSet given a ConfigManager and bigv.io/client/lib Client.
func NewCommandSet(config ConfigManager, client bigv.Client) *CommandSet {
	commandSet := new(CommandSet)
	commandSet.config = config
	commandSet.bigv = client
	return commandSet
}

// EnsureAuth authenticates with the BigV authentication server, prompting for credentials if necessary.
func (cmds *CommandSet) EnsureAuth() {
	err := cmds.bigv.AuthWithToken(cmds.config.Get("token"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to use token, trying credentials.\r\n\r\n")
		attempts := 3

		for err != nil {
			attempts--

			cmds.PromptForCredentials()
			credents := map[string]string{
				"username": cmds.config.Get("user"),
				"password": cmds.config.Get("pass"),
			}
			if cmds.config.Get("yubikey") != "" {
				credents["yubikey"] = cmds.config.Get("yubikey-otp")
			}

			err = cmds.bigv.AuthWithCredentials(credents)
			if err == nil {
				// sucess!
				cmds.config.SetPersistent("token", cmds.bigv.GetSessionToken(), "AUTH")
				break
			} else {
				if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") {
					if attempts > 0 {
						fmt.Fprintf(os.Stderr, "Invalid credentials, please try again\r\n")
						cmds.config.Set("user", "", "INVALID")
						cmds.config.Set("pass", "", "INVALID")
						cmds.config.Set("yubikey-otp", "", "INVALID")
					} else {
						exit(err, "Invalid credentials, giving up after three attempts.")
					}
				} else {
					exit(err)
				}

			}
		}
	}

}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
func (cmds *CommandSet) PromptForCredentials() {
	buf := bufio.NewReader(os.Stdin)
	for cmds.config.Get("user") == "" {
		fmt.Fprintf(os.Stderr, "User: ")
		user, _ := buf.ReadString('\n')
		cmds.config.Set("user", strings.TrimSpace(user), "INTERACTION")
	}

	for cmds.config.Get("pass") == "" {
		pass, err := speakeasy.Ask("Pass: ")

		if err != nil {
			exit(err, "Couldn't read password - are you sure you're using a terminal?")
		}
		cmds.config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
	}

	if cmds.config.Get("yubikey") != "" {
		for cmds.config.Get("yubikey-otp") == "" {
			fmt.Fprintf(os.Stderr, "Press yubikey: ")
			yubikey, _ := buf.ReadString('\n')
			cmds.config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}
	fmt.Fprintf(os.Stderr, "\r\n")

}
