package cmd

import (
	bigv "bigv.io/client/lib"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Commands interface {
	EnsureAuth()

	Debug([]string)
	Help([]string)
	Set([]string)
	ShowAccount([]string)
	ShowVM([]string)
	Unset([]string)

	HelpForHelp()
	HelpForSet()
	HelpForShow()
}

type CommandSet struct {
	bigv   bigv.Client
	config *Config
}

func NewCommandSet(config *Config, client bigv.Client) *CommandSet {
	commandSet := new(CommandSet)
	commandSet.config = config
	commandSet.bigv = client
	return commandSet
}

func (cmds *CommandSet) EnsureAuth() {
	err := cmds.bigv.AuthWithToken(cmds.config.Get("token"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to use token, trying credentials.\r\n\r\n")
		cmds.PromptForCredentials()
		credents := map[string]string{
			"username": cmds.config.Get("user"),
			"password": cmds.config.Get("pass"),
		}
		if cmds.config.Get("yubikey") != "" {
			credents["yubikey"] = cmds.config.Get("yubikey-otp")
		}

		err = cmds.bigv.AuthWithCredentials(credents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to use credentials.\r\n")
			panic(err)
		}
	}

	cmds.config.SetPersistent("token", cmds.bigv.GetSessionToken())
}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
func (cmds *CommandSet) PromptForCredentials() {
	buf := bufio.NewReader(os.Stdin)
	for cmds.config.Get("user") == "" {
		fmt.Fprintf(os.Stderr, "User: ")
		user, _ := buf.ReadString('\n')
		cmds.config.Set("user", strings.TrimSpace(user))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	for cmds.config.Get("pass") == "" {
		fmt.Fprintf(os.Stderr, "Pass: ")
		pass, _ := buf.ReadString('\n')
		cmds.config.Set("pass", strings.TrimSpace(pass))
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	if cmds.config.Get("yubikey") != "" {
		for cmds.config.Get("yubikey-otp") == "" {
			fmt.Fprintf(os.Stderr, "Press yubikey: ")
			yubikey, _ := buf.ReadString('\n')
			cmds.config.Set("yubikey-otp", strings.TrimSpace(yubikey))
		}
	}

}
