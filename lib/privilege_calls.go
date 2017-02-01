package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// user is optional, will get all privileges you can see without it
func (c *bytemarkClient) GetPrivileges(user string) (privileges brain.Privileges, err error) {
	req, err := c.BuildRequest("GET", BrainEndpoint, "/users/%s/privileges", user)
	if err != nil {
		return
	}
	if user == "" {
		req, err = c.BuildRequest("GET", BrainEndpoint, "/privileges")
		if err != nil {
			return
		}
	}
	_, _, err = req.Run(nil, &privileges)
	return
}

func (c *bytemarkClient) GrantPrivilege(privilege brain.Privilege) (err error) {
	username := privilege.Username
	req, err := c.BuildRequest("POST", BrainEndpoint, "/users/%s/privileges", username)
	if err != nil {
		return
	}
	privilege.Username = ""

	js, err := json.Marshal(privilege)
	if err != nil {
		return
	}
	_, _, err = req.Run(bytes.NewBuffer(js), nil)
	return
}

func (c *bytemarkClient) RevokePrivilege(privilege brain.Privilege) (err error) {
	if privilege.ID == 0 {
		// ok annoying, we have to go find out the privilege's id first
		privs, err := c.GetPrivileges(privilege.Username)
		if err != nil {
			return err
		}
		index := privs.IndexOf(privilege)
		if index == -1 {
			return fmt.Errorf("No such privilege found - %s", privilege.String())
		}
		privilege.ID = privs[index].ID
	}

	req, err := c.BuildRequest("DELETE", BrainEndpoint, "/privileges/%s", fmt.Sprintf("%d", privilege.ID))
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, nil)
	return
}
