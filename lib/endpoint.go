package lib

//go:generate stringer -type Endpoint

// Endpoint is an enum-style type to avoid people using endpoints like ints
type Endpoint int

const (
	// AuthEndpoint means "make the connection to auth!"
	AuthEndpoint Endpoint = iota
	// BrainEndpoint means "make the connection to the brain!"
	BrainEndpoint
	// BillingEndpoint means "make the connection to bmbilling!"
	BillingEndpoint
	// SPPEndpoint means "make the connection to SPP!"
	SPPEndpoint
	// APIEndpoint means "make the connection to the general API endpoint!" (api.bytemark.co.uk - atm only used for domains?)
	APIEndpoint
)
