package lib_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

// TestGetUser tests the behaviour of GetUser in a success case, as well as when the brain returns a 404
func TestGetUser(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/users/test-user",
		Response: testutil.FixtureUser,
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
		assert.Equal(t, testName, "ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user", user.AuthorizedKeys[0].String())
		assert.Equal(t, testName, "ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am", user.AuthorizedKeys[1].String())
	})
}
