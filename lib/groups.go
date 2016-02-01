package lib

import (
	"encoding/json"
)

// CreateGroup sends a request to the API server to create a group with the given name.
func (c *bytemarkClient) CreateGroup(name GroupName) error {
	err := c.validateGroupName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups", name.Account)

	obj := map[string]string{
		"name": name.Group,
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, _, err = c.Request(true, "POST", path, string(bytes))
	return err
}

// DeleteGroup requests that a given group be deleted. Will return an error if there are VMs in the group.
func (c *bytemarkClient) DeleteGroup(name GroupName) error {
	err := c.validateGroupName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s", name.Account, name.Group)

	_, err = c.RequestAndRead(true, "DELETE", path, "")
	return err
}

// GetGroup requests an overview of the group with the given name
func (c *bytemarkClient) GetGroup(name GroupName) (*Group, error) {
	err := c.validateGroupName(&name)
	if err != nil {
		return nil, err
	}
	group := new(Group)

	path := BuildURL("/accounts/%s/groups/%s?view=overview&include_deleted=true", name.Account, name.Group)

	err = c.RequestAndUnmarshal(true, "GET", path, "", group)
	if err != nil {
		return nil, err
	}
	return group, err
}
