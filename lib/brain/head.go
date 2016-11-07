package brain

import (
	"net"
)

// Head represents a Bytemark Cloud Servers head server.
type Head struct {
	ID       int    `json:"id,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Label    string `json:"label,omitempty"`
	ZoneName string `json:"zone,omit_empty"`

	// descriptive information that can be modified

	Architecture  string   `json:"architecture"`
	CCAddress     *net.IP  `json:"c_and_c_address"`
	Note          string   `json:"note"`
	Memory        int      `json:"memory,omitempty"`
	UsageStrategy string   `json:"usage_strategy,omitempty"`
	Models        []string `json:"models,omitempty"`

	// state

	MemoryFree int  `json:"memory_free,omitempty"`
	IsOnline   bool `json:"is_online,omitempty"`
	UsedCores  int  `json:"used_cores"`

	// You may have one or the other.

	VirtualMachineCount int               `json:"virtual_machines_count,omitempty"`
	VirtualMachines     []*VirtualMachine `json:"virtual_machines,omitempty"`
}

// CountVirtualMachines returns the number of virtual machines running on this head
func (h *Head) CountVirtualMachines() int {
	if h.VirtualMachines != nil {
		return len(h.VirtualMachines)
	}
	return h.VirtualMachineCount
}
