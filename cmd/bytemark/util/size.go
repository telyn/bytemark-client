package util

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// SizeSpecError indicates that the size spec could not be parsed.
type SizeSpecError struct {
	// Position is the index in Spec at which the the unexpected character was found.
	Position int
	// Character is the character that was unexpected.
	Character rune
	// Spec is the full size spec that caused the error.
	Spec string
}

func (e *SizeSpecError) Error() string {
	if e.Character == '\x00' {
		return "Size specification was empty."
	}
	return fmt.Sprintf("Invalid size '%s': unexpected '%c' at character %d.", e.Spec, e.Character, e.Position)
}

// ParseSize will take a size as a string like <num>[G|M][[i]B] and output a size.
// It's actually more permissive than that but w/e
func ParseSize(spec string) (int, error) {
	const (
		_num = iota
		_numGM
		_iB
		_B
	)
	pos := 0
	curSize := ""
	curMultiplier := 1024
	expecting := _num
	for {
		if pos >= len(spec) {
			break
		}
		c, w := utf8.DecodeRuneInString(spec[pos:])
		if c == ' ' {
			pos += w
			continue
		}
		switch expecting {
		case _num:
			if c >= '0' && c <= '9' {
				curSize += spec[pos : pos+w]
				expecting = _numGM
				pos += w
			} else {
				return -1, &SizeSpecError{pos, c, spec}
			}
		case _numGM:
			if c >= '0' && c <= '9' {
				curSize += spec[pos : pos+w]
				expecting = _numGM
				pos += w
			} else if c == 'm' || c == 'M' {
				curMultiplier = 1
				expecting = _iB
				pos += w
			} else if c == 'g' || c == 'G' {
				expecting = _iB
				pos += w
			} else {
				return -1, &SizeSpecError{pos, c, spec}
			}
		case _iB:
			if c == 'i' {
				expecting = _B
				pos += w
			} else if c == 'b' || c == 'B' {
				pos += w
			} else {
				return -1, &SizeSpecError{pos, c, spec}
			}
		case _B:
			if c == 'b' || c == 'B' {
				pos += w
			} else {
				return -1, &SizeSpecError{pos, c, spec}
			}

		}

	}
	size, err := strconv.ParseInt(curSize, 10, 32)
	if err != nil {
		return -1, &SizeSpecError{Position: 0}
	}
	return int(size) * curMultiplier, nil
}
