package testutil

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// MuxHandlers is like a HandlerMap, but for Mux objects instead of
// http.Handler.
// MuxHandlers is deprecated (replaced by MultiRequestTestSpec), will go away eventually.
type MuxHandlers struct {
	Auth    Mux
	Brain   Mux
	Billing Mux
	SPP     Mux
	API     Mux
}

// MakeServers creates a Servers whose httptest.Server elements are handled by these Muxes
func (mh MuxHandlers) MakeServers(t *testing.T) (s Servers) {
	h := HandlerMap{
		lib.AuthEndpoint:    mh.Auth.ToHandler(),
		lib.BrainEndpoint:   mh.Brain.ToHandler(),
		lib.BillingEndpoint: mh.Billing.ToHandler(),
		lib.SPPEndpoint:     mh.SPP.ToHandler(),
		lib.APIEndpoint:     mh.API.ToHandler(),
	}
	if mh.Auth == nil {
		delete(h, lib.AuthEndpoint)
	}
	return h.MakeServers(t)
}

func closeBodyAfter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(wr, r)
		// ignore the error cause we can assume it was already closed if it
		// errors. Or we just don't care - it's only a test. At least we tried.
		_ = r.Body.Close()
	})
}
