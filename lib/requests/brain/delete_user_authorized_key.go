package brain

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// DeleteUserAuthorizedKey removes a key from a user. The key may be specified in full or just the comment part (as long as it's unique)
func DeleteUserAuthorizedKey(client lib.Client, username string, key string) error {
	user, err := client.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	var newKeys []brain.Key
	potentiallyAmbiguous := false
	for _, k := range user.AuthorizedKeys {
		if strings.TrimSpace(k.String()) == key {
			continue
		} else {
			parts := strings.SplitN(k.String(), " ", 3)
			if len(parts) == 3 && strings.TrimSpace(parts[2]) == key {
				potentiallyAmbiguous = true
				continue
			}
		}
		newKeys = append(newKeys, k)
	}
	// if there's a difference of more than one then the key was ambiguous
	if len(newKeys) < len(user.AuthorizedKeys)-1 && potentiallyAmbiguous {
		err := lib.AmbiguousKeyError{}
		return err
	}

	user.AuthorizedKeys = newKeys

	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/users/%s", username)
	if err != nil {
		return err
	}
	_, _, err = r.MarshalAndRun(user, nil)
	return err

}
