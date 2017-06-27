package lib

import (
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func Test400BadRequestError(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", 400)
			_, err := w.Write([]byte(`{"something": "is not allowed"}`))
			if err != nil {
				t.Fatal(err)
			}
		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadDefinitions()
	is.NotNil(err)

	brErr, ok := err.(BadRequestError)
	is.NotNil(brErr)
	is.True(ok)

	is.OK(len(brErr.Problems))
}

func Test401UnauthorizedError(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", 401)
		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadDefinitions()
	is.NotNil(err)

	brErr, ok := err.(UnauthorizedError)
	is.NotNil(brErr)
	is.True(ok)

}

func Test403ForbiddenError(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", 403)
		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadDefinitions()
	is.NotNil(err)

	brErr, ok := err.(ForbiddenError)
	is.NotNil(brErr)
	is.True(ok)

}

func Test500InternalServerError(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", 500)
		}),
	})

	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadDefinitions()
	is.NotNil(err)

	brErr, ok := err.(InternalServerError)
	is.NotNil(brErr)
	is.True(ok)

}

func Test503ServiceUnavailableError(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service Temporarily Unavailable", 503)
		}),
	})

	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadDefinitions()
	is.NotNil(err)

	brErr, ok := err.(ServiceUnavailableError)
	is.NotNil(brErr)
	is.True(ok)

}
