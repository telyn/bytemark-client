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

// GetGroup requests an overview of the group with the given name
func (bigv *bigvClient) GetGroup(name GroupName) (*Group, error) {
	group := new(Group)

	path := BuildURL("/accounts/%s/groups/%s?view=overview", name.Account, name.Group)

	err := bigv.RequestAndUnmarshal(true, "GET", path, "", group)
	return group, err
}
