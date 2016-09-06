package brain

import (
	"encoding/json"
	"strings"
)

// JSONUser is used as an intermediate type that gets processed into a User. It should not have been exported.
type JSONUser struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	AuthorizedKeys string `json:"authorized_keys"`

	// passwords are handled by auth these days
	//Password       string `json:"password"`

	// "users can be created (using POST) without authentication. If the
	// request has no authentication, it will also accept an account_name
	// parameter and create an account at the same time."
	// this is almost certainly never going to be useful
	//AccountName string `json:"account_name"`
}

// Process turns this JSONUser into a User.
func (jsonUser *JSONUser) Process(into *User) {
	into.Username = jsonUser.Username
	into.Email = jsonUser.Email
	into.AuthorizedKeys = strings.Split(jsonUser.AuthorizedKeys, "\n")
}

// User represents a Bytemark user.
type User struct {
	Username       string
	Email          string
	AuthorizedKeys []string
}

// MarshalJSON marshals the User into a JSON bytestream.
func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONUser{
		Username:       user.Username,
		Email:          user.Email,
		AuthorizedKeys: strings.Join(user.AuthorizedKeys, "\n"),
	})
}
