package with

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// promptForCredentials ensures that user, pass and yubikey-otp are defined, by prompting the user for them.
// needs a for loop to ensure that they don't stay empty.
// returns nil on success or an error on failure
func (a Authenticator) promptForCredentials() error {
	userVar, _ := a.config.GetV("user")
	for userVar.Value == "" || userVar.Source != "INTERACTION" {
		if userVar.Value != "" {
			user := a.prompter.Prompt(fmt.Sprintf("User [%s]: ", userVar.Value))
			if strings.TrimSpace(user) == "" {
				a.config.Set("user", userVar.Value, "INTERACTION")
			} else {
				a.config.Set("user", strings.TrimSpace(user), "INTERACTION")
			}
		} else {
			user := a.prompter.Prompt("User: ")
			a.config.Set("user", strings.TrimSpace(user), "INTERACTION")
		}
		userVar, _ = a.config.GetV("user")
	}

	for a.config.GetIgnoreErr("pass") == "" {
		pass, err := a.passPrompter.Ask("Pass: ")

		if err != nil {
			return err
		}
		a.config.Set("pass", strings.TrimSpace(pass), "INTERACTION")
	}

	if a.config.GetIgnoreErr("yubikey") != "" {
		for a.config.GetIgnoreErr("yubikey-otp") == "" {
			yubikey := a.prompter.Prompt("Press yubikey: ")
			a.config.Set("yubikey-otp", strings.TrimSpace(yubikey), "INTERACTION")
		}
	}
	log.Log("")
	return nil
}
