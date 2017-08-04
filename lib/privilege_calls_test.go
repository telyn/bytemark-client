package lib

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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
		assertMethod(t, r, "GET")
		privs := testPrivileges
		if user != "" {
			privs = brain.Privileges{}
			for _, p := range testPrivileges {
				if p.Username == user {
					privs = append(privs, p)
				}
			}
		}
		writeJSON(t, wr, privs)
	}
}

func TestGetPrivileges(t *testing.T) {
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			// simple test of getting all privileges - hence the empty string
			"/privileges": mkGetPrivilegesHandler(t, ""),
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	// simple test of getting all privileges - hence the empty string
	privileges, err := client.GetPrivileges("")
	if err != nil {
		t.Error(err)
	}
	if len(privileges) != 2 {
		t.Errorf("Wrong number of privileges: %d", len(privileges))
	}

}

func TestGetPrivilegesForAccount(t *testing.T) {
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/accounts/mycoolaccount/privileges": mkGetPrivilegesHandler(t, ""),
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	privileges, err := client.GetPrivilegesForAccount("mycoolaccount")
	if err != nil {
		t.Error(err)
	}
	if len(privileges) != 2 {
		t.Errorf("Wrong number of privileges: %d", len(privileges))
	}
}

func TestGetPrivilegesForGroup(t *testing.T) {
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/accounts/test-account/groups/test-group/privileges": mkGetPrivilegesHandler(t, ""),
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	privileges, err := client.GetPrivilegesForGroup(GroupName{Group: "test-group", Account: "test-account"})
	if err != nil {
		t.Error(err)
	}
	if len(privileges) != 2 {
		t.Errorf("Wrong number of privileges: %d", len(privileges))
	}
}

func TestGetPrivilegesForServer(t *testing.T) {
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/accounts/test-account/groups/test-group/virtual_machines/test-vm/privileges": mkGetPrivilegesHandler(t, ""),
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	privileges, err := client.GetPrivilegesForVirtualMachine(VirtualMachineName{
		Account:        "test-account",
		Group:          "test-group",
		VirtualMachine: "test-vm",
	})
	if err != nil {
		t.Error(err)
	}
	if len(privileges) != 2 {
		t.Errorf("Wrong number of privileges: %d", len(privileges))
	}
}

func TestGrantPrivilege(t *testing.T) {
	done := false
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/users/satan/privileges": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Wrong method %s", r.Method)
				}
				// test the data gets there?
				done = true
			},
		},
	})
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
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/privileges/999": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Wrong method %s", r.Method)
				}
				done = true
			},
		},
	})
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

	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/users/satan/privileges": mkGetPrivilegesHandler(t, "satan"),
			"/privileges/999": func(wr http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Wrong method %s", r.Method)
				}
				done = true
			},
		},
	})
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
