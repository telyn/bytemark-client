package brain

// AccountPather is a type which can provide a URL path for an account
type AccountPather interface {
	// AccountPath returns the path component of a URL, starting with /, which
	// uniquely references an account with the brain, or an error if it cannot
	AccountPath() (string, error)
}

// DiscPather is a type which can provide a URL path for a disc
type DiscPather interface {
	// DiscPath returns the path component of a URL, starting with /, which
	// uniquely references the disc with the brain, or an error if it cannot
	DiscPath() (string, error)
}

// GroupPather is a type which can provide a URL path to a group
type GroupPather interface {
	// GroupPath returns the path component of a URL, starting with /, which
	// uniquely references the group with the brain, or an error if it cannot.
	GroupPath() (string, error)
}

// VirtualMachinePather is a type which can provide a URL path to a virtual machine
type VirtualMachinePather interface {
	// VirtualMachinePath returns the path component of a URL, starting with /,
	// which uniquely references the virtual machine with the brain, or an
	// error if it cannot.
	VirtualMachinePath() (string, error)
}
