package brain_test

import (
	"encoding/json"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetServersOnHead(t *testing.T) {
	testName := testutil.Name(0)

	testServers := brain.VirtualMachines{{
		Cores:					4,
		Memory:					4,
		Name:					"test_server",
		ID:                		123,
		Hostname:				"test_hostname",
		Head:					"test_head",
	}}

	rts := testutil.RequestTestSpec{
		Method:        "GET",
		URL:           "/admin/heads/123/virtual_machines",
		Endpoint:      lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"cores": 4,
		"memory": 4,
		"name": "test_server",
		"id": 123,
		"hostname": "test_hostname",
		"head": "test_head"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		servers, err := brainMethods.GetServersOnHead(client, "123", "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, servers, testServers)
	})
}

func TestGetServersOnTail(t *testing.T) {
	testName := testutil.Name(0)

	testServers := brain.VirtualMachines{{
		Cores:		4,
		Memory:		4,
		Name:		"test_server",
		ID:			123,
		Hostname:	"test_hostname",
		Head:		"test_head",
	}}

	rts := testutil.RequestTestSpec{
		Method:        "GET",
		URL:           "/admin/tails/123/virtual_machines",
		Endpoint:      lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"cores": 4,
		"memory": 4,
		"name": "test_server",
		"id": 123,
		"hostname": "test_hostname",
		"head": "test_head"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		servers, err := brainMethods.GetServersOnTail(client, "123", "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, servers, testServers)
	})
}

func TestGetServersOnStoragePool(t *testing.T) {
	testName := testutil.Name(0)

	testServers := brain.VirtualMachines{{
		Cores:		4,
		Memory:		4,
		Name:		"test_server",
		ID:			123,
		Hostname:	"test_hostname",
		Head:		"test_head",
	}}

	rts := testutil.RequestTestSpec{
		Method:        "GET",
		URL:           "/admin/storage_pools/123/virtual_machines",
		Endpoint:      lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"cores": 4,
		"memory": 4,
		"name": "test_server",
		"id": 123,
		"hostname": "test_hostname",
		"head": "test_head"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		servers, err := brainMethods.GetServersOnStoragePool(client, "123", "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, servers, testServers)
	})
}

