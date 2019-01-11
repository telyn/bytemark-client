package lib

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GroupName is the double-form of the name of a Group, which should be enough to find the group.
type GroupName struct {
	brain.AccountPather
	Group string
}

// DefaultGroup is the default group name (just the group part - don't add dots!). Defaults to "default". Wow.
var DefaultGroup = "default"

func (g GroupName) String() string {
	if g.Group == "" {
		g.Group = DefaultGroup
	}
	if g.Account == "" {
		return g.Group
	}
	return g.Group + "." + g.Account
}

// Path returns the URL path for this group, if possible.
// If the group is not full specified (i.e. does not have an account and group)
// it instead returns an error.
func (g GroupName) GroupPath() (string, error) {
	if g.Group == "" || g.Account == "" {
		return "", fmt.Errorf("Group %q was not fully specified so cannot make a URL", g)
	}
	path := fmt.Sprintf("/accounts/%s/groups/%s", g.Account, g.Group)
	return path, nil
}
