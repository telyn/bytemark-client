package assert

import (
	"reflect"
	"testing"
)

// Equal asserts that expected and actual are equal according to reflect.DeepEqual
func Equal(t *testing.T, testName string, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s objects weren't the same.\nExpected: %#v\nActual:   %#v", testName, expected, actual)
	}
}

// NotEqual asserts that expected and actual are not equal according to reflect.DeepEqual
func NotEqual(t *testing.T, testName string, unexpected interface{}, actual interface{}) {
	if reflect.DeepEqual(unexpected, actual) {
		t.Errorf("%s objects were not supposed to be equal.\nUnexpected: %#v\nActual:     %#v", testName, unexpected, actual)
	}
}
