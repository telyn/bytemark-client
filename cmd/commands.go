package cmd

import (
	bigv "bigv.io/client/lib"
	"bufio"
	"fmt"
	"github.com/bgentry/speakeasy"
	"os"
	"strings"
)

type Commands interface {
	EnsureAuth()

	Config([]string)
	Debug([]string)
	Help([]string)
	ShowAccount([]string)
	ShowVM([]string)

	HelpForConfig()
	HelpForHelp()
	HelpForShow()
}

type CommandSet struct {
	bigv   bigv.Client
	config ConfigManager
}

func NewCommandSet(config ConfigManager, client bigv.Client) *CommandSet {
	commandSet := new(CommandSet)
	commandSet.config = config
	commandSet.bigv = client
	return commandSet
}

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
						fmt.Fprintf(os.Stderr, "Invalid credentials, giving up after three attempts.\r\n")
						// TODO(telyn): define exit codes
						os.Exit(1)
					}
				} else {
					panic(err)
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
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	for cmds.config.Get("pass") == "" {
		pass, err := speakeasy.Ask("Pass: ")

		if err != nil {
			panic(err)
		}
		cmds.config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
		fmt.Fprintf(os.Stderr, "\r\n")
	}

	if cmds.config.Get("yubikey") != "" {
		for cmds.config.Get("yubikey-otp") == "" {
			fmt.Fprintf(os.Stderr, "Press yubikey: ")
			yubikey, _ := buf.ReadString('\n')
			cmds.config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}

}
