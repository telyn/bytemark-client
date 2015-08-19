package lib

import "fmt"

func labelDiscs(discs []*Disc) {
	for i, disc := range discs {
		if disc.Label == "" {
			disc.Label = fmt.Sprintf("%c", 'a'+i)
		}
	}

}

func generateDiscLabel(discIdx int) string {
	return fmt.Sprintf("vd%c", 'a'+discIdx)
}
