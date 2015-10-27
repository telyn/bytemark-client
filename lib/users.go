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

func (bigv *bigvClient) GetUser(name string) (*User, error) {
	url := BuildURL("/users/%s", name)

	var user JSONUser
	var realUser User
	err := bigv.RequestAndUnmarshal(true, "GET", url, "", &user)
	user.Process(&realUser)
	return &realUser, err
}

// AddUserAuthorizedKey
func (bigv *bigvClient) AddUserAuthorizedKey(username string, key string) error {
	user, err := bigv.GetUser(username)
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
	_, err = bigv.RequestAndRead(true, "PUT", url, string(userjs))
	return err

}

// DeleteUserAuthorizedKey
func (bigv *bigvClient) DeleteUserAuthorizedKey(username string, key string) error {
	user, err := bigv.GetUser(username)
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	newKeys := make([]string, 0)
	for _, k := range user.AuthorizedKeys {
		if strings.TrimSpace(k) == strings.TrimSpace(key) {
			continue
		}
		newKeys = append(newKeys, k)
	}

	user.AuthorizedKeys = newKeys
	userjs, err := json.Marshal(user)
	if err != nil {
		return err
	}

	url := BuildURL("/users/%s", username)
	_, err = bigv.RequestAndRead(true, "PUT", url, string(userjs))
	return err

}
