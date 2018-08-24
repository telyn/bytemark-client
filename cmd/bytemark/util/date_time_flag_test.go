package util

import (
	"testing"
	"time"

	"github.com/cheekybits/is"
)

func TestDateTimeFlag(t *testing.T) {
	is := is.New(t)

	var flag DateTimeFlag

	err := flag.Set("1/2/2018 3pm")

	is.Nil(err)
	is.OK(flag.String())

	// check format
	formatted, err := time.Parse("2006-01-02T15:04:05-0700", flag.String())
	is.Nil(err)
	is.OK(formatted)
}
