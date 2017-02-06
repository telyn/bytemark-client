package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ok this shit is a bit weird / crap

// basically, Handlers is a collection of http.Handlers, one per endpoint.
// downside is you need a complete http.Handler for each one which means that if your test touches both GET /virtual_machines/23 and POST /virtual_machines/23/discs you have to write that all as one http.HandlerFunc.

// what would be much nicer is to use http.ServeMuxes in a short-hand fashion.
// that's what Mux and MuxHandlers are for. The Mux.ToHandler into a full-on http.ServeMux ready to add to a Handlers
// MuxHandlers is used to make a collection of Muxes in the same way that Handlers is for a collection of http.Handlers
// and both these things can be turned into a Servers ready to test with using MakeServers.

// Mux is a map of URL paths to http.HandlerFuncs
type Mux map[string]func(wr http.ResponseWriter, r *http.Request)

// ToHandler turns the Mux into an http.ServeMux
func (m Mux) ToHandler() (serveMux *http.ServeMux) {

	serveMux = http.NewServeMux()
	for p, f := range m {
		serveMux.HandleFunc(p, f)
	}
	return
}

// MuxHandlers is the equivalent of Handlers, but for Mux objects instead of http.Handler.
type MuxHandlers struct {
	auth    Mux
	brain   Mux
	billing Mux
	spp     Mux
	api     Mux
}

func (mh MuxHandlers) MakeServers(t *testing.T) (s Servers) {
	h := Handlers{
		auth:    mh.auth.ToHandler(),
		brain:   mh.brain.ToHandler(),
		billing: mh.billing.ToHandler(),
		spp:     mh.spp.ToHandler(),
		api:     mh.api.ToHandler(),
	}
	if mh.auth == nil {
		h.auth = nil
	}
	return h.MakeServers(t)
}

type Handlers struct {
	auth    http.Handler
	brain   http.Handler
	billing http.Handler
	spp     http.Handler
	api     http.Handler
}

func (h *Handlers) Fill(t *testing.T) {
	if h.brain == nil {
		h.brain = mkNilHandler(t)
	}
	if h.billing == nil {
		h.billing = mkNilHandler(t)
	}
	if h.spp == nil {
		h.spp = mkNilHandler(t)
	}
	if h.api == nil {
		h.api = mkNilHandler(t)
	}
}

func (h Handlers) MakeServers(t *testing.T) (s Servers) {
	h.Fill(t)

	if h.auth != nil {
		s.auth = mkTestServer(h.auth)
	} else {
		s.auth = mkTestAuthServer()
	}
	s.brain = mkTestServer(h.brain)
	s.billing = mkTestServer(h.billing)
	s.api = mkTestServer(h.api)
	s.spp = mkTestServer(h.spp)

	return
}

type Servers struct {
	auth    *httptest.Server
	brain   *httptest.Server
	billing *httptest.Server
	spp     *httptest.Server
	api     *httptest.Server
}

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

// URLs creates an EndpointURLs
func (s Servers) URLs() (urls EndpointURLs) {
	urls.API = s.api.URL
	urls.Auth = s.auth.URL
	urls.Billing = s.billing.URL
	urls.Brain = s.brain.URL
	urls.SPP = s.spp.URL
	return
}

// Client makes a bytemarkClient for these Servers
func (s Servers) Client() (c Client, err error) {
	urls := s.URLs()
	c, err = NewWithURLs(urls)
	if err != nil {
		return nil, err
	}
	c.AllowInsecureRequests()
	return
}

type ServersFactory interface {
	MakeServers(t *testing.T) (s Servers)
}

// mkTestClientAndServers constructs httptest Servers for a pretend auth and API endpoint, then constructs a Client that uses those servers.
// Used to test that the right URLs are being requested and such.
func mkTestClientAndServers(t *testing.T, factory ServersFactory) (c Client, s Servers, err error) {
	s = factory.MakeServers(t)
	c, err = s.Client()
	return
}

// mkTestServer creates an httptest.Server for the given http.Handler. It's basically an alias for httptest.NewServer. Why did I write it?
func mkTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func mkNilHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Unexpected request to a nil server\r\n%s %s", r.Method, r.URL.String())
	})
}

func mkTestAuthServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`{
    "token": "working-auth-token",
    "username": "account",
    "factors": []
}`)
	}))

}
