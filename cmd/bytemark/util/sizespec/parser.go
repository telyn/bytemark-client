// Package sizespec implements a parser for size specifications.
// Yep, we could use regexes, but then we don't get to tell people where they failed.
package sizespec

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Error indicates that the size spec could not be parsed.
type Error struct {
	// Position is the index in Spec at which the the unexpected character was found.
	Position int
	// Character is the character that was unexpected.
	Character rune
	// Spec is the full size spec that caused the error.
	Spec string
}

func (e *Error) Error() string {
	if e.Character == '\x00' {
		return "Size specification was empty."
	}
	return fmt.Sprintf("Invalid size '%s': unexpected '%c' at character %d.", e.Spec, e.Character, e.Position)
}

type expectation int

const (
	// when we expect a number
	_num expectation = iota
	// when we expect a number or G or M or nothing
	_numGM
	// when we expect an i or a B or nothing
	_iB
	// when we expect a B or nothing
	_B
)

type parserState struct {
	pos        int
	buf        string
	multiplier int
	expecting  expectation
	spec       string
}

// create an error for our spec at the current position & given rune (too lazy to re-decode the rune)
func (st *parserState) err(c rune) error {
	return &Error{Position: st.pos, Character: c, Spec: st.spec}
}

// the rune and int are the next char and its width in bytes, and are args rather than part of state because they should not be edited by the function
type parserFn func(*parserState, rune, int) error

// parserFnMap is where the main meat of the parser is - it maps expectations to which function to try, to continue parsing.
var parserFnMap = map[expectation]parserFn{
	// at the beginning, we expect a digit.
	_num: func(st *parserState, c rune, w int) error {
		if c >= '0' && c <= '9' {
			st.buf += st.spec[st.pos : st.pos+w]
			st.expecting = _numGM
			return nil
		}

		return st.err(c)
	},
	// once we have a digit, we expect a digit or a G or an M
	_numGM: func(st *parserState, c rune, w int) (err error) {
		if c >= '0' && c <= '9' {
			st.buf += st.spec[st.pos : st.pos+w]
			st.expecting = _numGM
		} else if c == 'm' || c == 'M' {
			st.multiplier = 1
			st.expecting = _iB
		} else if c == 'g' || c == 'G' {
			st.expecting = _iB
		} else {
			return st.err(c)
		}
		return nil

	},
	// once we've had a G or an M, we expect the next thing to be an i or a B.
	_iB: func(st *parserState, c rune, w int) (err error) {
		if c == 'i' {
			// if we get an i, we expect the next to be a B.
			// There's nothing stopping the sizespec from ending here though, so weird constructions like 14Gi are valid according to this parser.
			st.expecting = _B
			return nil
		} else if c == 'b' || c == 'B' {
			return nil
		}
		return st.err(c)

	},
	// if we get an i, we expect the next to be a B.
	_B: func(st *parserState, c rune, w int) (err error) {
		if c == 'b' || c == 'B' {
			return nil
		}
		return st.err(c)

	},
}

// continueParse picks a function from parserFnMap to run, and runs it, then increments the state's position
func continueParse(st *parserState, c rune, w int) (err error) {
	err = parserFnMap[st.expecting](st, c, w)
	st.pos += w
	return
}

// Parse will take a size as a string like <num>[G|M][[i]B] and output an int - the size in megabytes.
// It's actually more permissive than that but w/e
func Parse(spec string) (size int, err error) {
	st := parserState{
		pos:        0,
		buf:        "",
		multiplier: 1024,
		expecting:  _num,
		spec:       spec,
	}
	for {
		// if we've read to the end, stop
		if st.pos >= len(st.spec) {
			break
		}
		// otherwise decode the next rune
		c, w := utf8.DecodeRuneInString(spec[st.pos:])
		// ignore whitespace. TODO(telyn): why?
		if c == ' ' {
			st.pos += w
			continue
		}
		// and continue parsing
		err = continueParse(&st, c, w)
		if err != nil {
			return
		}

	}

	size64, err := strconv.ParseInt(st.buf, 10, 32)
	size = int(size64) * st.multiplier
	return
}
