package lib

import ()

// GroupName is the double-form of the name of a Group, which should be enough to find the group.
type GroupName struct {
	Group   string
	Account string
}

func (g GroupName) String() string {
	if g.Group == "" {
		g.Group = "default"
	}
	if g.Account == "" {
		return g.Group
	}
	return g.Group + "." + g.Account
}
