package util

import (
	bigv "bytemark.co.uk/client/lib"
	"fmt"
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

// ParseDiscSpec takes a disc spec and returns a slice of Discs (from bytemark.co.uk/client/lib)
/*func ParseDiscSpec(spec string, trace bool) ([]bigv.Disc, error) {
	// parser!
	// this really needs to be rewritten with a lexer.
	pos := 0

	discs := make([]bigv.Disc, 0, 4)

	const (
		_either = 0
		_grade  = 1
		_size   = 2
	)

	curDisc := 0
	curGrade := ""
	curSize := ""

	sizePos := 0

	expecting := _either

	// parser!!
	for true {
		if pos >= len(spec) {
			break
		}
		c, w := utf8.DecodeRuneInString(spec[pos:])
		if c == ' ' {
			pos += w
			continue
		}
		switch expecting {
		case _grade:
			if c == ',' {
				if len(curSize) == 0 {
					return nil, &DiscSpecError{pos, c}
				}

				size, err := ParseSize(curSize)
				if err != nil {
					if ssErr, ok := err.(*SizeSpecError); ok {
						return nil, &DiscSpecError{ssErr.Position - sizePos, ssErr.Character}
					} else {
						return discs, err
					}
				}
				discs = append(discs, bigv.Disc{
					StorageGrade: curGrade,
					Size:         int(size),
				})

				curGrade = ""
				curSize = ""
				curDisc++
				expecting = _either
				pos += w
			} else if c >= 'a' && c <= 'z' {
				curGrade += spec[pos : pos+w]
				pos += w
			} else if c == ':' {
				expecting = _size
				pos += w
				sizePos = pos
			} else {
				return nil, &DiscSpecError{pos, c}
			}
		case _size:
			if c == ',' {
				size, err := ParseSize(curSize)
				if err != nil {
					if ssErr, ok := err.(*SizeSpecError); ok {
						return nil, &DiscSpecError{ssErr.Position - sizePos, ssErr.Character}
					} else {
						return discs, err
					}
				}
				discs = append(discs, bigv.Disc{
					StorageGrade: curGrade,
					Size:         int(size),
				})

				curGrade = ""
				curSize = ""
				curDisc++
				expecting = _either
				pos += w
			} else if c == ':' && curGrade == "" {
				expecting = _grade
			} else {
				curSize += spec[pos : pos+w]
				pos += w
			}
		default: // _either
			if c >= 'a' && c <= 'z' {
				expecting = _grade
			} else if c >= '0' && c <= '9' {
				expecting = _size
				sizePos = pos
			} else {
				return nil, &DiscSpecError{pos, c}
			}
		}
	}
	size, err := ParseSize(curSize)
	if err != nil {
		return nil, &DiscSpecError{Position: pos - len(curSize)}
	}
	discs = append(discs, bigv.Disc{
		StorageGrade: curGrade,
		Size:         int(size),
	})
	return discs, nil
}
*/
