package assert

import (
	"net/http"
	"testing"
)

// Method reads the request's method and checks that it's the same as expected
func Method(expected string) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		if r.Method != expected {
			t.Errorf("%s - Wrong HTTP method %s, expecting %s", testName, r.Method, expected)
		}
	}
}
