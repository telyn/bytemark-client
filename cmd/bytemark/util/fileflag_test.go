package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cheekybits/is"
)

func TestFileFlag(t *testing.T) {
	is := is.New(t)
	err := ioutil.WriteFile("test-fileflag", []byte("contents here yay!"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	flag := FileFlag{}

	err = flag.Set("test-fileflag")
	if err != nil {
		t.Fatal(err)
	}
	is.Equal("contents here yay!", flag.String())
	_ = os.Remove("test-fileflag")
}
