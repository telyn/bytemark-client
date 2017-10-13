package testutil

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

type TestLogWriter struct {
	t *testing.T
}

func (tlw TestLogWriter) Write(p []byte) (n int, err error) {
	tlw.t.Log(string(p))
	return len(p), nil
}

func OverrideLogWriters(t *testing.T) {
	log.Writer = TestLogWriter{t}
	log.ErrWriter = TestLogWriter{t}
}
