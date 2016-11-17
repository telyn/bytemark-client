package lib

import (
	"bytes"
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"strings"
)

// Getbigv.User grabs the named user from the brain
func (c *bytemarkClient) GetUser(name string) (user *brain.User, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/users/%s", name)
	if err != nil {
		return
	}

	var jsUser brain.JSONUser
	_, _, err = r.Run(nil, &jsUser)
	if err != nil {
		return
	}
	user = new(brain.User)
	jsUser.Process(user)
	return
}

// Addbigv.UserAuthorizedKey adds a key to the named user. These keys are used for logging into the management IPs for cloud servers
func (c *bytemarkClient) AddUserAuthorizedKey(username string, key string) error {
	user, err := c.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	user.AuthorizedKeys = append(user.AuthorizedKeys, key)

	userjs, err := json.Marshal(user)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("PUT", BrainEndpoint, "/users/%s", username)
	if err != nil {
		return err
	}
	_, _, err = r.Run(bytes.NewBuffer(userjs), nil)
	return err

}

// Deletebigv.UserAuthorizedKey removes a key from a user. The key may be specified in full or just the comment part (as long as it's unique)
func (c *bytemarkClient) DeleteUserAuthorizedKey(username string, key string) error {
	user, err := c.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	var newKeys []string
	potentiallyAmbiguous := false
	for _, k := range user.AuthorizedKeys {
		if strings.TrimSpace(k) == key {
			continue
		} else {
			parts := strings.SplitN(k, " ", 3)
			if len(parts) == 3 && strings.TrimSpace(parts[2]) == key {
				potentiallyAmbiguous = true
				continue
			}
		}
		newKeys = append(newKeys, k)
	}
	// if there's a difference of more than one then the key was ambiguous
	if len(newKeys) < len(user.AuthorizedKeys)-1 && potentiallyAmbiguous {
		err := AmbiguousKeyError{}
		return err
	}

	user.AuthorizedKeys = newKeys
	userjs, err := json.Marshal(user)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("PUT", BrainEndpoint, "/users/%s", username)
	if err != nil {
		return err
	}
	_, _, err = r.Run(bytes.NewBuffer(userjs), nil)
	return err

}
