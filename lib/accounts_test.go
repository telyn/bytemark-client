package lib

import (
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/valid-account" {
			w.Write([]byte(`{
    "name": "valid-account",
    "id": 1
}`))
		} else if req.URL.Path == "/accounts/invalid-account" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}))
	defer authServer.Close()
	defer brain.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	acc, err := client.GetAccount("invalid-account")
	is.Nil(acc)
	is.NotNil(err)

	acc, err = client.GetAccount("valid-account")
	is.NotNil(acc)
	is.Equal("valid-account", acc.Name)

}
