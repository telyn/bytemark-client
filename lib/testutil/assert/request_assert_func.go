package assert

import (
	"net/http"
	"testing"
)

// A RequestAssertFunc takes a request object and makes assertions
type RequestAssertFunc func(t *testing.T, testName string, r *http.Request)
