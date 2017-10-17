package assert

import (
	"net/http"
	"testing"
)

// RequestAssertFuncs take a request object and make assertions
type RequestAssertFunc func(t *testing.T, testName string, r *http.Request)
