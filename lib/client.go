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

func New(endpoint string) (bigv *BigVClient, err error) {
	bigv = new(BigVClient)
	bigv.endpoint = endpoint
	bigv.debugLevel = 0

	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}
	bigv.auth = auth
	return bigv, nil
}

func NewWithAuth(endpoint string, auth *auth3.Client) (bigv *BigVClient) {
	bigv = new(BigVClient)
	bigv.endpoint = endpoint
	bigv.debugLevel = 0
	bigv.auth = auth
	return bigv
}

func (bigv *BigVClient) AuthWithCredentials(credentials auth3.Credentials) error {
	session, err := bigv.auth.CreateSession(credentials)
	if err == nil {
		bigv.authSession = session
	}
	return err
}

func (bigv *BigVClient) AuthWithToken(token string) error {

	session, err := bigv.auth.ReadSession(token)
	if err == nil {
		bigv.authSession = session
	}
	return err

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
