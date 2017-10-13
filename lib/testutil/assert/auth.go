package assert

import (
	"net/http"
	"strings"
	"testing"
)

// Auth asserts that the right authorization header was put on the request
func Auth(tokenType string) RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		auths := r.Header["Authorization"]
		if len(auths) == 0 {
			t.Errorf("%s had no auth header", testName)
		}
		switch tokenType {
		case "bearer":
			if !strings.HasPrefix(auths[0], "Bearer ") {
				t.Errorf("%s had the wrong kind of authorization header. Expecting a 'Bearer ' prefix but got %s", testName, auths[0])
			}
		case "token":
			if !strings.HasPrefix(auths[0], "Token token=") {
				t.Errorf("%s had the wrong kind of authorization header. Expecting a 'Token token=' prefix but got %s", testName, auths[0])
			}
		}
	}
}

// Unauthed asserts that the request had no authorization header
func Unauthed() RequestAssertFunc {
	return func(t *testing.T, testName string, r *http.Request) {
		auths := r.Header["Authorization"]
		if len(auths) > 0 {
			t.Errorf("%s request had authorization when it shouldn't have.", testName)
		}
	}
}
