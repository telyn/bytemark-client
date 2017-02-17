package lib

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetPrivileges gets all privileges for the given user (if you are that user or are cluster-admin)
// is user is blank, assumed to be your user
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

// GetPrivilegesForAccount gets all privileges lower than your privilege on the given account
func (c *bytemarkClient) GetPrivilegesForAccount(account string) (privileges brain.Privileges, err error) {
	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/privileges", account)
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, &privileges)
	return
}

// GetPrivilegesForGroup gets all privileges lower than your privilege on the given group
func (c *bytemarkClient) GetPrivilegesForGroup(group GroupName) (privileges brain.Privileges, err error) {
	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/privileges", group.Account, group.Group)
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, &privileges)
	return
}

// GetPrivilegesForVirtualMachine gets all privileges lower than your privilege on the given virtual machine
func (c *bytemarkClient) GetPrivilegesForVirtualMachine(vm VirtualMachineName) (privileges brain.Privileges, err error) {
	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/privileges", vm.Account, vm.Group, vm.VirtualMachine)
	if err != nil {
		return
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

	_, _, err = req.MarshalAndRun(privilege, nil)
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
