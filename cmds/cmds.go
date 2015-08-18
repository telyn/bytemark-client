package cmds

import (
	util "bigv.io/client/cmds/util"
	bigv "bigv.io/client/lib"
	"bufio"
	auth3 "bytemark.co.uk/auth3/client"
	"fmt"
	"github.com/bgentry/speakeasy"
	"net/url"
	"os"
	"strings"
)

type CommandFunc func([]string) util.ExitCode

// Commands represents the available commands in the BigV client. Each command should have its own function defined here, with a corresponding HelpFor* function too.
type CommandManager interface {
	EnsureAuth() error

	Config([]string) util.ExitCode
	Console([]string) util.ExitCode
	CreateGroup([]string) util.ExitCode
	CreateVM([]string) util.ExitCode
	DeleteGroup([]string) util.ExitCode
	DeleteVM([]string) util.ExitCode
	Debug([]string) util.ExitCode
	Help([]string) util.ExitCode
	LockHWProfile([]string) util.ExitCode
	UnlockHWProfile([]string) util.ExitCode
	ResetVM([]string) util.ExitCode
	Restart([]string) util.ExitCode
	SetCores([]string) util.ExitCode
	SetHWProfile([]string) util.ExitCode
	SetMemory([]string) util.ExitCode
	Start([]string) util.ExitCode
	Stop([]string) util.ExitCode
	Shutdown([]string) util.ExitCode
	ShowAccount([]string) util.ExitCode
	ShowGroup([]string) util.ExitCode
	ShowVM([]string) util.ExitCode
	UndeleteVM([]string) util.ExitCode

	HelpForConfig() util.ExitCode
	HelpForCreate() util.ExitCode
	HelpForDebug() util.ExitCode
	HelpForDelete() util.ExitCode
	HelpForHelp() util.ExitCode
	HelpForLocks() util.ExitCode
	HelpForPower() util.ExitCode
	HelpForSet() util.ExitCode
	HelpForShow() util.ExitCode
}

// CommandSet is the main implementation of the Commands interface
type CommandSet struct {
	bigv   bigv.Client
	config util.ConfigManager
}

// NewCommandSet creates a CommandSet given a ConfigManager and bigv.io/client/lib Client.
func NewCommandSet(config util.ConfigManager, client bigv.Client) *CommandSet {
	commandSet := new(CommandSet)
	commandSet.config = config
	commandSet.bigv = client
	return commandSet
}

// EnsureAuth authenticates with the BigV authentication server, prompting for credentials if necessary.
func (cmds *CommandSet) EnsureAuth() error {
	token, err := cmds.config.Get("token")

	err = cmds.bigv.AuthWithToken(token)
	if err != nil {
		if aErr, ok := err.(*auth3.Error); ok {
			if _, ok := aErr.Err.(*url.Error); ok {
				return aErr
			}
		}
		fmt.Fprintf(os.Stderr, "Please log in to BigV\r\n\r\n")
		attempts := 3

		for err != nil {
			attempts--

			cmds.PromptForCredentials()
			credents := map[string]string{
				"username": cmds.config.GetIgnoreErr("user"),
				"password": cmds.config.GetIgnoreErr("pass"),
			}
			if useKey, _ := cmds.config.GetBool("yubikey"); useKey {
				credents["yubikey"] = cmds.config.GetIgnoreErr("yubikey-otp")
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
						return err
					}
				} else {
					return err
				}

			}
		}
	}
	return nil

}

// PromptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
// returns nil on success or an error on failure
func (cmds *CommandSet) PromptForCredentials() error {
	buf := bufio.NewReader(os.Stdin)
	for cmds.config.GetIgnoreErr("user") == "" {
		fmt.Fprintf(os.Stderr, "User: ")
		user, _ := buf.ReadString('\n')
		cmds.config.Set("user", strings.TrimSpace(user), "INTERACTION")
	}

	for cmds.config.GetIgnoreErr("pass") == "" {
		pass, err := speakeasy.Ask("Pass: ")

		if err != nil {
			return err
		}
		cmds.config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
	}

	if cmds.config.GetIgnoreErr("yubikey") != "" {
		for cmds.config.GetIgnoreErr("yubikey-otp") == "" {
			fmt.Fprintf(os.Stderr, "Press yubikey: ")
			yubikey, _ := buf.ReadString('\n')
			cmds.config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}
	fmt.Fprintf(os.Stderr, "\r\n")
	return nil
}
