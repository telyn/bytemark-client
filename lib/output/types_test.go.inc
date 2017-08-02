package output_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"reflect"
	"testing"
)

func TestEveryTypeIsDefaultFieldsHaver(t *testing.T) {
	dfhType := reflect.TypeOf((*DefaultFieldsHaver)(nil))
	// p is for pointer
	for p := range []interface{}{
		(*lib.Account)(nil),
		%s
	} {
		if IsArrayType(p) {
			continue
		}
		typ := reflect.TypeOf(p)
		if !typ.Implements(dhfType) {
			t.Fail(fmt.Sprintf("%%s is not a DefaultsFieldHaver. Implement the PrettyPrint and DefaultFields methods.", typ.Name()))
		}

	}
}

// vim: ft=go
