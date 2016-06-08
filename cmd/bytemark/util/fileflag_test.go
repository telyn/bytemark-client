package util

import (
	"github.com/cheekybits/is"
	"io/ioutil"
	"testing"
)

func TestFileFlag(t *testing.T) {
	is := is.New(t)
	err := ioutil.WriteFile("test-fileflag", []byte("contents here yay!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test-fileflag")

	flag := FileFlag{}

	flag.Set("test-fileflag")
	is.Equal("contents here yay!", flag.String())
}
