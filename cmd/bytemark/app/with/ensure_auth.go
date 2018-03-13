package with

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
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
	authenticator := NewAuthenticator(client, config)
	return authenticator.Authenticate()
}

func factorExists(factors []string, factor string) bool {
	for _, f := range factors {
		if f == factor {
			return true
		}
	}

	return false
}
