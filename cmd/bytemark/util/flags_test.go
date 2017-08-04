package util

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestIPFlag(t *testing.T) {
	is := is.New(t)

	flag := IPFlag{}
	// IPFlag.Set never returns an error
	_ = flag.Set("192.168.1.1")
	_ = flag.Set("2000::1")
	is.Equal("192.168.1.1", flag[0].String())
	is.Equal("192.168.1.1", flag[0].To4().String())
	is.Equal("2000::1", flag[1].String())
	is.NotEqual("2000::1", flag[1].To4().String())
}
