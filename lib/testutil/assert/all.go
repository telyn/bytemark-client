package assert

import (
	"net/http"
	"testing"
)

// All runs all the assertions in sequence
func All(funcs ...RequestAssertFunc) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		for _, assertFunc := range funcs {
			assertFunc(t, testName, r)
		}
	}
}
