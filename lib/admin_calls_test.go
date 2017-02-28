package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"net/http"
	"reflect"
	"testing"
)

func TestGetVLANS(t *testing.T) {
	testVLANs := []*brain.VLAN{
		{
			ID:        90210,
			Num:       123,
			UsageType: "recipes",
			IPRanges: []*brain.IPRange{
				{
					ID:      1234,
					Spec:    "192.168.13.0/24",
					VLANNum: 123,
					Zones: []string{
						"test-zone",
					},
					Available: 200.0,
				},
			},
		},
	}
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/admin/vlans": func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "GET")
				writeJSON(t, wr, testVLANs)
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

	vlans, err := client.GetVLANs()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vlans, testVLANs) {
		t.Errorf("VLANs returned from GetVLANs were not what was expected.\r\nExpected: %#v\r\nActual:%#v", testVLANs, vlans)
	}
}
