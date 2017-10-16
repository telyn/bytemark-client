package lib_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func getFixtureUser() (user *brain.User) {
	return &brain.User{
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
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/users/test-user",
		Response: getFixtureUser(),
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		user, err := client.GetUser("nonexistent-user")
		assert.NotEqual(t, testName, nil, err)

		if _, ok := err.(lib.NotFoundError); !ok {
			t.Fatalf("Expected NotFoundError, got %T", err)
		}

		user, err = client.GetUser("test-user")
		assert.Equal(t, testName, nil, err)
		assert.Equal(t, testName, "test-user", user.Username)
		assert.Equal(t, testName, "test@user.com", user.Email)
		assert.Equal(t, testName, 2, len(user.AuthorizedKeys))
		assert.Equal(t, testName, "ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user", user.AuthorizedKeys[0])
		assert.Equal(t, testName, "ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am", user.AuthorizedKeys[1])
	})
}
