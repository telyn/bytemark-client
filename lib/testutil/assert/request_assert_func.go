package assert

import (
	"net/http"
	"testing"
)

type RequestAssertFunc func(t *testing.T, testName string, r *http.Request)
