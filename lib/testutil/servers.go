package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// Servers holds an httptest.Server for each endpoint that the client could talk to
type Servers struct {
	auth    *httptest.Server
	brain   *httptest.Server
	billing *httptest.Server
	spp     *httptest.Server
	api     *httptest.Server
}

// Close ensures all the servers have been closed.
func (s *Servers) Close() {
	if s.auth != nil {
		s.auth.Close()
	}
	if s.brain != nil {
		s.brain.Close()
	}
	if s.billing != nil {
		s.billing.Close()
	}
	if s.spp != nil {
		s.spp.Close()
	}
	if s.api != nil {
		s.api.Close()
	}
	s.auth = nil
	s.brain = nil
	s.billing = nil
	s.spp = nil
	s.api = nil
}

// URLs creates an EndpointURLs filled with all the URLs for the Servers
func (s Servers) URLs() (urls lib.EndpointURLs) {
	urls.API = s.api.URL
	urls.Auth = s.auth.URL
	urls.Billing = s.billing.URL
	urls.Brain = s.brain.URL
	urls.SPP = s.spp.URL
	return
}

// Client makes a bytemarkClient for these Servers
func (s Servers) Client() (c lib.Client, err error) {
	urls := s.URLs()
	c, err = lib.NewWithURLs(urls)
	if err != nil {
		return nil, err
	}
	c.AllowInsecureRequests()
	return
}

// ServersFactory is an interface used for convenience - so that Handlers or MuxHandlers can be passed to NewClientAndServers
type ServersFactory interface {
	MakeServers(t *testing.T) Servers
}

// NewClientAndServers constructs httptest Servers for a pretend auth and API endpoint, then constructs a Client that uses those servers.
// Used to test that the right URLs are being requested and such.
// because this is used in nearly all of the tests in lib, this also does some weird magic to set up a writer for log such that all the test output comes out attached to the test it's from
func NewClientAndServers(t *testing.T, factory ServersFactory) (c lib.Client, s Servers, err error) {
	log.Writer = TestLogWriter{t}
	log.ErrWriter = TestLogWriter{t}
	s = factory.MakeServers(t)
	c, err = s.Client()
	return
}

// NewAuthServer creates a fake auth server that responds to any request with a session.
// just for convenience when writing tests that require auth.
func NewAuthServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w,
			`{
    "token": "working-auth-token",
    "username": "account",
    "factors": []
}`)
	}))

}
