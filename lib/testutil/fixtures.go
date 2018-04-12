package testutil

import "github.com/BytemarkHosting/bytemark-client/lib/brain"

// FixtureUser is a simple User with a couple of keys.
// TODO(telyn): get rid of this in favour of making User objects in each test
var FixtureUser = brain.User{
	Username: "test-user",
	Email:    "test@user.com",
	AuthorizedKeys: []brain.Key{
		{Key: "ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user"},
		{Key: "ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am"},
	},
}
