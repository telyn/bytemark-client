package testutil

import (
	"net/http"
	"testing"
)

// NilHandler creates an http.Handler that fails and ends the test when called.
// It's the default for MuxHandlers and Handlers - so that you only need specify the endpoints you expect to talk to.
func NilHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Unexpected request to a nil server\r\n%s %s", r.Method, r.URL.String())
	})
}
