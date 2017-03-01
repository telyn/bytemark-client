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

func TestGetIPRanges(t *testing.T) {
	testIPRanges := []*brain.IPRange{
		{
			ID:      1234,
			Spec:    "192.168.13.0/24",
			VLANNum: 123,
			Zones: []string{
				"test-zone",
			},
			Available: 200.0,
		},
	}
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/admin/ip_ranges": func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "GET")
				writeJSON(t, wr, testIPRanges)
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

	ipranges, err := client.GetIPRanges()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ipranges, testIPRanges) {
		t.Errorf("IPRanges returned from GetIPRanges were not what was expected.\r\nExpected: %#v\r\nActual:%#v", testIPRanges, ipranges)
	}
}

func TestGetIPRange(t *testing.T) {
	testIPRange := brain.IPRange{
		ID:      1234,
		Spec:    "192.168.13.0/24",
		VLANNum: 123,
		Zones: []string{
			"test-zone",
		},
		Available: 200.0,
	}
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/admin/ip_ranges/1234": func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "GET")
				writeJSON(t, wr, testIPRange)
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

	iprange, err := client.GetIPRange(1234)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(iprange, &testIPRange) {
		t.Errorf("IPRange returned from GetIPRange was not what was expected.\r\nExpected: %#v\r\nActual:%#v", testIPRange, iprange)
	}
}
func TestGetHeads(t *testing.T) {
	testHeads := []*brain.Head{
		{
			ID:       315,
			UUID:     "234833-2493-3423-324235",
			Label:    "test-head315",
			ZoneName: "awesomecoolguyzone",

			Architecture: "x86_64",
			// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
			// CCAddress:     &net.IP{214, 233, 32, 31},
			Note:          "melons",
			Memory:        241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			MemoryFree:          123400,
			IsOnline:            true,
			UsedCores:           9,
			VirtualMachineCount: 3,
		}, {
			ID:       239,
			UUID:     "235670-2493-3423-324235",
			Label:    "test-head239",
			ZoneName: "awesomecoolguyzone",

			Architecture: "x86_64",
			// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
			// CCAddress:     &net.IP{24, 43, 32, 49},
			Note:          "more than a hundred years old",
			Memory:        241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			MemoryFree:          234000,
			IsOnline:            true,
			UsedCores:           1,
			VirtualMachineCount: 1,
		},
	}
	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/admin/heads": func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "GET")
				writeJSON(t, wr, testHeads)
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

	heads, err := client.GetHeads()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(heads, testHeads) {
		t.Errorf("Heads returned from GetHeads was not what was expected.\r\nExpected: %#v\r\nActual:   %#v", testHeads, heads)
	}
}
