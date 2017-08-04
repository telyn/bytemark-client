package lib

// GroupName is the double-form of the name of a Group, which should be enough to find the group.
type GroupName struct {
	Group   string
	Account string
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
