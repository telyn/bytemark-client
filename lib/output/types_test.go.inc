package output_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"reflect"
	"testing"
)

func IsArrayType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Slice,reflect.Array:
		return true
	}
	return false
}

func TestEveryTypeIsDefaultFieldsHaver(t *testing.T) {
	dfhType := reflect.TypeOf((*output.DefaultFieldsHaver)(nil))
	// p is for pointer
	for p := range []interface{}{
		(*lib.Account)(nil),
		%s
	} {
		typ := reflect.TypeOf(p)
		if IsArrayType(typ) {
			if !typ.Elem().Implements(dfhType) {
				t.Errorf("%%s is a slice type, but its elem is not DefaultsFieldHaver. Implement the PrettyPrint and DefaultFields methods.", typ.Name())
			}
		}
		if !typ.Implements(dfhType) {
			t.Errorf("%%s is not a DefaultsFieldHaver. Implement the PrettyPrint and DefaultFields methods.", typ.Name())
		}

	}
}

// vim: ft=go
