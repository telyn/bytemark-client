package mocks_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestRequest(t *testing.T) {
	var testMap map[string]string
	r := mocks.Request{
		T:              t,
		ResponseObject: map[string]string{"hello": "hi"},
	}
	_, _, err := r.MarshalAndRun(nil, &testMap)
	if err != nil {
		t.Error(err)
	}
	if testMap["hello"] != "hi" {
		t.Errorf("Assignment failed")
	}
}
