// +build quality

package output_test

// this file gets used to generate a valid go source file which is
// used to test whether or not types that might get output have
// implement the PrettyPrinter and DefaultFieldsHaver interfaces

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"reflect"
	"testing"
)

var IgnoredTypes = []string {
}

func isIgnoredType(t string) bool {
	switch t {
	// the brain will never send us an IPSpec - we will never Output it.
	case "brain.IPSpec":
		return true
	// IPs are always part of a NetworkInterface
	case "brain.IPs":
		return true
	// JSONUser is only public because it needs to be so we can write JSON to it
	// no public method in lib.Client returns a JSONUser
	case "brain.JSONUser":
		return true
        // PrivilegeLevel is never gonna come out of the lib, except as part of a Privilege
	case "brain.PrivilegeLevel":
		return true
	// billing.Accounts should ALWAYS be dressed as lib.Accounts
	case "billing.Account":
		return true
	// billing.SPPTokenRequests will only be sent, never received.
	case "billing.SPPTokenRequest":
		return true
	// as if spp would return credit card details to us :]
	case "spp.CreditCard":
		return true
	}
	return false
}

func isArrayType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Ptr:
		return isArrayType(typ.Elem())
	case reflect.Slice,reflect.Array:
		return true
	case reflect.Invalid:
		fmt.Println("oops, type is invalid. this shouldn't happen")
		return false
	}
	return false
}

func TestQualityEveryTypeIsDefaultFieldsHaver(t *testing.T) {
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
		switch typ.Kind() {
		case reflect.Ptr:
			typ = typ.Elem()
		}
		orig := typ
		
		if isIgnoredType(typ.String()) {
			continue
		}

		if isArrayType(typ) {
			typ = typ.Elem()
			switch typ.Kind() {
			case reflect.Ptr:
				typ = typ.Elem()
			}

			if !typ.Implements(dfhType){
				t.Errorf("%%s is a slice type, but its elem (%%s) is not DefaultsFieldHaver. Implement the DefaultFields method.", orig, typ)
			}
			if !typ.Implements(ppType) {
				t.Errorf("%%s is a slice type, but its elem (%%s) is not PrettyPrinter. Implement the PrettyPrint method.", orig, typ)
			}
		} else {
			if !typ.Implements(dfhType) {
				t.Errorf("%%s is not a DefaultsFieldHaver. Implement the DefaultFields method.", typ)
			}
			if !typ.Implements(ppType) {
				t.Errorf("%%s is not a PrettyPrinter. Implement the PrettyPrint method.", typ)
			}
		}

	}
}

// vim: ft=go
