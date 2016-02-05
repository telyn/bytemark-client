package lib

import (
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func TestGetAccount(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account" {
			w.Write([]byte(`{
			    "name": "account",
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
	client.AllowInsecureRequests()

	acc, err := client.GetAccount("invalid-account")
	is.NotNil(err)

	acc, err = client.GetAccount("")
	is.Nil(err)
	is.Equal("account", acc.Name)

	acc, err = client.GetAccount("account")
	is.Nil(err)
	is.Equal("account", acc.Name)

}

func TestGetAccounts(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts" {
			w.Write([]byte(`[
			{
			    "name": "account",
			    "id": 1
			}, {
			    "name": "Dr. Evil",
			    "suspended": true,
			    "id": 10
			}
			]`))
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
	client.AllowInsecureRequests()

	acc, err := client.GetAccounts()
	is.Nil(err)
	is.Equal(2, len(acc))
	is.Equal("Dr. Evil", acc[1].Name)

}
