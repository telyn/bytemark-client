package lib

import (
	auth3 "bytemark.co.uk/client/lib/auth"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// mkTestClientAndServers constructs httptest Servers for a pretend auth and BigV endpoint, then constructs a BigVClient that uses those servers.
// The http.Handler passed is for the BigV endpoint - see the definition of mkTestAuthServer for the auth handler.
// Used to test that the right URLs are being requested and such.
func mkTestClientAndServers(handler http.Handler) (bigv *bigvClient, authServer *httptest.Server, brain *httptest.Server, err error) {
	authServer = mkTestAuthServer()
	brain = mkTestBrain(handler)

	auth, err := auth3.New(authServer.URL)
	return NewWithAuth(brain.URL, auth), authServer, brain, err
}

// mkTestBrain creates an httptest.Server for the given http.Handler. It's basically an alias for httptest.NewServer. Why did I write it?
func mkTestBrain(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func mkTestAuthServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`{
    "token": "fffffffffffffff",
    "username": "account",
    "factors": []
}`)
	}))

}
