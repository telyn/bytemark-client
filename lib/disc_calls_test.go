package lib_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/cheekybits/is"
)

func getFixtureDisc() brain.Disc {
	return brain.Disc{
		Label:            "",
		StorageGrade:     "sata",
		Size:             26400,
		ID:               1,
		VirtualMachineID: 1,
		StoragePool:      "fakepool",
	}
}

func getFixtureDiscSet() []brain.Disc {
	return []brain.Disc{
		getFixtureDisc(),
		brain.Disc{
			ID:           2,
			StorageGrade: "archive",
			Label:        "arch",
			Size:         1024000,
		},
		brain.Disc{
			ID:           3,
			StorageGrade: "",
			Size:         2048,
		},
	}
}

func TestCreateDisc(t *testing.T) {
	testName := testutil.Name(0)

	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/group/virtual_machines/vm/discs": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("POST"),
						//TODO assert.BodyUnmarshal
						//TODO to whoever has this MR: if you see this, make me write an ACTUAL test here.
					)
					w.Write([]byte(`{}`))
				},
				"/accounts/account/groups/group/virtual_machines/vm": func(w http.ResponseWriter, r *http.Request) {
					// TODO: request test
					// TODO: meaningful response
					// TODO to whoever has this MR: if you see this, make me write an ACTUAL test here.
					w.Write([]byte(`{}`))
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.CreateDisc(lib.VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, getFixtureDisc())
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
		// TODO response validation
		//TODO to whoever has this MR: if you see this, make me write an ACTUAL test here.
	})
}

func TestDeleteDisc(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:        "DELETE",
		Endpoint:      lib.BrainEndpoint,
		URL:           "/accounts/account/groups/group/virtual_machines/vm/discs/666",
		AssertRequest: assert.QueryValue("purge", "true"),
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.DeleteDisc(lib.VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})

}

func TestResizeDisc(t *testing.T) {
	testName := testutil.Name(0)

	actualDisc := make(map[string]interface{})
	expectedDisc := map[string]interface{}{
		"size": 35.0, // json library treats all numbers as float64
	}
	rts := testutil.RequestTestSpec{
		Method:   "PUT",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/account/groups/group/virtual_machines/vm/discs/666",
		AssertRequest: assert.BodyUnmarshal(&actualDisc, func(_ *testing.T, _ string) {
			if !reflect.DeepEqual(actualDisc, expectedDisc) {
				t.Errorf("Resize disc request wasn't as expected.\r\nExpected: %#v\r\nActual: %#v", expectedDisc, actualDisc)
			}
		}),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.ResizeDisc(lib.VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666", 35)
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})

}

func TestShowDisc(t *testing.T) {
	is := is.New(t)
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/account/groups/group/virtual_machines/vm/discs/666",
		Response: getFixtureDisc(),
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		disc, err := client.GetDisc(lib.VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
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
	})
}

func TestShowDiscByID(t *testing.T) {
	is := is.New(t)
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/discs/667",
		Response: getFixtureDisc(),
	}

	rts.Run(t, testName, true, func(client lib.Client) {

		disc, err := client.GetDiscByID(667)
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
	})
}
