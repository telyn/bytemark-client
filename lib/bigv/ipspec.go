package bigv

// IPSpec represents one v4 and one v6 address to assign to a server during creation.
type IPSpec struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}
