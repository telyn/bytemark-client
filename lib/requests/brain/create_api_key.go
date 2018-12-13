package brain

import (
	"errors"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateAPIKey creates an API key for the given user, then returns it.
// Neither the ID nor APIKey field should be specified in the spec.
// username may be blank if the spec.UserID is set.
func CreateAPIKey(client lib.Client, username string, spec brain.APIKey) (apiKey brain.APIKey, err error) {
	if len(spec.Privileges) > 0 {
		err = errors.New("CreateAPIKey does not support creating an APIKey with privileges already attached - please use GrantPrivilege after CreateAPIKey, and provide an APIKey with Privileges unset to CreateAPIKey")
		return
	}
	if spec.UserID != 0 && username != "" {
		err = errors.New("only specify one of username and spec.UserID")
		return
	}
	if spec.UserID == 0 && username == "" {
		err = errors.New("one of user and spec.UserID must be specified")
		return
	}
	if spec.UserID == 0 {
		var user brain.User
		user, err = client.GetUser(username)
		if err != nil {
			return
		}
		spec.UserID = user.ID
	}
	r, err := client.BuildRequest("POST", lib.BrainEndpoint, "/api_keys")
	if err != nil {
		return
	}
	_, _, err = r.MarshalAndRun(spec, &apiKey)
	return
}
