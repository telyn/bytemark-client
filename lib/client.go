package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
)

// bigvClient is the main type in the BigV client library
type bigvClient struct {
	endpoint    string
	auth        *auth3.Client
	authSession *auth3.SessionData
	debugLevel  int
}

// New creates a new BigV client using the given BigV endpoint and the default Bytemark auth endpoint
func New(endpoint string) (bigv *bigvClient, err error) {
	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}
	return NewWithAuth(endpoint, auth), nil
}

// NewWithAuth creates a new BigV client using the given BigV endpoint and bytemark.co.uk/auth3/client Client
func NewWithAuth(endpoint string, auth *auth3.Client) (bigv *bigvClient) {
	bigv = new(bigvClient)
	bigv.endpoint = endpoint
	bigv.debugLevel = 0
	bigv.auth = auth
	return bigv
}

// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
func (bigv *bigvClient) AuthWithCredentials(credentials auth3.Credentials) error {
	session, err := bigv.auth.CreateSession(credentials)
	if err == nil {
		bigv.authSession = session
	}
	return err
}

// AuthWithToken attempts to read sessiondata from auth for the given token. Returns nil on success or an error otherwise.
func (bigv *bigvClient) AuthWithToken(token string) error {

	session, err := bigv.auth.ReadSession(token)
	if err == nil {
		bigv.authSession = session
	}
	return err

}

// GetEndpoint returns the BigV endpoint currently in use.
func (bigv *bigvClient) GetEndpoint() string {
	return bigv.endpoint
}

// SetDebugLevel sets the debug level / verbosity of the BigV client. 0 (default) is silent.
func (bigv *bigvClient) SetDebugLevel(debugLevel int) {
	bigv.debugLevel = debugLevel
}

// GetSessionToken returns the token for the current auth session
func (bigv *bigvClient) GetSessionToken() string {
	if bigv.authSession == nil {
		return ""
	}
	return bigv.authSession.Token
}

func (bigv *bigvClient) GetSessionUser() string {
	if bigv.authSession == nil {
		return ""
	}
	return bigv.authSession.Username
}
