package lib_test

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

type errTestSpec struct {
	Handler   http.HandlerFunc
	Assertion func(err error)
}

func errorTest(t *testing.T, testName string, spec errTestSpec) {
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/definitions": spec.Handler,
			},
		},
	}
	rts.Run(t, testName, false, func(client lib.Client) {
		_, err := client.ReadDefinitions()
		assert.Equal(t, testName, false, err == nil)
		spec.Assertion(err)
	})
}

func Test400BadRequestError(t *testing.T) {
	testName := testutil.Name(0)
	errorTest(t, testName, errTestSpec{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			testutil.WriteJSON(t, w, map[string]string{"something": "is not allowed"})
		},
		Assertion: func(err error) {
			brErr, ok := err.(lib.BadRequestError)
			assert.Equal(t, testName, true, ok)

			assert.Equal(t, testName, true, 0 < len(brErr.Problems))
		},
	})
}

func Test401UnauthorizedError(t *testing.T) {
	testName := testutil.Name(0)
	errorTest(t, testName, errTestSpec{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		},
		Assertion: func(err error) {
			_, ok := err.(lib.UnauthorizedError)
			assert.Equal(t, testName, true, ok)
		},
	})
}

func Test403ForbiddenError(t *testing.T) {
	testName := testutil.Name(0)
	errorTest(t, testName, errTestSpec{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		},
		Assertion: func(err error) {
			_, ok := err.(lib.ForbiddenError)
			assert.Equal(t, testName, true, ok)
		},
	})
}

func Test500InternalServerError(t *testing.T) {
	testName := testutil.Name(0)
	errorTest(t, testName, errTestSpec{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		},
		Assertion: func(err error) {
			_, ok := err.(lib.InternalServerError)
			assert.Equal(t, testName, true, ok)
		},
	})
}

func Test503ServiceUnavailableError(t *testing.T) {
	testName := testutil.Name(0)
	errorTest(t, testName, errTestSpec{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		},
		Assertion: func(err error) {
			_, ok := err.(lib.ServiceUnavailableError)
			assert.Equal(t, testName, true, ok)
		},
	})
}
