package flags

// This file was automatically generate using
// cmd/bytemark/app/flags/gen/slice_flags - do not edit it by hand!

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// GroupNameSliceFlag is used for GroupNameFlags that may be specified more than
// once. It's a slice of GroupNameFlag in order to avoid rewriting parsing
// logic.
type GroupNameSliceFlag []GroupNameFlag

func (sf *GroupNameSliceFlag) Preprocess(ctx *app.Context) error {
	for i := range *sf {
		err := (*sf)[i].Preprocess(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set appends a GroupNameFlag (created for you) to the slice
func (sf *GroupNameSliceFlag) Set(value string) error {
	flag := GroupNameFlag{}
	err := flag.Set(value)
	if err != nil {
		return err
	}
	*sf = append(*sf, flag)
	return nil
}

// String returns all values in the slice, comma-delimeted
func (sf GroupNameSliceFlag) String() string {
	strs := make([]string, len(sf))
	for i, value := range sf {
		strs[i] = value.String()
	}
	return strings.Join(strs, ", ")
}

// GroupNameSlice returns the named flag as a GroupNameSliceFlag,
// if it was one in the first place.
func GroupNameSlice(ctx *app.Context, name string) GroupNameSliceFlag {
	if sf, ok := ctx.Context.Generic(name).(*GroupNameSliceFlag); ok {
		return *sf
	}
	return GroupNameSliceFlag{}
}
