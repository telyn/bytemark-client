package lib

import (
	"encoding/json"
	"strings"
)

func (jsonUser *JSONUser) Process(into *User) {
	into.Username = jsonUser.Username
	into.Email = jsonUser.Email
	into.AuthorizedKeys = strings.Split(jsonUser.AuthorizedKeys, "\n")
}

func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONUser{
		Username:       user.Username,
		Email:          user.Email,
		AuthorizedKeys: strings.Join(user.AuthorizedKeys, "\n"),
	})
}

func (c *bytemarkClient) GetUser(name string) (*User, error) {
	url := BuildURL("/users/%s", name)

	var user JSONUser
	var realUser User
	err := c.RequestAndUnmarshal(true, "GET", url, "", &user)
	user.Process(&realUser)
	return &realUser, err
}

// AddUserAuthorizedKey
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

	url := BuildURL("/users/%s", username)
	_, err = c.RequestAndRead(true, "PUT", url, string(userjs))
	return err

}

// DeleteUserAuthorizedKey removes a key from a user. The key may be specified in full or just the comment part
func (c *bytemarkClient) DeleteUserAuthorizedKey(username string, key string) error {
	user, err := c.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	newKeys := make([]string, 0)
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

	url := BuildURL("/users/%s", username)
	_, err = c.RequestAndRead(true, "PUT", url, string(userjs))
	return err

}
