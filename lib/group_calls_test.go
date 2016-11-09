package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func getFixtureGroup() brain.Group {
	vm := getFixtureVM()
	return brain.Group{
		ID:   1,
		Name: "default",
		VirtualMachines: []*brain.VirtualMachine{
			&vm,
		},
	}

}

func TestCreateGroup(t *testing.T) {
	is := is.New(t)
	groupHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups" && req.Method == "POST" {
			// TODO: unmarshal the groupname and check
			str, err := json.Marshal(getFixtureGroup())
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.Write(str)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	is.Nil(err)

	err = client.CreateGroup(&GroupName{Group: "invalid-group", Account: "account"})
	is.Nil(err)
}

func TestDeleteGroup(t *testing.T) {
	is := is.New(t)
	groupHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/default" && req.Method == "DELETE" {
			_, err := w.Write([]byte(""))
			if err != nil {
				t.Fatal(err)
			}
		} else if req.URL.Path == "/accounts/account/groups/default" && req.Method == "GET" {
			str, err := json.Marshal(getFixtureGroup())
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.Write(str)
			if err != nil {
				t.Fatal(err)
			}
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})

	is.Nil(err)
	err = client.DeleteGroup(&GroupName{Group: "default", Account: "account"})
	is.Nil(err)

}

func TestGetGroup(t *testing.T) {
	is := is.New(t)
	groupHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/default" {
			str, err := json.Marshal(getFixtureGroup())
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.Write(str)
			if err != nil {
				t.Fatal(err)
			}
		} else if req.URL.Path == "/accounts/account/groups/invalid-group" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	billingHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/accounts" {
			_, err := w.Write([]byte(`[{ "bigv_account_subscription": "account" }]`))
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	group, err := client.GetGroup(&GroupName{Group: "invalid-group", Account: "account"})
	is.NotNil(err)
	is.Equal("", group.Name)

	group, err = client.GetGroup(&GroupName{Group: "default", Account: "account"})
	is.NotNil(group)
	is.Nil(err)

	authServer.Close()
	brain.Close()
	billing.Close()

	client, authServer, brain, billing, err =
		mkTestClientAndServers(http.HandlerFunc(groupHandler), http.HandlerFunc(billingHandler))
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	group, err = client.GetGroup(&GroupName{Group: "", Account: ""})
	is.NotNil(group)
	is.Nil(err)
}
