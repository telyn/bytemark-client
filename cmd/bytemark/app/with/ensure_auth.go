package with

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/auth"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// EnsureAuth authenticates with the Bytemark authentication server, prompting for credentials if necessary.
func EnsureAuth(client lib.Client, config config.Manager) error {
	authenticator := auth.NewAuthenticator(client, config)
	return authenticator.Authenticate()
}
