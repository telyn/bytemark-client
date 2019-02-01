package lib_test

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func getFixtureGroup() brain.Group {
	vm := getFixtureVM()
	return brain.Group{
		ID:   1,
		Name: "default",
		VirtualMachines: []brain.VirtualMachine{
			vm,
		},
	}

}

func TestCreateGroup(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "POST",
		URL:      "/accounts/account/groups",
		Endpoint: lib.BrainEndpoint,
		Response: getFixtureGroup(),
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.CreateGroup(pathers.GroupName{Group: "invalid-group", Account: "account"})
		if err != nil {
			t.Errorf("%s err %s", testName, err)
		}
	})
}

func TestDeleteGroup(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/default": func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case "DELETE":
						_, err := w.Write([]byte(""))
						if err != nil {
							t.Fatal(err)
						}
					case "GET":
						testutil.WriteJSON(t, w, getFixtureGroup())
					default:
						t.Errorf("%s wrong method %s", testName, r.Method)
					}
				},
			},
		},
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		err := client.DeleteGroup(pathers.GroupName{Group: "default", Account: "account"})
		if err != nil {
			t.Errorf("%s: %s", testName, err)
		}
	})

}

func TestGetGroup(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/default": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, getFixtureGroup())
				},
			},
			Billing: testutil.Mux{
				"/api/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, []billing.Account{{ID: 233, Name: "account"}})
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		t.Logf("Trying invalid-group")

		group, err := client.GetGroup(pathers.GroupName{Group: "invalid-group", Account: "account"})
		assert.NotEqual(t, testName, nil, err)
		assert.Equal(t, testName, "", group.Name)

		t.Logf("Trying default.account")
		group, err = client.GetGroup(pathers.GroupName{Group: "default", Account: "account"})
		assert.Equal(t, testName, nil, err)
		assert.Equal(t, testName, "default", group.Name)
		assert.Equal(t, testName, 1, len(group.VirtualMachines))

		t.Logf("Trying blank group")

		group, err = client.GetGroup(pathers.GroupName{Group: "", Account: ""})
		assert.Equal(t, testName, nil, err)
		assert.Equal(t, testName, "default", group.Name)
		assert.Equal(t, testName, 1, len(group.VirtualMachines))
	})

}
