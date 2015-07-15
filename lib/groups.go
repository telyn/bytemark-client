package lib

import (
	"encoding/json"
)

// CreateGroup sends a request to the BigV server to create a group with the given name.
func (bigv *bigvClient) CreateGroup(name GroupName) error {
	path := BuildURL("/accounts/%s/groups", name.Account)

	obj := map[string]string{
		"name": name.Group,
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, _, err = bigv.Request(true, "POST", path, string(bytes))
	return err
}

// DeleteGroup requests that a given group be deleted. Will return an error if there are VMs in the group.
func (bigv *bigvClient) DeleteGroup(name GroupName) error {
	path := BuildURL("/accounts/%s/groups/%s", name.Account, name.Group)

	_, err := bigv.RequestAndRead(true, "DELETE", path, "")
	return err
}

// GetGroup requests an overview of the group with the given name
func (bigv *bigvClient) GetGroup(name GroupName) (*Group, error) {
	group := new(Group)

	path := BuildURL("/accounts/%s/groups/%s?view=overview&include_deleted=true", name.Account, name.Group)

	err := bigv.RequestAndUnmarshal(true, "GET", path, "", group)
	return group, err
}
