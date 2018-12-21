package lib

import (
	"context"
	"errors"
	"strings"

	auth3 "github.com/BytemarkHosting/auth-client"
)

// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
func (c *bytemarkClient) AuthWithCredentials(credentials auth3.Credentials) error {
	session, err := c.auth.CreateSession(context.TODO(), credentials)
	if err == nil {
		c.authSession = session
	}
	return err
}

// AuthWithToken attempts to read sessiondata from auth for the given token. Returns nil on success or an error otherwise.
func (c *bytemarkClient) AuthWithToken(token string) error {
	if token == "" {
		return errors.New("No token provided")
	}

	if strings.HasPrefix(token, "apikey.") {
		c.authSession = &auth3.SessionData{
			Username: "apikeyuser",
			Factors: []string{
				"apikey",
			},
			Token: token,
		}
		c.urls.Billing = ""
		return nil
	}

	session, err := c.auth.ReadSession(context.TODO(), token)
	if err == nil {
		c.authSession = session
	}
	return err

}

// Impersonate creates a session for the given user (assuming the client has already authenticated as someone who can)
func (c *bytemarkClient) Impersonate(user string) (err error) {
	c.authSession, err = c.auth.CreateImpersonatedSession(context.TODO(), c.authSession.Token, user)

	return
}
