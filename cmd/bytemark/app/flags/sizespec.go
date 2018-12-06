package flags

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
)

// SizeSpecFlag represents a capacity as an integer number of megabytes.
type SizeSpecFlag int

// Set sets the value to the size specified. Users can add "M" or "G" as a suffix to specify that they are talking about megabytes/gigabytes. Gigabytes are assumed by default.
func (ssf *SizeSpecFlag) Set(spec string) error {
	sz, err := sizespec.Parse(spec)
	if err != nil {
		return err
	}
	*ssf = SizeSpecFlag(sz)
	return nil
}

func (ssf *SizeSpecFlag) String() string {
	// default value is 0, but is checked for in code that uses SizeSpecFlag and changed to 1
	if *ssf == 0 {
		return ""
	}
	return fmt.Sprintf("%d", *ssf)
}
