package testutil

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// MultiRequestTestSpec is used for more complex operations that call out to
// multiple API endpoints.
// MultiRequestTestSpec will combine all the specs into one MuxHandlers,
// supporting multiple methods for the same endpoint and url path
type MultiRequestTestSpec struct {
	// Specs is the RequestTestSpecs that define requests, responses and
	// assertions. In order to work with MultiRequestTestSpec, the MuxHandlers
	// property of each Spec must be nil.
	Specs []RequestTestSpec
	// Auth is an override for all the RequestTestSpecs' Auth bools - set to
	// true to set all the RequestTestSpecs Auth to true, leave false to use the
	// RequestTestSpecs' Auth bools
	Auth bool
}

// setAuth sets all the Specs' Auth bools to true, if mrts.Auth == true
func (mrts MultiRequestTestSpec) setAuth() {
	if mrts.Auth {
		for i := range mrts.Specs {
			mrts.Specs[i].Auth = true
		}
	}
}

// Run performs the setup and for the test, calls fn, then verifies all the
// requests were requested (according to their NoVerify values)
func (mrts MultiRequestTestSpec) Run(t *testing.T, fn RequestTestFunc) {
	mrts.setAuth()
	client, servers, err := NewClientAndServers(t, mrts)
	defer servers.Close()
	if err != nil {
		t.Fatalf("NewClientAndServers failed: %s", err)
	}

	if mrts.Auth {
		err = client.AuthWithCredentials(map[string]string{})
		if err != nil {
			t.Fatalf("AuthWithCredentials failed: %s", err)
		}
	}

	fn(client)
	for _, spec := range mrts.Specs {
		spec.Verify(t)
	}
}

// MakeServers implements ServerFactory, so MultiRequestTestSpec can be fed to
// NewClientAndServers (as is done in MultiRequestTestSpec.Run)
func (mrts MultiRequestTestSpec) MakeServers(t *testing.T) Servers {
	hm := HandlerMap{}
	for _, spec := range mrts.Specs {
		if _, ok := hm[spec.Endpoint]; !ok {
			hm[spec.Endpoint] = mrts.Handler(t, spec.Endpoint)
		}
	}
	return hm.MakeServers(t)
}

// Handler returns an http.Handler for the specified endpoint - which will
// handle any incoming requests, falling back to NilHandler if no
// RequestTestSpec can be found matching the endpoint / request.
func (mrts MultiRequestTestSpec) Handler(t *testing.T, endpoint lib.Endpoint) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		for i := range mrts.Specs {
			// we need to be acting on a pointer to the spec rather than a copy
			// (which is what range gives us) since the Handler modifies the
			// spec.
			spec := &mrts.Specs[i]
			if spec.Endpoint != endpoint {
				continue
			}
			if spec.Method != r.Method {
				continue
			}
			if spec.URL != r.URL.Path {
				continue
			}
			spec.Handler(t).ServeHTTP(wr, r)
			return
		}
		NilHandler(t).ServeHTTP(wr, r)
	})
}
