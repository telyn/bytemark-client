package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"	
)

// Function to get array of servers on specified head
func GetServersOnHead(client lib.Client, id string) (servers brain.VirtualMachines, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/heads/%s/virtual_machines", id)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

// Function to get array of servers on specified tail
func GetServersOnTail(client lib.Client, id string) (servers brain.VirtualMachines, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/tails/%s/virtual_machines", id)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

// Function to get array of servers on specified storage pool
func GetServersOnStoragePool(client lib.Client, id string) (servers brain.VirtualMachines, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/storage_pools/%s/virtual_machines", id)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}
