package lib

import (
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func getFixtureDisc() Disc {
	return Disc{
		Label:            "",
		StorageGrade:     "sata",
		Size:             26400,
		ID:               1,
		VirtualMachineID: 1,
		StoragePool:      "fakepool",
	}
}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
			if req.URL.Query().Get("purge") != "true" {
				http.NotFound(w, req)
			}

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

	err = client.DeleteDisc(VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, 666)
	is.Nil(err)

}
