package lib

import (
	"encoding/json"
	"fmt"
	"github.com/cheekybits/is"
	"io/ioutil"
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

func getFixtureDiscSet() []Disc {
	return []Disc{
		getFixtureDisc(),
		Disc{
			ID:           2,
			StorageGrade: "archive",
			Label:        "arch",
			Size:         1024000,
		},
		Disc{
			ID:           3,
			StorageGrade: "",
			Size:         2048,
		},
	}
}

func TestLabelDisc(t *testing.T) {
	is := is.New(t)
	discs := getFixtureDiscSet()
	labelDiscs(discs)
	for _, d := range discs {
		switch d.ID {
		case 1:
			is.Equal("vda", d.Label)
		case 2:
			is.Equal("arch", d.Label)
		case 3:
			is.Equal("vdc", d.Label)
		default:
			fmt.Printf("Unexpected disc ID %d\r\n", d.ID)
			t.Fail()
		}
	}
}

func TestValidateDisc(t *testing.T) {
	is := is.New(t)
	discs := getFixtureDiscSet()
	for _, d := range discs {
		d2, err := d.Validate()
		is.Nil(err)

		is.Equal(d.ID, d2.ID)
		is.Equal(d.Label, d2.Label)
		is.Equal(d.Size, d2.Size)
		is.Equal(d.StoragePool, d2.StoragePool)
		is.Equal(d.VirtualMachineID, d2.VirtualMachineID)
		switch d.ID {
		case 1, 3:
			is.Equal("sata", d2.StorageGrade)
		case 2:
			is.Equal("archive", d2.StorageGrade)
		}
	}
}

func TestCreateDisc(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs" && req.Method == "POST" {
			// TODO: unmarshal the disc
			// then test for sanity, equality to disk put in
			w.Write([]byte("{}"))
		} else if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm" {
			// TODO: return a VM that has some discs
			w.Write([]byte("{}"))
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), mkNilHandler(t))
	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()

	is.Nil(err)
	err = client.AuthWithCredentials(map[string]string{})
	is.Nil(err)
	if err != nil {
		t.Fatal(err)
	}

	err = client.CreateDisc(&VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, getFixtureDisc())

	is.Nil(err)

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
			if req.URL.Query().Get("purge") != "true" {
				http.NotFound(w, req)
			}

		} else if req.URL.Path == "/accounts/invalid-account" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), mkNilHandler(t))
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

	err = client.DeleteDisc(&VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
	is.Nil(err)

}

func TestResizeDisc(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
			bytes, err := ioutil.ReadAll(req.Body)
			is.Nil(err)
			var disc Disc
			err = json.Unmarshal(bytes, &disc)
			is.Nil(err)

		} else if req.URL.Path == "/accounts/invalid-account" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), mkNilHandler(t))
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

	err = client.ResizeDisc(&VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666", 35)
	is.Nil(err)

}

func TestShowDisc(t *testing.T) {
	is := is.New(t)

	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
			bytes, err := json.Marshal(getFixtureDisc())
			is.Nil(err)
			w.Write(bytes)
		} else if req.URL.Path == "/accounts/invalid-account" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), mkNilHandler(t))
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

	disc, err := client.GetDisc(&VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
	if err != nil {
		t.Fatal(err)
	}
	is.Nil(err)
	fx := getFixtureDisc()

	is.Equal(fx.ID, disc.ID)
	is.Equal(fx.Label, disc.Label)
	is.Equal(fx.Size, disc.Size)
	is.Equal(fx.StorageGrade, disc.StorageGrade)
	is.Equal(fx.StoragePool, disc.StoragePool)

}
