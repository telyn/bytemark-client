package assert

import (
	"net/url"
	"testing"
)

func URLValue(t *testing.T, testName string, urlValues url.Values, key, expectedValue string) {
	values := urlValues[key]
	if values == nil || len(values) == 0 {
		t.Errorf("%s %s parameter was not set", testName, key)
		return
	}

	if values[0] != expectedValue {
		t.Errorf("%s %s parameter was %q when it should've been %q", testName, key, values[0], expectedValue)
	}
}
