package assert

import (
	"net/http"
	"testing"
)

// QueryValue asserts that the requests' query string has key=value.
func QueryValue(key, expectedValue string) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		URLValue(t, testName, r.URL.Query(), key, expectedValue)
	}
}
