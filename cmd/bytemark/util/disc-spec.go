package util

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/bigv"
	"strings"
)

// DiscSpecError represents an error during parse.
type DiscSpecError struct {
	Position  int
	Character rune
}

func (e *DiscSpecError) Error() string {
	return fmt.Sprintf("Disc spec error: Unexpected %c at character %d.", e.Character, e.Position)
}

// ParseDiscSpec reads the given string and attempts to interpret it as a disc spec.
func ParseDiscSpec(spec string) (*bigv.Disc, error) {
	bits := strings.Split(spec, ":")
	size, err := ParseSize(bits[len(bits)-1])
	if err != nil {
		return nil, err
	}
	switch {
	case len(bits) >= 4:
		return nil, &DiscSpecError{}
	case len(bits) == 3:
		return &bigv.Disc{Label: bits[0], StorageGrade: bits[1], Size: size}, nil
	case len(bits) == 2:
		return &bigv.Disc{StorageGrade: bits[0], Size: size}, nil
	case len(bits) == 1:
		return &bigv.Disc{Size: size}, nil
	case len(bits) == 0:
		return nil, &DiscSpecError{}
	}
	return nil, nil
}
