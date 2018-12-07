package flags

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// genericSliceFlag allows any flag type to be converted into a slicey version
// that can be specified more than once. It supports the app.Preprocesser
// interface too, and passes that call on to each of its flags.
// The purpose is to ease implementing slice flags, particularly on
// app.Preprocessers.
//
// To use it you should be implementing a new flag type which embeds it and
// contains a function to return a slice of whatever underlying flag type it
// contains. See GroupNameSliceFlag for example.
type genericSliceFlag struct {
	// template is the default value for each flag, if you like. It MUST be set
	// and MUST NOT be nil.
	template flag.Value
	Values   []flag.Value
}

func (gsf *genericSliceFlag) Preprocess(c *app.Context) error {
	for _, value := range gsf.Values {
		if preprocesser, ok := value.(app.Preprocesser); ok {
			err = preprocesser.Preprocess(c)
			if err != nil {
				return
			}
		}
	}
	return
}

func (gsf *genericSliceFlag) Set(value string) error {
	if gsf.template == nil {
		return fmt.Errorf("Can't set flag - template is nil. This is a bug")
	}
	newValue := cloneFlagValue(gsf.template)
	err := newValue.Set(value)
	if err != nil {
		return err
	}
	gsf.Values = append(gsf.Values, newValue)
	return nil
}

func (gsf genericSliceFlag) String() string {
	strs := make([]string, len(gsf.Values))
	for i, value := range gsf.Values {
		strs[i] = value.String()
	}
	return strings.Join(strs, ", ")
}

func (gsf genericSliceFlag) copyValues(dst interface{}) {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Slice {
		panic("copyValues dst must be a slice")
	}
	tmplType := reflect.TypeOf(gsf.template)
	if !tmplType.AssignableTo(dstType.Elem()) {
		panic(fmt.Sprintf("copyValues dst must be something that a %T can be assigned to", gsf.template))
	}

	dstValue := reflect.ValueOf(dst)
	for _, val := range gsf.Values {
		reflect.Append(dstValue, reflect.ValueOf(val))
	}
}
