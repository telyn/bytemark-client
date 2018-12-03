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
	// StatusCode is the status code that will be returned
	StatusCode int
	// AssertRequest is an optional func which will get called to check the
	// request object further - for example to check the URL has particular
	// query string keys
	AssertRequest func(t *testing.T, testName string, r *http.Request)

	// visits is how many times the generated handler has been called
	// TODO(telyn): maybe refactor this out to MuxHandlers so we get visits-tracking across all endpoints for custom MuxHandlers.
	visits int
}

// handlerFunc creates a http.HandlerFunc which validates the Method & ExpectedRequestBody
// then writes the Response
func (rts *RequestTestSpec) handlerFunc(t *testing.T, testName string, auth bool) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		// TODO(telyn): refactor the visits thing out to MuxHandlers
		rts.visits++

		assert.Method(rts.Method)(t, testName, r)
		if auth {
			assert.Auth(lib.TokenType(rts.Endpoint))(t, testName, r)
		}
		if rts.AssertRequest != nil {
			rts.AssertRequest(t, testName, r)
		}
		WriteJSON(t, wr, rts.Response)
	}
}

// mkMuxHandlers makes a MuxHandlers which contains one endpoint (specified by this RequestTestSpec)
// which validates the Method and ExpectedRequest, then writes the Response
func (rts *RequestTestSpec) mkMuxHandlers(t *testing.T, testName string, auth bool) (mh MuxHandlers, err error) {
	return NewMuxHandlers(rts.Endpoint, rts.URL, rts.handlerFunc(t, testName, auth))
}

// Run sets up fake servers, creates a client that talks to them, then passes the client to fn.
// fn should run some request method using the client & test the results of that function.
func (rts *RequestTestSpec) Run(t *testing.T, testName string, auth bool, fn RequestTestFunc) {
	if rts.MuxHandlers == nil {
		mh, err := rts.mkMuxHandlers(t, testName, auth)
		if err != nil {
			t.Fatalf("Couldn't create MuxHandlers - %s", err)
		}
		rts.MuxHandlers = &mh
	}

	client, servers, err := NewClientAndServers(t, rts.MuxHandlers)
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
	if rts.MuxHandlers == nil {
		if rts.visits == 0 {
			t.Errorf("%s never called the HTTP endpoint", testName)
		}
	}
}
