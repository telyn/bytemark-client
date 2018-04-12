package assert

import (
	"net/url"
	"testing"
)

// URLValue asserts that key->expectedValue is a key-value pair in the urlValues.
func URLValue(t *testing.T, testName string, urlValues url.Values, key, expectedValue string) {
	values := urlValues[key]
	if len(values) == 0 {
		t.Errorf("%s %s parameter was not set", testName, key)
		return
	}

	if values[0] != expectedValue {
		t.Errorf("%s %s parameter was %q when it should've been %q", testName, key, values[0], expectedValue)
	}
}
