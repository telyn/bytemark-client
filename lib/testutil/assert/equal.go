package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, testName string, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s objects weren't the same.\nExpected: %#v\nActual: %#v", testName, expected, actual)
	}
}
