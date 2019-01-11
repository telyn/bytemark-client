package brain

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// AddUserAuthorizedKey adds a key to the named user. These keys are used for logging into the management IPs for cloud servers
func AddUserAuthorizedKey(client lib.Client, username string, key string) error {
	user, err := client.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	user.AuthorizedKeys = append(user.AuthorizedKeys, brain.Key{Key: key})

	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/users/%s", username)
	if err != nil {
		return err
	}
	_, _, err = r.MarshalAndRun(user, nil)
	return err

}
