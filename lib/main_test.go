package lib

import (
	"flag"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.DebugLevel = 99
	os.Exit(m.Run())
}
