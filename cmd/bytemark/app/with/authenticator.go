package with

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	auth3 "gitlab.bytemark.co.uk/auth/client"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

type retryErr string

func (r retryErr) Error() string {
	return string(r)
}

// Authenticator is an object which is temporarily created during the authentication phase of the client
// (i.e. during with.Auth) and is responsible for making requests to the auth server,
// making sure the user is logged in as who they want to be, performing impersonation,
// and prompting the user for any missing credentials.
// TODO(telyn): ensure prompts / output comes out on app.ErrWriter.
type Authenticator struct {
	client       lib.Client
	config       util.ConfigManager
	prompter     util.Prompter
	passPrompter passPrompter
}

// NewAuthenticator creates a new authenticator which will prompt on stderr (expecting input on stdin) and using speakeasy for password prompts.
func NewAuthenticator(client lib.Client, config util.ConfigManager) Authenticator {
	return Authenticator{
		client:       client,
		config:       config,
		prompter:     util.NewPrompter(),
		passPrompter: speakeasyWrapper{},
	}
}

func (a Authenticator) get2FAOTP() (otp string) {
	otp = a.config.GetIgnoreErr("2fa-otp")
	for otp == "" {
		token := a.prompter.Prompt("Enter 2FA token: ")
		a.config.Set("2fa-otp", strings.TrimSpace(token), "INTERACTION")
		otp = a.config.GetIgnoreErr("2fa-otp")
	}
	return otp
}

func (a Authenticator) tryCredentialsAttempt() error {
	credents, err := a.makeCredentials()

	if err != nil {
		return err
	}
	err = a.client.AuthWithCredentials(credents)

	// Handle the special case here where we just need to prompt for 2FA and try again
	if err != nil && strings.Contains(err.Error(), "Missing 2FA") {
		otp := a.get2FAOTP()

		credents["2fa"] = otp

		err = a.client.AuthWithCredentials(credents)
	}
	return err
}

func (a Authenticator) tryCredentials() (err error) {
	attempts := 3
	err = errors.New("fake error")

	for err != nil {
		attempts--
		err = a.tryCredentialsAttempt()

		if err != nil {
			if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") || strings.Contains(err.Error(), "Authentication failed") {
				// if the credentials are bad in some way, make another attempt.
				if attempts <= 0 {
					return
				}
				log.Errorf("Invalid credentials, please try again\r\n")
				// reset all credentials and set the default user to
				// whoever the last login attempt was for to make the prompt nicer
				a.config.Set("user", a.config.GetIgnoreErr("user"), "PRIOR INTERACTION")
				a.config.Set("pass", "", "INVALID")
				a.config.Set("yubikey-otp", "", "INVALID")
				a.config.Set("2fa-otp", "", "INVALID")
				continue
			} else {
				// if the credentials were okay and login still failed, let the user know
				return
			}
		} else {
			// we have successfully authenticated!

			// TODO(telyn): warn on failure to write to token
			_ = a.config.SetPersistent("token", a.client.GetSessionToken(), "AUTH")

			// Check that the 2fa factor was set if --2fa-otp was specified.
			// Checking here rather than in checkSession as it is only relevant
			// during the initial login, not subsequent validations of the
			// token (as opposed to yubikey)
			if a.config.GetIgnoreErr("2fa-otp") != "" {
				factors := a.client.GetSessionFactors()

				if !factorExists(factors, "2fa") {
					// Should never happen, as long as auth correctly returns the factors
					return errors.New("Unexpected error with 2FA login. Please report this as a bug")
				}
			}
		}
	}
	return
}

func (a Authenticator) tryToken() error {
	token := a.config.GetIgnoreErr("token")
	if token == "" {
		return errors.New("blank token")
	}

	return a.client.AuthWithToken(token)
}

func (a Authenticator) checkSession(shortCircuit bool) error {
	factors := a.client.GetSessionFactors()

	// make sure that we authenticated with a yubikey if we requested to do so
	if a.config.GetIgnoreErr("yubikey") != "" {
		// yubikey factor isn't included when we impersonate, so && !factorExists("impersonated")
		if !factorExists(factors, "yubikey") && !factorExists(factors, "impersonated") {
			// prompt the user to login again with yubikey

			// This happens when someone has logged in already,
			// but then tries to run a command with the
			// "yubikey" flag set

			a.config.Set("token", "", "FLAG yubikey")

			return EnsureAuth(a.client, a.config)
		}
	}

	currentUser := a.client.GetSessionUser()
	requestedUser := a.config.GetIgnoreErr("impersonate")

	// if we want to impersonate someone and we're not currently them
	if requestedUser != "" && currentUser != requestedUser {
		// if we already tried impersonating we should just give up
		if shortCircuit {
			err := a.config.Unset("token")
			if err != nil {
				return fmt.Errorf("Couldn't unset token: %v", err)
			}
			return fmt.Errorf("Impersonation as %s requested, but unable to impersonate as them - got %s instead", requestedUser, currentUser)
		}
		// if our token is already an impersonated one then we need to unset it
		// and start over
		if factorExists(factors, "impersonated") {
			err := a.config.Unset("token")
			if err != nil {
				return fmt.Errorf("Couldn't unset token: %v", err)
			}
			return retryErr(fmt.Sprintf("Impersonation as %s requested but already impersonating %s", requestedUser, currentUser))
		}
		// if not, run impersonation and see
		err := a.client.Impersonate(requestedUser)
		if err != nil {
			return err
		}
		// check that we got the right user this time
		return a.checkSession(true)
	} else if requestedUser == "" {
		// we didn't want to impersonate
		if factorExists(factors, "impersonated") {
			// but we got an impersonated token anyway, so unset token and retry
			err := a.config.Unset("token")
			if err != nil {
				return fmt.Errorf("Couldn't unset token: %v", err)
			}
			return retryErr("Impersonation was not requested but impersonation still happened")
		}
		// and we didn't impersonate but we aren't logged in as who we want to be
		if currentUser != a.config.GetIgnoreErr("user") {
			// we didn't want to impersonate anyone and we're not ourselves
			// so unset the token and cry about it
			err := a.config.Unset("token")
			if err != nil {
				return fmt.Errorf("Couldn't unset token: %v", err)
			}
			return fmt.Errorf("Expected to log in as %s but logged in as %s", a.config.GetIgnoreErr("user"), currentUser)
		}
	}

	return nil
}

// Authenticate performs authentication, checks and impersonation.
func (a Authenticator) Authenticate() error {
	err := a.tryToken()
	if err != nil {
		// check for url.Error cause that indicates something worse than a simple auth fail.
		if aErr, ok := err.(*auth3.Error); ok {
			if _, ok := aErr.Err.(*url.Error); ok {
				return aErr
			}
		}

		log.Error("Please log in to Bytemark\r\n")

		err = a.tryCredentials()
		if err != nil {
			return err
		}
	}

	err = a.checkSession(false)
	if _, ok := err.(retryErr); ok {
		err = a.Authenticate()
	}
	return err
}
