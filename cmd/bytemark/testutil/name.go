package testutil

import (
	ltu "github.com/BytemarkHosting/bytemark-client/lib/testutil"
)

// Name returns a sensible name for this test
func Name(iteration int) string {
	return ltu.Name(iteration)
}
