package flags

// This file was automatically generate using
// cmd/bytemark/app/flags/gen/slice_flags - do not edit it by hand!

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// AccountNameSliceFlag is used for AccountNameFlags that may be specified more than
// once. It's a slice of AccountNameFlag in order to avoid rewriting parsing
// logic.
type AccountNameSliceFlag []AccountNameFlag

func (sf *AccountNameSliceFlag) Preprocess(ctx *app.Context) error {
	for i := range *sf {
		err := (*sf)[i].Preprocess(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set appends a AccountNameFlag (created for you) to the slice
func (sf *AccountNameSliceFlag) Set(value string) error {
	flag := AccountNameFlag{}
	err := flag.Set(value)
	if err != nil {
		return err
	}
	*sf = append(*sf, flag)
	return nil
}

// String returns all values in the slice, comma-delimeted
func (sf AccountNameSliceFlag) String() string {
	strs := make([]string, len(sf))
	for i, value := range sf {
		strs[i] = value.String()
	}
	return strings.Join(strs, ", ")
}

// AccountNameSlice returns the named flag as a AccountNameSliceFlag,
// if it was one in the first place.
func AccountNameSlice(ctx *app.Context, name string) AccountNameSliceFlag {
	if sf, ok := ctx.Context.Generic(name).(AccountNameSliceFlag); ok {
		return sf
	}
	return AccountNameSliceFlag{}
}
