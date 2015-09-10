package util

import (
	bigv "bigv.io/client/lib"
	"fmt"
	"unicode/utf8"
)

// DiscSpecError represents an error during parse.
type DiscSpecError struct {
	Position  int
	Character rune
}

func (e *DiscSpecError) Error() string {
	return fmt.Sprintf("Disc spec error: Unexpected %c at character %d.", e.Character, e.Position)
}

// ParseDiscSpec takes a disc spec and returns a slice of Discs (from bigv.io/client/lib)
func ParseDiscSpec(spec string, trace bool) ([]bigv.Disc, error) {
	// parser!
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
			if c >= 'a' && c <= 'z' {
				curGrade += spec[pos : pos+w]
				pos += w
			} else if c == ':' {
				expecting = _size
				pos += w
			} else {
				return nil, &DiscSpecError{pos, c}
			}
		case _size:
			if c == ',' {
				size, err := ParseSize(curSize)
				if err != nil {
					return discs, err
					// this should logically be impossible - curSize should be a string solely containing characters from 0-9
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
			} else {
				curSize += spec[pos : pos+w]
				pos += w
			}
		default: // _expecting
			if c >= 'a' && c <= 'z' {
				expecting = _grade
			} else if c >= '0' && c <= '9' {
				expecting = _size
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
