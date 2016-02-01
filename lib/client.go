package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
	"bytemark.co.uk/client/util/log"
	"errors"
)

// bytemarkClient is the main type in the Bytemark API client library
type bytemarkClient struct {
	endpoint    string
	auth        *auth3.Client
	authSession *auth3.SessionData
	debugLevel  int
}

// New creates a new Bytemark API client using the given Bytemark API endpoint and the default Bytemark auth endpoint
func New(endpoint string) (c *bytemarkClient, err error) {
	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}
	return NewWithAuth(endpoint, auth), nil
}

// NewWithAuth creates a new Bytemark API client using the given Bytemark API endpoint and bytemark.co.uk/auth3/client Client
func NewWithAuth(endpoint string, auth *auth3.Client) (c *bytemarkClient) {
	c = new(bytemarkClient)
	c.endpoint = endpoint
	c.debugLevel = 0
	c.auth = auth
	return c
}

// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
func (c *bytemarkClient) AuthWithCredentials(credentials auth3.Credentials) error {
	session, err := c.auth.CreateSession(credentials)
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

	session, err := c.auth.ReadSession(token)
	if err == nil {
		c.authSession = session
	}
	log.Debug(5, session)
	return err

}

// GetEndpoint returns the Bytemark API endpoint currently in use.
func (c *bytemarkClient) GetEndpoint() string {
	return c.endpoint
}

// SetDebugLevel sets the debug level / verbosity of the Bytemark API client. 0 (default) is silent.
func (c *bytemarkClient) SetDebugLevel(debugLevel int) {
	c.debugLevel = debugLevel
}

// GetSessionFactors returns the factors provided when the current auth session was set up
func (c *bytemarkClient) GetSessionFactors() []string {
	if c.authSession == nil {
		return []string{}
	}
	return c.authSession.Factors
}

// GetSessionToken returns the token for the current auth session
func (c *bytemarkClient) GetSessionToken() string {
	if c.authSession == nil {
		return ""
	}
	return c.authSession.Token
}

func (c *bytemarkClient) GetSessionUser() string {
	if c.authSession == nil {
		return ""
	}
	return c.authSession.Username
}
