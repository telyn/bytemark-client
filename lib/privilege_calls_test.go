package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"net/http"
	"testing"
)

var testPrivileges = brain.Privileges{
	{
		ID:               6,
		VirtualMachineID: 23,
		Level:            brain.VMAdminPrivilege,
		Username:         "test-user",
	}, {
		ID:        999,
		Username:  "satan",
		AccountID: 433224,
		Level:     brain.AccountAdminPrivilege,
	},
}

func mkGetPrivilegesHandler(t *testing.T, user string) func(http.ResponseWriter, *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Wrong method %s", r.Method)
		}
		privs := testPrivileges
		if user != "" {
			privs = brain.Privileges{}
			for _, p := range testPrivileges {
				if p.Username == user {
					privs = append(privs, p)
				}
			}
		}
		js, err := json.Marshal(privs)
		if err != nil {
			t.Fatalf("couldn't marshal testPrivileges: %s", err.Error())
		}
		_, err = wr.Write(js)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetPrivileges(t *testing.T) {
	handlers := MuxHandlers{
		brain: Mux{
			"/privileges": mkGetPrivilegesHandler(t, ""),
		},
	}
	client, servers, err := mkTestClientAndServers(t, handlers.ToHandlers())
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	privileges, err := client.GetPrivileges("")
	if err != nil {
		t.Error(err)
	}
	if len(privileges) != 2 {
		t.Errorf("Wrong number of privileges: %d", len(privileges))
	}

}

func TestGrantPrivilege(t *testing.T) {
	done := false
	handlers := MuxHandlers{
		brain: Mux{
			"/users/satan/privileges": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Wrong method %s", r.Method)
				}
				// test the data gets there?
				done = true
			},
		},
	}
	client, servers, err := mkTestClientAndServers(t, handlers.ToHandlers())
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.GrantPrivilege(brain.Privilege{
		Username:  "satan",
		AccountID: 433224,
		Level:     brain.AccountAdminPrivilege,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !done {
		t.Errorf("Privilege was never added")
	}
}

func TestRevokePrivilegeWithID(t *testing.T) {
	done := false
	handlers := MuxHandlers{
		brain: Mux{
			"/privileges/999": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Wrong method %s", r.Method)
				}
				done = true
			},
		},
	}
	client, servers, err := mkTestClientAndServers(t, handlers.ToHandlers())
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.RevokePrivilege(brain.Privilege{
		ID:        999,
		Username:  "satan",
		AccountID: 433224,
		Level:     brain.AccountAdminPrivilege,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !done {
		t.Errorf("Privilege 999 was never deleted")
	}
}

func TestRevokePrivilegeWithoutID(t *testing.T) {
	done := false
	handlers := MuxHandlers{
		brain: Mux{
			"/users/satan/privileges": mkGetPrivilegesHandler(t, "satan"),
			"/privileges/999": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Wrong method %s", r.Method)
				}
				done = true
			},
		},
	}
	client, servers, err := mkTestClientAndServers(t, handlers.ToHandlers())
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.RevokePrivilege(brain.Privilege{
		Username:  "satan",
		AccountID: 433224,
		Level:     brain.AccountAdminPrivilege,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !done {
		t.Errorf("Privilege 999 was never deleted")
	}
}
