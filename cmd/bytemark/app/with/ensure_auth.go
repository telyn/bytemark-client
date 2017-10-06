package with

import (
	"fmt"
	"net/url"
	"strings"

	auth3 "gitlab.bytemark.co.uk/auth/client"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

func makeCredentials(config util.ConfigManager) (credents map[string]string, err error) {
	err = PromptForCredentials(config)
	if err != nil {
		return
	}
	credents = map[string]string{
		"username": config.GetIgnoreErr("user"),
		"password": config.GetIgnoreErr("pass"),
		"validity": config.GetIgnoreErr("session-validity"),
	}
	if useKey, _ := config.GetBool("yubikey"); useKey {
		credents["yubikey"] = config.GetIgnoreErr("yubikey-otp")
	}
	return
}

// EnsureAuth authenticates with the Bytemark authentication server, prompting for credentials if necessary.
// TODO(telyn): This REALLY, REALLY needs breaking apart into more manageable chunks
func EnsureAuth(client lib.Client, config util.ConfigManager) error {
	token := config.GetIgnoreErr("token")

	err := client.AuthWithToken(token)
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

			credents, err := makeCredentials(config)

			if err != nil {
				return err
			}
			err = client.AuthWithCredentials(credents)

			// Handle the special case here where we just need to prompt for 2FA and try again
			if err != nil && strings.Contains(err.Error(), "Missing 2FA") {
				for config.GetIgnoreErr("2fa-otp") == "" {
					token := util.Prompt("Enter 2FA token: ")
					config.Set("2fa-otp", strings.TrimSpace(token), "INTERACTION")
				}

				credents["2fa"] = config.GetIgnoreErr("2fa-otp")

				err = client.AuthWithCredentials(credents)
			}

			if err == nil {
				// success!
				// TODO(telyn): warn on failure to write to token
				_ = config.SetPersistent("token", client.GetSessionToken(), "AUTH")

				// Check this here, as it is only relevant the initial login,
				// not subsequent validations of the token (as opposed to yubikey)
				if config.GetIgnoreErr("2fa-otp") != "" {
					factors := client.GetSessionFactors()

					if config.GetIgnoreErr("2fa-otp") != "" {
						if !factorExists(factors, "2fa") {
							// Should never happen, as long as auth correctly returns the factors
							return fmt.Errorf("Unexpected error with 2FA login. Please report this as a bug")
						}
					}
				}

				break
			} else {
				if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") {
					if attempts > 0 {
						log.Errorf("Invalid credentials, please try again\r\n")
						config.Set("user", config.GetIgnoreErr("user"), "PRIOR INTERACTION")
						config.Set("pass", "", "INVALID")
						config.Set("yubikey-otp", "", "INVALID")
						config.Set("2fa-otp", "", "INVALID")
					} else {
						return err
					}
				} else {
					return err
				}

			}
		}
	}
	if config.GetIgnoreErr("yubikey") != "" {
		factors := client.GetSessionFactors()

		if config.GetIgnoreErr("yubikey") != "" {
			if !factorExists(factors, "yubikey") {
				// Current auth token doesn't have a yubikey,
				// so prompt the user to login again with yubikey

				// This happens when someone has logged in already,
				// but then tries to run a command with the
				// "yubikey" flag set

				config.Set("token", "", "FLAG yubikey")

				return EnsureAuth(client, config)
			}
		}
	}
	return nil
}

func factorExists(factors []string, factor string) bool {
	for _, f := range factors {
		if f == factor {
			return true
		}
	}

	return false
}
