package lib

import (
	"flag"
	"os"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.DebugLevel = 99
	os.Exit(m.Run())
}
