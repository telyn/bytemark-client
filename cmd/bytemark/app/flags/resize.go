package flags

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
)

// ResizeMode represents whether to increment a size or just to set it.
type ResizeMode int

const (
	// ResizeModeSet will cause resize disk to set the disc size to the one specified
	ResizeModeSet ResizeMode = iota
	// ResizeModeIncrease will cause resize disk to increase the disc size by the one specified
	ResizeModeIncrease
)

// Resize is effectively an extension of SizeSpecFlag which has a ResizeMode. The Size stored in the flag is the size to set to or increase by depending on the Mode
type Resize struct {
	Mode ResizeMode
	Size int
}

// Set parses the string into a Resize. If it starts with +, Mode is set to ResizeModeIncrease. Otherwise, it's set to ResizeModeSet. The rest of the string is parsed as a sizespec using sizespec.Parse
func (rf *Resize) Set(value string) (err error) {
	rf.Mode = ResizeModeSet
	if strings.HasPrefix(value, "+") {
		rf.Mode = ResizeModeIncrease
		value = value[1:]
	}

	sz, err := sizespec.Parse(value)
	if err != nil {
		return
	}
	rf.Size = sz
	return
}

// String returns the size, in GiB or TiB (if the size is > 1TIB) with the unit used as a suffix. If Mode is ResizeModeIncrease, the string is prefixed with '+'
func (rf Resize) String() string {
	plus := ""
	if rf.Mode == ResizeModeIncrease {
		plus += "+"
	}
	sz := rf.Size
	units := "GiB"
	sz /= 1024
	if sz > 1024 {
		sz /= 1024
		units = "TiB"
	}
	return fmt.Sprintf("%s%d%s", plus, sz, units)
}
