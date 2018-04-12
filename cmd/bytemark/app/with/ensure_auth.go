package with

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

func (a Authenticator) makeCredentials() (credents map[string]string, err error) {
	err = a.promptForCredentials()
	if err != nil {
		return
	}
	credents = map[string]string{
		"username": a.config.GetIgnoreErr("user"),
		"password": a.config.GetIgnoreErr("pass"),
		"validity": a.config.GetIgnoreErr("session-validity"),
	}
	if useKey, _ := a.config.GetBool("yubikey"); useKey {
		credents["yubikey"] = a.config.GetIgnoreErr("yubikey-otp")
	}
	return
}

// EnsureAuth authenticates with the Bytemark authentication server, prompting for credentials if necessary.
func EnsureAuth(client lib.Client, config config.Manager) error {
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
