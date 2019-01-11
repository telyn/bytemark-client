package assert

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func cmpoptions() []cmp.Option {
	return []cmp.Option{
		cmp.AllowUnexported(big.Int{}),
	}
}

// Equal asserts that expected and actual are equal according to reflect.DeepEqual
func Equal(t *testing.T, testName string, expected interface{}, actual interface{}) {
	if !cmp.Equal(expected, actual, cmpoptions()...) {
		t.Errorf("%s objects weren't the same.\nExpected: %#v\nActual:   %#v\nDiff:\n%s", testName, expected, actual, cmp.Diff(expected, actual, cmpoptions()...))
	}
}

// NotEqual asserts that expected and actual are not equal according to reflect.DeepEqual
func NotEqual(t *testing.T, testName string, unexpected interface{}, actual interface{}) {
	if cmp.Equal(unexpected, actual, cmpoptions()...) {
		t.Errorf("%s objects were not supposed to be equal.\nUnexpected: %#v\n    Actual:     %#v", testName, unexpected, actual)
	}
}
