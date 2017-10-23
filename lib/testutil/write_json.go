package testutil

import (
	"encoding/json"
	"io"
	"testing"
)

// WriteJSON marshals the object as JSON and writes it to the writer
func WriteJSON(t *testing.T, wr io.Writer, object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		t.Fatalf("Couldn't marshal. %s", err.Error())
	}
	_, err = wr.Write(js)
	if err != nil {
		t.Fatalf("Couldn't write JSON out: %s", err.Error())
	}
}
