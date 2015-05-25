package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func mkTestClientAndServers(handler http.Handler) (bigv *BigVClient, authServer *httptest.Server, brain *httptest.Server, err error) {
	authServer = mkTestAuthServer()
	brain = mkTestBrain(handler)

	auth, err := auth3.New(authServer.URL)
	return NewWithAuth(brain.URL, auth), authServer, brain, err
}

func mkTestBrain(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func mkTestAuthServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`{
    "token": "fffffffffffffff",
    "username": "valid-user",
    "factors": []
}`)
	}))

}
