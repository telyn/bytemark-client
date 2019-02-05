package brain_test

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestMoveDisc(t *testing.T) {
	testName := testutil.Name(0)

	disc := brain.Disc{}
	expectedDisc := brain.Disc{
		VirtualMachineID: 999,
	}

	oldVM := brain.VirtualMachine{
		ID: 999,
	}
	oldVMName := pathers.VirtualMachineName{
		VirtualMachine: "vm",
		GroupName: pathers.GroupName{
			Group:   "group",
			Account: "account",
		},
	}
	newVMName := pathers.VirtualMachineName{
		VirtualMachine: "new-vm",
		GroupName: pathers.GroupName{
			Group:   "group",
			Account: "account",
		},
	}

	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/group/virtual_machines/vm/discs/666": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("PUT"),
						assert.BodyUnmarshal(&disc, func(_ *testing.T, _ string) {
							assert.Equal(t, testName, expectedDisc, disc)
						}),
					)(t, testName, r)
					testutil.WriteJSON(t, w, map[string]interface{}{})
				},
				"/accounts/account/groups/group/virtual_machines/new-vm": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, oldVM)
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		err := brainRequests.MoveDisc(client, oldVMName, "666", newVMName)
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})

}
