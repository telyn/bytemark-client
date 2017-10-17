package testutil

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// TestLogWriter is a writer that writes to the test log using testing.T.Log
// it's a bit naff since things like fmt.Printf and text/template.Execute call
// Write multiple times and each Write adds a new line to the log.
type TestLogWriter struct {
	t *testing.T
}

// Write writes the bytes as a string to t.Log
func (tlw TestLogWriter) Write(p []byte) (n int, err error) {
	tlw.t.Log(string(p))
	return len(p), nil
}

// OverrideLogWriters overrides the bytemark-client/util/log package's writers with TestLogWriters
func OverrideLogWriters(t *testing.T) {
	log.Writer = TestLogWriter{t}
	log.ErrWriter = TestLogWriter{t}
}
