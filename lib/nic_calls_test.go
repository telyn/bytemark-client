package lib_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func getFixtureNic() brain.NetworkInterface {
	ip := net.IPv4(127, 0, 0, 2)
	return brain.NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []net.IP{ip},
		ExtraIPs:         map[string]net.IP{},
		VirtualMachineID: 1,
	}
}

func TestAddIP(t *testing.T) {
	local1 := net.IPv4(127, 0, 0, 1)
	local2 := net.IPv4(127, 0, 0, 2)
	tests := []struct {
		name       string
		serverName lib.VirtualMachineName
		nicID      int
		spec       brain.IPCreateRequest
		created    brain.IPCreateRequest
		shouldErr  bool
	}{
		{
			name:       "add one ip",
			serverName: lib.VirtualMachineName{Account: "test", Group: "testo", VirtualMachine: "testing"},
			nicID:      252,
			spec:       brain.IPCreateRequest{Addresses: 1, Family: "ipv4", Reason: "jeff", Contiguous: false},
			created:    brain.IPCreateRequest{IPs: brain.IPs{local1}},
		},
		{
			name:       "add two ips",
			serverName: lib.VirtualMachineName{Account: "borm", Group: "galp", VirtualMachine: "sklep"},
			nicID:      564,
			spec:       brain.IPCreateRequest{Addresses: 2, Family: "ipv4", Reason: "jeff", Contiguous: false},
			created:    brain.IPCreateRequest{IPs: brain.IPs{local1, local2}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			specMap := map[string]interface{}{
				"addresses":  float64(test.spec.Addresses),
				"family":     test.spec.Family,
				"reason":     test.spec.Reason,
				"contiguous": test.spec.Contiguous,
			}
			vmURL := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", test.serverName.Account, test.serverName.Group, test.serverName.VirtualMachine)
			ipcreateURL := vmURL + fmt.Sprintf("/nics/%d/ip_create", test.nicID)
			vm := brain.VirtualMachine{
				NetworkInterfaces: []brain.NetworkInterface{
					{
						ID: test.nicID,
					},
				},
			}

			rts := testutil.RequestTestSpec{
				MuxHandlers: &testutil.MuxHandlers{
					Brain: testutil.Mux{
						vmURL: func(wr http.ResponseWriter, r *http.Request) {
							assert.All(
								assert.Auth(lib.TokenType(lib.BrainEndpoint)),
								assert.Method("GET"),
							)(t, test.name, r)

							testutil.WriteJSON(t, wr, vm)
						},
						ipcreateURL: func(wr http.ResponseWriter, r *http.Request) {
							assert.All(
								assert.Auth(lib.TokenType(lib.BrainEndpoint)),
								assert.Method("POST"),
								assert.BodyUnmarshalEqual(specMap),
							)(t, test.name, r)

							testutil.WriteJSON(t, wr, test.created)
						},
					},
				},
			}

			rts.Run(t, test.name, true, func(client lib.Client) {
				ips, err := client.AddIP(test.serverName, test.spec)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Errorf("Error expected but not returned")
				}

				assert.Equal(t, test.name, test.created.IPs, ips)
			})
		})
	}
}
