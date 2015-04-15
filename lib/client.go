package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
)

// Client is the main type in the BigV client library.
type Client struct {
	Endpoint    string
	Auth        *auth3.Client
	AuthSession *auth3.SessionData
	DebugLevel  int
}

// NewWithSession creates a new Client using the bytemark auth.Client you specify.
func NewWithSession(auth *auth3.Client, session *auth3.SessionData) (bigv *Client, err error) {

	bigv = new(Client)
	bigv.Endpoint = "https://uk0.bigv.io"
	bigv.Auth = auth
	bigv.AuthSession = session
	bigv.DebugLevel = 0

	return bigv, nil
}

// NewWithCredentials creates a new Client authenticating against auth.bytemark.co.uk using the given credentials.
// If an alternative auth endpoint is needed, use NewWithSession
func NewWithCredentials(credentials auth3.Credentials) (bigv *Client, err error) {
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
func NewWithToken(token string) (bigv *Client, err error) {

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
