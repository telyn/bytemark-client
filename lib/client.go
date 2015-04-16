package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
)

// Client is the main type in the BigV client library.
type BigVClient struct {
	endpoint    string
	auth        *auth3.Client
	authSession *auth3.SessionData
	debugLevel  int
}

// NewWithSession creates a new Client using the bytemark auth.Client you specify.
func NewWithSession(auth *auth3.Client, session *auth3.SessionData) (bigv *BigVClient, err error) {

	bigv = new(BigVClient)
	bigv.endpoint = "https://uk0.bigv.io"
	bigv.auth = auth
	bigv.authSession = session
	bigv.debugLevel = 0

	return bigv, nil
}

// NewWithCredentials creates a new Client authenticating against auth.bytemark.co.uk using the given credentials.
// If an alternative auth endpoint is needed, use NewWithSession
func NewWithCredentials(credentials auth3.Credentials) (bigv *BigVClient, err error) {
	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}

	session, err := auth.CreateSession(credentials)
	if err != nil {
		return nil, err
	}

	return NewWithSession(auth, session)
}

// NewWithToken creates a new Client authenticating against auth.bytemark.co.uk using the given token
// If an alternative auth endpoint is needed, use NewWithSession
func NewWithToken(token string) (bigv *BigVClient, err error) {

	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}

	session, err := auth.ReadSession(token)
	if err != nil {
		return nil, err
	}

	return NewWithSession(auth, session)
}

func (bigv *BigVClient) GetEndpoint() string {
	return bigv.endpoint
}

func (bigv *BigVClient) SetDebugLevel(debugLevel int) {
	bigv.debugLevel = debugLevel
}

// TODO(telyn): remove GetSessionToken - Dispatcher should get the AuthSession and pass it to NewWithSession
func (bigv *BigVClient) GetSessionToken() string {
	return bigv.authSession.Token
}
