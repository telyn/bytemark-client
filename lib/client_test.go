package lib

import (
	"fmt"
	auth3 "github.com/BytemarkHosting/auth-client"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mkTestClientAndServers constructs httptest Servers for a pretend auth and API endpoint, then constructs a Client that uses those servers.
// The http.Handler passed is for the API endpoint - see the definition of mkTestAuthServer for the auth handler.
// Used to test that the right URLs are being requested and such.
func mkTestClientAndServers(brainHandler http.Handler, billingHandler http.Handler) (c *bytemarkClient, authServer *httptest.Server, brain *httptest.Server, billing *httptest.Server, err error) {
	authServer = mkTestAuthServer()
	brain = mkTestServer(brainHandler)
	billing = mkTestServer(billingHandler)

	auth, err := auth3.New(authServer.URL)
	//FIXME: the "" is the spp endpoint. at the moment there are no tests that hit it.
	client := NewWithAuth(brain.URL, billing.URL, "", auth)
	client.AllowInsecureRequests()
	return client, authServer, brain, billing, err
}

// mkTestServer creates an httptest.Server for the given http.Handler. It's basically an alias for httptest.NewServer. Why did I write it?
func mkTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func mkNilHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Unexpected request to a nil server\r\n%s %s", r.Method, r.URL.String())
	})
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
