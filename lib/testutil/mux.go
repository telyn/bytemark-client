package testutil

import "net/http"

// Mux is a map of URL paths to http.HandlerFuncs
// In general it's preferable to use a MultiRequestTestSpec instead, but we need
// to keep this type around for in the short term (used by
// RequestTestSpec.Handler)
type Mux map[string]http.HandlerFunc

// Handler makes a http.Handler for this Mux
func (m Mux) Handler() http.Handler {
	serveMux := http.NewServeMux()
	for p, f := range m {
		serveMux.HandleFunc(p, f)
	}
	return serveMux
}
