package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// HandlerMap is a collection of handlers to use when creating complex tests, or
// implicitly using MultiRequestTestSpec (which is the recommended way to write
// tests)
type HandlerMap map[lib.Endpoint]http.Handler

func (hm HandlerMap) getHandler(t *testing.T, ep lib.Endpoint) http.Handler {
	if h, ok := hm[ep]; ok {
		return closeBodyAfter(h)
	}
	return NilHandler(t)
}
func (hm HandlerMap) mkServer(t *testing.T, ep lib.Endpoint) *httptest.Server {
	if ep == lib.AuthEndpoint {
		if _, ok := hm[lib.AuthEndpoint]; !ok {
			return NewAuthServer()
		}
	}
	return httptest.NewServer(hm.getHandler(t, ep))
}

// MakeServers creates a Servers from this HandlerMap. If an auth handler is
// defined (i.e. if hm[lib.AuthEndpoint] exists) then that is used for auth,
// otherwise NewAuthServer() is used. You usually don't want to set
// hm[lib.AuthEndpoint] unless you are testing something related to
// authentication.
func (hm HandlerMap) MakeServers(t *testing.T) Servers {
	s := Servers{
		auth:    hm.mkServer(t, lib.AuthEndpoint),
		brain:   hm.mkServer(t, lib.BrainEndpoint),
		billing: hm.mkServer(t, lib.BillingEndpoint),
		api:     hm.mkServer(t, lib.APIEndpoint),
		spp:     hm.mkServer(t, lib.SPPEndpoint),
	}
	return s
}
