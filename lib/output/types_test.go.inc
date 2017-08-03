package output_test

// this file gets used to generate a valid go source file which is
// used to test whether or not types that might get output have
// implement the PrettyPrinter and DefaultFieldsHaver interfaces

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
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
	// weird trick we have to do to get a reflect.Type for output.DefaultFieldsHaver and prettyprint.PrettyPrinter
	dfhType := reflect.TypeOf((*output.DefaultFieldsHaver)(nil)).Elem()
	ppType := reflect.TypeOf((*prettyprint.PrettyPrinter)(nil)).Elem()
	// p is for pointer
	for _, p := range []interface{}{
		(*lib.Account)(nil),
		// BEGIN GENERATED LIST OF TYPE POINTERS
		%s
		// END GENERATED LIST OF TYPE POINTERS
	} {
		typ := reflect.TypeOf(p)
		
		if IsArrayType(typ) {
			if !typ.Elem().Implements(dfhType) {
				t.Errorf("%%s is a slice type, but its elem is not DefaultsFieldHaver. Implement the DefaultFields method.", typ)
			}
			if !typ.Elem().Implements(ppType) {
				t.Errorf("%%s is a slice type, but its elem is not PrettyPrinter. Implement the PrettyPrint method.", typ)
			}
		}
		if !typ.Implements(dfhType) {
			t.Errorf("%%s is not a DefaultsFieldHaver. Implement the DefaultFields method.", typ)
		}
		if !typ.Implements(ppType) {
			t.Errorf("%%s is not a PrettyPrinter. Implement the PrettyPrint method.", typ)
		}

	}
}

// vim: ft=go
