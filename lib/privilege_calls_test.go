package lib_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
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

func testPrivilegesForUser(user string) (privs brain.Privileges) {
	privs = testPrivileges
	if user != "" {
		privs = brain.Privileges{}
		for _, p := range testPrivileges {
			if p.Username == user {
				privs = append(privs, p)
			}
		}
	}
	return
}

func TestGetPrivileges(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/privileges",
		Response: testPrivileges,
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		// simple test of getting all privileges - hence the empty string
		privileges, err := client.GetPrivileges("")
		if err != nil {
			t.Error(err)
		}
		if len(privileges) != 2 {
			t.Errorf("Wrong number of privileges: %d\nfull list of privs: %#v", len(privileges), privileges)
		}
	})

}

func TestGetPrivilegesForAccount(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/mycoolaccount/privileges",
		Response: testPrivileges,
	}

	rts.Run(t, testName, true, func(client lib.Client) {

		privileges, err := client.GetPrivilegesForAccount("mycoolaccount")
		if err != nil {
			t.Error(err)
		}
		if len(privileges) != 2 {
			t.Errorf("Wrong number of privileges: %d", len(privileges))
		}
	})
}

func TestGetPrivilegesForGroup(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/test-account/groups/test-group/privileges",
		Response: testPrivileges,
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		privileges, err := client.GetPrivilegesForGroup(lib.GroupName{Group: "test-group", Account: "test-account"})
		if err != nil {
			t.Error(err)
		}
		if len(privileges) != 2 {
			t.Errorf("Wrong number of privileges: %d", len(privileges))
		}
	})
}

func TestGetPrivilegesForServer(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/test-account/groups/test-group/virtual_machines/test-vm/privileges",
		Response: testPrivileges,
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		privileges, err := client.GetPrivilegesForVirtualMachine(lib.VirtualMachineName{
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
	})
}

func TestGrantPrivilege(t *testing.T) {
	testName := testutil.Name(0)
	tests := []struct {
		name string

		input    brain.Privilege
		expected map[string]interface{}
	}{
		{
			name: "account admin with password",
			input: brain.Privilege{
				Username:         "satan",
				AccountID:        433224.0,
				PasswordRequired: true,
				Level:            brain.AccountAdminPrivilege,
			},
			expected: map[string]interface{}{
				"account_id":        433224.0,
				"level":             "account_admin",
				"password_required": true,
				"yubikey_required":  false,
			},
		}, {
			name: "yubikey",
			input: brain.Privilege{
				Username:         "azriel",
				PasswordRequired: true,
				YubikeyRequired:  true,
				Level:            brain.ClusterAdminPrivilege,
			},
			expected: map[string]interface{}{
				"level":             "cluster_admin",
				"password_required": true,
				"yubikey_required":  true,
			},
		}, {
			name: "api key",
			input: brain.Privilege{
				APIKeyID: 2013,
				Username: "dr-robotnik",
				GroupID:  2145,
				Level:    brain.GroupAdminPrivilege,
			},
			expected: map[string]interface{}{
				"api_key_id":        2013.0,
				"level":             "group_admin",
				"group_id":          2145.0,
				"password_required": false,
				"yubikey_required":  false,
			},
		},
	}

	for _, test := range tests {
		rts := testutil.RequestTestSpec{
			Method:        "POST",
			Endpoint:      lib.BrainEndpoint,
			URL:           fmt.Sprintf("/users/%s/privileges", test.input.Username),
			AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := client.GrantPrivilege(test.input)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRevokePrivilegeWithID(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "DELETE",
		Endpoint: lib.BrainEndpoint,
		URL:      "/privileges/999",
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.RevokePrivilege(brain.Privilege{
			ID:        999,
			Username:  "satan",
			AccountID: 433224,
			Level:     brain.AccountAdminPrivilege,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

/*
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
*/
