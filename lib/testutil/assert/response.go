package assert

import (
	"reflect"
	"testing"
)

// Response asserts that response deepequals expected
func Response(t *testing.T, testName string, response interface{}, expected interface{}) {
	if !reflect.DeepEqual(expected, response) {
		t.Errorf("%s - unexpected response.\nExpected: %#v\nActual: %#v", testName, expected, response)
	}
}
