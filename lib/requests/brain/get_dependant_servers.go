package brain

import (
	"time"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// Function to get array of servers on specified head
func GetServersOnHead(client lib.Client, id string, at string) (servers brain.VirtualMachines, err error) {
	var r lib.Request

	if at != "" {
		at, err = ParseDateTime(at)

		if err != nil {
			return
		}

		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/heads/%s/virtual_machines?at=%s", id, at)
	} else {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/heads/%s/virtual_machines", id)
	}

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

// Function to get array of servers on specified tail
func GetServersOnTail(client lib.Client, id string, at string) (servers brain.VirtualMachines, err error) {
	var r lib.Request

	if at != "" {
		at, err = ParseDateTime(at)

		if err != nil {
			return
		}

		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/tails/%s/virtual_machines?at=%s", id, at)
	} else {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/tails/%s/virtual_machines", id)
	}

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

// Function to get array of servers on specified storage pool
func GetServersOnStoragePool(client lib.Client, id string, at string) (servers brain.VirtualMachines, err error) {
	var r lib.Request

	if at != "" {
		at, err = ParseDateTime(at)

		if err != nil {
			return
		}

		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/storage_pools/%s/virtual_machines?at=%s", id, at)
	} else {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/storage_pools/%s/virtual_machines", id)
	}

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

func ParseDateTime(at string) (r string, err error) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	when, err := w.Parse(at, time.Now())

	r = when.Time.Format("2006-01-02T15:04:05-0700")

	return
}
