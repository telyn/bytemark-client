package prettyprint

import (
	"bytes"
	"testing"
)

// Test represents a test that can be used with RunPrettyPrintTest - it's not for general use outside of testing.
type Test struct {
	Object   PrettyPrinter
	Detail   DetailLevel
	Expected string
}

// RunTests runs the prettyprint tests provided. testName should be the name of the calling function (TestSomethingSomethingPrettyPrint normally)
func RunTests(t *testing.T, testName string, tests []Test) {
	var seenFull, seenMedium, seenSingleLine bool
	for i, test := range tests {
		var b bytes.Buffer

		err := test.Object.PrettyPrint(&b, test.Detail)
		if err != nil {
			t.Errorf("%s %d ERROR: %s", testName, i, err.Error())
		}
		str := b.String()
		if str != test.Expected {
			t.Errorf("%s %d FAIL: expected '%s', got '%s'", testName, i, test.Expected, str)
		}

		switch test.Detail {
		case Full:
			seenFull = true
		case Medium:
			seenMedium = true
		case SingleLine:
			seenSingleLine = true
		}
	}

	if !seenFull {
		t.Errorf("%s FAIL - didn't see a test with Detail: prettyprint.Full", testName)
	}
	if !seenMedium {
		t.Errorf("%s FAIL - didn't see a test with Detail: prettyprint.Medium", testName)
	}
	if !seenSingleLine {
		t.Errorf("%s FAIL - didn't see a test with Detail: prettyprint.Single", testName)
	}
}
