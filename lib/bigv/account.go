package bigv

// Account represents an account object that's returned by the brain
type Account struct {
	Name string `json:"name"`

	// the following cannot be set
	ID        int      `json:"id"`
	Suspended bool     `json:"suspended"`
	Groups    []*Group `json:"groups"`
}
