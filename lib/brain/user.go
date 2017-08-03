package brain

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
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

// User turns this JSONUser into a User.
func (jsonUser *JSONUser) User() (user User) {
	user.Username = jsonUser.Username
	user.Email = jsonUser.Email
	user.AuthorizedKeys = strings.Split(jsonUser.AuthorizedKeys, "\n")
	return
}

// User represents a Bytemark user.
type User struct {
	Username       string
	Email          string
	AuthorizedKeys []string
}

func (user User) DefaultFields(f output.Format) string {
	return "Username, Email"
}

func (user User) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	userTpl := `
{{ define "user_sgl" }}{{ .Username }}{{ end }}
{{ define "user_medium" }}{{ .Username }} - {{ .Email }}{{ end }}
{{ define "user_full " }}{{ template "user_medium" }}

Authorized keys:
{{ for .AuthorizedKeys }}
{{ . }}	
{{ end }}
`
	return prettyprint.Run(wr, userTpl, "user"+string(detail), user)
}

// MarshalJSON marshals the User into a JSON bytestream.
func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONUser{
		Username:       user.Username,
		Email:          user.Email,
		AuthorizedKeys: strings.Join(user.AuthorizedKeys, "\n"),
	})
}
