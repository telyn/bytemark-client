package testutil_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/cheekybits/is"
)

func call(fn func() string) string {
	return fn()
}

func indirectName() string {
	return testutil.Name(309)
}

func TestName(t *testing.T) {
	is := is.New(t)
	is.Equal("TestName 0", testutil.Name(0))
	call(func() string {
		is.Equal("TestName 400", testutil.Name(400))
		return ""
	})
	is.Equal("TestName 309", indirectName())
	is.Equal("TestName 309", call(indirectName))

}
