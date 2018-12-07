package flags

import (
	"flag"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// GroupNameSliceFlag is for any --group flags that can be specified multiple
// times
type GroupNameSliceFlag interface {
	flag.Value
	app.Preprocesser
	GroupNames() []lib.GroupName
}

type groupNameSliceFlag struct {
	GenericSliceFlag
}

// NewGroupNameSliceFlag creates a new
func NewGroupNameSliceFlag() GroupNameSliceFlag {
	return groupNameSliceFlag{
		template: GroupNameFlag{},
	}
}

// GroupNames returns the GroupNameSlice for which this flag is named
func (gnsf GroupNameSliceFlag) GroupNames() (groupNames []lib.GroupName) {
	gnsf.copyValues(groupNames)
	return
}
