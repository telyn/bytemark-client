package lib

import (
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

// New creates a new Bytemark API client using the default bytemark endpoints.
// This function will be renamed to New in 3.0
func New() (Client, error) {
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

// GetEndpoint returns the Bytemark API endpoint currently in use.
func (c *bytemarkClient) GetEndpoint() string {
	return c.urls.Brain
}

// GetBillingEndpoint returns the Bytemark Billing API endpoint in use.
// This function is deprecated and will be removed in a point release.
// DO NOT DEPEND ON IT
// TODO(telyn): remove this
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
