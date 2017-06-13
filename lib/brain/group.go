package brain

import (
	"fmt"
)

// Group represents a group
type Group struct {
	Name string `json:"name"`

	// the following cannot be set
	AccountID       int               `json:"account_id"`
	ID              int               `json:"id"`
	VirtualMachines []*VirtualMachine `json:"virtual_machines"`
}

func (g Group) String() string {
	return fmt.Sprintf("group %d %q - has %d servers", g.ID, g.Name, len(g.VirtualMachines))
}
