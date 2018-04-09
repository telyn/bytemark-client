package testutil

import "github.com/BytemarkHosting/bytemark-client/lib/brain"

var FixtureUser = brain.User{
	Username: "test-user",
	Email:    "test@user.com",
	AuthorizedKeys: []brain.Key{
		{Key: "ssh-rsa AAAAAAAAAAAAAAAAAIm a scary test user"},
		{Key: "ssh-dsa AAAAAAAI even use DSA keys, that's how scary I am"},
	},
}
