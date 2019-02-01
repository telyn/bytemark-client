package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// CreateGroup sends a request to the API server to create a group with the given name.
func (c *bytemarkClient) CreateGroup(name pathers.GroupName) (err error) {
	err = c.EnsureGroupName(&name)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups", string(name.Account))
	if err != nil {
		return
	}

	obj := map[string]string{
		"name": name.Group,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

// DeleteGroup requests that a given group be deleted. Will return an error if there are VMs in the group.
func (c *bytemarkClient) DeleteGroup(name pathers.GroupName) (err error) {
	err = c.EnsureGroupName(&name)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s", string(name.Account), name.Group)
	if err != nil {
		return
	}
	_, _, err = r.Run(nil, nil)
	return
}

// GetGroup requests an overview of the group with the given name
func (c *bytemarkClient) GetGroup(name pathers.GroupName) (group brain.Group, err error) {
	err = c.EnsureGroupName(&name)
	if err != nil {
		return
	}

	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s?view=overview&include_deleted=true", string(name.Account), name.Group)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &group)
	if err != nil {
		return
	}
	return
}
