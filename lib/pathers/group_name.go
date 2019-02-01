package pathers

import (
	"fmt"
)

// GroupName is the double-form of the name of a Group, which should be enough to find the group.
type GroupName struct {
	Group   string
	Account AccountName
}

// DefaultGroup is the default group name (just the group part - don't add dots!). Defaults to "default". Wow.
var DefaultGroup = "default"

func (g GroupName) defaultIfNeeded() {
	if g.Group == "" {
		g.Group = DefaultGroup
	}
}

func (g GroupName) String() string {
	g.defaultIfNeeded()
	if g.Account == "" {
		return g.Group
	}
	return g.Group + "." + string(g.Account)
}

// GroupPath returns a Brain URL for this group, or an error if the group is
// invalid
func (g GroupName) GroupPath() (string, error) {
	g.defaultIfNeeded()
	base, err := g.AccountPath()
	return base + fmt.Sprintf("/groups/%s", g.Group), err
}

// AccountPath returns a Brain URL for the account specified in this GroupName,
// or an error if it is blank
func (g GroupName) AccountPath() (string, error) {
	return g.Account.AccountPath()
}
