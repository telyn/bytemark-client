package flags

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// GroupNameFlag is used for all --group flags, including the global one.
type GroupNameFlag struct {
	GroupName *lib.GroupName
	Value     string
}

// Set runs lib.ParseGroupName to make sure we have a valid group name
func (name *GroupNameFlag) Set(value string) error {
	name.Value = value
	return nil
}

// Preprocess defaults the value of this flag to the default group from the
// config attached to the context and then runs lib.ParseGroupName
// This is an implementation of `app.Preprocessor`, which is detected and
// called automatically by actions created with `app.Action`
func (name *GroupNameFlag) Preprocess(c *app.Context) (err error) {
	if name.GroupName != nil {
		c.Debug("GroupName.Preprocess before %#v", *name.GroupName)
	}
	if name.Value == "" {
		return
	}
	groupName := lib.ParseGroupName(name.Value, c.Config().GetGroup())
	name.GroupName = &groupName
	c.Debug("GroupName.Preprocess after %#v", *name.GroupName)
	return
}

// String returns the GroupName as a string.
func (name GroupNameFlag) String() string {
	if name.GroupName != nil {
		return name.GroupName.String()
	}
	return ""
}
