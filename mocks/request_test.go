package mocks_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestRequest(t *testing.T) {
	var testMap map[string]string
	r := mocks.Request{
		t:              t,
		ResponseObject: map[string]string{"hello": "hi"},
	}
	r.MarshalAndRun(nil, &testMap)
	if testMap["hello"] != "hi" {
		t.Errorf("Assignment failed")
	}
}
