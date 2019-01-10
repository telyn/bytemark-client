package testutil

import "net/http"

// Mux is a map of URL paths to http.HandlerFuncs
// In general it's preferable to use a MultiRequestTestSpec instead, but we need
// to keep this type around for in the short term (used by
// RequestTestSpec.Handler)
type Mux map[string]http.HandlerFunc

// ToHandler turns the Mux into an http.ServeMux
func (m Mux) ToHandler() (serveMux *http.ServeMux) {

	serveMux = http.NewServeMux()
	for p, f := range m {
		serveMux.HandleFunc(p, f)
	}
	return
}
