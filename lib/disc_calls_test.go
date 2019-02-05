package lib_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
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

func TestCreateDisc(t *testing.T) {
	testName := testutil.Name(0)

	server := brain.VirtualMachine{
		Discs: []brain.Disc{
			{
				Label: "disc-3",
			},
		},
	}
	expected := brain.Disc{
		Label:        "disc-4",
		StorageGrade: "sata",
		Size:         26400,
	}
	disc := brain.Disc{}
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/group/virtual_machines/vm/discs": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("POST"),
						assert.BodyUnmarshal(&disc, func(_ *testing.T, _ string) {
							assert.Equal(t, testName, expected, disc)
						}),
					)(t, testName, r)
					testutil.WriteJSON(t, w, map[string]interface{}{})
				},
				"/accounts/account/groups/group/virtual_machines/vm": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, server)
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.CreateDisc(pathers.VirtualMachineName{VirtualMachine: "vm", GroupName: pathers.GroupName{Group: "group", Account: "account"}}, brain.Disc{
			StorageGrade: "sata",
			Size:         26400,
		})
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
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
		err := client.DeleteDisc(pathers.VirtualMachineName{VirtualMachine: "vm", GroupName: pathers.GroupName{Group: "group", Account: "account"}}, "666")
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
		err := client.ResizeDisc(pathers.VirtualMachineName{VirtualMachine: "vm", GroupName: pathers.GroupName{Group: "group", Account: "account"}}, "666", 35)
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})

}

func TestGetDisc(t *testing.T) {
	is := is.New(t)
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/discs/666",
		Response: getFixtureDisc(),
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		disc, err := client.GetDisc(pathers.VirtualMachineName{VirtualMachine: "vm", GroupName: pathers.GroupName{Group: "group", Account: "account"}}, "666")
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

func TestGetDiscByID(t *testing.T) {
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
