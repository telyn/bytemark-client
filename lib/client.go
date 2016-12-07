package lib

import (
	"errors"
	auth3 "github.com/BytemarkHosting/auth-client"
)

// EndpointURLs are the URLs stored by the client for the various API endpoints the client touches.
// The key endpoints that you may wish to alter are Auth and Brain. When using an auth server and brain
// that doesn't have a matching bmbilling API, Billing should be set to ""
type EndpointURLs struct {
	API     string
	Auth    string
	Billing string
	Brain   string
	SPP     string
}

// DefaultURLs returns an EndpointURLs for the usual customer-facing Bytemark APIs.
func DefaultURLs() EndpointURLs {
	return EndpointURLs{
		API:     "https://api.bytemark.co.uk",
		Auth:    "https://auth.bytemark.co.uk",
		Billing: "https://bmbilling.bytemark.co.uk",
		Brain:   "https://uk0.bigv.io",
		SPP:     "https://spp-submissions.bytemark.co.uk",
	}
}

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

// bytemarkClient is the main type in the Bytemark API client library
type bytemarkClient struct {
	allowInsecure bool
	auth          *auth3.Client
	authSession   *auth3.SessionData
	debugLevel    int
	urls          EndpointURLs
}

// NewSimple creates a new Bytemark API client using the default bytemark endpoints.
// This function will be renamed to New in 3.0
func NewSimple() (Client, error) {
	return NewWithURLs(DefaultURLs())
}

// NewWithURLs creates a new Bytemark API client using the given endpoints.
func NewWithURLs(urls EndpointURLs) (c Client, err error) {
	if urls.Auth == "" {
		urls.Auth = "https://auth.bytemark.co.uk"
	}
	auth, err := auth3.New(urls.Auth)
	if err != nil {
		return nil, err
	}
	client := bytemarkClient{
		urls:       urls,
		auth:       auth,
		debugLevel: 0,
	}
	return &client, nil
}

// New creates a new Bytemark API client using the given Bytemark API endpoint and the default Bytemark auth endpoint, and fills the rest in with defaults.
// This function will be replaced with NewSimple in 3.0
func New(brainEndpoint, billingEndpoint, sppEndpoint string) (c Client, err error) {
	auth, err := auth3.New("https://auth.bytemark.co.uk")
	if err != nil {
		return nil, err
	}
	return NewWithAuth(brainEndpoint, billingEndpoint, sppEndpoint, auth), nil
}

// NewWithAuth creates a new Bytemark API client using the given Bytemark API endpoint and github.com/BytemarkHosting/auth-client Client
// This function is deprecated and will be removed in 3.0
func NewWithAuth(brainEndpoint, billingEndpoint, sppEndpoint string, auth *auth3.Client) Client {
	urls := DefaultURLs()
	urls.Brain = brainEndpoint
	urls.Billing = billingEndpoint
	urls.SPP = sppEndpoint
	client := bytemarkClient{
		urls:       urls,
		auth:       auth,
		debugLevel: 0,
	}
	return &client
}

// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
func (c *bytemarkClient) AuthWithCredentials(credentials auth3.Credentials) error {
	session, err := c.auth.CreateSession(credentials)
	if err == nil {
		c.authSession = session
	}
	return err
}

// AuthWithToken attempts to read sessiondata from auth for the given token. Returns nil on success or an error otherwise.
func (c *bytemarkClient) AuthWithToken(token string) error {
	if token == "" {
		return errors.New("No token provided")
	}

	session, err := c.auth.ReadSession(token)
	if err == nil {
		c.authSession = session
	}
	return err

}

// GetEndpoint returns the Bytemark API endpoint currently in use.
func (c *bytemarkClient) GetEndpoint() string {
	return c.urls.Brain
}

// GetBillingEndpoint returns the Bytemark Billing API endpoint in use.
// This function is deprecated and will be removed in a point release.
// DO NOT DEPEND ON IT
func (c *bytemarkClient) GetBillingEndpoint() string {
	return c.urls.Billing
}

// SetDebugLevel sets the debug level / verbosity of the Bytemark API client. 0 (default) is silent.
func (c *bytemarkClient) SetDebugLevel(debugLevel int) {
	c.debugLevel = debugLevel
}

// GetSessionFactors returns the factors provided when the current auth session was set up
func (c *bytemarkClient) GetSessionFactors() []string {
	if c.authSession == nil {
		return []string{}
	}
	return c.authSession.Factors
}

// GetSessionToken returns the token for the current auth session
func (c *bytemarkClient) GetSessionToken() string {
	if c.authSession == nil {
		return ""
	}
	return c.authSession.Token
}

func (c *bytemarkClient) GetSessionUser() string {
	if c.authSession == nil {
		return ""
	}
	return c.authSession.Username
}

func (c *bytemarkClient) AllowInsecureRequests() {
	c.allowInsecure = true
}

func (c *bytemarkClient) validateVirtualMachineName(vm *VirtualMachineName) error {
	if vm.Account == "" {
		if err := c.validateAccountName(&vm.Account); err != nil {
			return err
		}
	}
	if vm.Group == "" {
		vm.Group = DefaultGroup
	}

	if vm.VirtualMachine == "" {
		return BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return nil
}

func (c *bytemarkClient) validateGroupName(group *GroupName) error {
	if group.Account == "" {
		if err := c.validateAccountName(&group.Account); err != nil {
			return err
		}
	}
	if group.Group == "" {
		group.Group = DefaultGroup
	}
	return nil
}

func (c *bytemarkClient) validateAccountName(account *string) error {
	if *account == "" && c.authSession != nil {
		billAcc, err := c.getDefaultBillingAccount()
		if err == nil && billAcc != nil {
			*account = billAcc.Name
		} else {
			return err
		}
	}
	if *account == "" {
		return NoDefaultAccountError{}
	}
	return nil
}
