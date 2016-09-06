package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/bigv"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func getFixtureUser() (user *bigv.User) {
	return &bigv.User{
		Username: "test-user",
		Email:    "test@user.com",
		AuthorizedKeys: []string{
			"ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user",
			"ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am",
		},
	}
}

// TestGetUser tests the behaviour of GetUser in a success case, as well as when the brain returns a 404
func TestGetUser(t *testing.T) {
	is := is.New(t)
	client, auth, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/users/test-user" {
			str, err := json.Marshal(getFixtureUser())
			if err != nil {
				t.Fatal(err)
			}
			w.Write(str)
		} else if req.URL.Path == "/users/nonexistent-user" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}
	}), mkNilHandler(t))
	defer auth.Close()
	defer brain.Close()
	defer billing.Close()
	log.DebugLevel = 9
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	user, err := client.GetUser("nonexistent-user")
	is.Nil(user)
	is.NotNil(err)

	if _, ok := err.(NotFoundError); !ok {
		t.Fatalf("Expected NotFoundError, got %T", err)
	}

	user, err = client.GetUser("test-user")
	is.Nil(err)
	is.NotNil(user)

	is.Equal("test-user", user.Username)
	is.Equal("test@user.com", user.Email)
	if 2 != len(user.AuthorizedKeys) {
		t.Fatalf("User didn't have enough authorized keys - %d instead of 2", len(user.AuthorizedKeys))
	}
	is.Equal("ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user", user.AuthorizedKeys[0])
	is.Equal("ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am", user.AuthorizedKeys[1])
}
