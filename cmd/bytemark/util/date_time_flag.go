package util

import (
	"github.com/bcampbell/fuzzytime"
)

// DateTimeFlag holds datatime in iso8601 format
type DateTimeFlag string

// Set takes user input and attempts to parse the datetime from any format to iso8601
func (dtf *DateTimeFlag) Set(value string) (err error) {
	// WesternContext will not raise and ambiguous error and expects dd/mm/yyyy
	dt, _, err := fuzzytime.WesternContext.Extract(value)

	if err != nil {
		return
	}

	*dtf = DateTimeFlag(dt.ISOFormat())
	return
}

func (dtf *DateTimeFlag) String() string {
	return string(*dtf)
}
