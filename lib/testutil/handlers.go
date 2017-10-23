package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// ok this shit is a bit weird / crap

// basically, Handlers is a collection of http.Handlers, one per endpoint.
// downside is you need a complete http.Handler for each one which means that if your test touches both GET /virtual_machines/23 and POST /virtual_machines/23/discs you have to write that all as one http.HandlerFunc.

// what would be much nicer is to use http.ServeMuxes in a short-hand fashion.
// that's what Mux and MuxHandlers are for. The Mux.ToHandler into a full-on http.ServeMux ready to add to a Handlers
// MuxHandlers is used to make a collection of Muxes in the same way that Handlers is for a collection of http.Handlers
// and both these things can be turned into a Servers ready to test with using MakeServers.

// Mux is a map of URL paths to http.HandlerFuncs
type Mux map[string]http.HandlerFunc

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
	Auth    Mux
	Brain   Mux
	Billing Mux
	SPP     Mux
	API     Mux
}

// MakeServers creates a Servers whose httptest.Server elements are handled by these Muxes
func (mh MuxHandlers) MakeServers(t *testing.T) (s Servers) {
	h := Handlers{
		auth:    mh.Auth.ToHandler(),
		brain:   mh.Brain.ToHandler(),
		billing: mh.Billing.ToHandler(),
		spp:     mh.SPP.ToHandler(),
		api:     mh.API.ToHandler(),
	}
	if mh.Auth == nil {
		h.auth = nil
	}
	return h.MakeServers(t)
}

// AddMux adds a Mux for the endpoint passed. Returns an error if it didn't recognise that endpoint
func (mh *MuxHandlers) AddMux(ep lib.Endpoint, m Mux) (err error) {
	switch ep {
	case lib.AuthEndpoint:
		mh.Auth = m
	case lib.BrainEndpoint:
		mh.Brain = m
	case lib.BillingEndpoint:
		mh.Billing = m
	case lib.SPPEndpoint:
		mh.SPP = m
	case lib.APIEndpoint:
		mh.API = m
	default:
		return fmt.Errorf("'%d' is not a known endpoint const. Take another look at lib/client.go's Endpoint type", ep)
	}
	return nil
}

func closeBodyAfter(h http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(wr, r)
		_ = r.Body.Close()
	}
}

// NewMuxHandlers creates a MuxHandler which will respond on the given endpoint URL with the handler provided, after which the request body will be automatically closed.
func NewMuxHandlers(endpoint lib.Endpoint, url string, h http.HandlerFunc) (mh MuxHandlers, err error) {
	err = mh.AddMux(endpoint, Mux{
		url: closeBodyAfter(h),
	})
	return
}

// Handlers is a struct which holds a http.Handler for each endpoint the client
// could possibly talk to. It gets converted to a Servers by MakeServers, which
// is the structure that's used in tests.
// Due to the existence of RequestTestSpec, most if not all the tests in lib
// and lib/requests should use RequestTestSpec instead of directly getting a Servers
type Handlers struct {
	auth    http.Handler
	brain   http.Handler
	billing http.Handler
	spp     http.Handler
	api     http.Handler
}

// Fill ensures each endpoint has an http.Handler - except auth.
// auth is allowed to be nil.
func (h *Handlers) Fill(t *testing.T) {
	if h.brain == nil {
		h.brain = NilHandler(t)
	}
	if h.billing == nil {
		h.billing = NilHandler(t)
	}
	if h.spp == nil {
		h.spp = NilHandler(t)
	}
	if h.api == nil {
		h.api = NilHandler(t)
	}
}

// MakeServers creates a Servers object from this Handlers.
// if auth is nil, creates the default auth server with NewAuthServer.
func (h Handlers) MakeServers(t *testing.T) (s Servers) {
	h.Fill(t)

	if h.auth != nil {
		s.auth = httptest.NewServer(h.auth)
	} else {
		s.auth = NewAuthServer()
	}
	s.brain = httptest.NewServer(h.brain)
	s.billing = httptest.NewServer(h.billing)
	s.api = httptest.NewServer(h.api)
	s.spp = httptest.NewServer(h.spp)

	return
}
