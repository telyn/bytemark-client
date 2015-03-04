package client

import (
	auth3 "bytemark.co.uk/auth3/client"
)

type Client struct {
	Endpoint    string
	Auth        *auth3.Client
	AuthSession *auth3.SessionData
	DebugLevel  int
}

// Creates a new Client using the auth.Client you specify.
func NewWithSession(auth *auth3.Client, session *auth3.SessionData) (bigv *Client, err error) {
	// check your session's good or smth

	bigv = new(Client)
	bigv.Endpoint = "https://uk0.bigv.io"
	bigv.Auth = auth
	bigv.AuthSession = session
	bigv.DebugLevel = 0

	return bigv, nil
}

// Creates a new Client authenticating against auth.bytemark.co.uk
func NewWithCredentials(credentials auth3.Credentials) (bigv *Client, err error) {
	// sort out all the auth shit yo
	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}

	session, err := auth.CreateSession(credentials)

	return NewWithSession(auth, session)
}

// Creates a new Client authenticating against auth.bytemark.co.uk
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
