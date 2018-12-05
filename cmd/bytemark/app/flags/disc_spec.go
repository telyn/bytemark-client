package flags

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// DiscSpecFlag is a flag which reads its argument as a disc spec. It can be specified multiple times to add multiple discs.
type DiscSpecFlag []brain.Disc

// Set adds all the defined discs to this flag's value
func (discsFlag *DiscSpecFlag) Set(value string) error {
	for _, val := range strings.Split(value, " ") {
		disc, err := util.ParseDiscSpec(val)
		if err != nil {
			return err
		}

		*discsFlag = append(*discsFlag, *disc)

	}
	return nil
}

func (discsFlag *DiscSpecFlag) String() string {
	var discs []string
	for _, d := range *discsFlag {
		if d.Label == "" {
			discs = append(discs, fmt.Sprintf("%s:%dGiB", d.StorageGrade, d.Size/1024))
		} else {
			discs = append(discs, fmt.Sprintf("%s:%s:%dGiB", d.Label, d.StorageGrade, d.Size/1024))
		}
	}
	return strings.Join(discs, ",")
}
