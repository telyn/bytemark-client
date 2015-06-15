package main

import (
	bigv "bigv.io/client/lib"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// DiscSpecError represents an error during parse.
type DiscSpecError struct {
	Position  int
	Character rune
}

func (e *DiscSpecError) Error() string {
	return fmt.Sprintf("Disc spec error: Unexpected %s at %d.")
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
			if c >= '0' && c <= '9' {
				curSize += spec[pos : pos+w]
				pos += w
			} else if c == ',' {
				size, err := strconv.ParseInt(curSize, 10, 32)
				if err != nil {
					panic(err)
				}
				discs = append(discs, bigv.Disc{
					StorageGrade: curGrade,
					Size:         int(size) * 1024,
				})

				curGrade = ""
				curSize = ""
				curDisc++
				expecting = _either
				pos += w
			} else {
				return nil, &DiscSpecError{pos, c}
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
	size, err := strconv.ParseInt(curSize, 10, 32)
	if err != nil {
		panic(err)
	}
	discs = append(discs, bigv.Disc{
		StorageGrade: curGrade,
		Size:         int(size) * 1024,
	})
	return discs, nil
}
