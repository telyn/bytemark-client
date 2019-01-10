package testutil

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

// RequestTestFunc is a function which will be run by RequestTestSpec.Run.
// This function should generally call some request method (such as those built-in
// to the lib.Client interface or external ones, usually in the lib/requests/* packages.
//
// Any validation of the results of that request method should be done as part of the
// RequestTestFunc
type RequestTestFunc func(lib.Client)

// RequestTestSpec is used to build up httptest servers for a test, then run a
// function and check it makes the correct request.
type RequestTestSpec struct {
	// MuxHandlers will be used if defined - this allows for the test to support
	// multiple endpoints, URLs, methods, etc. while still keeping as DRY as
	// possible. Otherwise, set the Method, Endpoint, URL, AssertRequest,
	// Response and StatusCode.
	MuxHandlers *MuxHandlers

	// Method is used to assert that the request was given the correct type
	// it is only used if MuxHandlers is nil
	Method string
	// Endpoint is used to build a MuxHandlers from when MuxHandlers is nil
	Endpoint lib.Endpoint
	// URL is used to build a MuxHandlers when MuxHandlers is nil
	URL string
	// Response will be JSON marshalled
	// to use a raw string (i.e. if you don't want to use JSON) cast it to
	// encoding/json.RawMessage - this will be reproduced verbatim
	Response interface{}
	// StatusCode is the status code that will be returned. If unset, will
	// default to whatever http.ResponseWriter.Write defaults to.
	// Only used if MuxHandlers is nil
	StatusCode int
	// AssertRequest is an optional func which will get called to check the
	// request object further - for example to check the URL has particular
	// query string keys
	AssertRequest func(t *testing.T, testName string, r *http.Request)

	// Auth is used to determine whether authentication should be performed - if
	// it is set to true it will be. Auth is automatically set by Run and only
	// used
	Auth bool
	// NoVerify is used to disable visit count verification in case the test is
	// expected not to get as far as calling out to this API call.
	NoVerify bool
	// visits is how many times the generated handler has been called
	visits int
}

// Run sets up fake servers, creates a client that talks to them, then passes the client to fn.
// fn should run some request method using the client & test the results of that function.
func (rts *RequestTestSpec) Run(t *testing.T, testName string, auth bool, fn RequestTestFunc) {
	client, servers, err := NewClientAndServers(t, rts)
	defer servers.Close()
	if err != nil {
		t.Fatalf("%s NewClientAndServers failed - %s", testName, err)
	}
	if auth {
		err = client.AuthWithCredentials(map[string]string{})
		if err != nil {
			t.Fatalf("%s AuthWithCredentials failed - %s", testName, err)
		}
	}

	fn(client)
	rts.Verify(t)
}

// MakeServers creates a Servers object for this RequestTestSpec. If MuxHandlers
// is set, that is used as the entirety of the basis for the Servers. Otherwise,
// makes a HandlerMap and calls MakeServers on that.
func (rts *RequestTestSpec) MakeServers(t *testing.T) Servers {
	if rts.MuxHandlers != nil {
		return rts.MuxHandlers.MakeServers(t)
	}
	return HandlerMap{
		rts.Endpoint: rts.Handler(t),
	}.MakeServers(t)
}

// Verify ensures that the request was visited at least once, as long as
// NoVerify is false and this RequestTestSpec is not the MuxHandlers style.
func (rts *RequestTestSpec) Verify(t *testing.T) {
	if rts.NoVerify {
		return
	}
	if rts.MuxHandlers == nil {
		if rts.visits == 0 {
			t.Error("never called the HTTP endpoint")
		}
	}
}

// Handler returns an http.Handler for the request(s) expected by this
// RequestTestSpec.
func (rts *RequestTestSpec) Handler(t *testing.T) http.Handler {
	return Mux{
		rts.URL: func(wr http.ResponseWriter, r *http.Request) {
			assert.Method(rts.Method)(t, "", r)
			rts.visits++
			if rts.Auth {
				assert.Auth(lib.TokenType(rts.Endpoint))(t, "", r)
			}
			if rts.AssertRequest != nil {
				rts.AssertRequest(t, "", r)
			}
			if rts.StatusCode != 0 {
				wr.WriteHeader(rts.StatusCode)
			}
			WriteJSON(t, wr, rts.Response)
		},
	}.Handler()
}
