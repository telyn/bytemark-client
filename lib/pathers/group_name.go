package pathers

import (
	"errors"
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
	if g.Account == "" {
		return "", errors.New("no account specified - use GetDefaultAccount to set the Account if none is otherwise specified")
	}
	url := fmt.Sprintf("/accounts/%s/groups/%s", string(g.Account), g.Group)
	return url, nil
}
