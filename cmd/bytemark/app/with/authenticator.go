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

type retryErr string

func (r retryErr) Error() string {
	return string(r)
}

type Authenticator struct {
	client lib.Client
	config util.ConfigManager
}

func NewAuthenticator(client lib.Client, config util.ConfigManager) Authenticator {
	return Authenticator{client: client, config: config}
}

func (a Authenticator) get2FAOTP() (otp string) {
	otp = a.config.GetIgnoreErr("2fa-otp")
	for otp == "" {
		token := util.Prompt("Enter 2FA token: ")
		a.config.Set("2fa-otp", strings.TrimSpace(token), "INTERACTION")
		otp = a.config.GetIgnoreErr("2fa-otp")
	}
	return otp
}

func (a Authenticator) tryCredentialsAttempt() error {
	credents, err := makeCredentials(a.config)

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
	err = fmt.Errorf("fake error")

	for err != nil {
		attempts--
		err = a.tryCredentialsAttempt()

		if err != nil {
			if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") {
				// if the credentials are bad in some way, make another attempt.
				if attempts <= 0 {
					return err
				}
				log.Errorf("Invalid credentials, please try again\r\n")
				a.config.Set("user", a.config.GetIgnoreErr("user"), "PRIOR INTERACTION")
				a.config.Set("pass", "", "INVALID")
				a.config.Set("yubikey-otp", "", "INVALID")
				a.config.Set("2fa-otp", "", "INVALID")
				continue
			} else {
				// if the credentials were okay and login still failed, let the user know
				return err
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
					return fmt.Errorf("Unexpected error with 2FA login. Please report this as a bug")
				}
			}
		}
	}
	return
}

func (a Authenticator) tryToken() error {
	token := a.config.GetIgnoreErr("token")

	return a.client.AuthWithToken(token)
}

func (a Authenticator) checkSession(shortCircuit bool) error {
	factors := a.client.GetSessionFactors()

	// make sure that we authenticated with a yubikey if we requested to do so
	if a.config.GetIgnoreErr("yubikey") != "" {
		// yubikey factor isn't included when we impersonate, so && !factorExists("impersonate")
		if !factorExists(factors, "yubikey") && !factorExists(factors, "impersonate") {
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
	if requestedUser != "" && currentUser != requestedUser {
		// if our token is already an impersonated one then we need to unset it
		// and start over
		if factorExists(factors, "impersonate") {
			return retryErr(fmt.Sprintf("Impersonation as %s requested but already impersonating %s", requestedUser, currentUser))
		} else {
			// if not then we need to run the impersonation - but if we already tried it we should just give up
			if shortCircuit {
				return fmt.Errorf("Impersonation as %s requested, but unable to impersonate as them - got %s instead", requestedUser, currentUser)
			}
			err := a.client.Impersonate(requestedUser)
			if err != nil {
				return err
			}
			// check that we got the right user this time
			return a.checkSession(true)
		}
	}

	if a.config.GetIgnoreErr("impersonate") == "" {
		if factorExists(factors, "impersonate") {
			err := a.config.Unset("impersonate")
			if err != nil {
				return fmt.Errorf("Couldn't edit config directory")
			}
			return retryErr("Impersonation was not requested but impersonation still happened")
		}
	}
	return nil
}

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
		fmt.Printf("%s. Retrying\n\n", err.Error())
		err = a.Authenticate()
	}
	return err
}
