package lib

import (
	"encoding/json"
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func getFixtureGroup() Group {
	vm := getFixtureVM()
	return Group{
		ID:   1,
		Name: "default",
		VirtualMachines: []*VirtualMachine{
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
			w.Write(str)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	client.AllowInsecureRequests()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	is.Nil(err)

	err = client.CreateGroup(GroupName{Group: "invalid-group", Account: "account"})
	is.Nil(err)
}

func TestDeleteGroup(t *testing.T) {
	is := is.New(t)
	groupHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/default" && req.Method == "DELETE" {
			w.Write([]byte(""))
		} else if req.URL.Path == "/accounts/account/groups/default" && req.Method == "GET" {
			str, err := json.Marshal(getFixtureGroup())
			if err != nil {
				t.Fatal(err)
			}
			w.Write(str)
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	client.AllowInsecureRequests()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})

	is.Nil(err)
	err = client.DeleteGroup(GroupName{Group: "default", Account: "account"})
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
			w.Write(str)
		} else if req.URL.Path == "/accounts/account/groups/invalid-group" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}
	client, authServer, brain, billing, err :=
		mkTestClientAndServers(http.HandlerFunc(groupHandler), mkNilHandler(t))

	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()
	client.AllowInsecureRequests()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})

	group, err := client.GetGroup(GroupName{Group: "invalid-group", Account: "account"})
	is.NotNil(err)

	group, err = client.GetGroup(GroupName{Group: "default", Account: "account"})
	is.NotNil(group)
	is.Nil(err)

	group, err = client.GetGroup(GroupName{Group: "", Account: ""})
	is.NotNil(group)
	is.Nil(err)
}
